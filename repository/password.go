// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package repository

import "time"

var PasswordTable = "password"

// PasswordModel password table data model
type PasswordModel struct {
	ID        string    `db:"id"`
	Type      string    `db:"type"`
	Hash      string    `db:"hash"`
	Secret    string    `db:"secret"`
	CreatedAt time.Time `db:"created_at"`
}

// Repository password table repository interface
type PasswordRepository interface {
	// FindByID get one password from storage object
	FindByID(id string) (*PasswordModel, error)
	// Create add new password
	Create(data *PasswordModel) error
	// Update update password
	Update(data *PasswordModel) error
	// DeleteById delete one password from storage
	Delete(data *PasswordModel) error
}
