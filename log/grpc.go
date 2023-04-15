package log

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	return grpc_logrus.UnaryServerInterceptor(e)
}
