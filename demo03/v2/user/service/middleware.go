package service

import (
	"context"
	"fmt"
	"go-kit-demo/demo03/v2/user/pb"
	"go.uber.org/zap"
)

const ContextReqUUid = "req_uuid"

type SMiddlewareServer func(service Service) Service

type logMiddlewareServer struct {
	logger *zap.Logger
	next Service
}

func (l logMiddlewareServer) Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Login logMiddlewareServer", "Login"), zap.Any("req", in), zap.Any("res", ack), zap.Any("err", err))
	}()
	return l.next.Login(ctx, in)
}

func NewLogMiddleWareServer(log *zap.Logger) SMiddlewareServer {
	return func(service Service) Service {
		return logMiddlewareServer{
			logger: log,
			next:   service,
		}
	}
}
