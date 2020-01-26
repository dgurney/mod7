package main

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
	"time"

	"github.com/dgurney/mod7/elevendigit"
	"github.com/dgurney/mod7/oem"
	"github.com/dgurney/mod7/tendigit"
	"github.com/dgurney/mod7/validation"
)

// generationBenchmark generates 3000000 keys and shows the elapsed time. It's meant to be much more understandable and user-accessible than "make bench"
func generationBenchmark() []string {
	och := make(chan string)
	dch := make(chan string)
	keys := make([]string, 0)
	started := time.Now()
	count := 0
	for i := 0; i < 1000000; i++ {
		count++
		go oem.GenerateOEM(och)
		keys = append(keys, <-och)
		go tendigit.Generate10digit(dch)
		keys = append(keys, <-dch)
		go elevendigit.Generate11digit(dch)
		keys = append(keys, <-dch)
	}

	fmt.Printf("Took %s to generate %d keys.\n", time.Since(started).Round(time.Millisecond), count*3)
	return keys
}

// validationBenchmark validates 3000000 keys and shows the elapsed time. It's meant to be much more understandable and user-accessible than "make bench"
func validationBenchmark(keys []string) {
	vch := make(chan bool)
	started := time.Now()
	for _, v := range keys {
		go validation.BatchValidate(v, vch)
		<-vch
	}
	fmt.Printf("Took %s to validate %d keys.\n", time.Since(started).Round(time.Millisecond), len(keys))
	return
}
