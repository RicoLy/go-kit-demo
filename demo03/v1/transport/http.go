package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
	"go-kit-demo/demo03/v1/endpoint"
	"go-kit-demo/demo03/v1/service"
	"go-kit-demo/demo03/v1/utils"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

func NewHttpHandler(log *zap.Logger, endpoints endpoint.UserEndpoints) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
			log.Warn(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Error(err))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(errorWrapper{Error:err.Error()})
		}),
		httptransport.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			log.Debug("给请求添加uuid", zap.Any("UUID", UUID))
			ctx = context.WithValue(ctx, service.ContextReqUUid, UUID)
			log.Debug("把请求中的token发到Context中", zap.Any("Token", request.Header.Get("Authorization")))

			return ctx
		}),
	}

	m := http.NewServeMux()
	m.Handle("/login", httptransport.NewServer(
		endpoints.LoginEndpoint,
		decodeHTTPLoginRequest,
		encodeHTTPGenericResponse,
		options...
	))

	m.Handle("/register", httptransport.NewServer(
		endpoints.RegisterEndpoint,
		decodeHTTPRegisterRequest,
		encodeHTTPGenericResponse,
		options...
	))

	return m
}

func decodeHTTPRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var vo service.RegisterUserVO
	vo.Password = r.FormValue("password")
	vo.Username = r.FormValue("username")
	vo.Email = r.FormValue("email")
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Any(" 开始解析请求数据", vo))
	if vo.Email == "" || vo.Username == "" || vo.Password == "" {
		return nil, ErrorBadRequest
	}
	return &vo, nil
}

func decodeHTTPLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var vo service.LoginVO
	err := json.NewDecoder(r.Body).Decode(&vo)
	if err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Any(" 开始解析请求数据", vo))
	if vo.Email == "" ||  vo.Password == "" {
		return nil, ErrorBadRequest
	}
	return &vo, nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Any("请求结束封装返回值", response))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"errors"`
}