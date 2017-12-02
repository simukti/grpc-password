// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package repository

import (
	"errors"

	"github.com/gocraft/dbr"
)

type passwordMySQLStorage struct {
	db *dbr.Session
}

// PasswordMysqlStorage initiate password mysql storage
func PasswordMysqlStorage(db *dbr.Session) *passwordMySQLStorage {
	return &passwordMySQLStorage{db: db}
}

// FindByID get one password from storage object
func (c *passwordMySQLStorage) FindByID(id string) (*PasswordModel, error) {
	result := &PasswordModel{}
	err := c.db.Select("*").
		From(PasswordTable).
		Where("id = ?", id).
		LoadStruct(result)

	if err != nil {
		if err == dbr.ErrNotFound {
			return result, errors.New("NOT_FOUND")
		}

		return result, err
	}

	return result, nil
}

// Create add new password
func (c *passwordMySQLStorage) Create(data *PasswordModel) error {
	_, err := c.db.InsertInto(PasswordTable).
		Columns("id", "type", "hash", "secret").
		Record(data).
		Exec()

	return err
}

// Update update password
func (c *passwordMySQLStorage) Update(data *PasswordModel) error {
	_, err := c.db.Update(PasswordTable).
		Set("hash", data.Hash).
		Where("id = ?", data.ID).
		Exec()

	return err
}

// DeleteById delete one password from storage
func (c *passwordMySQLStorage) Delete(data *PasswordModel) error {
	_, err := c.db.DeleteFrom(PasswordTable).
		Where("id = ?", data.ID).
		Exec()

	return err
}
