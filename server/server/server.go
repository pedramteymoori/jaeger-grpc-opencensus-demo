package server

import (
	"context"

	"github.com/pedramteymoori/grpc-jaeger-demo/protocols"
)

type DemoServer struct {
	protocols.UnimplementedDemoServer
}

func (ds DemoServer) SayHello(ctx context.Context, req *protocols.SayHelloRequest) (*protocols.SayHelloResponse, error) {
	// _, span := trace.StartSpan(ctx, "say-hello-server-span")
	// defer span.End()

	// fmt.Println(span.SpanContext().SpanID.String())
	// fmt.Println(span.SpanContext().TraceID.String())

	ret := &protocols.SayHelloResponse{Greeting: "hello" + req.Name}
	return ret, nil
}
