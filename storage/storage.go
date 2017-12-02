// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package storage

import (
	"github.com/simukti/grpc-password/repository"
)

// Storage storage object and its repositories
type Storage interface {
	// Password get password repository
	Password() repository.PasswordRepository
}
