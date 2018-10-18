// Package tendigit handles the generation and validation of an "XXX-XXXXXXX"-style product key.
package tendigit

import (
	"fmt"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var site int
var serial [7]int

func checkSite(s int) bool {
	// Technically we could omit 999 as we don't generate a number that high, but we include it for posterity anyway.
	invalidSites := []int{333, 444, 555, 666, 777, 888, 999}
	for _, v := range invalidSites {
		if v == s || s == 0 {
			// Site number is invalid
			return false
		}
	}
	return true
}

// Generate the so-called "site" number, which is the first segment of the key.
func genSite() int {
	for site < 100 {
		site = r.Intn(998)
	}
	return site
}

// Generate the second segment of the key. The digit sum of the seven numbers must be divisible by seven.
// The last digit is the "check digit". The check digit cannot be 0 or >=8.
func genSeven() [7]int {
	for i := 0; i < 7; i++ {
		serial[i] = r.Intn(9)
		if i == 6 {
			for serial[i] == 0 || serial[i] >= 8 {
				serial[i] = r.Intn(7)
			}
		}
	}
	return serial
}

// Perform the actual validation
func validateSeven(serial [7]int) bool {
	sum := 0
	for _, dig := range serial {
		sum += dig
	}
	if sum%7 == 0 {
		return true
	}
	return false
}

// Generate10digit generates a 10-digit product key.
func Generate10digit() {
	for !checkSite(site) {
		genSite()
	}
	for !validateSeven(genSeven()) {
		genSeven()
	}
	fmt.Printf("%d-", site)
	for _, digits := range serial {
		fmt.Print(digits)
	}
}
