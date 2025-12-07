package cqrs

import (
	"context"
)

func ApplyQueryDecorator[Q any, R any](handler QueryHandler[Q, R]) QueryHandler[Q, R] {
	return queryLoggingDecorator[Q, R]{
		base: handler,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, cmd Q) (R, error)
}
