package logging

import (
	"context"
	// "fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GrpcCliLogger struct {
	Headers map[string]string
}

func (inte *GrpcCliLogger) attach(ctx context.Context) context.Context {
	kv := make([]string, 0, len(inte.Headers)*2)

	for k, v := range inte.Headers {
		kv = append(kv, k, v)
	}

	return metadata.AppendToOutgoingContext(ctx, kv...)
}

func (inte *GrpcCliLogger) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) (err error) {

		return invoker(inte.attach(ctx), method, req, reply, cc, opts...)
	}
}

func (inte *GrpcCliLogger) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption,
	) (client grpc.ClientStream, err error) {

		return streamer(inte.attach(ctx), desc, cc, method, opts...)
	}
}
