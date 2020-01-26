package oem

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
	"testing"

	"github.com/dgurney/mod7/validation"
)

func TestOEM(t *testing.T) {
	ka := make([]string, 0)
	och := make(chan string)
	vch := make(chan bool)
	for i := 0; i < 500000; i++ {
		go GenerateOEM(och)
		ka = append(ka, <-och)
	}
	for i := 0; i < len(ka); i++ {
		go validation.BatchValidate(ka[i], vch)
		if !<-vch {
			t.Errorf("Generated key %s is invalid!", ka[i])
		}

	}
}

func BenchmarkOEM100(b *testing.B) {
	b.StopTimer()
	och := make(chan string)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			go GenerateOEM(och)
			<-och
		}
	}
}
