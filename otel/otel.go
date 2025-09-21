package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func Init() (*sdktrace.TracerProvider, error) {
	// Create stdout exporter
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())

	if err != nil {
		return nil, err
	}

	// Create a tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("bookstore-api"),
			),
		),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	return tp, nil
}
