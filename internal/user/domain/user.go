// Package domain defines the core user domain entities and types.
package domain

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name" fieldtag:"insert,update"`
	Email     string    `db:"email" fieldtag:"insert,update"`
	Password  string    `db:"password" fieldtag:"insert,update"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
