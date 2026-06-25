package leads

import "billionmail-core/api/leads"

type ControllerV1 struct{}

func NewV1() leads.ILeadsV1 {
	return &ControllerV1{}
}
