package store

import "time"

type Invitation struct {
	ID     int64     `json:"id"`
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
	UserID int64     `json:"user_id"`
}
