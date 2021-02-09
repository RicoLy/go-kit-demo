package client

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	uuid "github.com/satori/go.uuid"
	endpoint2 "go-kit-demo/demo03/v2/user/endpoint"
	"go-kit-demo/demo03/v2/user/pb"
	"go-kit-demo/demo03/v2/user/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewGrpcClient(conn *grpc.ClientConn, log *zap.Logger) service.Service {
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(func(ctx context.Context, md *metadata.MD) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			log.Debug("给请求添加uuid", zap.Any("UUID", UUID))
			md.Set(service.ContextReqUUid, UUID)
			ctx = metadata.NewOutgoingContext(context.Background(), *md)
			return ctx
		}),
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"RpcUserLogin",
			RequestLogin,
			ResponseLogin,
			pb.LoginAck{},
			options...).Endpoint()
	}
	return endpoint2.EndPointServer{
		LoginEndPoint: loginEndpoint,
	}
}

func RequestLogin(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.Login)
	return &pb.Login{Account: req.Account, Password: req.Password}, nil
}

func ResponseLogin(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.LoginAck)
	return &pb.LoginAck{Token: resp.Token}, nil
}
