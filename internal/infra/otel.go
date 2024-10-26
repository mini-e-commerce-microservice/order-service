package infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"order-service/generated/proto/secret_proto"
	"time"
)

func NewOtel(cred *secret_proto.Otel, tracerName string) collection.CloseFnCtx {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(cred.Username+":"+cred.Password))
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cred.Endpoint),
		otlptracegrpc.WithHeaders(map[string]string{
			"Authorization": authHeader,
		}),
	)

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		panic(err)
	}

	traceProvider, closeFn, err := startOtelProvider(traceExp, tracerName)
	if err != nil {
		panic(err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(traceProvider)

	log.Info().Msg("initialization opentelemetry successfully")
	return closeFn
}

func startOtelProvider(exp trace.SpanExporter, tracerName string) (*trace.TracerProvider, collection.CloseFnCtx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(tracerName),
		),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
		resource.WithFromEnv(),
	)

	if err != nil {
		err = fmt.Errorf("failed to created resource: %w", err)
		return nil, nil, err
	}

	bsp := trace.NewBatchSpanProcessor(exp)

	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	closeFn := func(ctx context.Context) error {
		log.Info().Msg("starting shutdown export and provider")
		ctxClosure, cancelClosure := context.WithTimeout(ctx, 5*time.Second)
		defer cancelClosure()

		if err = exp.Shutdown(ctxClosure); err != nil {
			return err
		}

		if err = provider.Shutdown(ctxClosure); err != nil {
			return err
		}

		log.Info().Msg("shutdown export and provider successfully")
		return nil
	}

	return provider, closeFn, nil
}
