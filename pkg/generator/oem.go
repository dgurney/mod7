package generator

/*
   Copyright (C) 2020 Daniel Gurney
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"fmt"
	"math/rand"
	"strconv"
)

// Generate generates an OEM key
func (o OEM) Generate(ch chan string) {
	// Generate the first segment of the key. The first three digits represent the julian date the COA was printed (001 to 366), and the last two are the year.
	// The year cannot be below 95 or above 03 (not Y2K-compliant D:).
	var d int
	first := ""
	nonzero := false
	for !nonzero {
		switch {
		case d != 0:
			nonzero = true
		default:
			d = rand.Intn(366)
		}
	}
	date := fmt.Sprintf("%03d", d)
	years := []string{"95", "96", "97", "98", "99", "00", "01", "02", "03"}
	year := years[rand.Intn(len(years))]
	first = date + year

	// The third segment (OEM is the second) must begin with a zero, but otherwise it follows the same rule as the second segment of 10-digit keys:
	// The digit sum must be divisible by seven, and the check digit cannot be 0 or >=8.
	serial := make([]int, 6)
	third := ""
	for {
		// We generate only 6 digits because of the "first digit must be 0" rule
		for i := 0; i < 6; i++ {
			serial[i] = rand.Intn(9)
			if i == 5 {
				// We must also generate a valid check digit
				for serial[i] == 0 || serial[i] >= 8 {
					serial[i] = rand.Intn(7)
				}
			}
		}
		sum := 0
		for _, dig := range serial {
			sum += dig
		}
		if sum%7 == 0 {
			for _, digits := range serial {
				third += strconv.Itoa(digits)
			}
			break
		}
	}

	// The fourth segment is truly irrelevant
	f := rand.Intn(99999)
	fourth := fmt.Sprintf("%05d", f)
	ch <- first + "-OEM-0" + third + "-" + fourth
}
