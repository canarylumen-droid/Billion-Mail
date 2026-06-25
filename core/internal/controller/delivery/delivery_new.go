package delivery

import (
	"billionmail-core/api/delivery"
)

type ControllerV1 struct{}

func NewV1() delivery.IDeliveryV1 {
	return &ControllerV1{}
}
