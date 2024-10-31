package primitive

type SagaStateStatus string

const (
	SagaStateStatusOnProcess SagaStateStatus = "ON PROCESS"
	SagaStateStatusSuccess   SagaStateStatus = "SUCCESS"
	SagaStateStatusFailed    SagaStateStatus = "FAILED"
)
