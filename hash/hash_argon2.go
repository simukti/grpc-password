// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.
// hashing related method copied from https://github.com/simukti/passwd/blob/master/passwd.go

package hash

// #cgo pkg-config: libsodium
// #include <stdlib.h>
// #include <sodium.h>
import "C"
import (
	"unsafe"

	"github.com/simukti/grpc-password/util"
	"golang.org/x/crypto/blake2b"
)

// Argon2Name hash name
const Argon2Name = "ARGON2"

var (
	OpsLimitInteractive = ArgonLimit(C.crypto_pwhash_opslimit_interactive())
	MemLimitInteractive = ArgonLimit(C.crypto_pwhash_memlimit_interactive())
	OpsLimitModerate    = ArgonLimit(C.crypto_pwhash_opslimit_moderate())
	MemLimitModerate    = ArgonLimit(C.crypto_pwhash_memlimit_moderate())
	OpsLimitSensitive   = ArgonLimit(C.crypto_pwhash_opslimit_sensitive())
	MemLimitSensitive   = ArgonLimit(C.crypto_pwhash_memlimit_sensitive())

	argon2iStringBytesLen = int(C.crypto_pwhash_strbytes())
)

// ArgonLimit operation and memory limit type
type ArgonLimit int

type argon2Hash struct {
	name          string
	plainPassword string
	secretString  string
	operation     ArgonLimit
	memory        ArgonLimit
}

// Argon2 create new argon2Hash object
func Argon2(plainPassword, secretString string) *argon2Hash {
	return &argon2Hash{
		name:          Argon2Name,
		plainPassword: plainPassword,
		secretString:  secretString,
		operation:     OpsLimitInteractive,
		memory:        MemLimitInteractive,
	}
}

func (c *argon2Hash) argon2iHash(plaintext []byte) ([]byte, bool) {
	length := len(plaintext)
	result := make([]C.char, argon2iStringBytesLen)
	isOk := int(C.crypto_pwhash_str(
		(*C.char)(unsafe.Pointer(&result[0])),
		(*C.char)(unsafe.Pointer(&plaintext[0])),
		(C.ulonglong)(length),
		(C.ulonglong)(c.operation),
		(C.size_t)(c.memory),
	)) == 0

	C.sodium_memzero(unsafe.Pointer(&plaintext[0]), C.size_t(length))

	return []byte(C.GoString(&result[0])), isOk
}

func (c *argon2Hash) argon2iPasswordVerify(plaintext, hash []byte) bool {
	if len(hash) < argon2iStringBytesLen {
		hash = append(hash, make([]byte, argon2iStringBytesLen-len(hash))...)
	}

	length := len(plaintext)
	isOk := int(C.crypto_pwhash_str_verify(
		(*C.char)(unsafe.Pointer(&hash[0])),
		(*C.char)(unsafe.Pointer(&plaintext[0])),
		(C.ulonglong)(length),
	)) == 0

	C.sodium_memzero(unsafe.Pointer(&plaintext[0]), C.size_t(length))

	return isOk
}

func (c *argon2Hash) prehash() string {
	passSum := blake2b.Sum512(util.ASCII85Byte([]byte(c.plainPassword)))
	secretSum := blake2b.Sum512(util.ASCII85Byte([]byte(c.secretString)))
	prehash, _ := blake2b.New512(secretSum[:])
	prehash.Write(passSum[:])
	pwd := util.ASCII85String(prehash.Sum(nil))

	return pwd
}

// Name hash name
func (c *argon2Hash) Name() string {
	return c.name
}

// Value argon2id hash from internally-prehashed password string.
func (c *argon2Hash) Value() string {
	b, _ := c.argon2iHash([]byte(c.prehash()))

	return string(b)
}

// IsValid validate stored hashed password against internally-prehashed password string.
func (c *argon2Hash) IsValid(storedHash string) bool {
	return c.argon2iPasswordVerify([]byte(c.prehash()), []byte(storedHash))
}
