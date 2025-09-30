package otel

import (
	"context"
	"os"

	"github.com/BULLKNIGHT/bookstore/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc/credentials"
)

var tp *sdktrace.TracerProvider

func Init() (*sdktrace.TracerProvider, error) {
	// Create stdout exporter
	// exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())

	//OTLP exporter to NewRelics
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint("otlp.nr-data.net:4317"),
		otlptracegrpc.WithHeaders(map[string]string{
			"api-key": os.Getenv("NEW_RELIC_LICENSE_KEY"),
		}),
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)

	if err != nil {
		return nil, err
	}

	// Create a tracer provider
	tp = sdktrace.NewTracerProvider(
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

func ShutDown() {
	if err := tp.Shutdown(context.Background()); err != nil {
		logger.Log.WithError(err).Error("Failed to shutdown tracer provider!! 👎")
	} else {
		logger.Log.Info("Tracer provider shutdown gracefully!! 👎")
	}
}
