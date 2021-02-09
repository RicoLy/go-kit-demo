package src

import (
	"context"
	"fmt"
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
	h.logger.Warn(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Error(err))
}
