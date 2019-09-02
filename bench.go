package main

import (
	"fmt"
	"mod7/oem"
	"mod7/tendigit"
	"mod7/validation"
	"time"
)

// generationBenchmark generates 6000000 keys and shows the elapsed time. It's meant to be much more understandable and user-accessible than "make bench"
func generationBenchmark() {
	och := make(chan string)
	dch := make(chan string)
	started := time.Now()
	count := 0
	for i := 0; i < 3000000; i++ {
		count++
		go oem.GenerateOEM(och)
		<-och
		go tendigit.Generate10digit(dch)
		<-dch
	}

	fmt.Printf("Took %s to generate %d keys.\n", time.Since(started).Round(time.Millisecond), count*2)
	return
}

// validationBenchmark validates 1000000 keys and shows the elapsed time. It's meant to be much more understandable and user-accessible than "make bench"
func validationBenchmark() {
	och := make(chan string)
	dch := make(chan string)
	count := 0
	fmt.Println("Generating keys to validate (does not affect the result)...")
	keys := make([]string, 0)
	for i := 0; i < 500000; i++ {
		count++
		go oem.GenerateOEM(och)
		go tendigit.Generate10digit(dch)
		keys = append(keys, <-dch, <-och)
	}
	vch := make(chan bool)
	started := time.Now()
	for _, v := range keys {
		go validation.BatchValidate(v, vch)
		<-vch
	}
	fmt.Printf("Took %s to validate %d keys.\n", time.Since(started).Round(time.Millisecond), count*2)
	return
}
