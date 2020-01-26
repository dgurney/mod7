// Package elevendigit handles generation of 11-digit CD keys (XXXX-XXXXXXX).
package elevendigit

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
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate the first segment of the key.
// Formula for last digit: third digit + 1 or 2. If the result is more than 9, it's 0 or 1.
func genSite(ch chan string, m *sync.Mutex) {
	m.Lock()
	s := r.Intn(999)
	site := fmt.Sprintf("%03d", s)
	last, _ := strconv.Atoi(site[len(site)-1:])
	fourth := 0
	switch {
	default:
		fourth = last + 1
	case r.Intn(2) == 1:
		fourth = last + 2
	}
	switch {
	case fourth == 10:
		fourth = 0
	case fourth > 10:
		fourth = 1
	}
	m.Unlock()
	ch <- fmt.Sprintf("%s%d", site, fourth)
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

// Generate11digit generates an 11-digit CD key.
func Generate11digit(ch chan string) {
	var m sync.Mutex
	sch := make(chan string)
	dch := make(chan string)
	go genSite(sch, &m)
	go genSeven(dch, &m)
	ch <- <-sch + "-" + <-dch
}
