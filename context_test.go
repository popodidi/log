package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	require.Equal(t, Null, GetFromCtx(ctx))

	ctx, l := Context(context.Background(), "test")
	require.NotEqual(t, Null, l)
	id := l.GetID()

	_, ctxl := Context(ctx, "test")
	require.NotEqual(t, Null, ctxl)
	require.Equal(t, id, ctxl.GetID())

	ctxl = GetFromCtx(ctx)
	require.NotEqual(t, Null, ctxl)
	require.Equal(t, id, ctxl.GetID())
}
