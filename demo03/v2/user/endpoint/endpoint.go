package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go-kit-demo/demo03/v2/user/pb"
	"go-kit-demo/demo03/v2/user/service"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type EndPointServer struct {
	LoginEndPoint endpoint.Endpoint
}

func NewEndPointServer(srv service.Service, log *zap.Logger, limit *rate.Limiter) EndPointServer {
	var loginEndPoint endpoint.Endpoint
	{
		loginEndPoint = MakeLoginEndPoint(srv)
		loginEndPoint = LoggingMiddleware(log)(loginEndPoint)
		loginEndPoint = NewGolangRateAllowMiddleware(limit)(loginEndPoint)
	}
	return EndPointServer{LoginEndPoint: loginEndPoint}
}

func (e EndPointServer) Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error) {
	res, err := e.LoginEndPoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.LoginAck), nil
}

func MakeLoginEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.Login)
		return s.Login(ctx, req)
	}
}


