// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package hash

// PasswordHash password hashing interface
type PasswordHash interface {
	Name() string
	Value() string
	IsValid(storedHash string) bool
}
