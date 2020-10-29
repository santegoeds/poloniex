package test

import (
	"context"
	"testing"
)

type contextKey string

const (
	ctKeyT contextKey = "ctx-key-poloniex-T"
)

func WithT(ctx context.Context, t *testing.T) context.Context {
	return context.WithValue(ctx, ctKeyT, t)
}

func T(ctx context.Context) *testing.T {
	val := ctx.Value(ctKeyT)
	return val.(*testing.T)
}
