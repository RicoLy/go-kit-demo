package transpoint

import (
	"context"
	"encoding/json"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
	"go-kit-demo/demo03/utils"
	"go-kit-demo/demo03/v4/user_agent/src/endpoint"
	"go-kit-demo/demo03/v4/user_agent/src/service"
	"go.uber.org/zap"
	"net/http"
)

func NewHttpHandler(endpoint endpoint.EndPointServer, log *zap.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
			log.Warn(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Error(err))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
		}), //程序中的全部报错都会走这里面
		httptransport.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), "req_uuid").String()
			log.Debug("给请求添加uuid", zap.Any("UUID", UUID))
			ctx = context.WithValue(ctx, service.ContextReqUUid, UUID)
			ctx = context.WithValue(ctx, utils.JWT_CONTEXT_KEY, request.Header.Get("Authorization"))
			log.Debug("把请求中的token发到Context中", zap.Any("Token", request.Header.Get("Authorization")))

			return ctx
		}),
	}
	m := http.NewServeMux()
	m.Handle("/login", httptransport.NewServer(
		endpoint.LoginEndPoint,
		decodeHTTPLoginRequest,
		encodeHTTPGenericResponse,
		options...,
	))

	return m
}

func decodeHTTPLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var login service.Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		return nil, err
	}
	return &login, nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"errors"`
}
