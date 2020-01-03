package main

import (
	"fmt"
	"mod7/validation"
)

// Calculate the total amount of valid CD keys.
func total10digit(ch chan int) {
	valid := 0
	for main := 0; main < 9999999; main++ {
		test := fmt.Sprintf("000-%07d", main)
		vch := make(chan bool)
		go validation.BatchValidate(test, vch)
		if <-vch {
			valid++
		}
	}
	// There are 993 valid site numbers
	ch <- valid * 993
}
func total11digit(ch chan int) {
	valid := 0
	for main := 0; main < 9999999; main++ {
		test := fmt.Sprintf("0001-%07d", main)
		vch := make(chan bool)
		go validation.BatchValidate(test, vch)
		if <-vch {
			valid++
		}
	}
	// There are 1998 first segments instead of 1000 due to the whole +1 or +2 thing.
	ch <- valid * 1998
}

// Calculate the total amount of valid OEM keys.
func totaloem(ch chan int) {
	valid := 0
	for main := 0; main < 999999; main++ {
		test := fmt.Sprintf("00196-OEM-0%06d-11111", main)
		vch := make(chan bool)
		go validation.BatchValidate(test, vch)
		if <-vch {
			valid++
		}
	}
	// 3294 valid dates + 100000 valid last segments
	ch <- valid * 103294
}
