package service

import (
	"context"
	"errors"
	"go-kit-demo/version1/v3/discovery"
	"log"
)

const ContextReqUUid = "req_uuid"

type HealthCheckRequest struct {
	ServiceName string 	`json:"service_name"`
}

type HealthCheckResponse struct {
	HealthResponse string `json:"health_response"`
}

type DiscoveryResponse struct {
	Instances []*discovery.InstanceInfo `json:"instances"`
}

type DiscoveryRequest struct {
	ServiceName string 	`json:"service_name"`
}

type Service interface {
	HealthCheck(ctx context.Context, checkReq HealthCheckRequest) (HealthCheckResponse, error)
	DiscoveryService(ctx context.Context, discoveryRequest DiscoveryRequest) (DiscoveryResponse, error)
}

var ErrNotServiceInstances = errors.New("instances are not existed")

func NewRegisterServiceImpl(discoveryClient *discovery.DiscoveryClient) Service {
	return &RegisterServiceImpl{disconveryClient:discoveryClient}
}

type RegisterServiceImpl struct {
	disconveryClient *discovery.DiscoveryClient
}

func (s *RegisterServiceImpl) HealthCheck(ctx context.Context, checkReq HealthCheckRequest) (HealthCheckResponse, error) {
	return HealthCheckResponse{HealthResponse:"ok"}, nil
}

func (s *RegisterServiceImpl) DiscoveryService(ctx context.Context, discoveryRequest DiscoveryRequest) (DiscoveryResponse, error) {
	instances, err := s.disconveryClient.DiscoverServices(ctx, discoveryRequest)
	if err != nil {
		log.Printf("get service info err: %s", err)
	}
	if instances == nil || len(instances) == 0 {
		return DiscoveryResponse{}, ErrNotServiceInstances
	}
	return DiscoveryResponse{Instances:instances}, nil
}

