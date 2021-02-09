package transport

import (
	"context"
	"fmt"
	"go-kit-demo/demo03/v2/user/service"
	"go.uber.org/zap"
)

type LogErrorHandler struct {
	logger *zap.Logger
}

func NewZapLogErrorHandler(logger *zap.Logger) *LogErrorHandler {
	return &LogErrorHandler{
		logger: logger,
	}
}

func (h *LogErrorHandler) Handle(ctx context.Context, err error) {
	h.logger.Warn(fmt.Sprint(ctx.Value(service.ContextReqUUid)), zap.Error(err))
}
