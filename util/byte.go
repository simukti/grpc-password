// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package util

import (
	"bytes"
	"encoding/ascii85"
)

// ASCII85Byte encode bytes to ascii-85 bytes
func ASCII85Byte(b []byte) []byte {
	bb := ascii85Encode(b)

	return bb.Bytes()
}

// https://github.com/tuupola/base85
func ascii85Encode(b []byte) *bytes.Buffer {
	bb := &bytes.Buffer{}
	encoder := ascii85.NewEncoder(bb)
	encoder.Write(b)
	encoder.Close()

	return bb
}
