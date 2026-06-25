package relay_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// ProviderRow is a loaded DB row for a relay provider slot.
type ProviderRow struct {
	ID               int64
	Name             string
	ProviderType     string
	Slot             int
	APIKey           string
	APIKey2          string
	DailyLimit       int
	MonthlyLimit     int
	DailySent        int
	MonthlySent      int
	LastDailyReset   int64
	LastMonthlyReset int64
	RetryAfter       int64
	Priority         int
	Status           string
	IsActive         bool
}

// SendViaPool picks the first available configured provider and sends the
// message. It falls back through all providers in priority order.
// When all configured providers are exhausted for today, it returns
// ErrPoolExhausted so the caller can queue for retry next day.
var ErrPoolExhausted = fmt.Errorf("relay_pool: all configured providers are exhausted — email queued for retry after midnight UTC")

// SendViaPool routes a message through the highest-priority available slot.
// Returns the SendResult from the first successful provider or ErrPoolExhausted
// if every slot has hit its daily/monthly limit.
func SendViaPool(ctx context.Context, msg *EmailMessage) (SendResult, error) {
	providers, err := loadActiveProviders(ctx)
	if err != nil {
		return SendResult{}, fmt.Errorf("relay_pool: load providers: %w", err)
	}
	if len(providers) == 0 {
		return SendResult{}, fmt.Errorf("relay_pool: no active providers configured — add API keys in Settings → Relay Providers")
	}

	now := time.Now().Unix()
	allExhausted := true
	var lastErr error

	for _, row := range providers {
		// Honour backoff window from previous exhaustion
		if row.RetryAfter > 0 && now < row.RetryAfter {
			continue
		}

		// Lazy daily reset: if the stored reset timestamp is from a previous
		// calendar day (or was never set), zero the counter and record today.
		if needsDailyReset(row) {
			_ = resetDaily(ctx, row.ID)
			row.DailySent = 0
		}
		// Lazy monthly reset
		if needsMonthlyReset(row) {
			_ = resetMonthly(ctx, row.ID)
			row.MonthlySent = 0
		}

		// Skip exhausted slots
		if row.DailySent >= row.DailyLimit || row.MonthlySent >= row.MonthlyLimit {
			_ = markExhausted(ctx, row.ID)
			continue
		}

		allExhausted = false

		impl, ok := GetProvider(row.ProviderType)
		if !ok {
			continue
		}

		result := impl.Send(ctx, row.APIKey, row.APIKey2, msg)
		if result.Err != nil {
			lastErr = result.Err
			_ = updateStatus(ctx, row.ID, "error", result.Err.Error())
			continue
		}

		// Increment counters and mark active
		_ = incrementCounters(ctx, row.ID)
		return result, nil
	}

	if allExhausted {
		return SendResult{}, ErrPoolExhausted
	}
	if lastErr != nil {
		return SendResult{}, fmt.Errorf("relay_pool: all providers failed, last error: %w", lastErr)
	}
	return SendResult{}, ErrPoolExhausted
}

// IsPoolAvailable returns true if at least one active, non-exhausted provider is configured.
func IsPoolAvailable(ctx context.Context) bool {
	providers, err := loadActiveProviders(ctx)
	if err != nil || len(providers) == 0 {
		return false
	}
	now := time.Now()
	nowUnixVal := now.Unix()
	for _, row := range providers {
		if row.RetryAfter > 0 && nowUnixVal < row.RetryAfter {
			continue
		}
		dailySent := row.DailySent
		if needsDailyReset(row) {
			dailySent = 0
		}
		monthlySent := row.MonthlySent
		if needsMonthlyReset(row) {
			monthlySent = 0
		}
		if dailySent < row.DailyLimit && monthlySent < row.MonthlyLimit {
			return true
		}
	}
	return false
}

// ── DB helpers ────────────────────────────────────────────────────────────────

func loadActiveProviders(ctx context.Context) ([]ProviderRow, error) {
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
		LastDailyReset   int64  `json:"last_daily_reset"`
		LastMonthlyReset int64  `json:"last_monthly_reset"`
		RetryAfter       int64  `json:"retry_after"`
		Priority         int    `json:"priority"`
		Status           string `json:"status"`
		IsActive         bool   `json:"is_active"`
	}

	err := g.DB().Model("bm_relay_providers").
		Where("is_active", true).
		WhereNot("api_key", "").
		Order("priority DESC, id ASC").
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	result := make([]ProviderRow, 0, len(rows))
	for _, r := range rows {
		result = append(result, ProviderRow{
			ID: r.ID, Name: r.Name, ProviderType: r.ProviderType, Slot: r.Slot,
			APIKey: r.APIKey, APIKey2: r.APIKey2,
			DailyLimit: r.DailyLimit, MonthlyLimit: r.MonthlyLimit,
			DailySent: r.DailySent, MonthlySent: r.MonthlySent,
			LastDailyReset: r.LastDailyReset, LastMonthlyReset: r.LastMonthlyReset,
			RetryAfter: r.RetryAfter,
			Priority:   r.Priority, Status: r.Status, IsActive: r.IsActive,
		})
	}
	return result, nil
}

// needsDailyReset returns true when the daily counter should be zeroed.
// Treats last_daily_reset == 0 as "first use" and resets to initialize.
func needsDailyReset(row ProviderRow) bool {
	if row.LastDailyReset == 0 {
		// Never been reset — initialize the timestamp to today so counters
		// track correctly from the very first send.
		return true
	}
	last := time.Unix(row.LastDailyReset, 0).UTC()
	now := time.Now().UTC()
	return last.Year() != now.Year() || last.YearDay() != now.YearDay()
}

// needsMonthlyReset returns true when the monthly counter should be zeroed.
func needsMonthlyReset(row ProviderRow) bool {
	if row.LastMonthlyReset == 0 {
		return true
	}
	last := time.Unix(row.LastMonthlyReset, 0).UTC()
	now := time.Now().UTC()
	return last.Year() != now.Year() || last.Month() != now.Month()
}

func resetDaily(ctx context.Context, id int64) error {
	_, err := g.DB().Model("bm_relay_providers").Where("id", id).Data(g.Map{
		"daily_sent":       0,
		"last_daily_reset": time.Now().Unix(),
		"retry_after":      0,
		"updated_at":       time.Now().Unix(),
	}).Update()
	return err
}

func resetMonthly(ctx context.Context, id int64) error {
	_, err := g.DB().Model("bm_relay_providers").Where("id", id).Data(g.Map{
		"monthly_sent":       0,
		"last_monthly_reset": time.Now().Unix(),
		"updated_at":         time.Now().Unix(),
	}).Update()
	return err
}

// markExhausted sets status to exhausted and sets retry_after to tomorrow midnight UTC.
func markExhausted(ctx context.Context, id int64) error {
	now := time.Now().UTC()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	_, err := g.DB().Model("bm_relay_providers").Where("id", id).Data(g.Map{
		"status":      "exhausted",
		"retry_after": tomorrow.Unix(),
		"updated_at":  time.Now().Unix(),
	}).Update()
	return err
}

func incrementCounters(ctx context.Context, id int64) error {
	_, err := g.DB().Model("bm_relay_providers").Where("id", id).Increment("daily_sent", 1)
	if err != nil {
		return err
	}
	_, err = g.DB().Model("bm_relay_providers").Where("id", id).Increment("monthly_sent", 1)
	if err != nil {
		return err
	}
	_, err = g.DB().Model("bm_relay_providers").Where("id", id).Data(g.Map{
		"last_used_at": time.Now().Unix(),
		"status":       "active",
		"last_error":   "",
		"retry_after":  0,
	}).Update()
	return err
}

func updateStatus(ctx context.Context, id int64, status, lastError string) error {
	data := g.Map{"status": status, "updated_at": time.Now().Unix()}
	if lastError != "" {
		data["last_error"] = lastError
	}
	_, err := g.DB().Model("bm_relay_providers").Where("id", id).Data(data).Update()
	return err
}
