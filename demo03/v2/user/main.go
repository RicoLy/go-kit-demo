package main

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"go-kit-demo/demo03/utils"
	"go-kit-demo/demo03/v2/user/endpoint"
	"go-kit-demo/demo03/v2/user/pb"
	"go-kit-demo/demo03/v2/user/service"
	"go-kit-demo/demo03/v2/user/transport"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	utils.NewLoggerServer()
	golanglimit := rate.NewLimiter(10, 1)
	server := service.NewService(utils.GetLogger())
	endpoints := endpoint.NewEndPointServer(server, utils.GetLogger(), golanglimit)
	grpcServer := transport.NewGRPCServer(endpoints, utils.GetLogger())
	utils.GetLogger().Info("server run : 8881")
	grpcListener, err := net.Listen("tcp", ":8881")
	if err != nil {
		utils.GetLogger().Warn("Listen", zap.Error(err))
		os.Exit(0)
	}
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpctransport.Interceptor))
	pb.RegisterUserServer(baseServer, grpcServer)
	if err = baseServer.Serve(grpcListener); err != nil {
		utils.GetLogger().Warn("Serve", zap.Error(err))
		os.Exit(0)
	}
}