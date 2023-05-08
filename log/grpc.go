package log

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"strings"
)

type GrpcLogger struct {
	*logrus.Entry
}

func Grpc() GrpcLogger {
	return GrpcLogger{e}
}

func (g GrpcLogger) V(l int) bool {
	return false
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	healthCheckDecider := func(fullMethodName string, err error) bool {
		isHealthCheck := strings.Contains(fullMethodName, "Health") && strings.Contains(fullMethodName, "Check")
		return err != nil || !isHealthCheck
	}
	healthCheckOption := grpc_logrus.WithDecider(healthCheckDecider)
	return grpc_logrus.UnaryServerInterceptor(e, healthCheckOption)
}
