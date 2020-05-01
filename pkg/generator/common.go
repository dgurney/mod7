// Package generator handles generation of all key types
package generator

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// KeyGenerator is what it says on the tin
type KeyGenerator interface {
	Generate(chan string)
}

// OEM key
type OEM struct {
}

// ElevenCD is an 11-digit CD key
type ElevenCD struct {
}

// CD Key
type CD struct {
}
