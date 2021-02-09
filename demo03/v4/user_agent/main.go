package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	metricsprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/sd/etcdv3"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-kit-demo/demo03/utils"
	"go-kit-demo/demo03/v4/user_agent/pb"
	"go-kit-demo/demo03/v4/user_agent/src/endpoint"
	"go-kit-demo/demo03/v4/user_agent/src/service"
	"go-kit-demo/demo03/v4/user_agent/src/transpoint"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"hash/crc32"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var grpcAddr = flag.String("g", "127.0.0.1:8881", "grpcAddr")
var httpAddr = flag.String("h", "127.0.0.1:8882", "httpAddr")
var prometheusAddr = flag.String("p", "127.0.0.1:9100", "prometheus addr")

var quitChan = make(chan error, 1)

func main() {
	flag.Parse()
	var (
		etcdAddr = []string{"127.0.0.1:2379"}
		serName  = "svc.user.agent"
		ttl      = 5 * time.Second
	)
	utils.NewLoggerServer()
	options := etcdv3.ClientOptions{
		DialTimeout:   ttl,
		DialKeepAlive: ttl,
	}
	etcdClient, err := etcdv3.NewClient(context.Background(), etcdAddr, options)
	if err != nil {
		utils.GetLogger().Error("[user_agent]  NewClient", zap.Error(err))
		return
	}
	Register := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
		Key:   fmt.Sprintf("%s/%d", serName, crc32.ChecksumIEEE([]byte(*grpcAddr))),
		Value: *grpcAddr,
	}, log.NewNopLogger())
	count := metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
		Subsystem: "user_agent",
		Name:      "request_count",
		Help:      "Number of requests",
	}, []string{"method"})

	histogram := metricsprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Subsystem: "user_agent",
		Name:      "request_consume",
		Help:      "Request consumes time",
	}, []string{"method"})
	golangLimit := rate.NewLimiter(10, 1)
	server := service.NewService(utils.GetLogger(), count, histogram)
	endpoints := endpoint.NewEndPointServer(server, golangLimit)

	go func() {
		grpcServer := transpoint.NewGRPCServer(endpoints, utils.GetLogger())
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			utils.GetLogger().Warn("[user_agent] Listen", zap.Error(err))
			quitChan <- err
			return
		}
		Register.Register()
		utils.GetLogger().Info("[user_agent] grpc run " + *grpcAddr)
		baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpctransport.Interceptor))
		pb.RegisterUserServer(baseServer, grpcServer)
		quitChan <- baseServer.Serve(grpcListener)
	}()
	go func() {
		httpServer := transpoint.NewHttpHandler(endpoints, utils.GetLogger())
		_ = http.ListenAndServe(*httpAddr, httpServer)
	}()
	go func() {
		utils.GetLogger().Info("[user_agent] prometheus run " + *prometheusAddr)
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		quitChan <- http.ListenAndServe(*prometheusAddr, m)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		quitChan <- fmt.Errorf("%s", <-c)
	}()
	err = <-quitChan
	Register.Deregister()
	utils.GetLogger().Info("[user_agent] quit", zap.Any("info", err))
}
