package metrics

import (
	"context"
	"log"
	"math/rand"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

// InitMeter initializes and returns an instance of the MeterProvider.
// The MeterProvider is configured with a resource describing the service
// and an OTLP exporter to send metrics to a remote collector.
func InitMeter() *metricsdk.MeterProvider {
	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	// if len(insecure) > 0 {
	// 	secureOption = otlpmetricgrpc.WithInsecure()
	// }
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		secureOption,
		otlpmetricgrpc.WithEndpoint(collectorURL),
		otlpmetricgrpc.WithHeaders(headers),
	)
	if err != nil {
		log.Fatalf("Error creating OTLP exporter: %v", err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(res),
		metricsdk.WithReader(metricsdk.NewPeriodicReader(exporter)),
	)
	return provider
}

// Metrics Generator can be used to all types of meter
// This is only used for testing, it can be modified for necessary changes
func MetricsGenerator(meter metric.Meter) {
	exceptionsCounter(meter)
	rqstDurationHistogram(meter)
	countItemsGauge(meter)
}

// Exceptions counter ideally will catch the http error, or various kinds of errors
// Presently on call we are simulating a error for test purposes
func exceptionsCounter(meter metric.Meter) {
	counter, err := meter.Int64Counter(
		"http_errors",
		metric.WithUnit("1"),
		metric.WithDescription("Counts exceptions in the system"),
	)
	if err != nil {
		log.Fatal("Error creating counter: ", err)
	}

	counter.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("endpoint", "/items"),
			attribute.String("error_type", "NullPointerException"),
		),
	)
}

// requestDurationHistogram records the simulated duration of HTTP requests.
// Metric is recorded with random durations for testing purposes.
func rqstDurationHistogram(meter metric.Meter) {
	histogram, err := meter.Int64Histogram("request_duration", metric.WithUnit("ms"), metric.WithDescription("HTTP request duration"))
	if err != nil {
		log.Fatal("Error creating histogram: ", err)
	}

	histogram.Record(context.Background(), rand.Int63n(1000), metric.WithAttributes(attribute.String("path", "/api")))
}

// countItemsGauge creates an observable gauge to report the count of items.
// Value is updated via a callback function with random data for testing purposes.
func countItemsGauge(meter metric.Meter) {
	gauge, err := meter.Float64ObservableGauge("items_count", metric.WithUnit("1"), metric.WithDescription("Duration of HTTP requests"))
	if err != nil {
		log.Fatal("Error creating gauge: ", err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			attrSet := attribute.NewSet(attribute.String("process", "data"))
			withAttrSet := metric.WithAttributeSet(attrSet)
			o.ObserveFloat64(gauge, rand.Float64()*100, withAttrSet)
			return nil
		},
		gauge,
	)
	if err != nil {
		log.Fatal("Error in registering callback: ", err)
	}
}

func RecordSpanError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err, trace.WithAttributes(attribute.String("error.type", "simulated_error")))
	}
}
