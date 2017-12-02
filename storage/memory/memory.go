// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package memory

import (
	"log"

	"github.com/hashicorp/go-memdb"
	"github.com/simukti/grpc-password/repository"
	"github.com/simukti/grpc-password/storage"
)

type memoryStorage struct {
	db *memdb.MemDB
}

// Storage create memory storage object
func Storage() storage.Storage {
	schema := migrations()
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("In-memory storage")
	return &memoryStorage{
		db: db,
	}
}

// Password in-memory PasswordRepository
func (c *memoryStorage) Password() repository.PasswordRepository {
	return repository.PasswordMemoryStorage(c.db)
}
