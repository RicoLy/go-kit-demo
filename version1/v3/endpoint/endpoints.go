package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go-kit-demo/version1/v3/service"
)

type RegisterEndpoints struct {
	DiscoveryEndPoint   endpoint.Endpoint
	HealthCheckEndPoint endpoint.Endpoint
}

func NewRegisterEndPoints(svc service.Service) RegisterEndpoints {
	return RegisterEndpoints{
		DiscoveryEndPoint:   MakeDiscoveryEndpoint(svc),
		HealthCheckEndPoint: MakeHealthCheckEndpoint(svc),
	}
}

func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.DiscoveryRequest)
		return svc.DiscoveryService(ctx, req)
	}
}

func (s RegisterEndpoints) HealthCheck(ctx context.Context, checkReq service.HealthCheckRequest) (service.HealthCheckResponse, error) {
	resp , err := s.HealthCheckEndPoint(ctx, checkReq)
	return resp.(service.HealthCheckResponse), err
}

func MakeDiscoveryEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.DiscoveryRequest)
		return svc.DiscoveryService(ctx, req)
	}
}

func (s RegisterEndpoints) DiscoveryService(ctx context.Context, discoveryRequest service.DiscoveryRequest) (service.DiscoveryResponse, error) {
	resp, err := s.DiscoveryEndPoint(ctx, discoveryRequest)
	return resp.(service.DiscoveryResponse), err
}

