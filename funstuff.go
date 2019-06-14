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
	// 3294 valid dates + 100000 valid last segments
	ch <- valid * 103294
}

func generateAllCD(ch chan bool) {
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
	ch <- true
}

func generateAllOEM(ch chan bool) {
	dates := []string{}
	years := []string{"95", "96", "97", "98", "99", "00", "01", "02", "03"}
	mainSegments := []string{}
	lastSegments := []string{}
	for _, y := range years {
		for day := 1; day <= 366; day++ {
			dates = append(dates, fmt.Sprintf("%03d%s", day, y))
		}
	}
	for ls := 0; ls <= 99999; ls++ {
		lastSegments = append(lastSegments, fmt.Sprintf("%05d", ls))
	}
	valid := make(chan bool)
	for main := 0; main < 999999; main++ {
		test := fmt.Sprintf("00195-OEM-%07d-00000", main)
		go validation.BatchValidate(test, valid)
		if <-valid {
			mainSegments = append(mainSegments, fmt.Sprintf("%07d", main))
		}
	}

	// And finally print all the actual keys
	for _, d := range dates {
		for _, m := range mainSegments {
			for _, l := range lastSegments {
				fmt.Printf("%s-OEM-%s-%s\n", d, m, l)
			}
		}
	}
	ch <- true
}
