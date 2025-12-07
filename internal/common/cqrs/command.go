package cqrs

import (
	"context"
	"fmt"
	"strings"
)

func ApplyCommandDecorators[C any](handler CommandHandler[C]) CommandHandler[C] {
	return commandLoggingDecorator[C]{
		base: handler,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
