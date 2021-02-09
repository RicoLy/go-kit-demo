package service

import (
	"context"
	"errors"
	"fmt"
	"go-kit-demo/demo03/utils"
	"go-kit-demo/demo03/v2/user/pb"
	"go.uber.org/zap"
)

type Service interface {
	Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error)
}

type baseServer struct {
	logger *zap.Logger
}

func NewService(log *zap.Logger) Service {
	var server Service
	server = &baseServer{log}
	server = NewLogMiddleWareServer(log)(server)
	return server
}

func (s baseServer) Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v5_service Service", "Login 处理请求"))
	if in.Account != "ryan" || in.Password != "123" {
		err = errors.New("用户信息错误")
		return
	}
	ack = &pb.LoginAck{}
	ack.Token, err = utils.CreateJwtToken(in.Account, 1)
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 v5_service Service", "Login 处理请求"), zap.Any("处理返回值", ack))
	return
}

