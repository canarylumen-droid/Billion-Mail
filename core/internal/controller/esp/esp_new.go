package esp

import (
	"billionmail-core/api/esp"
)

type ControllerV1 struct{}

func NewV1() esp.IEspV1 {
	return &ControllerV1{}
}
