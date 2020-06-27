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
	second := ""
	for {
		s := rand.Intn(9999999)
		// Perform the actual validation
		sum := digitsum(s)
		if sum%7 == 0 {
			second = fmt.Sprintf("%07d", s)
			if checkdigitCheck(s) {
				break
			}
		}
	}
	ch <- first + "-" + second
}
