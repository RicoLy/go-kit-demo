package service

import (
	"context"
	"errors"
	"github.com/go-kit/kit/metrics"
	"go-kit-demo/demo03/utils"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type Service interface {
	Login(ctx context.Context, in *Login) (ack *LoginAck, err error)
}

type Login struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginAck struct {
	Token string `json:"token"`
}

type baseServer struct {
	logger *zap.Logger
}

func NewService(log *zap.Logger, counter metrics.Counter, histogram metrics.Histogram) Service {
	var server Service
	server = &baseServer{log}
	server = NewLogMiddlewareServer(log)(server)
	server = NewMetricsMiddlewareServer(counter, histogram)(server)
	return server
}

func (s baseServer) Login(ctx context.Context, in *Login) (ack *LoginAck, err error) {
	if in.Account != "ryan" || in.Password != "123" {
		err = errors.New("用户信息错误")
		return
	}
	//模拟耗时
	rand.Seed(time.Now().UnixNano())
	sl := rand.Int31n(10-1) + 1
	time.Sleep(time.Duration(sl) * time.Millisecond * 100)
	ack = &LoginAck{}
	ack.Token, err = utils.CreateJwtToken(in.Account, 1)
	return
}



