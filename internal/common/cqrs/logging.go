package cqrs

import (
	"context"
	"sumni-finance-backend/internal/common/logs"
)

type commandLoggingDecorator[C any] struct {
	base CommandHandler[C]
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	logger := logs.FromContext(ctx).With(
		"command", handlerType,
	)

	logger.Info("Execute command")

	defer func() {
		if err != nil {
			logger.Error("Failed to execute command", "error", err)
		} else {
			logger.Info("Command executed successfully")
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[Q any, R any] struct {
	base QueryHandler[Q, R]
}

func (d queryLoggingDecorator[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {
	handlerType := generateActionName(query)

	logger := logs.FromContext(ctx).With(
		"query", handlerType,
	)

	logger.Debug("Execute query")

	defer func() {
		if err != nil {
			logger.Error("Failed to execute query", "error", err)
		} else {
			logger.Info("Query executed successfully")
		}
	}()

	return d.base.Handle(ctx, query)
}
