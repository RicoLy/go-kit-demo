package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go-kit-demo/version1/utils"
	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
	"go-kit-demo/version1/v3/endpoint"
	"go-kit-demo/version1/v3/service"
	"go.uber.org/zap"
	"net/http"
)

var ErrorBadRequest = errors.New("invalid request parameter")

func MakeHttpHandler(endpoints endpoint.RegisterEndpoints, log *zap.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
			log.Warn(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Error(err))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
		}),
		httptransport.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			log.Debug("给请求添加uuid", zap.Any("UUID", UUID))
			ctx = context.WithValue(ctx, service.ContextReqUUid, UUID)
			ctx = context.WithValue(ctx, utils.JWT_CONTEXT_KEY, request.Header.Get("Authorization"))
			log.Debug("把请求中的token发到Context中", zap.Any("Token", request.Header.Get("Authorization")))

			return ctx
		}),
	}

	r := mux.NewRouter()
	r.Methods("GET").Path("./health").Handler(httptransport.NewServer(
			endpoints.HealthCheckEndPoint,
			decodeHTTPHealthCheckRequest,
			encodeHttpGenericResponse,
			options...,
		))
	r.Methods("GET").Path("./discovery/name").Handler(httptransport.NewServer(
		endpoints.DiscoveryEndPoint,
		decodeHTTPDiscoveryRequest,
		encodeHttpGenericResponse,
		options...,
	))
	return r
}

func decodeHTTPHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req service.HealthCheckRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Any(" 开始解析请求数据", req))
	return req, nil
}

func decodeHTTPDiscoveryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req service.DiscoveryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Any(" 开始解析请求数据", req))
	return req, nil
}

func encodeHttpGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Any("请求结束封装返回值", response))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"errors"`
}