// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package repository

import (
	"errors"

	"github.com/hashicorp/go-memdb"
)

type passwordMemoryStorage struct {
	db *memdb.MemDB
}

// PasswordMemoryStorage initiate password memory storage
func PasswordMemoryStorage(db *memdb.MemDB) *passwordMemoryStorage {
	return &passwordMemoryStorage{db: db}
}

// FindByID find credential by its id from memory storage
func (c *passwordMemoryStorage) FindByID(id string) (*PasswordModel, error) {
	tx := c.db.Txn(false)
	defer tx.Abort()
	data, err := tx.First(PasswordTable, "id", id)
	if err != nil {
		return &PasswordModel{}, err
	}
	tx.Commit()

	if data == nil {
		return &PasswordModel{}, errors.New("NOT_FOUND")
	}

	return data.(*PasswordModel), nil
}

// Create save credential to memory storage
func (c *passwordMemoryStorage) Create(data *PasswordModel) error {
	tx := c.db.Txn(true)
	if err := tx.Insert(PasswordTable, data); err != nil {
		return err
	}
	tx.Commit()

	return nil
}

// Update delete and re-create credential to memory storage
func (c *passwordMemoryStorage) Update(data *PasswordModel) error {
	update := &PasswordModel{}
	*update = *data
	if err := c.Delete(data); err != nil {
		return err
	}

	return c.Create(update)
}

// Delete delete credential from memory storage
func (c *passwordMemoryStorage) Delete(data *PasswordModel) error {
	tx := c.db.Txn(true)
	defer tx.Abort()
	if err := tx.Delete(PasswordTable, data); err != nil {
		return err
	}
	tx.Commit()

	return nil
}
