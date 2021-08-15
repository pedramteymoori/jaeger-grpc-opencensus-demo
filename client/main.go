package main

import (
	"context"
	"log"
	"time"

	"github.com/pedramteymoori/grpc-jaeger-demo/client/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

func main() {

	shutdown := initProvider()
	defer shutdown()

	// oce, err := ocagent.NewExporter(
	// 	ocagent.WithInsecure(),
	// 	ocagent.WithReconnectionPeriod(5*time.Second),
	// 	ocagent.WithAddress("otel-collector.jaeger:55678"),
	// 	ocagent.WithServiceName("pedram-client"))
	// if err != nil {
	// 	log.Fatalf("Failed to create ocagent-exporter: %v", err)
	// }
	// defer oce.Stop()
	// trace.RegisterExporter(oce)
	// trace.ApplyConfig(trace.Config{
	// 	DefaultSampler: trace.AlwaysSample(),
	// })

	// agentEndpointURI := "localhost:6831"
	// collectorEndpointURI := "http://localhost:14268/api/traces"

	// je, err := jaeger.NewExporter(jaeger.Options{
	// 	AgentEndpoint:          agentEndpointURI,
	// 	CollectorEndpoint:      collectorEndpointURI,
	// 	ServiceName:            "demo",
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to create the Jaeger exporter: %v", err)
	// }

	// time.Sleep(10 * time.Second)

	cli := client.Client{}
	for {
		cli.Run()
		time.Sleep(30 * time.Second)
	}
}

func initProvider() func() {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("pedram-client"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("otel-collector.jaeger:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		err := tracerProvider.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}
}
