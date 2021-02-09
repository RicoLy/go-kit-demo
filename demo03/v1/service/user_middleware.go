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

func (l logMiddlewareServer) Login(ctx context.Context, vo *LoginVO) (userInfoDTO *UserInfoDTO, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Login logMiddlewareServer", "Login"), zap.Any("req", vo), zap.Any("res", userInfoDTO))
	}()
	return l.next.Login(ctx, vo)
}

func (l logMiddlewareServer) Register(ctx context.Context, vo *RegisterUserVO) (userInfoDTO *UserInfoDTO, err error) {
	defer func() {
		l.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Register logMiddlewareServer", "Login"), zap.Any("req", vo), zap.Any("res", userInfoDTO))
	}()
	return l.next.Register(ctx, vo)
}

func NewLogMiddlewareServer(log *zap.Logger) NewMiddlewareServer {
	return func(service UserService) UserService {
		return logMiddlewareServer{
			logger: log,
			next:   service,
		}
	}
}
