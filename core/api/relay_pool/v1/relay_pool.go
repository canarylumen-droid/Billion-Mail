package v1

import (
        "billionmail-core/utility/types/api_v1"
        "github.com/gogf/gf/v2/frame/g"
)

// RelayProvider represents one API slot in the provider pool.
type RelayProvider struct {
        ID               int64  `json:"id"`
        Name             string `json:"name"`
        ProviderType     string `json:"provider_type"`
        Slot             int    `json:"slot"`
        APIKey           string `json:"api_key"`   // masked on list
        APIKey2          string `json:"api_key_2"` // masked on list
        DailyLimit       int    `json:"daily_limit"`
        MonthlyLimit     int    `json:"monthly_limit"`
        DailySent        int    `json:"daily_sent"`
        MonthlySent      int    `json:"monthly_sent"`
        Priority         int    `json:"priority"`
        Status           string `json:"status"`
        IsActive         bool   `json:"is_active"`
        LastError        string `json:"last_error"`
        LastUsedAt       int64  `json:"last_used_at"`
}

// ProviderMeta provides metadata about each provider type (logo, limits, etc.)
type ProviderMeta struct {
        Type          string `json:"type"`
        Label         string `json:"label"`
        DailyFree     int    `json:"daily_free"`
        MonthlyFree   int    `json:"monthly_free"`
        SignupURL      string `json:"signup_url"`
        APIKeyLabel   string `json:"api_key_label"`
        APIKey2Label  string `json:"api_key_2_label"`
        DocsURL       string `json:"docs_url"`
}

// Aggregate stats across all 30 slots.
type PoolStats struct {
        TotalMonthlyCapacity int `json:"total_monthly_capacity"`
        TotalMonthlySent     int `json:"total_monthly_sent"`
        TotalDailyCapacity   int `json:"total_daily_capacity"`
        TotalDailySent       int `json:"total_daily_sent"`
        ActiveSlots          int `json:"active_slots"`
        ExhaustedSlots       int `json:"exhausted_slots"`
        ErrorSlots           int `json:"error_slots"`
        UnconfiguredSlots    int `json:"unconfigured_slots"`
}

// ListProvidersReq
type ListProvidersReq struct {
        g.Meta        `path:"/relay_pool/list" method:"get" tags:"Relay Pool" summary:"List all 30 provider slots"`
        Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
}
type ListProvidersRes struct {
        api_v1.StandardRes
        Data struct {
                Providers []RelayProvider `json:"providers"`
                Stats     PoolStats       `json:"stats"`
                Meta      []ProviderMeta  `json:"meta"`
        } `json:"data"`
}

// UpdateProviderReq — update API key and settings for one slot.
type UpdateProviderReq struct {
        g.Meta        `path:"/relay_pool/update" method:"post" tags:"Relay Pool" summary:"Update a provider slot (set API key, enable/disable)"`
        Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
        ID            int64  `json:"id" v:"required|min:1"`
        APIKey        string `json:"api_key"`
        APIKey2       string `json:"api_key_2"`
        IsActive      *bool  `json:"is_active"`
        DailyLimit    int    `json:"daily_limit"`
        MonthlyLimit  int    `json:"monthly_limit"`
        Priority      int    `json:"priority"`
}
type UpdateProviderRes struct {
        api_v1.StandardRes
}

// TestProviderReq — send a test email via a specific provider slot.
type TestProviderReq struct {
        g.Meta        `path:"/relay_pool/test" method:"post" tags:"Relay Pool" summary:"Test a specific provider slot by sending a test email"`
        Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
        ID            int64  `json:"id" v:"required|min:1"`
        TestToEmail   string `json:"test_to_email" v:"required|email"`
}
type TestProviderRes struct {
        api_v1.StandardRes
        Data struct {
                OK                bool   `json:"ok"`
                ProviderName      string `json:"provider_name"`
                ProviderMessageID string `json:"provider_message_id"`
                LatencyMs         int64  `json:"latency_ms"`
                Error             string `json:"error,omitempty"`
        } `json:"data"`
}

// TestPoolReq — send via the pool (auto-selects best provider).
type TestPoolReq struct {
        g.Meta        `path:"/relay_pool/test_pool" method:"post" tags:"Relay Pool" summary:"Send a test email via the pool (auto-selects provider)"`
        Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
        TestToEmail   string `json:"test_to_email" v:"required|email"`
}
type TestPoolRes struct {
        api_v1.StandardRes
        Data struct {
                OK                bool   `json:"ok"`
                ProviderName      string `json:"provider_name"`
                ProviderMessageID string `json:"provider_message_id"`
                LatencyMs         int64  `json:"latency_ms"`
                Error             string `json:"error,omitempty"`
        } `json:"data"`
}

// ResetCountersReq — admin: reset daily/monthly sent counters.
type ResetCountersReq struct {
        g.Meta        `path:"/relay_pool/reset_counters" method:"post" tags:"Relay Pool" summary:"Reset daily/monthly sent counters for all or one provider"`
        Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
        ID            int64  `json:"id"` // 0 = reset all
        ResetType     string `json:"reset_type" v:"in:daily,monthly,both" d:"daily"`
}
type ResetCountersRes struct {
        api_v1.StandardRes
}

// GetPoolConfigReq — returns the global pool config (auto-route toggle etc.).
type GetPoolConfigReq struct {
        g.Meta        `path:"/relay_pool/config" method:"get" tags:"Relay Pool" summary:"Get relay pool global config"`
        Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
}
type GetPoolConfigRes struct {
        api_v1.StandardRes
        Data struct {
                AutoRouteEnabled bool `json:"auto_route_enabled"`
        } `json:"data"`
}

// UpdatePoolConfigReq — toggle auto-route and other global pool settings.
type UpdatePoolConfigReq struct {
        g.Meta           `path:"/relay_pool/config/update" method:"post" tags:"Relay Pool" summary:"Update relay pool global config"`
        Authorization    string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
        AutoRouteEnabled *bool  `json:"auto_route_enabled"`
}
type UpdatePoolConfigRes struct {
        api_v1.StandardRes
}

// CreateProviderReq — add a custom (non-seeded) provider slot.
type CreateProviderReq struct {
        g.Meta        `path:"/relay_pool/create" method:"post" tags:"Relay Pool" summary:"Add a custom provider slot"`
        Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
        Name         string `json:"name" v:"required"`
        ProviderType string `json:"provider_type" v:"required"`
        Slot         int    `json:"slot" v:"required|min:1" d:"1"`
        APIKey       string `json:"api_key" v:"required"`
        APIKey2      string `json:"api_key_2"`
        DailyLimit   int    `json:"daily_limit" d:"100"`
        MonthlyLimit int    `json:"monthly_limit" d:"1000"`
        Priority     int    `json:"priority" d:"50"`
        IsActive     bool   `json:"is_active" d:"true"`
}
type CreateProviderRes struct {
        api_v1.StandardRes
        Data struct {
                ID int64 `json:"id"`
        } `json:"data"`
}

// DeleteProviderReq — delete a custom provider slot (seeded slots are reset instead).
type DeleteProviderReq struct {
        g.Meta        `path:"/relay_pool/delete" method:"post" tags:"Relay Pool" summary:"Delete a custom provider slot (or clear API key from a seeded slot)"`
        Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
        ID            int64  `json:"id" v:"required|min:1"`
}
type DeleteProviderRes struct {
        api_v1.StandardRes
}
