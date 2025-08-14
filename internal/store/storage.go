package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Post interface {
		Create(context.Context) error
	}
	User interface {
		Create(context.Context) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Post: NewPostsStore(db),
		User: NewUsersStore(db),
	}
}
