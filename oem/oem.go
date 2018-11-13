// Package oem handles generating and validating OEM keys (XXXXX-OEM-XXXXXXX-XXXXX).
package oem

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var serial [6]int
var eggChance int

// Generate the first segment of the key. The first three digits represent the julian date the COA was printed (001 to 366), and the last two are the year.
// The year cannot be below 95 or above 03 (not Y2K-compliant D:).
func generateFirst(ch chan string, wg *sync.WaitGroup, m *sync.Mutex) {
	wg.Add(1)
	defer wg.Done()
	m.Lock()
	d := r.Intn(366)
	var date string
	switch {
	case d == 0:
		date = "123"
	case d < 100 && d > 9:
		date = "0" + strconv.Itoa(d)
	case d < 10:
		date = "00" + strconv.Itoa(d)
	default:
		date = strconv.Itoa(d)
	}
	years := []string{"95", "96", "97", "98", "99", "00", "01", "02", "03"}
	year := years[r.Intn(len(years))]
	m.Unlock()
	ch <- date + year
}

// The third segment (OEM is the second) must begin with a zero, but otherwise it follows the same rule as the second segment of 10-digit keys:
// The digit sum must be divisible by seven, and the check digit cannot be 0 or >=8.
func generateThird(ch chan string, wg *sync.WaitGroup, m *sync.Mutex) {
	wg.Add(1)
	defer wg.Done()
	m.Lock()
	var final string
	var valid bool
	for !valid {
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

// The fourth segment is truly irrelevant
func generateFourth(ch chan string, wg *sync.WaitGroup, m *sync.Mutex) {
	wg.Add(1)
	defer wg.Done()
	m.Lock()
	f := r.Intn(99999)
	var fourth string
	switch {
	case f < 10:
		fourth = "0000" + strconv.Itoa(f)
	case f < 100 && f > 9:
		fourth = "000" + strconv.Itoa(f)
	case f < 1000 && f > 99:
		fourth = "00" + strconv.Itoa(f)
	case f < 10000 && f > 999:
		fourth = "0" + strconv.Itoa(f)
	default:
		fourth = strconv.Itoa(f)
	}
	m.Unlock()
	ch <- fourth
}

// GenerateOEM generates an OEM key (duh).
func GenerateOEM(ch chan string) {
	var wg sync.WaitGroup
	var m sync.Mutex
	dch := make(chan string)
	tch := make(chan string)
	fch := make(chan string)
	go generateFirst(dch, &wg, &m)
	go generateThird(tch, &wg, &m)
	go generateFourth(fch, &wg, &m)
	ch <- <-dch + "-OEM-0" + <-tch + "-" + <-fch
}
