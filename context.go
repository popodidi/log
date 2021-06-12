package log

import (
	"context"
)

type ctxKey struct{}

// CtxKey is the key for logger stored in context.
func CtxKey() interface{} { return ctxKey{} }

// Context gets from context or creates a logger with tags.
func Context(ctx context.Context, tags ...string) (context.Context, Logger) {
	l := GetFromCtx(ctx)
	if l != Null {
		return ctx, l
	}
	l = New(tags...)
	ctx = context.WithValue(ctx, CtxKey(), l)
	return ctx, l
}

// GetFromCtx returns logger from context, nil for no logger found.
func GetFromCtx(ctx context.Context) Logger {
	val := ctx.Value(CtxKey())
	if val == nil {
		return Null
	}
	return val.(Logger)
}
