package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go-kit-demo/demo03/v1/service"
	"go.uber.org/zap"
)

type UserEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint endpoint.Endpoint
}

func NewUserEndpoints(srv service.UserService, log *zap.Logger) UserEndpoints {
	var registerEndpoint endpoint.Endpoint
	{
		registerEndpoint = MakeRegisterEndpoint(srv)
		registerEndpoint = LoggingMiddleware(log)(registerEndpoint)
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndpoint(srv)
		loginEndpoint = LoggingMiddleware(log)(loginEndpoint)
	}

	return UserEndpoints{
		RegisterEndpoint: registerEndpoint,
		LoginEndpoint:    loginEndpoint,
	}
}

func (e *UserEndpoints) Login(ctx context.Context, vo *service.LoginVO) (userInfoDTO *service.UserInfoDTO, err error) {
	user, err := e.LoginEndpoint(ctx, vo)
	return user.(*service.UserInfoDTO), err
}

func MakeLoginEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*service.LoginVO)
		return svc.Login(ctx, req)
	}
}

func (e *UserEndpoints) Register(ctx context.Context, vo *service.RegisterUserVO) (userInfoDTO *service.UserInfoDTO, err error) {
	user, err := e.RegisterEndpoint(ctx, vo)
	return user.(*service.UserInfoDTO), err
}

func MakeRegisterEndpoint(srv service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*service.RegisterUserVO)
		return srv.Register(ctx, req)
	}
}


