package main

import (
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"github.com/pedramteymoori/grpc-jaeger-demo/client/client"
	"go.opencensus.io/trace"
)

func main() {
	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress("collector.linkerd-jaeger:55678"),
		ocagent.WithServiceName("pedram-client"))
	if err != nil {
		log.Fatalf("Failed to create ocagent-exporter: %v", err)
	}
	defer oce.Stop()
	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

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
