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
		GetByID(context.Context, string) (*Post, error)
		Update(context.Context, int64, *Post) error
		DeleteByID(context.Context, string) error
	}
	User interface {
		Create(context.Context, *User) error
	}
	Comment interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Post:    NewPostsStore(db),
		User:    NewUsersStore(db),
		Comment: NewCommentsStore(db),
	}
}
