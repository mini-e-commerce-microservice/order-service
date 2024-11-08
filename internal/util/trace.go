package util

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

func GetTraceParent(ctx context.Context) *string {
	traceParent := whttp.GetTraceParent(ctx)
	if traceParent != "" {
		return &traceParent
	}

	return nil
}

func SpanRecordErrorWIthEnd(span trace.Span, err error, errType string) {
	span.RecordError(collection.Err(err))
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(semconv.ErrorTypeKey.String(errType))
	span.End()
}
