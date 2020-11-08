package handlers

import (
	"balancer-api/balancer"
)

const (
	Asset     = "ASSET"
	Liability = "LIABILITY"
)

type IDPayload struct {
	ID uint64 `json:"recordId"`
}

type ValuePayload struct {
	Value float64 `json:"value"`
}

type Handler struct {
	RecordService balancer.RecordService
}
