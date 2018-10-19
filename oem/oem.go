// Package oem handles generating and validating OEM keys (XXXXX-OEM-XXXXXXX-XXXXX).
package oem

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var serial [6]int

// Generate the first segment of the key. The first three digits represent the julian date the COA was printed (001 to 366), and the last two are the year.
// The year cannot be below 95 or above 03 (not Y2K-compliant D:).
func generateFirst() string {
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
	r.Shuffle(len(years), func(i, j int) {
		years[i], years[j] = years[j], years[i]
	})
	year := years[0]
	return date + year
}

// The third segment (OEM is the second) must begin with a zero, but otherwise it follows the same rule as the second segment of 10-digit keys:
// The digit sum must be seven, and the check digit cannot be 0 or >=8.
func generateThird() [6]int {
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
	return serial
}

// The fourth segment is truly irrelevant
func generateFourth() int {
	var fourth int
	for fourth < 10000 {
		fourth = r.Intn(99999)
	}
	return fourth
}

func validateKey(serial [6]int) bool {
	sum := 0
	for _, dig := range serial {
		sum += dig
	}
	if sum%7 == 0 {
		return true
	}
	return false
}

// GenerateOEM generates an OEM key (duh).
func GenerateOEM() {
	for !validateKey(generateThird()) {
		generateThird()
	}
	fmt.Printf("%s-OEM-0", generateFirst())
	for _, digits := range serial {
		fmt.Print(digits)
	}
	fmt.Printf("-%d\n", generateFourth())
}
