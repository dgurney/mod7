// Package generator handles generation of all key types
package generator

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func checkdigitCheck(k int) bool {
	// Check digit cannot be 0 or >= 8.
	if k%10 == 0 || k%10 >= 8 {
		return false
	}
	return true
}

func digitsum(num int) int {
	s := 0
	for num != 0 {
		digit := num % 10
		s += digit
		num /= 10
	}
	return s
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
