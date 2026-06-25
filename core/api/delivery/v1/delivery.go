package v1

import (
	"billionmail-core/utility/types/api_v1"
	"github.com/gogf/gf/v2/frame/g"
)

type DeliveryBackend struct {
	Id          int64                  `json:"id"`
	Name        string                 `json:"name"`
	BackendType string                 `json:"backend_type"`
	Config      map[string]interface{} `json:"config"`
	Status      string                 `json:"status"`
	Priority    int                    `json:"priority"`
	DailySent   int                    `json:"daily_sent"`
	DailyLimit  int                    `json:"daily_limit"`
	IsDefault   bool                   `json:"is_default"`
	LastTested  int64                  `json:"last_tested"`
	LastError   string                 `json:"last_error"`
	CreatedAt   int64                  `json:"created_at"`
}

// List
type ListBackendsReq struct {
	g.Meta        `path:"/delivery/backends/list" method:"get" tags:"Delivery" summary:"List delivery backends"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
}
type ListBackendsRes struct {
	api_v1.StandardRes
	Data []*DeliveryBackend `json:"data"`
}

// Create
type CreateBackendReq struct {
	g.Meta        `path:"/delivery/backends/create" method:"post" tags:"Delivery" summary:"Create delivery backend"`
	Authorization string                 `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Name          string                 `json:"name" v:"required|max-length:100"`
	BackendType   string                 `json:"backend_type" v:"required|in:cloudflare_worker,haraka,oracle_vm,smtp_relay"`
	Config        map[string]interface{} `json:"config"`
	DailyLimit    int                    `json:"daily_limit" d:"100000"`
	IsDefault     bool                   `json:"is_default"`
}
type CreateBackendRes struct {
	api_v1.StandardRes
	Data *DeliveryBackend `json:"data"`
}

// Update
type UpdateBackendReq struct {
	g.Meta        `path:"/delivery/backends/update" method:"post" tags:"Delivery" summary:"Update delivery backend"`
	Authorization string                 `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64                  `json:"id" v:"required|min:1"`
	Name          string                 `json:"name" v:"max-length:100"`
	Config        map[string]interface{} `json:"config"`
	DailyLimit    int                    `json:"daily_limit"`
	IsDefault     bool                   `json:"is_default"`
	Status        string                 `json:"status" v:"in:active,disabled"`
}
type UpdateBackendRes struct {
	api_v1.StandardRes
}

// Delete
type DeleteBackendReq struct {
	g.Meta        `path:"/delivery/backends/delete" method:"post" tags:"Delivery" summary:"Delete delivery backend"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type DeleteBackendRes struct {
	api_v1.StandardRes
}

// Test connection
type TestBackendReq struct {
	g.Meta        `path:"/delivery/backends/test" method:"post" tags:"Delivery" summary:"Test delivery backend connection"`
	Authorization string                 `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	BackendType   string                 `json:"backend_type" v:"required"`
	Config        map[string]interface{} `json:"config"`
}
type TestBackendRes struct {
	api_v1.StandardRes
	Data struct {
		Ok          bool   `json:"ok"`
		Latency     int64  `json:"latency_ms"`
		Message     string `json:"message"`
		BackendType string `json:"backend_type"`
	} `json:"data"`
}

// Set default
type SetDefaultBackendReq struct {
	g.Meta        `path:"/delivery/backends/set_default" method:"post" tags:"Delivery" summary:"Set default delivery backend"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type SetDefaultBackendRes struct {
	api_v1.StandardRes
}
