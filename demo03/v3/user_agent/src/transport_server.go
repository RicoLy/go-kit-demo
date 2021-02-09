package src

import (
	"context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"go-kit-demo/demo03/v3/user_agent/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	login grpctransport.Handler
}

func NewGRPCServer(endpoint EndPointServer, log *zap.Logger) pb.UserServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerBefore(func(ctx context.Context, md metadata.MD) context.Context {
			ctx = context.WithValue(ctx, ContextReqUUid, md.Get(ContextReqUUid))
			return ctx
		}),
		grpctransport.ServerErrorHandler(NewZapLogErrorHandler(log)),
	}
	return &grpcServer{login: grpctransport.NewServer(
		endpoint.LoginEndPoint,
		func(_ context.Context, grpcReq interface{}) (request interface{}, err error) {
			req := grpcReq.(*pb.Login)
			return req, nil
		},
		func(_ context.Context, grpcResponse interface{}) (response interface{}, err error) {
			resp := grpcResponse.(*pb.LoginAck)
			return resp, nil
		},
		options...)}
}

func (s *grpcServer) RpcUserLogin(ctx context.Context, req *pb.Login) (*pb.LoginAck, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginAck), nil
}

//func RequestGrpcLogin(_ context.Context, grpcReq interface{}) (interface{}, error) {
//	req := grpcReq.(*pb.Login)
//	//return &pb.Login{Account: req.GetAccount(), Password: req.GetPassword()}, nil
//	return req, nil
//}

//func ResponseGrpcLogin(_ context.Context, response interface{}) (interface{}, error) {
//	resp := response.(*pb.LoginAck)
//	//return &pb.LoginAck{Token: resp.Token}, nil
//	return resp, nil
//}