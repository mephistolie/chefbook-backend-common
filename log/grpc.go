package log

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type GrpcLogger struct {
}

func Grpc() GrpcLogger {
	return GrpcLogger{}
}

func (g GrpcLogger) V(l int) bool {
	return false
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		code := status.Code(err)
		event := Event{
			Event:      "grpc.request.completed",
			Message:    "gRPC request completed",
			Component:  ComponentGRPC,
			Operation:  info.FullMethod,
			Duration:   time.Since(start),
			GRPCMethod: info.FullMethod,
			GRPCCode:   code.String(),
		}

		if err != nil {
			LogError(ctx, event, err)
		} else if !isHealthCheck(info.FullMethod) {
			Log(ctx, event)
		}

		return resp, err
	}
}

func isHealthCheck(fullMethodName string) bool {
	return strings.Contains(fullMethodName, "Health") && strings.Contains(fullMethodName, "Check")
}
