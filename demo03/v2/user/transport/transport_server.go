package transport

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"go-kit-demo/demo03/v2/user/endpoint"
	"go-kit-demo/demo03/v2/user/pb"
	"go-kit-demo/demo03/v2/user/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	login grpctransport.Handler
}

func NewGRPCServer(endpoint endpoint.EndPointServer, log *zap.Logger) pb.UserServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			ctx = context.WithValue(ctx, service.ContextReqUUid, md.Get(service.ContextReqUUid))
			return ctx
		}),
		grpctransport.ServerErrorHandler(NewZapLogErrorHandler(log)),
	}

	return &grpcServer{login: grpctransport.NewServer(
		endpoint.LoginEndPoint,
		RequestGrpcLogin,
		ResponseGrpcLogin,
		options...,
		)}
}

func (s *grpcServer) RpcUserLogin(ctx context.Context, req *pb.Login) (*pb.LoginAck, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginAck), nil
}

func RequestGrpcLogin(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.Login)
	return &pb.Login{Account: req.GetAccount(), Password: req.GetPassword()}, nil
}

func ResponseGrpcLogin(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.LoginAck)
	return &pb.LoginAck{
		Token: resp.Token,
	}, nil
}