package src

import (
	"context"
	"errors"
	"go-kit-demo/demo03/utils"
	"go-kit-demo/demo03/v3/user_agent/pb"
	"go.uber.org/zap"
)

type Service interface {
	Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error)
}

type baseServer struct {
	logger *zap.Logger
}

func NewService(log *zap.Logger) Service {
	var service Service
	service = &baseServer{log}
	service = NewLogMiddlewareServer(log)(service)
	return service
}

func (s baseServer) Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error) {
	if in.Account != "ryan" || in.Password != "123" {
		err = errors.New("用户信息错误")
		return
	}
	ack = &pb.LoginAck{}
	ack.Token, err = utils.CreateJwtToken(in.Account, 1)
	return
}


















