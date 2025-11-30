package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var (
	ErrNotFound          = fmt.Errorf("sql row not found in the database")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Post interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, string) (*Post, error)
		Update(context.Context, int64, int64, *Post) error
		DeleteByID(context.Context, string) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]*PostWithMetadata, error)
	}
	User interface {
		Create(context.Context, *User) error
		CreateAndInvite(context.Context, *User) error
		CreateAndInviteTx(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		DeleteByID(context.Context, int64) error
		Activate(context.Context, string) error
	}
	Comment interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
	Follow interface {
		CreateFollow(context.Context, int64, int64) error
		DeleteFollow(context.Context, int64, int64) error
	}
	Invitation interface {
		CleanByID(context.Context, int64) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Post:       NewPostsStore(db),
		User:       NewUsersStore(db),
		Comment:    NewCommentsStore(db),
		Follow:     NewFollowsStore(db),
		Invitation: NewInvitationStore(db),
	}
}
