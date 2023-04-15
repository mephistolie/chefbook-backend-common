package log

import "github.com/sirupsen/logrus"

type GrpcLogger struct {
	*logrus.Entry
}

func Grpc() GrpcLogger {
	return GrpcLogger{e}
}
