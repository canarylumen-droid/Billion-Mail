package relay_pool

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// PoolConfig holds global relay pool behaviour settings.
type PoolConfig struct {
	AutoRouteEnabled bool `json:"auto_route_enabled"` // use pool-first routing when true
}

// GetConfig returns the current relay pool configuration from the DB.
// Returns safe defaults (auto_route disabled) on any error.
func GetConfig(ctx context.Context) PoolConfig {
	var row struct {
		AutoRouteEnabled bool `json:"auto_route_enabled"`
	}
	if err := g.DB().Model("bm_relay_pool_config").Where("1=1").Limit(1).Scan(&row); err != nil {
		return PoolConfig{}
	}
	return PoolConfig{AutoRouteEnabled: row.AutoRouteEnabled}
}

// SetAutoRoute enables or disables pool-first routing.
func SetAutoRoute(ctx context.Context, enabled bool) error {
	count, _ := g.DB().Model("bm_relay_pool_config").Count()
	if count == 0 {
		_, err := g.DB().Model("bm_relay_pool_config").Data(g.Map{
			"auto_route_enabled": enabled,
		}).Insert()
		return err
	}
	_, err := g.DB().Model("bm_relay_pool_config").Where("1=1").Data(g.Map{
		"auto_route_enabled": enabled,
	}).Update()
	return err
}

// IsAutoRouteEnabled is a convenience wrapper used by outbound send paths.
func IsAutoRouteEnabled(ctx context.Context) bool {
	return GetConfig(ctx).AutoRouteEnabled
}
