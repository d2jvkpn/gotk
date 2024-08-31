package cloud

import (
	"context"
	// "fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type LogServer struct {
	logger *zap.Logger
}

func NewLogServer(logger *zap.Logger) *LogServer {
	return &LogServer{logger: logger}
}

func (inte *LogServer) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		resp any, err error) {

		var (
			ok    bool
			start time.Time
			md    metadata.MD
		)

		start = time.Now()
		if md, ok = metadata.FromIncomingContext(ctx); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "no metadata")
		}

		inte.logger.Debug(info.FullMethod, zap.Any("metadata", md))

		resp, err = handler(ctx, req)
		latency := zap.Float64("latencyMs", float64(time.Since(start).Microseconds())/1e3)

		if err == nil {
			inte.logger.Info(info.FullMethod, latency)
		} else {
			inte.logger.Error(info.FullMethod, latency, zap.Any("error", err))
		}

		return resp, err
	}

	// return grpc.UnaryInterceptor(call) // grpc.ServerOption
}

func (inte *LogServer) Stream() grpc.StreamServerInterceptor {
	return func(
		srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) (err error) {

		var (
			ok    bool
			start time.Time
			md    metadata.MD
		)

		start = time.Now()
		if md, ok = metadata.FromIncomingContext(ss.Context()); !ok {
			return status.Errorf(codes.Unauthenticated, "no metadata")
		}

		inte.logger.Debug(info.FullMethod, zap.Any("metadata", md))

		err = handler(srv, ss)
		latency := zap.Float64("latencyMs", float64(time.Since(start).Microseconds())/1e3)

		if err == nil {
			inte.logger.Info(info.FullMethod, latency)
		} else {
			inte.logger.Error(info.FullMethod, latency, zap.Any("error", err))
		}

		return err
	}

	// return grpc.StreamInterceptor(call) // grpc.ServerOption
}
