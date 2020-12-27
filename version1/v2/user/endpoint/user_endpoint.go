package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go-kit-demo/version1/v2/user/service"
	"go.uber.org/zap"
)

type UserEndpoints struct {
	RegisterEndpoint  endpoint.Endpoint
	LoginEndpoint endpoint.Endpoint
}

func NewUserEndpoints(svc service.UserService, log *zap.Logger) UserEndpoints {
	var registerEndpoint endpoint.Endpoint
	{
		registerEndpoint = MakeRegisterEndpoint(svc)
		registerEndpoint = LoggingMiddleware(log)(registerEndpoint)
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndpoint(svc)
		loginEndpoint = LoggingMiddleware(log)(loginEndpoint)
	}

	return UserEndpoints{
		RegisterEndpoint: registerEndpoint,
		LoginEndpoint:    loginEndpoint,
	}
}

func (e *UserEndpoints) Login(ctx context.Context, vo *service.LoginVO) (*service.UserInfoDTO, error) {
	user, err := e.LoginEndpoint(ctx, vo)
	return user.(*service.UserInfoDTO), err
}

func MakeLoginEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*service.LoginVO)
		return svc.Login(ctx, req)
	}
}

func (e *UserEndpoints) Register(ctx context.Context, vo *service.RegisterUserVO) (*service.UserInfoDTO, error) {
	user, err := e.RegisterEndpoint(ctx, vo)
	return user.(*service.UserInfoDTO), err
}

func MakeRegisterEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*service.RegisterUserVO)
		return svc.Register(ctx, req)
	}
}
