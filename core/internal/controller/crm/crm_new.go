package crm

import "billionmail-core/api/crm"

type ControllerV1 struct{}

func NewV1() crm.ICrmV1 {
	return &ControllerV1{}
}
