package main

import (
	"context"
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go-kit-demo/version1/utils"
	"go-kit-demo/version1/v3/discovery"
	"go-kit-demo/version1/v3/endpoint"
	"go-kit-demo/version1/v3/service"
	"go-kit-demo/version1/v3/transport"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	consulAddr := flag.String("consul.addr", "localhost", "consul address")
	consulPort := flag.Int("consul.port", 8500, "consul port")
	serviceName := flag.String("service.name", "register", "service name")
	serviceAddr := flag.String("service.addr", "localhost", "service addr")
	servicePort := flag.Int("service.port", 12312, "service port")

	flag.Parse()

	client := discovery.NewDiscoveryClient(*consulAddr, *consulPort)

	errChan := make(chan error)
	srv := service.NewRegisterServiceImpl(client)

	endPoints := endpoint.NewRegisterEndPoints(srv)
	handler := transport.MakeHttpHandler(endPoints, utils.GetLogger())
	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	instanceId := *serviceName + "-" + uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()

	err := client.Register(context.Background(), *serviceName, instanceId, "/health", *serviceAddr, *servicePort, nil, nil)
	if err != nil {
		utils.GetLogger().Debug("register service err : ", zap.Any("err", err))
		os.Exit(-1)
	}

	err = <-errChan
	utils.GetLogger().Debug("listen", zap.Any(" err :", err))
	client.Deregister(context.Background(), instanceId)
}
