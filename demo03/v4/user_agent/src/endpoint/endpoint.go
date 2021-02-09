package endpoint

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"go-kit-demo/demo03/v4/user_agent/src/service"
	"golang.org/x/time/rate"
)

type EndPointServer struct {
	LoginEndPoint endpoint.Endpoint
}

func NewEndPointServer(svc service.Service, limit *rate.Limiter) EndPointServer {
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndPoint(svc)
		loginEndpoint = NewGolangRateAllowMiddleware(limit)(loginEndpoint)
	}
	return EndPointServer{loginEndpoint}
}

func (s EndPointServer) Login(ctx context.Context, in *service.Login) (ack *service.LoginAck, err error) {
	res, err := s.LoginEndPoint(ctx, in)
	if err != nil {
		fmt.Println("s.LoginEndPoint", err)
		return nil, err
	}
	return res.(*service.LoginAck), nil
}

func MakeLoginEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*service.Login)
		return s.Login(ctx, req)
	}
}

