package otel

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc/credentials"
)

var tracerProvider *sdktrace.TracerProvider
var loggerProvider *sdklog.LoggerProvider
var shutDownFuncs []func()

func Init() error {
	nrEndPoint := "otlp.nr-data.net:4317"
	headers := map[string]string{"api-key": os.Getenv("NEW_RELIC_LICENSE_KEY")}
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("bookstore-api"),
	)

	if err := initLogProvider(res, nrEndPoint, headers); err != nil {
		return err
	}

	shutDownFuncs = append(shutDownFuncs, logShutDown)

	if err := initTraceProvider(res, nrEndPoint, headers); err != nil {
		return err
	}

	shutDownFuncs = append(shutDownFuncs, traceShutDown)

	return nil
}

func initTraceProvider(res *resource.Resource, endPoint string, headers map[string]string) error {
	// Create stdout exporter
	// exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())

	//OTLP trace exporter to NewRelics
	traceExporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(endPoint),
		otlptracegrpc.WithHeaders(headers),
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)

	if err != nil {
		return err
	}

	// Create a tracer provider
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)

	return nil
}

func initLogProvider(res *resource.Resource, endPoint string, headers map[string]string) error {
	//OTLP log exporter to NewRelics
	logExporter, err := otlploggrpc.New(
		context.Background(),
		otlploggrpc.WithEndpoint(endPoint),
		otlploggrpc.WithHeaders(headers),
	)

	if err != nil {
		return err
	}

	// Create a log provider
	loggerProvider = sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)

	// Set global logger provider
	global.SetLoggerProvider(loggerProvider)

	return nil
}

func initMetricProvider(res *resource.Resource, endPoint string, headers map[string]string) error {
	metricExporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithEndpoint(endPoint),
		otlpmetricgrpc.WithHeaders(headers),
	)

	if err != nil {
		return err
	}

	metricReader := sdkmetric.NewPeriodicReader(
		metricExporter,
		sdkmetric.WithInterval(10 * time.Second),
	)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(metricReader),
	)

	otel.SetMeterProvider(meterProvider)

	return nil
}

func traceShutDown() {
	if err := tracerProvider.Shutdown(context.Background()); err != nil {
		log.Printf("Failed to shutdown tracer provider: %v", err)
	} else {
		log.Println("Tracer provider shutdown gracefully!! 👍")
	}
}

func logShutDown() {
	if err := loggerProvider.Shutdown(context.Background()); err != nil {
		log.Printf("Failed to shutdown logger provider: %v", err)
	} else {
		log.Println("Logger provider shutdown gracefully!! 👍")
	}
}

func ShutDown() {
	for _, shutDown := range shutDownFuncs {
		shutDown()
	}

	log.Println("OpenTelemetry shut down successfully.")
}
