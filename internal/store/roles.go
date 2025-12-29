package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RolesStore struct {
	db *sql.DB
}

func NewRolesStore(db *sql.DB) *RolesStore {
	return &RolesStore{db: db}
}

func (s *RolesStore) GetByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	query := `
	SELECT id, name, level, description
	FROM roles
	WHERE name = $1
	`

	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &role, nil
}

func (s *RolesStore) IsPrecedent(ctx context.Context, roleID int, requiredRole string) (bool, error) {
	var allowed bool
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM roles
		WHERE id = $1 AND level >= (
			SELECT level
			FROM roles
			WHERE name = $2
		)
	)
	`

	err := s.db.QueryRowContext(ctx, query, int64(roleID), requiredRole).Scan(&allowed)
	if err != nil {
		return false, err
	}
	return allowed, nil
}
