package store

import (
	"context"
	"fmt"
)

type CredStore interface {
	Get(ctx context.Context, unitName, credID string) (string, error)
}

type MockCredStore struct {
}

func (c *MockCredStore) Get(_ context.Context, unitName, credID string) (string, error) {
	return fmt.Sprint(unitName, credID), nil
}
