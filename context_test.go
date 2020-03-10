package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	require.Nil(t, GetFromCtx(ctx))

	ctx, l := Context(context.Background(), "test")
	require.NotNil(t, l)
	id := l.GetID()

	_, ctxl := Context(ctx, "test")
	require.NotNil(t, ctxl)
	require.Equal(t, id, ctxl.GetID())

	ctxl = GetFromCtx(ctx)
	require.NotNil(t, ctxl)
	require.Equal(t, id, ctxl.GetID())
}
