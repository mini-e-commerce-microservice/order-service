package cdc

import (
	"github.com/SyaibanAhmadRamadhan/event-bus/debezium"
)

type DebeziumPayload[T any] struct {
	Payload T `json:"payload"`
}

type ProductItemData struct {
	ID            int64              `json:"id"`
	AggregateID   int64              `json:"aggregate_id"`
	AggregateType string             `json:"aggregate_type"`
	Payload       string             `json:"payload"`
	TraceParent   *string            `json:"trace_parent"`
	Op            debezium.Operation `json:"op"`
}
