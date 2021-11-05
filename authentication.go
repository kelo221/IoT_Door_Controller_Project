package main

import (
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func HashPassword(password, salt []byte) []byte {
	defer clear(password)
	return pbkdf2.Key(password, salt, 4096, sha256.Size, sha256.New)
}
