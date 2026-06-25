package esp

import (
	"context"
	"strings"
	"time"

	v1 "billionmail-core/api/esp/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) ListSuppression(ctx context.Context, req *v1.ListSuppressionReq) (res *v1.ListSuppressionRes, err error) {
	res = &v1.ListSuppressionRes{}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	m := g.DB().Model("bm_suppression_list")
	if req.Reason != "" {
		m = m.Where("reason", req.Reason)
	}
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		m = m.Where("email LIKE ? OR domain LIKE ?", kw, kw)
	}

	count, err := m.Count()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}

	var list []*v1.SuppressionEntry
	err = m.Order("created_at DESC").Page(page, pageSize).Scan(&list)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}
	if list == nil {
		list = []*v1.SuppressionEntry{}
	}
	res.Data.Total = count
	res.Data.List = list
	return
}

func (c *ControllerV1) AddSuppression(ctx context.Context, req *v1.AddSuppressionReq) (res *v1.AddSuppressionRes, err error) {
	res = &v1.AddSuppressionRes{}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	domain := ""
	if parts := strings.Split(email, "@"); len(parts) == 2 {
		domain = parts[1]
	}
	reason := req.Reason
	if reason == "" {
		reason = "manual"
	}

	_, err = g.DB().Model("bm_suppression_list").Data(g.Map{
		"email":      email,
		"reason":     reason,
		"domain":     domain,
		"created_at": time.Now().Unix(),
	}).Insert()
	if err != nil {
		res.Code = 1
		res.Message = "Already exists or DB error: " + err.Error()
	}
	return
}

func (c *ControllerV1) DeleteSuppression(ctx context.Context, req *v1.DeleteSuppressionReq) (res *v1.DeleteSuppressionRes, err error) {
	res = &v1.DeleteSuppressionRes{}
	_, err = g.DB().Model("bm_suppression_list").Where("id", req.Id).Delete()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	return
}
