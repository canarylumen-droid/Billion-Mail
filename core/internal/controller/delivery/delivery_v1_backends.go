package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	v1 "billionmail-core/api/delivery/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) ListBackends(ctx context.Context, req *v1.ListBackendsReq) (res *v1.ListBackendsRes, err error) {
	res = &v1.ListBackendsRes{}

	var rows []struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		BackendType string `json:"backend_type"`
		Config      string `json:"config"`
		Status      string `json:"status"`
		Priority    int    `json:"priority"`
		DailySent   int    `json:"daily_sent"`
		DailyLimit  int    `json:"daily_limit"`
		IsDefault   bool   `json:"is_default"`
		LastTested  int64  `json:"last_tested"`
		LastError   string `json:"last_error"`
		CreatedAt   int64  `json:"created_at"`
	}

	err = g.DB().Model("bm_delivery_backends").Order("priority DESC, id ASC").Scan(&rows)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}

	res.Data = make([]*v1.DeliveryBackend, 0, len(rows))
	for _, row := range rows {
		var cfg map[string]interface{}
		_ = json.Unmarshal([]byte(row.Config), &cfg)
		if cfg == nil {
			cfg = map[string]interface{}{}
		}
		// Mask secrets in config
		if secret, ok := cfg["secret"].(string); ok && len(secret) > 6 {
			cfg["secret"] = secret[:3] + "****" + secret[len(secret)-3:]
		}
		res.Data = append(res.Data, &v1.DeliveryBackend{
			Id:          row.Id,
			Name:        row.Name,
			BackendType: row.BackendType,
			Config:      cfg,
			Status:      row.Status,
			Priority:    row.Priority,
			DailySent:   row.DailySent,
			DailyLimit:  row.DailyLimit,
			IsDefault:   row.IsDefault,
			LastTested:  row.LastTested,
			LastError:   row.LastError,
			CreatedAt:   row.CreatedAt,
		})
	}
	return
}

func (c *ControllerV1) CreateBackend(ctx context.Context, req *v1.CreateBackendReq) (res *v1.CreateBackendRes, err error) {
	res = &v1.CreateBackendRes{}

	cfg := req.Config
	if cfg == nil {
		cfg = map[string]interface{}{}
	}
	cfgJSON, _ := json.Marshal(cfg)
	now := time.Now().Unix()
	dailyLimit := req.DailyLimit
	if dailyLimit <= 0 {
		dailyLimit = 100000
	}

	// If set as default, clear other defaults first
	if req.IsDefault {
		_, _ = g.DB().Model("bm_delivery_backends").Data(g.Map{"is_default": false}).Where("1=1").Update()
	}

	result, err := g.DB().Model("bm_delivery_backends").Data(g.Map{
		"name":         req.Name,
		"backend_type": req.BackendType,
		"config":       string(cfgJSON),
		"status":       "active",
		"priority":     0,
		"daily_sent":   0,
		"daily_limit":  dailyLimit,
		"is_default":   req.IsDefault,
		"created_at":   now,
		"updated_at":   now,
	}).Insert()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}

	id, _ := result.LastInsertId()
	res.Data = &v1.DeliveryBackend{
		Id:          id,
		Name:        req.Name,
		BackendType: req.BackendType,
		Config:      cfg,
		Status:      "active",
		DailyLimit:  dailyLimit,
		IsDefault:   req.IsDefault,
		CreatedAt:   now,
	}
	return
}

func (c *ControllerV1) UpdateBackend(ctx context.Context, req *v1.UpdateBackendReq) (res *v1.UpdateBackendRes, err error) {
	res = &v1.UpdateBackendRes{}

	data := g.Map{"updated_at": time.Now().Unix()}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.Config != nil {
		cfgJSON, _ := json.Marshal(req.Config)
		data["config"] = string(cfgJSON)
	}
	if req.DailyLimit > 0 {
		data["daily_limit"] = req.DailyLimit
	}
	if req.Status != "" {
		data["status"] = req.Status
	}
	if req.IsDefault {
		_, _ = g.DB().Model("bm_delivery_backends").Data(g.Map{"is_default": false}).Where("1=1").Update()
		data["is_default"] = true
	}

	_, err = g.DB().Model("bm_delivery_backends").Where("id", req.Id).Data(data).Update()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	return
}

func (c *ControllerV1) DeleteBackend(ctx context.Context, req *v1.DeleteBackendReq) (res *v1.DeleteBackendRes, err error) {
	res = &v1.DeleteBackendRes{}
	_, err = g.DB().Model("bm_delivery_backends").Where("id", req.Id).Delete()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	return
}

func (c *ControllerV1) TestBackend(ctx context.Context, req *v1.TestBackendReq) (res *v1.TestBackendRes, err error) {
	res = &v1.TestBackendRes{}
	res.Data.BackendType = req.BackendType

	start := time.Now()

	switch req.BackendType {
	case "cloudflare_worker":
		workerURL, _ := req.Config["worker_url"].(string)
		secret, _ := req.Config["secret"].(string)

		if workerURL == "" {
			res.Data.Ok = false
			res.Data.Message = "worker_url is required"
			return
		}

		// Send a test ping to the worker
		testPayload := map[string]interface{}{
			"ping": true,
			"to":   "test@billionmail.internal",
			"from": "noreply@billionmail.internal",
			"subject": "BillionMail connectivity test",
			"text": "This is a connectivity test from BillionMail.",
		}
		body, _ := json.Marshal(testPayload)

		httpReq, _ := http.NewRequestWithContext(ctx, "POST", workerURL, bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		if secret != "" {
			httpReq.Header.Set("X-BillionMail-Secret", secret)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, httpErr := client.Do(httpReq)
		res.Data.Latency = time.Since(start).Milliseconds()

		if httpErr != nil {
			res.Data.Ok = false
			res.Data.Message = fmt.Sprintf("Connection failed: %s", httpErr.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 401 {
			res.Data.Ok = false
			res.Data.Message = "Authentication failed — check your secret matches the BILLIONMAIL_SECRET environment variable in your Cloudflare Worker"
			return
		}
		if resp.StatusCode >= 500 {
			res.Data.Ok = false
			res.Data.Message = fmt.Sprintf("Worker returned error %d — check that the EMAIL binding is configured in your Worker settings", resp.StatusCode)
			return
		}

		res.Data.Ok = true
		res.Data.Message = fmt.Sprintf("Worker reachable in %dms. Ready to deliver email through Cloudflare's IP network.", res.Data.Latency)

	case "haraka", "oracle_vm":
		host, _ := req.Config["host"].(string)
		portFloat, _ := req.Config["port"].(float64)
		port := int(portFloat)
		if port == 0 {
			port = 587
		}

		if host == "" {
			res.Data.Ok = false
			res.Data.Message = "host is required"
			return
		}

		// TCP connectivity check on submission port
		conn, dialErr := (&http.Transport{}).Clone().DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", host, port))
		res.Data.Latency = time.Since(start).Milliseconds()

		if dialErr != nil {
			res.Data.Ok = false
			res.Data.Message = fmt.Sprintf("Cannot connect to %s:%d — %s. Ensure Haraka is running and port %d is open.", host, port, dialErr.Error(), port)
			return
		}
		conn.Close()

		res.Data.Ok = true
		res.Data.Message = fmt.Sprintf("Connected to %s:%d in %dms. SMTP server is reachable.", host, port, res.Data.Latency)

	default:
		res.Data.Ok = false
		res.Data.Message = "Unknown backend type: " + req.BackendType
	}

	return
}

func (c *ControllerV1) SetDefaultBackend(ctx context.Context, req *v1.SetDefaultBackendReq) (res *v1.SetDefaultBackendRes, err error) {
	res = &v1.SetDefaultBackendRes{}

	// Clear all defaults, then set the selected one
	_, err = g.DB().Model("bm_delivery_backends").Data(g.Map{"is_default": false}).Where("1=1").Update()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}
	_, err = g.DB().Model("bm_delivery_backends").Where("id", req.Id).Data(g.Map{
		"is_default": true,
		"updated_at": time.Now().Unix(),
	}).Update()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	return
}
