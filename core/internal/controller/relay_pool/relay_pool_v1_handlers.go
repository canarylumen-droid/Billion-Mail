package relay_pool_ctrl

import (
        "context"
        "strings"
        "time"

        v1 "billionmail-core/api/relay_pool/v1"
        "billionmail-core/internal/service/relay_pool"
        "github.com/gogf/gf/v2/frame/g"
)

// providerMeta defines metadata for all 15 providers shown in the UI.
var providerMeta = []v1.ProviderMeta{
        {Type: "resend", Label: "Resend", DailyFree: 100, MonthlyFree: 3000, SignupURL: "https://resend.com", APIKeyLabel: "API Key (re_...)", DocsURL: "https://resend.com/docs/api-reference/emails/send-email"},
        {Type: "sendgrid", Label: "SendGrid", DailyFree: 100, MonthlyFree: 3000, SignupURL: "https://signup.sendgrid.com", APIKeyLabel: "API Key (SG....)", DocsURL: "https://docs.sendgrid.com/api-reference/mail-send/mail-send"},
        {Type: "brevo", Label: "Brevo", DailyFree: 300, MonthlyFree: 9000, SignupURL: "https://app.brevo.com/account/register", APIKeyLabel: "API Key", DocsURL: "https://developers.brevo.com/reference/sendtransacemail"},
        {Type: "mailjet", Label: "Mailjet", DailyFree: 200, MonthlyFree: 6000, SignupURL: "https://app.mailjet.com/signup", APIKeyLabel: "API Key (Public)", APIKey2Label: "Secret Key", DocsURL: "https://dev.mailjet.com/email/guides/send-api-v31/"},
        {Type: "mailersend", Label: "Mailersend", DailyFree: 100, MonthlyFree: 3000, SignupURL: "https://app.mailersend.com/register", APIKeyLabel: "API Token", DocsURL: "https://developers.mailersend.com/api/v1/email.html"},
        {Type: "sendpulse", Label: "Sendpulse", DailyFree: 500, MonthlyFree: 15000, SignupURL: "https://sendpulse.com/en/signup", APIKeyLabel: "Client ID", APIKey2Label: "Client Secret", DocsURL: "https://sendpulse.com/integrations/api/smtp"},
        {Type: "zeptomail", Label: "ZeptoMail", DailyFree: 333, MonthlyFree: 10000, SignupURL: "https://www.zoho.com/zeptomail/signup.html", APIKeyLabel: "Send Mail Token", DocsURL: "https://www.zoho.com/zeptomail/help/api/email-sending.html"},
        {Type: "smtp2go", Label: "SMTP2GO", DailyFree: 33, MonthlyFree: 1000, SignupURL: "https://www.smtp2go.com/pricing/", APIKeyLabel: "API Key", DocsURL: "https://apidoc.smtp2go.com/documentation/#/POST%20/email/send"},
        {Type: "elasticemail", Label: "Elastic Email", DailyFree: 100, MonthlyFree: 3000, SignupURL: "https://app.elasticemail.com/register", APIKeyLabel: "API Key", DocsURL: "https://elasticemail.com/developers/api-documentation/rest-api"},
        {Type: "mailgun", Label: "Mailgun", DailyFree: 33, MonthlyFree: 1000, SignupURL: "https://signup.mailgun.com", APIKeyLabel: "API Key", APIKey2Label: "Sending Domain (e.g. mg.yourdomain.com)", DocsURL: "https://documentation.mailgun.com/en/latest/api-sending.html"},
        {Type: "postmark", Label: "Postmark", DailyFree: 4, MonthlyFree: 100, SignupURL: "https://account.postmarkapp.com/sign_up", APIKeyLabel: "Server API Token", DocsURL: "https://postmarkapp.com/developer/api/email-api"},
        {Type: "sparkpost", Label: "SparkPost", DailyFree: 17, MonthlyFree: 500, SignupURL: "https://app.sparkpost.com/join", APIKeyLabel: "API Key", DocsURL: "https://developers.sparkpost.com/api/transmissions/"},
        {Type: "socketlabs", Label: "SocketLabs", DailyFree: 1333, MonthlyFree: 40000, SignupURL: "https://www.socketlabs.com/signup/", APIKeyLabel: "API Key", APIKey2Label: "Server ID", DocsURL: "https://www.socketlabs.com/docs/inject/"},
        {Type: "netcore", Label: "Netcore", DailyFree: 100, MonthlyFree: 3000, SignupURL: "https://netcorecloud.com/email-api/", APIKeyLabel: "API Key", DocsURL: "https://support.netcorecloud.com/support/solutions/articles/43000551726"},
        {Type: "loops", Label: "Loops", DailyFree: 67, MonthlyFree: 2000, SignupURL: "https://loops.so", APIKeyLabel: "API Key", APIKey2Label: "Transactional Template ID", DocsURL: "https://loops.so/docs/api-reference/send-transactional-email"},
}

func (c *ControllerV1) ListProviders(ctx context.Context, req *v1.ListProvidersReq) (res *v1.ListProvidersRes, err error) {
        res = &v1.ListProvidersRes{}

        var rows []struct {
                ID               int64  `json:"id"`
                Name             string `json:"name"`
                ProviderType     string `json:"provider_type"`
                Slot             int    `json:"slot"`
                APIKey           string `json:"api_key"`
                APIKey2          string `json:"api_key_2"`
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

        if err = g.DB().Model("bm_relay_providers").Order("priority DESC, slot ASC, id ASC").Scan(&rows); err != nil {
                res.Code = 1
                res.Message = err.Error()
                return
        }

        stats := v1.PoolStats{}
        providers := make([]v1.RelayProvider, 0, len(rows))
        for _, r := range rows {
                // Mask API keys
                maskedKey := maskKey(r.APIKey)
                maskedKey2 := maskKey(r.APIKey2)

                providers = append(providers, v1.RelayProvider{
                        ID: r.ID, Name: r.Name, ProviderType: r.ProviderType, Slot: r.Slot,
                        APIKey: maskedKey, APIKey2: maskedKey2,
                        DailyLimit: r.DailyLimit, MonthlyLimit: r.MonthlyLimit,
                        DailySent: r.DailySent, MonthlySent: r.MonthlySent,
                        Priority: r.Priority, Status: r.Status, IsActive: r.IsActive,
                        LastError: r.LastError, LastUsedAt: r.LastUsedAt,
                })

                if r.IsActive && r.APIKey != "" {
                        stats.ActiveSlots++
                        stats.TotalDailyCapacity += r.DailyLimit
                        stats.TotalMonthlyCapacity += r.MonthlyLimit
                        stats.TotalDailySent += r.DailySent
                        stats.TotalMonthlySent += r.MonthlySent
                }
                switch r.Status {
                case "exhausted":
                        stats.ExhaustedSlots++
                case "error":
                        stats.ErrorSlots++
                case "unconfigured":
                        stats.UnconfiguredSlots++
                }
        }

        res.Data.Providers = providers
        res.Data.Stats = stats
        res.Data.Meta = providerMeta
        return
}

func (c *ControllerV1) UpdateProvider(ctx context.Context, req *v1.UpdateProviderReq) (res *v1.UpdateProviderRes, err error) {
        res = &v1.UpdateProviderRes{}

        data := g.Map{"updated_at": time.Now().Unix()}
        if req.APIKey != "" {
                data["api_key"] = req.APIKey
                data["status"] = "active"
        }
        if req.APIKey2 != "" {
                data["api_key_2"] = req.APIKey2
        }
        if req.IsActive != nil {
                data["is_active"] = *req.IsActive
                if !*req.IsActive {
                        data["status"] = "unconfigured"
                }
        }
        if req.DailyLimit > 0 {
                data["daily_limit"] = req.DailyLimit
        }
        if req.MonthlyLimit > 0 {
                data["monthly_limit"] = req.MonthlyLimit
        }
        if req.Priority > 0 {
                data["priority"] = req.Priority
        }

        if _, err = g.DB().Model("bm_relay_providers").Where("id", req.ID).Data(data).Update(); err != nil {
                res.Code = 1
                res.Message = err.Error()
        }
        return
}

func (c *ControllerV1) TestProvider(ctx context.Context, req *v1.TestProviderReq) (res *v1.TestProviderRes, err error) {
        res = &v1.TestProviderRes{}

        var row struct {
                ProviderType string `json:"provider_type"`
                APIKey       string `json:"api_key"`
                APIKey2      string `json:"api_key_2"`
                Name         string `json:"name"`
        }
        if err = g.DB().Model("bm_relay_providers").Where("id", req.ID).Scan(&row); err != nil || row.ProviderType == "" {
                res.Code = 1
                res.Message = "Provider not found"
                return
        }
        if row.APIKey == "" {
                res.Code = 1
                res.Message = "API key not configured for this slot"
                return
        }

        impl, ok := relay_pool.GetProvider(row.ProviderType)
        if !ok {
                res.Code = 1
                res.Message = "Provider type not supported: " + row.ProviderType
                return
        }

        msg := &relay_pool.EmailMessage{
                FromEmail: "test@billionmail.app",
                FromName:  "BillionMail Test",
                ToEmail:   req.TestToEmail,
                ToName:    "Test Recipient",
                Subject:   "BillionMail Provider Test — " + row.Name,
                HTML:      "<h2>✅ Provider test successful</h2><p>This test email was sent via <strong>" + row.Name + "</strong> from BillionMail's relay pool.</p>",
                Text:      "Provider test successful. Sent via " + row.Name + " from BillionMail's relay pool.",
        }

        result := impl.Send(ctx, row.APIKey, row.APIKey2, msg)
        if result.Err != nil {
                res.Data.OK = false
                res.Data.Error = result.Err.Error()
                _, _ = g.DB().Model("bm_relay_providers").Where("id", req.ID).Data(g.Map{
                        "status":     "error",
                        "last_error": result.Err.Error(),
                        "updated_at": time.Now().Unix(),
                }).Update()
        } else {
                res.Data.OK = true
                res.Data.ProviderMessageID = result.ProviderMessageID
                _, _ = g.DB().Model("bm_relay_providers").Where("id", req.ID).Data(g.Map{
                        "status":       "active",
                        "last_error":   "",
                        "last_used_at": time.Now().Unix(),
                        "updated_at":   time.Now().Unix(),
                }).Update()
        }
        res.Data.ProviderName = row.Name
        res.Data.LatencyMs = result.LatencyMs
        return
}

func (c *ControllerV1) TestPool(ctx context.Context, req *v1.TestPoolReq) (res *v1.TestPoolRes, err error) {
        res = &v1.TestPoolRes{}

        msg := &relay_pool.EmailMessage{
                FromEmail: "test@billionmail.app",
                FromName:  "BillionMail Pool Test",
                ToEmail:   req.TestToEmail,
                ToName:    "Test Recipient",
                Subject:   "BillionMail Pool Test — Auto-routed",
                HTML:      "<h2>✅ Pool routing test successful</h2><p>This test email was auto-routed through your relay provider pool.</p>",
                Text:      "Pool routing test successful. Auto-routed through your relay provider pool.",
        }

        result, poolErr := relay_pool.SendViaPool(ctx, msg)
        if poolErr != nil {
                res.Data.OK = false
                res.Data.Error = poolErr.Error()
        } else {
                res.Data.OK = true
                res.Data.ProviderName = result.ProviderName
                res.Data.ProviderMessageID = result.ProviderMessageID
                res.Data.LatencyMs = result.LatencyMs
        }
        return
}

func (c *ControllerV1) ResetCounters(ctx context.Context, req *v1.ResetCountersReq) (res *v1.ResetCountersRes, err error) {
        res = &v1.ResetCountersRes{}

        data := g.Map{"updated_at": time.Now().Unix()}
        resetType := req.ResetType
        if resetType == "" {
                resetType = "daily"
        }

        if resetType == "daily" || resetType == "both" {
                data["daily_sent"] = 0
                data["last_daily_reset"] = time.Now().Unix()
        }
        if resetType == "monthly" || resetType == "both" {
                data["monthly_sent"] = 0
                data["last_monthly_reset"] = time.Now().Unix()
        }

        var query = g.DB().Model("bm_relay_providers")
        if req.ID > 0 {
                query = query.Where("id", req.ID)
        } else {
                query = query.Where("1=1")
        }

        if _, err = query.Data(data).Update(); err != nil {
                res.Code = 1
                res.Message = err.Error()
        }
        return
}

func (c *ControllerV1) GetPoolConfig(ctx context.Context, req *v1.GetPoolConfigReq) (res *v1.GetPoolConfigRes, err error) {
        res = &v1.GetPoolConfigRes{}
        cfg := relay_pool.GetConfig(ctx)
        res.Data.AutoRouteEnabled = cfg.AutoRouteEnabled
        return
}

func (c *ControllerV1) UpdatePoolConfig(ctx context.Context, req *v1.UpdatePoolConfigReq) (res *v1.UpdatePoolConfigRes, err error) {
        res = &v1.UpdatePoolConfigRes{}
        if req.AutoRouteEnabled != nil {
                if setErr := relay_pool.SetAutoRoute(ctx, *req.AutoRouteEnabled); setErr != nil {
                        res.Code = 1
                        res.Message = setErr.Error()
                }
        }
        return
}

func (c *ControllerV1) CreateProvider(ctx context.Context, req *v1.CreateProviderReq) (res *v1.CreateProviderRes, err error) {
        res = &v1.CreateProviderRes{}
        result, dbErr := g.DB().Model("bm_relay_providers").Data(g.Map{
                "name":          req.Name,
                "provider_type": req.ProviderType,
                "slot":          req.Slot,
                "api_key":       req.APIKey,
                "api_key_2":     req.APIKey2,
                "daily_limit":   req.DailyLimit,
                "monthly_limit": req.MonthlyLimit,
                "priority":      req.Priority,
                "is_active":     req.IsActive,
                "status":        "active",
                "created_at":    time.Now().Unix(),
                "updated_at":    time.Now().Unix(),
        }).Insert()
        if dbErr != nil {
                res.Code = 1
                res.Message = dbErr.Error()
                return
        }
        id, _ := result.LastInsertId()
        res.Data.ID = id
        return
}

func (c *ControllerV1) DeleteProvider(ctx context.Context, req *v1.DeleteProviderReq) (res *v1.DeleteProviderRes, err error) {
        res = &v1.DeleteProviderRes{}
        // Check if this is a seeded slot (provider_type + slot combo that came from migration).
        // For seeded slots: clear the key and deactivate instead of deleting the row.
        var row struct {
                ProviderType string `json:"provider_type"`
                Slot         int    `json:"slot"`
        }
        if dbErr := g.DB().Model("bm_relay_providers").Where("id", req.ID).Scan(&row); dbErr != nil || row.ProviderType == "" {
                res.Code = 1
                res.Message = "Provider not found"
                return
        }

        // Check if seeded (any of the 15 known types)
        seededTypes := map[string]bool{
                "resend": true, "sendgrid": true, "brevo": true, "mailjet": true,
                "mailersend": true, "sendpulse": true, "zeptomail": true, "smtp2go": true,
                "elasticemail": true, "mailgun": true, "postmark": true, "sparkpost": true,
                "socketlabs": true, "netcore": true, "loops": true,
        }
        if seededTypes[row.ProviderType] && row.Slot <= 2 {
                // Seeded slot — clear key and deactivate
                if _, dbErr := g.DB().Model("bm_relay_providers").Where("id", req.ID).Data(g.Map{
                        "api_key":    "",
                        "api_key_2":  "",
                        "is_active":  false,
                        "status":     "unconfigured",
                        "last_error": "",
                        "updated_at": time.Now().Unix(),
                }).Update(); dbErr != nil {
                        res.Code = 1
                        res.Message = dbErr.Error()
                }
                return
        }

        // Custom slot — hard delete
        if _, dbErr := g.DB().Model("bm_relay_providers").Where("id", req.ID).Delete(); dbErr != nil {
                res.Code = 1
                res.Message = dbErr.Error()
        }
        return
}

// maskKey masks an API key for display.
func maskKey(key string) string {
        if key == "" {
                return ""
        }
        if len(key) <= 8 {
                return strings.Repeat("*", len(key))
        }
        return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}
