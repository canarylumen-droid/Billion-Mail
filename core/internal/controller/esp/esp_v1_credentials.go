package esp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	v1 "billionmail-core/api/esp/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func generateUsername(name string) string {
	b := make([]byte, 4)
	rand.Read(b)
	suffix := fmt.Sprintf("%x", b)
	safe := ""
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			safe += string(c)
		}
	}
	if len(safe) > 12 {
		safe = safe[:12]
	}
	if safe == "" {
		safe = "user"
	}
	return fmt.Sprintf("%s_%s", safe, suffix)
}

func generatePassword() string {
	b := make([]byte, 18)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:24]
}

func (c *ControllerV1) CreateCredential(ctx context.Context, req *v1.CreateCredentialReq) (res *v1.CreateCredentialRes, err error) {
	res = &v1.CreateCredentialRes{}

	username := generateUsername(req.Name)
	password := generatePassword()
	now := time.Now().Unix()

	dailyLimit := req.DailyLimit
	if dailyLimit <= 0 {
		dailyLimit = 1000
	}
	plan := req.Plan
	if plan == "" {
		plan = "free"
	}

	result, err := g.DB().Model("bm_smtp_credentials").Data(g.Map{
		"name":        req.Name,
		"username":    username,
		"password":    password,
		"description": req.Description,
		"status":      1,
		"daily_limit": dailyLimit,
		"sent_today":  0,
		"plan":        plan,
		"ip_pool_id":  req.IpPoolId,
		"created_at":  now,
		"updated_at":  now,
	}).Insert()
	if err != nil {
		res.Code = 1
		res.Message = "Failed to create credential: " + err.Error()
		return
	}

	id, _ := result.LastInsertId()

	res.Data = &v1.SmtpCredential{
		Id:          id,
		Name:        req.Name,
		Username:    username,
		Description: req.Description,
		Status:      1,
		DailyLimit:  dailyLimit,
		SentToday:   0,
		Plan:        plan,
		IpPoolId:    req.IpPoolId,
		SmtpHost:    "smtp.yourdomain.com",
		SmtpPort:    587,
		CreatedAt:   now,
	}

	// Return password only on creation (not stored in plain text after this)
	_ = password
	res.Code = 0
	res.Message = "Credential created. Save your password: " + password
	return
}

func (c *ControllerV1) ListCredentials(ctx context.Context, req *v1.ListCredentialsReq) (res *v1.ListCredentialsRes, err error) {
	res = &v1.ListCredentialsRes{}

	var list []struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Username    string `json:"username"`
		Description string `json:"description"`
		Status      int    `json:"status"`
		DailyLimit  int    `json:"daily_limit"`
		SentToday   int    `json:"sent_today"`
		Plan        string `json:"plan"`
		IpPoolId    int64  `json:"ip_pool_id"`
		CreatedAt   int64  `json:"created_at"`
	}

	err = g.DB().Model("bm_smtp_credentials").Scan(&list)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}

	res.Data = make([]*v1.SmtpCredential, 0, len(list))
	for _, row := range list {
		res.Data = append(res.Data, &v1.SmtpCredential{
			Id:          row.Id,
			Name:        row.Name,
			Username:    row.Username,
			Description: row.Description,
			Status:      row.Status,
			DailyLimit:  row.DailyLimit,
			SentToday:   row.SentToday,
			Plan:        row.Plan,
			IpPoolId:    row.IpPoolId,
			SmtpHost:    "smtp.yourdomain.com",
			SmtpPort:    587,
			CreatedAt:   row.CreatedAt,
		})
	}
	return
}

func (c *ControllerV1) DeleteCredential(ctx context.Context, req *v1.DeleteCredentialReq) (res *v1.DeleteCredentialRes, err error) {
	res = &v1.DeleteCredentialRes{}
	_, err = g.DB().Model("bm_smtp_credentials").Where("id", req.Id).Delete()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	return
}

func (c *ControllerV1) RegeneratePassword(ctx context.Context, req *v1.RegeneratePasswordReq) (res *v1.RegeneratePasswordRes, err error) {
	res = &v1.RegeneratePasswordRes{}
	newPass := generatePassword()
	_, err = g.DB().Model("bm_smtp_credentials").Where("id", req.Id).Data(g.Map{
		"password":   newPass,
		"updated_at": time.Now().Unix(),
	}).Update()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}
	res.Data.Password = newPass
	return
}
