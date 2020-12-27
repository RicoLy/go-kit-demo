package main

import (
	"fmt"
	"go-kit-demo/version1/utils"
	"go-kit-demo/version1/v2/user/dao"
	"go-kit-demo/version1/v2/user/endpoint"
	"go-kit-demo/version1/v2/user/redis"
	"go-kit-demo/version1/v2/user/service"
	"go-kit-demo/version1/v2/user/transport"

	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	utils.NewLoggerServer()

	errChan := make(chan error)
	logger := utils.GetLogger()

	err := dao.InitMysql("192.168.190.139", "3306", "root", "mysqlly", "user")
	if err != nil{
		logger.Debug("err:" ,zap.Any(" mysql", err))
	}

	err = redis.InitRedis("192.168.190.139","6379", "" )
	if err != nil{
		logger.Debug("err:" ,zap.Any(" redis", err))
	}

	service := service.NewUserServiceImpl(&dao.UserDAOImpl{}, logger)

	userEndpoints := endpoint.NewUserEndpoints(service, logger)

	r := transport.NewHttpHandler(logger, userEndpoints)

	go func() {
		errChan <- http.ListenAndServe(":10081", r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	err = <-errChan
	logger.Debug("err: ", zap.Any("error", err))
}