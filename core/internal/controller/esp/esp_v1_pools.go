package esp

import (
	"context"
	"encoding/json"
	"time"

	v1 "billionmail-core/api/esp/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) CreateIpPool(ctx context.Context, req *v1.CreateIpPoolReq) (res *v1.CreateIpPoolRes, err error) {
	res = &v1.CreateIpPoolRes{}

	ips := req.Ips
	if ips == nil {
		ips = []string{}
	}
	ipsJSON, _ := json.Marshal(ips)
	now := time.Now().Unix()
	dailyLimit := req.DailyLimit
	if dailyLimit <= 0 {
		dailyLimit = 500
	}

	result, err := g.DB().Model("bm_ip_pools").Data(g.Map{
		"name":        req.Name,
		"pool_type":   req.PoolType,
		"ips":         string(ipsJSON),
		"status":      "warming",
		"warm_day":    0,
		"daily_limit": dailyLimit,
		"sent_today":  0,
		"description": req.Description,
		"created_at":  now,
		"updated_at":  now,
	}).Insert()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}

	id, _ := result.LastInsertId()
	res.Data = &v1.IpPool{
		Id:          id,
		Name:        req.Name,
		PoolType:    req.PoolType,
		Ips:         ips,
		Status:      "warming",
		WarmDay:     0,
		DailyLimit:  dailyLimit,
		SentToday:   0,
		Description: req.Description,
		CreatedAt:   now,
	}
	return
}

func (c *ControllerV1) ListIpPools(ctx context.Context, req *v1.ListIpPoolsReq) (res *v1.ListIpPoolsRes, err error) {
	res = &v1.ListIpPoolsRes{}

	var rows []struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		PoolType    string `json:"pool_type"`
		Ips         string `json:"ips"`
		Status      string `json:"status"`
		WarmDay     int    `json:"warm_day"`
		DailyLimit  int    `json:"daily_limit"`
		SentToday   int    `json:"sent_today"`
		Description string `json:"description"`
		CreatedAt   int64  `json:"created_at"`
	}
	err = g.DB().Model("bm_ip_pools").Scan(&rows)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}

	res.Data = make([]*v1.IpPool, 0, len(rows))
	for _, row := range rows {
		var ips []string
		_ = json.Unmarshal([]byte(row.Ips), &ips)
		if ips == nil {
			ips = []string{}
		}
		res.Data = append(res.Data, &v1.IpPool{
			Id:          row.Id,
			Name:        row.Name,
			PoolType:    row.PoolType,
			Ips:         ips,
			Status:      row.Status,
			WarmDay:     row.WarmDay,
			DailyLimit:  row.DailyLimit,
			SentToday:   row.SentToday,
			Description: row.Description,
			CreatedAt:   row.CreatedAt,
		})
	}
	return
}

func (c *ControllerV1) DeleteIpPool(ctx context.Context, req *v1.DeleteIpPoolReq) (res *v1.DeleteIpPoolRes, err error) {
	res = &v1.DeleteIpPoolRes{}
	_, err = g.DB().Model("bm_ip_pools").Where("id", req.Id).Delete()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	return
}
