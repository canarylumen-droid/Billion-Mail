package relay_pool

import (
        "context"
        "time"
)

// EmailMessage is the payload passed to every provider's Send method.
type EmailMessage struct {
        FromEmail string
        FromName  string
        ToEmail   string
        ToName    string
        Subject   string
        HTML      string
        Text      string
        ReplyTo   string
        MessageID string
        Headers   map[string]string
}

// SendResult is returned by every provider implementation.
type SendResult struct {
        ProviderMessageID string
        ProviderName      string
        LatencyMs         int64
        Err               error
}

// EmailProvider is the interface every provider must implement.
type EmailProvider interface {
        // Name returns a human-readable provider name.
        Name() string
        // Send transmits one email message. Returns SendResult with provider message ID.
        Send(ctx context.Context, apiKey string, apiKey2 string, msg *EmailMessage) SendResult
}

// providerRegistry maps provider_type to implementation.
var providerRegistry = map[string]EmailProvider{
        "resend":       &ResendProvider{},
        "sendgrid":     &SendGridProvider{},
        "brevo":        &BrevoProvider{},
        "mailjet":      &MailjetProvider{},
        "mailersend":   &MailersendProvider{},
        "sendpulse":    &SendpulseProvider{},
        "zeptomail":    &ZeptoMailProvider{},
        "smtp2go":      &SMTP2GOProvider{},
        "elasticemail": &ElasticEmailProvider{},
        "mailgun":      &MailgunProvider{},
        "postmark":     &PostmarkProvider{},
        "sparkpost":    &SparkPostProvider{},
        "socketlabs": &SocketLabsProvider{},
        "netcore":     &NetcoreProvider{},
        "loops":       &LoopsProvider{},
}

// GetProvider returns the provider implementation for a given type.
func GetProvider(providerType string) (EmailProvider, bool) {
        p, ok := providerRegistry[providerType]
        return p, ok
}

// AllProviderTypes returns all registered provider type strings.
func AllProviderTypes() []string {
        keys := make([]string, 0, len(providerRegistry))
        for k := range providerRegistry {
                keys = append(keys, k)
        }
        return keys
}

// nowUnix returns current Unix timestamp.
func nowUnix() int64 {
        return time.Now().Unix()
}
