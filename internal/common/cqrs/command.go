package cqrs

import (
	"context"
	"fmt"
	"strings"
)

func ApplyCommandDecorators[C any, R any](handler CommandHandler[C, R]) CommandHandler[C, R] {
	return commandLoggingDecorator[C, R]{
		base: handler,
	}
}

type CommandHandler[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
