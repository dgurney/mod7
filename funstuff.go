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
	// 3294 valid dates + 99999 valid last segments
	ch <- valid * 103293
}

func printFinalCD(site int, ms []string, done chan bool) {
	done <- true
}

func generateAllCD() {
	mainSegments := []string{}
	valid := make(chan bool)
	for main := 0; main < 9999999; main++ {
		test := fmt.Sprintf("000-%07d", main)
		go validation.BatchValidate(test, valid)
		if <-valid {
			mainSegments = append(mainSegments, fmt.Sprintf("%07d", main))
		}
	}
	invalidSites := map[int]int{333: 333, 444: 444, 555: 555, 666: 666, 777: 777, 888: 888, 999: 999}
	for site := 0; site < 999; site++ {
		_, invalid := invalidSites[site]
		if !invalid {
			for _, m := range mainSegments {
				fmt.Printf("%03d-%s\n", site, m)
			}
		}
	}
}
