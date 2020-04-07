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

// Generate generates an 10-digit CD key.
func (CD) Generate(ch chan string) {
	// Generate the so-called site number, which is the first segment of the key.
	first := ""
	s := rand.Intn(998)
	// Technically 999 could be omitted as we don't generate a number that high, but we include it for posterity anyway.
	invalidSites := []int{333, 444, 555, 666, 777, 888, 999}
	for _, v := range invalidSites {
		if v == s {
			// Site number is invalid, so we replace it with a guaranteed valid number
			s = rand.Intn(300)
		}
	}
	first = fmt.Sprintf("%03d", s)

	// Generate the second segment of the key. The digit sum of the seven numbers must be divisible by seven.
	// The last digit is the check digit. The check digit cannot be 0 or >=8.
	serial := make([]int, 7)
	second := ""
	for {
		for i := 0; i < 7; i++ {
			serial[i] = rand.Intn(9)
			if i == 6 {
				// We must also generate a valid check digit
				for serial[i] == 0 || serial[i] >= 8 {
					serial[i] = rand.Intn(7)
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
				second += strconv.Itoa(digits)
			}
			break
		}
	}
	ch <- first + "-" + second
}
