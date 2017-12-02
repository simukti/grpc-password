// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package util

import (
	"crypto/rand"
)

// ASCII printable chars except: <space> <double quote> <single quote> <backtick> <backslash>
// This chars intended for random string use.
var safeChars = "!#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_abcdefghijklmnopqrstuvwxyz{|}~"

// ASCII85String encode bytes to ascii-85 string
func ASCII85String(b []byte) string {
	bb := ascii85Encode(b)

	return bb.String()
}

// RandomString generate fixed length random string from bytes seeded by crypto/rand
// reference: https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb#file-main-go-L46
func RandomString(n int) string {
	buffer := make([]byte, n)
	rand.Read(buffer)

	for k, v := range buffer {
		buffer[k] = safeChars[v%byte(len(safeChars))]
	}

	return string(buffer)
}
