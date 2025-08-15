package store

import (
	"context"
	"database/sql"
	"fmt"
)

var ErrNotFound = fmt.Errorf("sql row not found in the database")

type Storage struct {
	Post interface {
		Create(context.Context, *Post) error
	}
	User interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Post: NewPostsStore(db),
		User: NewUsersStore(db),
	}
}
