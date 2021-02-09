package client

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	uuid "github.com/satori/go.uuid"
	"go-kit-demo/demo03/v4/user_agent/pb"
	endpoint2 "go-kit-demo/demo03/v4/user_agent/src/endpoint"
	"go-kit-demo/demo03/v4/user_agent/src/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"time"
)

type UserAgent struct {
	instance *etcdv3.Instancer
	logger   log.Logger
}

func NewUserAgentClient(addr []string, logger log.Logger) (*UserAgent, error) {
	var (
		etcdAddrs = addr
		serName = "svc.user.agent"
		ttl = 5 * time.Second
	)
	options := etcdv3.ClientOptions{
		DialKeepAlive: ttl,
		DialTimeout: ttl,
	}
	etcdClient, err := etcdv3.NewClient(context.Background(), etcdAddrs, options)
	if err != nil {
		return nil, err
	}
	instance, err := etcdv3.NewInstancer(etcdClient, serName, logger)
	if err != nil {
		return nil, err
	}
	return &UserAgent{
		instance: instance,
		logger:   logger,
	}, err
}

func (u *UserAgent) UserAgentClient() (service.Service, error) {
	var (
		retryMax = 3
		retryTimeout = 5 * time.Second
	)
	var (
		endpoints endpoint2.EndPointServer
	)
	{
		factory := u.factoryFor(endpoint2.MakeLoginEndPoint)
		endpoint := sd.NewEndpointer(u.instance, factory, u.logger)
		balancer := lb.NewRoundRobin(endpoint)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.LoginEndPoint = retry
	}
	return endpoints, nil
}

func (u *UserAgent) factoryFor(makeEndpoint func(service service.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		fmt.Println("instance >>>>>>>>>>>>>>>>   ",instance)
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		srv := u.NewGRPCClient(conn)
		endpoints := makeEndpoint(srv)
		return endpoints, conn, err
	}
}

func (u *UserAgent) NewGRPCClient(conn *grpc.ClientConn) service.Service {
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(func(ctx context.Context, md *metadata.MD) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			md.Set(service.ContextReqUUid, UUID)
			ctx = metadata.NewOutgoingContext(ctx, *md)
			return ctx
		}),
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpctransport.NewClient(
				conn,
				"pb.User",
				"RpcUserLogin",
				u.RequestLogin,
				u.ResponseLogin,
				pb.LoginAck{},
				options...
			).Endpoint()
	}
	return endpoint2.EndPointServer{LoginEndPoint: loginEndpoint}
}

func (u *UserAgent) RequestLogin(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*service.Login)
	return &pb.Login{Account: req.Account, Password: req.Password}, nil
}

func (u *UserAgent) ResponseLogin(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*service.LoginAck)
	return &pb.LoginAck{Token: resp.Token}, nil
}