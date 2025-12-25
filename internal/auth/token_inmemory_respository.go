package auth

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/oauth2"
)

var ErrTokenNotFound = errors.New("token not found")

type InMemoryTokenRepository struct {
	mu    sync.RWMutex
	store map[string]*oauth2.Token
}

func NewInMemoryTokenRepository() (*InMemoryTokenRepository, error) {
	return &InMemoryTokenRepository{
		store: make(map[string]*oauth2.Token, 5),
	}, nil
}

func (r *InMemoryTokenRepository) GetBySessionID(ctx context.Context, sessionID string) (*oauth2.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, exists := r.store[sessionID]

	if !exists {
		return nil, ErrTokenNotFound
	}

	return token, nil
}

func (r *InMemoryTokenRepository) Save(ctx context.Context, sessionID string, token *oauth2.Token) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[sessionID] = token
	return nil
}

func (r *InMemoryTokenRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[sessionID]; !exists {
		return ErrTokenNotFound
	}

	delete(r.store, sessionID)
	return nil
}
