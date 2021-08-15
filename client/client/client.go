package client

import (
	"context"
	"fmt"
	"os"

	"github.com/pedramteymoori/grpc-jaeger-demo/protocols"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Client struct{}

func (dc Client) Run() {
	grpcDialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}

	path := os.Getenv("SERVER_URL")
	conn, _ := grpc.Dial(path, grpcDialOpts...)
	cli := protocols.NewDemoClient(conn)

	ctx := context.Background()
	// newCtx, span := trace.StartSpan(ctx, "say-hello-client-span")
	// defer span.End()

	// fmt.Println(span.SpanContext().SpanID.String())
	// fmt.Println(span.SpanContext().TraceID.String())

	ret, err := cli.SayHello(ctx, &protocols.SayHelloRequest{Name: "Pedram"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Greeting)
}
