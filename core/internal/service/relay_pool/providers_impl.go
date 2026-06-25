package relay_pool

import (
        "context"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "net/url"
        "time"
)

// ── 1. Resend ─────────────────────────────────────────────────────────────────

type ResendProvider struct{}

func (p *ResendProvider) Name() string { return "Resend" }
func (p *ResendProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "from":    fmt.Sprintf("%s <%s>", msg.FromName, msg.FromEmail),
                "to":      []string{msg.ToEmail},
                "subject": msg.Subject,
                "html":    msg.HTML,
                "text":    msg.Text,
        }
        if msg.ReplyTo != "" {
                payload["reply_to"] = msg.ReplyTo
        }
        body, status, err := doPost(ctx, "https://api.resend.com/emails",
                map[string]string{"Authorization": "Bearer " + apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        var resp struct{ ID string `json:"id"` }
        _ = json.Unmarshal(body, &resp)
        return SendResult{ProviderName: p.Name(), ProviderMessageID: resp.ID, LatencyMs: ms(start)}
}

// ── 2. SendGrid ───────────────────────────────────────────────────────────────

type SendGridProvider struct{}

func (p *SendGridProvider) Name() string { return "SendGrid" }
func (p *SendGridProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "personalizations": []map[string]interface{}{
                        {"to": []map[string]string{{"email": msg.ToEmail, "name": msg.ToName}}},
                },
                "from":    map[string]string{"email": msg.FromEmail, "name": msg.FromName},
                "subject": msg.Subject,
                "content": []map[string]string{{"type": "text/html", "value": msg.HTML}},
        }
        if msg.Text != "" {
                payload["content"] = []map[string]string{
                        {"type": "text/plain", "value": msg.Text},
                        {"type": "text/html", "value": msg.HTML},
                }
        }
        body, status, err := doPost(ctx, "https://api.sendgrid.com/v3/mail/send",
                map[string]string{"Authorization": "Bearer " + apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 3. Brevo (Sendinblue) ────────────────────────────────────────────────────

type BrevoProvider struct{}

func (p *BrevoProvider) Name() string { return "Brevo" }
func (p *BrevoProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "sender":      map[string]string{"email": msg.FromEmail, "name": msg.FromName},
                "to":          []map[string]string{{"email": msg.ToEmail, "name": msg.ToName}},
                "subject":     msg.Subject,
                "htmlContent": msg.HTML,
                "textContent": msg.Text,
        }
        body, status, err := doPost(ctx, "https://api.brevo.com/v3/smtp/email",
                map[string]string{"api-key": apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        var resp struct{ MessageID string `json:"messageId"` }
        _ = json.Unmarshal(body, &resp)
        return SendResult{ProviderName: p.Name(), ProviderMessageID: resp.MessageID, LatencyMs: ms(start)}
}

// ── 4. Mailjet ────────────────────────────────────────────────────────────────

type MailjetProvider struct{}

func (p *MailjetProvider) Name() string { return "Mailjet" }
func (p *MailjetProvider) Send(ctx context.Context, apiKey, apiSecret string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "Messages": []map[string]interface{}{
                        {
                                "From":     map[string]string{"Email": msg.FromEmail, "Name": msg.FromName},
                                "To":       []map[string]string{{"Email": msg.ToEmail, "Name": msg.ToName}},
                                "Subject":  msg.Subject,
                                "HTMLPart": msg.HTML,
                                "TextPart": msg.Text,
                        },
                },
        }
        auth := base64.StdEncoding.EncodeToString([]byte(apiKey + ":" + apiSecret))
        body, status, err := doPost(ctx, "https://api.mailjet.com/v3.1/send",
                map[string]string{"Authorization": "Basic " + auth}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 5. Mailersend ─────────────────────────────────────────────────────────────

type MailersendProvider struct{}

func (p *MailersendProvider) Name() string { return "Mailersend" }
func (p *MailersendProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "from":    map[string]string{"email": msg.FromEmail, "name": msg.FromName},
                "to":      []map[string]string{{"email": msg.ToEmail, "name": msg.ToName}},
                "subject": msg.Subject,
                "html":    msg.HTML,
                "text":    msg.Text,
        }
        body, status, err := doPost(ctx, "https://api.mailersend.com/v1/email",
                map[string]string{"Authorization": "Bearer " + apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 6. Sendpulse ──────────────────────────────────────────────────────────────
// Sendpulse requires OAuth2 — apiKey = client_id, apiKey2 = client_secret

type SendpulseProvider struct{}

func (p *SendpulseProvider) Name() string { return "Sendpulse" }
func (p *SendpulseProvider) Send(ctx context.Context, clientID, clientSecret string, msg *EmailMessage) SendResult {
        start := time.Now()

        // Step 1: get access token
        tokenPayload := map[string]string{
                "grant_type":    "client_credentials",
                "client_id":     clientID,
                "client_secret": clientSecret,
        }
        tokenBody, tokenStatus, err := doPost(ctx, "https://api.sendpulse.com/oauth/access_token",
                nil, tokenPayload)
        if err != nil || tokenStatus >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: fmt.Errorf("sendpulse auth failed: %v", err)}
        }
        var tokenResp struct{ AccessToken string `json:"access_token"` }
        _ = json.Unmarshal(tokenBody, &tokenResp)

        // Step 2: send email
        payload := map[string]interface{}{
                "email": map[string]interface{}{
                        "html":    msg.HTML,
                        "text":    msg.Text,
                        "subject": msg.Subject,
                        "from":    map[string]string{"name": msg.FromName, "email": msg.FromEmail},
                        "to":      []map[string]string{{"name": msg.ToName, "email": msg.ToEmail}},
                },
        }
        body, status, err := doPost(ctx, "https://api.sendpulse.com/smtp/emails",
                map[string]string{"Authorization": "Bearer " + tokenResp.AccessToken}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 7. ZeptoMail (Zoho) ──────────────────────────────────────────────────────

type ZeptoMailProvider struct{}

func (p *ZeptoMailProvider) Name() string { return "ZeptoMail" }
func (p *ZeptoMailProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "from":     map[string]string{"address": msg.FromEmail, "name": msg.FromName},
                "to":       []map[string]interface{}{{"email_address": map[string]string{"address": msg.ToEmail, "name": msg.ToName}}},
                "subject":  msg.Subject,
                "htmlbody": msg.HTML,
                "textbody": msg.Text,
        }
        body, status, err := doPost(ctx, "https://api.zeptomail.com/v1.1/email",
                map[string]string{"Authorization": apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 8. SMTP2GO ────────────────────────────────────────────────────────────────

type SMTP2GOProvider struct{}

func (p *SMTP2GOProvider) Name() string { return "SMTP2GO" }
func (p *SMTP2GOProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "api_key":  apiKey,
                "to":       []string{fmt.Sprintf("%s <%s>", msg.ToName, msg.ToEmail)},
                "sender":   fmt.Sprintf("%s <%s>", msg.FromName, msg.FromEmail),
                "subject":  msg.Subject,
                "htmlbody": msg.HTML,
                "textbody": msg.Text,
        }
        body, status, err := doPost(ctx, "https://api.smtp2go.com/v3/email/send", nil, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 9. Elastic Email ──────────────────────────────────────────────────────────

type ElasticEmailProvider struct{}

func (p *ElasticEmailProvider) Name() string { return "Elastic Email" }
func (p *ElasticEmailProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "Recipients": map[string]interface{}{
                        "To": []map[string]string{{"Email": msg.ToEmail, "Name": msg.ToName}},
                },
                "Content": map[string]interface{}{
                        "From":    map[string]string{"Email": msg.FromEmail, "Name": msg.FromName},
                        "Subject": msg.Subject,
                        "Body": []map[string]string{
                                {"ContentType": "HTML", "Content": msg.HTML},
                                {"ContentType": "PlainText", "Content": msg.Text},
                        },
                },
        }
        body, status, err := doPost(ctx, "https://api.elasticemail.com/v4/emails",
                map[string]string{"X-ElasticEmail-ApiKey": apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 10. Mailgun ───────────────────────────────────────────────────────────────
// apiKey = Mailgun API key, apiKey2 = sending domain (e.g. mg.yourdomain.com)

type MailgunProvider struct{}

func (p *MailgunProvider) Name() string { return "Mailgun" }
func (p *MailgunProvider) Send(ctx context.Context, apiKey, domain string, msg *EmailMessage) SendResult {
        start := time.Now()
        if domain == "" {
                domain = "sandbox.mailgun.org"
        }
        formData := url.Values{}
        formData.Set("from", fmt.Sprintf("%s <%s>", msg.FromName, msg.FromEmail))
        formData.Set("to", msg.ToEmail)
        formData.Set("subject", msg.Subject)
        formData.Set("html", msg.HTML)
        if msg.Text != "" {
                formData.Set("text", msg.Text)
        }

        auth := base64.StdEncoding.EncodeToString([]byte("api:" + apiKey))
        apiURL := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", domain)
        body, status, err := doPostForm(ctx, apiURL,
                map[string]string{"Authorization": "Basic " + auth}, formData.Encode())
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        var resp struct{ ID string `json:"id"` }
        _ = json.Unmarshal(body, &resp)
        return SendResult{ProviderName: p.Name(), ProviderMessageID: resp.ID, LatencyMs: ms(start)}
}

// ── 11. Postmark ──────────────────────────────────────────────────────────────

type PostmarkProvider struct{}

func (p *PostmarkProvider) Name() string { return "Postmark" }
func (p *PostmarkProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "From":     fmt.Sprintf("%s <%s>", msg.FromName, msg.FromEmail),
                "To":       msg.ToEmail,
                "Subject":  msg.Subject,
                "HtmlBody": msg.HTML,
                "TextBody": msg.Text,
                "MessageStream": "outbound",
        }
        body, status, err := doPost(ctx, "https://api.postmarkapp.com/email",
                map[string]string{"X-Postmark-Server-Token": apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        var resp struct{ MessageID string `json:"MessageID"` }
        _ = json.Unmarshal(body, &resp)
        return SendResult{ProviderName: p.Name(), ProviderMessageID: resp.MessageID, LatencyMs: ms(start)}
}

// ── 12. SparkPost ─────────────────────────────────────────────────────────────

type SparkPostProvider struct{}

func (p *SparkPostProvider) Name() string { return "SparkPost" }
func (p *SparkPostProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "recipients": []map[string]interface{}{
                        {"address": map[string]string{"email": msg.ToEmail, "name": msg.ToName}},
                },
                "content": map[string]interface{}{
                        "from":    map[string]string{"email": msg.FromEmail, "name": msg.FromName},
                        "subject": msg.Subject,
                        "html":    msg.HTML,
                        "text":    msg.Text,
                },
        }
        body, status, err := doPost(ctx, "https://api.sparkpost.com/api/v1/transmissions",
                map[string]string{"Authorization": apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 13. SocketLabs ────────────────────────────────────────────────────────────
// apiKey = API key, apiKey2 = Server ID

type SocketLabsProvider struct{}

func (p *SocketLabsProvider) Name() string { return "SocketLabs" }
func (p *SocketLabsProvider) Send(ctx context.Context, apiKey, serverID string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "ServerId": serverID,
                "ApiKey":   apiKey,
                "Messages": []map[string]interface{}{
                        {
                                "To":      []map[string]string{{"EmailAddress": msg.ToEmail, "FriendlyName": msg.ToName}},
                                "From":    map[string]string{"EmailAddress": msg.FromEmail, "FriendlyName": msg.FromName},
                                "Subject": msg.Subject,
                                "HtmlBody": msg.HTML,
                                "TextBody": msg.Text,
                        },
                },
        }
        body, status, err := doPost(ctx, "https://inject.socketlabs.com/api/v1/email", nil, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 14. Netcore (Pepipost) ───────────────────────────────────────────────────

type NetcoreProvider struct{}

func (p *NetcoreProvider) Name() string { return "Netcore" }
func (p *NetcoreProvider) Send(ctx context.Context, apiKey, _ string, msg *EmailMessage) SendResult {
        start := time.Now()
        payload := map[string]interface{}{
                "from": map[string]string{"email": msg.FromEmail, "name": msg.FromName},
                "subject": msg.Subject,
                "content": []map[string]string{
                        {"type": "html", "value": msg.HTML},
                },
                "personalizations": []map[string]interface{}{
                        {"to": []map[string]string{{"email": msg.ToEmail, "name": msg.ToName}}},
                },
        }
        body, status, err := doPost(ctx, "https://emailapi.netcorecloud.net/v5/mail/send",
                map[string]string{"api_key": apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── 15. Loops.so ─────────────────────────────────────────────────────────────
// apiKey  = Loops API key (loops.so/settings → API)
// apiKey2 = Transactional template ID from Loops dashboard (required).
//
// Create a template in Loops with dataVariables: subject, html, text,
// fromName, fromEmail. Then paste the template ID into api_key_2.

type LoopsProvider struct{}

func (p *LoopsProvider) Name() string { return "Loops" }
func (p *LoopsProvider) Send(ctx context.Context, apiKey, templateID string, msg *EmailMessage) SendResult {
        start := time.Now()
        if templateID == "" {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start),
                        Err: fmt.Errorf("loops: api_key_2 must be set to your transactional template ID")}
        }
        payload := map[string]interface{}{
                "transactionalId": templateID,
                "email":           msg.ToEmail,
                "dataVariables": map[string]string{
                        "subject":   msg.Subject,
                        "html":      msg.HTML,
                        "text":      msg.Text,
                        "fromName":  msg.FromName,
                        "fromEmail": msg.FromEmail,
                },
        }
        body, status, err := doPost(ctx, "https://app.loops.so/api/v1/transactional",
                map[string]string{"Authorization": "Bearer " + apiKey}, payload)
        if err != nil {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: err}
        }
        if status >= 400 {
                return SendResult{ProviderName: p.Name(), LatencyMs: ms(start), Err: errFromStatus(p.Name(), status, body)}
        }
        return SendResult{ProviderName: p.Name(), LatencyMs: ms(start)}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func ms(start time.Time) int64 {
        return time.Since(start).Milliseconds()
}
