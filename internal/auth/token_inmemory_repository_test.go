package auth_test

import (
	"context"
	"fmt"
	"sumni-finance-backend/internal/auth"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestInMemoryAuthRepository(t *testing.T) {
	t.Run("should save token then get by sessionID successfully", func(t *testing.T) {
		t.Parallel()

		authRepo, _ := auth.NewInMemoryTokenRepository()
		sessionID := uuid.New().String()
		token := &oauth2.Token{
			AccessToken: "mock-access-token",
			Expiry:      time.Now().Add(time.Hour),
		}

		err := authRepo.Save(context.Background(), sessionID, token)
		require.NoError(t, err)

		got, err := authRepo.GetBySessionID(context.Background(), sessionID)
		assert.NoError(t, err)
		assertToken(t, token, got)
	})

	t.Run("should return error when sessionID does not exist", func(t *testing.T) {
		t.Parallel()

		authRepo, _ := auth.NewInMemoryTokenRepository()

		got, err := authRepo.GetBySessionID(context.Background(), "non-existent")
		assert.Nil(t, got)
		assert.ErrorIs(t, err, auth.ErrTokenNotFound)
	})

	t.Run("should delete token successfully", func(t *testing.T) {
		authRepo, _ := auth.NewInMemoryTokenRepository()
		sessionID := "session-to-delete"
		token := &oauth2.Token{AccessToken: "secret"}

		err := authRepo.Save(context.Background(), sessionID, token)
		assert.NoError(t, err)

		// Delete
		err = authRepo.DeleteBySessionID(context.Background(), sessionID)
		assert.NoError(t, err)

		// Verify it's gone
		got, err := authRepo.GetBySessionID(context.Background(), sessionID)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, auth.ErrTokenNotFound)
	})
}

func TestInmemoryAuthRepository_Race(t *testing.T) {
	authRepo, err := auth.NewInMemoryTokenRepository()
	require.NoError(t, err)

	ctx := context.Background()
	wg := sync.WaitGroup{}

	iterations := 50
	for i := 0; i < iterations; i++ {
		wg.Add(3)

		go func(id int) {
			defer wg.Done()

			sessionID := fmt.Sprintf("session-%d", id)
			_ = authRepo.Save(ctx, sessionID, &oauth2.Token{AccessToken: "secret"})
		}(i)

		go func(id int) {
			defer wg.Done()
			sessionID := fmt.Sprintf("session-%d", id)
			_, _ = authRepo.GetBySessionID(ctx, sessionID)
		}(i)

		go func(id int) {
			defer wg.Done()
			sessionID := fmt.Sprintf("session-%d", id)
			_ = authRepo.DeleteBySessionID(ctx, sessionID)
		}(i)
	}

	wg.Wait()
}

func assertToken(t *testing.T, expected, actual *oauth2.Token) {
	t.Helper()

	opts := []cmp.Option{
		cmpopts.EquateApproxTime(time.Second),
		cmp.AllowUnexported(oauth2.Token{}),
	}
	diff := cmp.Diff(expected, actual, opts...)
	if diff != "" {
		assert.Fail(t, "Tokens are not deeply equal", "Diff (-expected +actual):\n%s", diff)
	}
}
