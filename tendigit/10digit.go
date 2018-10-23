// Package tendigit handles the generation and validation of an "XXX-XXXXXXX"-style product key.
package tendigit

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var serial [7]int

// Generate the so-called "site" number, which is the first segment of the key.
func genSite(ch chan string, wg *sync.WaitGroup, m *sync.Mutex) {
	wg.Add(1)
	defer wg.Done()
	m.Lock()
	var site string
	s := r.Intn(998)
	// Technically we could omit 999 as we don't generate a number that high, but we include it for posterity anyway.
	invalidSites := []int{333, 444, 555, 666, 777, 888, 999}
	for _, v := range invalidSites {
		if v == s {
			// Site number is invalid, so we must obliterate it
			s = 69
		}
	}

	switch {
	case s < 10:
		site = "00" + strconv.Itoa(s)
	case s < 100 && s > 9:
		site = "0" + strconv.Itoa(s)
	default:
		site = strconv.Itoa(s)
	}
	m.Unlock()
	ch <- site
}

// Generate the second segment of the key. The digit sum of the seven numbers must be divisible by seven.
// The last digit is the "check digit". The check digit cannot be 0 or >=8.
func genSeven(ch chan string, wg *sync.WaitGroup, m *sync.Mutex) {
	wg.Add(1)
	defer wg.Done()
	m.Lock()
	var valid bool
	var final string
	for !valid {
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
		switch {
		case sum%7 == 0:
			valid = true
			for _, digits := range serial {
				final += strconv.Itoa(digits)
			}
		default:
			valid = false
		}
	}
	m.Unlock()
	ch <- final
}

// Generate10digit generates a 10-digit product key.
func Generate10digit(ch chan string) {
	var wg sync.WaitGroup
	var m sync.Mutex
	sch := make(chan string)
	dch := make(chan string)
	go genSite(sch, &wg, &m)
	go genSeven(dch, &wg, &m)
	ch <- <-sch + "-" + <-dch
	wg.Wait()
}
