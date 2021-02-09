package client

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"go-kit-demo/demo03/all_packaged_library/logtool"
	"go-kit-demo/demo03/v2/user/pb"
	"go-kit-demo/demo03/v2/user/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestNewGrpcClient(t *testing.T) {
	logger := logtool.NewLogger(
		logtool.SetAppName("go-kit"),
		logtool.SetDevelopment(true),
		logtool.SetLevel(zap.DebugLevel),
	)
	conn, err := grpc.Dial("127.0.0.1:8881", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()
	svr := NewGrpcClient(conn, logger)
	ack, err := svr.Login(context.Background(), &pb.Login{
		Account:  "ryan",
		Password: "123",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ack.Token)
}

func TestGrpc(t *testing.T) {
	serviceAddress := "127.0.0.1:8881"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	userClient := pb.NewUserClient(conn)
	UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
	md := metadata.Pairs(service.ContextReqUUid, UUID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	res, err := userClient.RpcUserLogin(ctx, &pb.Login{
		Account:  "ryan",
		Password: "123",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.Token)
}