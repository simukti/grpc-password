// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package hash

import (
	"github.com/simukti/grpc-password/util"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
)

// BcryptName hash name
const BcryptName = "BCRYPT"

type bcryptHash struct {
	name          string
	plainPassword string
	secretString  string
	cost          int
}

// Bcrypt create bcrypt password object.
func Bcrypt(plainPassword, secretString string) *bcryptHash {
	return &bcryptHash{
		name:          BcryptName,
		plainPassword: plainPassword,
		secretString:  secretString,
		// https://security.stackexchange.com/a/83382/151014
		// https://labs.clio.com/bcrypt-cost-factor-4ca0a9b03966
		cost: 13,
	}
}

// prehash create HMAC of plaintext password with secret as key and set result to fit in BLAKE2b-384 (60 chars of ascii-85).
func (c *bcryptHash) prehash() string {
	passSum := blake2b.Sum512(util.ASCII85Byte([]byte(c.plainPassword)))
	secretSum := blake2b.Sum512(util.ASCII85Byte([]byte(c.secretString)))
	prehash, _ := blake2b.New384(secretSum[:])
	prehash.Write(passSum[:])
	pwd := util.ASCII85String(prehash.Sum(nil))

	return pwd
}

// Name hash name
func (c *bcryptHash) Name() string {
	return c.name
}

// Value bcrypt hash from internally-prehashed password string.
func (c *bcryptHash) Value() string {
	b, _ := bcrypt.GenerateFromPassword([]byte(c.prehash()), c.cost)

	return string(b)
}

// IsValid validate stored hashed password against internally-prehashed password string.
func (c *bcryptHash) IsValid(storedHash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(c.prehash())); err != nil {
		return false
	}

	return true
}
