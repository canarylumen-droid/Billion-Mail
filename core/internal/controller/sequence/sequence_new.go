package sequence

import (
	"billionmail-core/api/sequence"
)

type ControllerV1 struct{}

func NewV1() sequence.ISequenceV1 {
	return &ControllerV1{}
}
