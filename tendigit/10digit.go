// Package tendigit handles generation of CD keys (XXX-XXXXXXX).
package tendigit

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate the so-called site number, which is the first segment of the key.
func genSite(ch chan string, m *sync.Mutex) {
	m.Lock()
	site := ""
	s := r.Intn(998)
	// Technically 999 could be omitted as we don't generate a number that high, but we include it for posterity anyway.
	invalidSites := []int{333, 444, 555, 666, 777, 888, 999}
	for _, v := range invalidSites {
		if v == s {
			// Site number is invalid, so we replace it with a guaranteed valid number
			s = r.Intn(300)
		}
	}

	site = fmt.Sprintf("%03d", s)
	m.Unlock()
	ch <- site
}

// Generate the second segment of the key. The digit sum of the seven numbers must be divisible by seven.
// The last digit is the check digit. The check digit cannot be 0 or >=8.
func genSeven(ch chan string, m *sync.Mutex) {
	serial := make([]int, 7)
	m.Lock()
	final := ""
	for {
		for i := 0; i < 7; i++ {
			serial[i] = r.Intn(9)
			if i == 6 {
				// We must also generate a valid check digit
				for serial[i] == 0 || serial[i] >= 8 {
					serial[i] = r.Intn(7)
				}
			}
		}
		// Perform the actual validation
		sum := 0
		for _, dig := range serial {
			sum += dig
		}
		if sum%7 == 0 {
			for _, digits := range serial {
				final += strconv.Itoa(digits)
			}
			break
		}
	}
	m.Unlock()
	ch <- final
}

// Generate10digit generates a 10-digit CD key.
func Generate10digit(ch chan string) {
	var m sync.Mutex
	sch := make(chan string)
	dch := make(chan string)
	go genSite(sch, &m)
	go genSeven(dch, &m)
	ch <- <-sch + "-" + <-dch
}
