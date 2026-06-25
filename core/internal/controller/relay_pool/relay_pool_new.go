package relay_pool_ctrl

import (
	relay_pool_api "billionmail-core/api/relay_pool"
)

type ControllerV1 struct{}

func NewV1() relay_pool_api.IRelayPoolV1 {
	return &ControllerV1{}
}
