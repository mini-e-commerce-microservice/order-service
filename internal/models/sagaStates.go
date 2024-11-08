package models

type SagaState struct {
	ID      int64  `db:"id"`
	Payload any    `db:"payload"`
	Status  string `db:"status"`
	Step    any    `db:"step"`
	Type    string `db:"type"`
	Version string `db:"version"`
}

type SagaStateCreateOrderProductStep struct {
	Initiated string `json:"initiated,omitempty"`
	Payment   string `json:"payment,omitempty"`
	Shipping  string `json:"shipping,omitempty"`
}
