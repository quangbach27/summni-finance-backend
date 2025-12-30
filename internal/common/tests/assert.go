package tests

import (
	"sumni-finance-backend/internal/common/server/httperr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertSlugError(
	t *testing.T,
	err error,
	slugMsg string,
	wrappedErr error,
) {
	t.Helper()

	assert.Error(t, err)

	var wantSlugErr httperr.SlugError
	assert.ErrorAs(t, err, &wantSlugErr)
	if slugMsg != "" {
		assert.Equal(t, wantSlugErr.Slug(), slugMsg)
	}
	if wrappedErr != nil {
		assert.ErrorIs(t, wantSlugErr.Unwrap(), wrappedErr)
	}
}
