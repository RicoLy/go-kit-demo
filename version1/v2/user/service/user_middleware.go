package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

const ContextReqUUid = "req_uuid"

type NewMiddlewareServer func(service UserService) UserService

type logMiddlewareServer struct {
	logger *zap.Logger
	next UserService
}

func NewLogMiddlewareServer(log *zap.Logger) NewMiddlewareServer {
	return func(service UserService) UserService {
		return logMiddlewareServer{
			logger: log,
			next:   service,
		}
	}
}

func (l logMiddlewareServer) Login(ctx context.Context, vo *LoginVO) (user *UserInfoDTO, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Login logMiddlewareServer", "Login"), zap.Any("req", vo), zap.Any("res", user))
	}()
	return l.next.Login(ctx, vo)
}

func (l logMiddlewareServer) Register(ctx context.Context, vo *RegisterUserVO) (user *UserInfoDTO, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Login logMiddlewareServer", "Register"), zap.Any("req", vo), zap.Any("res", user))
	}()
	return l.next.Register(ctx, vo)
}






