// Package oem handles generation of OEM keys (XXXXX-OEM-XXXXXXX-XXXXX).
package oem

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var serial [6]int

// Generate the first segment of the key. The first three digits represent the julian date the COA was printed (001 to 366), and the last two are the year.
// The year cannot be below 95 or above 03 (not Y2K-compliant D:).
func generateFirst(ch chan string, m *sync.Mutex) {
	m.Lock()
	var d int
	nonzero := false
	for !nonzero {
		switch {
		case d != 0:
			nonzero = true
		default:
			d = r.Intn(366)
		}
	}
	date := fmt.Sprintf("%03d", d)
	years := []string{"95", "96", "97", "98", "99", "00", "01", "02", "03"}
	year := years[r.Intn(len(years))]
	m.Unlock()
	ch <- date + year
}

// The third segment (OEM is the second) must begin with a zero, but otherwise it follows the same rule as the second segment of 10-digit keys:
// The digit sum must be divisible by seven, and the check digit cannot be 0 or >=8.
func generateThird(ch chan string, m *sync.Mutex) {
	m.Lock()
	final := ""
	for {
		// We generate only 6 digits because of the "first digit must be 0" rule
		for i := 0; i < 6; i++ {
			serial[i] = r.Intn(9)
			if i == 5 {
				// We must also generate a valid check digit
				for serial[i] == 0 || serial[i] >= 8 {
					serial[i] = r.Intn(7)
				}
			}
		}
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

// The fourth segment is truly irrelevant
func generateFourth(ch chan string, m *sync.Mutex) {
	m.Lock()
	f := r.Intn(99999)
	fourth := fmt.Sprintf("%05d", f)
	m.Unlock()
	ch <- fourth
}

// GenerateOEM generates an OEM key.
func GenerateOEM(ch chan string) {
	var m sync.Mutex
	dch := make(chan string)
	tch := make(chan string)
	fch := make(chan string)
	go generateFirst(dch, &m)
	go generateThird(tch, &m)
	go generateFourth(fch, &m)
	ch <- <-dch + "-OEM-0" + <-tch + "-" + <-fch
}
