package store

import (
	"context"
)

func NewMockStorage() *Storage {
	return &Storage{
		User: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (mus *MockUserStore) Create(ctx context.Context, u *User) error {
	return nil
}
func (mus *MockUserStore) CreateAndInvite(ctx context.Context, u *User) error {
	return nil
}
func (mus *MockUserStore) CreateAndInviteTx(ctx context.Context, u *User) error {
	return nil
}
func (mus *MockUserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	return &User{}, nil
}
func (mus *MockUserStore) GetByEmail(ctx context.Context, em string) (*User, error) {
	return &User{}, nil
}
func (mus *MockUserStore) DeleteByID(ctx context.Context, id int64) error {
	return nil
}
func (mus *MockUserStore) Activate(ctx context.Context, plainToken string) error {
	return nil
}
