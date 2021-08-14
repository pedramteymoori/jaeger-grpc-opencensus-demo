package main

import (
	"log"
	"net"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"

	"github.com/pedramteymoori/grpc-jaeger-demo/protocols"
	"github.com/pedramteymoori/grpc-jaeger-demo/server/server"
)

func main() {
	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress("otel-collector.jaeger:55678"),
		ocagent.WithServiceName("pedram-server"))
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

	grpcAddress := "0.0.0.0:5001"

	conn, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to listen on %q", grpcAddress)
		return
	}

	server := server.DemoServer{}
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	)

	protocols.RegisterDemoServer(grpcServer, server)

	logrus.Infof("GRPC Listening on %s", grpcAddress)
	err = grpcServer.Serve(conn)
	if err != nil {
		logrus.WithError(err).Fatal("failed to grpcServer.Serve(conn)")
	}
	logrus.Info("Stop listening.")
}
