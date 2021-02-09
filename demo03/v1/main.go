package main

import (
	"fmt"
	"go-kit-demo/demo03/v1/dao"
	"go-kit-demo/demo03/v1/endpoint"
	"go-kit-demo/demo03/v1/redis"
	"go-kit-demo/demo03/v1/service"
	"go-kit-demo/demo03/v1/transport"
	"go-kit-demo/demo03/v1/utils"
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

	err := dao.InitMysql("192.168.142.128", "3306", "root", "mysqlly", "user")
	if err != nil {
		logger.Debug("err:" ,zap.Any(" mysql", err))
	}

	err = redis.InitRedis("192.168.142.128","6379", "" )
	if err != nil{
		logger.Debug("err:" ,zap.Any(" redis", err))
	}
	userDao := dao.NewUserDaoImpl()
	service := service.NewUserServiceImpl(userDao, logger)

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