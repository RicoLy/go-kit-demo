package service

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

type metricsMiddlewareServer struct {
	next      Service
	counter   metrics.Counter
	histogram metrics.Histogram
}

func NewMetricsMiddlewareServer(counter metrics.Counter, histogram metrics.Histogram) NewMiddlewareServer {
	return func(service Service) Service {
		return metricsMiddlewareServer{
			next:      service,
			counter:   counter,
			histogram: histogram,
		}
	}
}

func (m metricsMiddlewareServer) Login(ctx context.Context, in *Login) (ack *LoginAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "login"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	return m.next.Login(ctx, in)
}



