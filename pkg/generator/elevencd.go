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

// Generate generates an 11-digit CD key.
func (ElevenCD) Generate(ch chan string) {
	// Generate the first segment of the key.
	// Formula for last digit: third digit + 1 or 2. If the result is more than 9, it's 0 or 1.
	s := rand.Intn(999)
	site := fmt.Sprintf("%03d", s)
	last, _ := strconv.Atoi(site[len(site)-1:])
	first := ""
	fourth := 0
	switch {
	default:
		fourth = last + 1
	case rand.Intn(2) == 1:
		fourth = last + 2
	}
	switch {
	case fourth == 10:
		fourth = 0
	case fourth > 10:
		fourth = 1
	}
	first = fmt.Sprintf("%s%d", site, fourth)

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
