package log

import "github.com/sirupsen/logrus"

type GrpcLogger struct {
	*logrus.Entry
}

func Grpc() GrpcLogger {
	return GrpcLogger{e}
}

func (g *GrpcLogger) V(l int) bool {
	return false
}
