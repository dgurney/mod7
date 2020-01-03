package main

import (
	"fmt"
	"mod7/elevendigit"
	"mod7/oem"
	"mod7/tendigit"
	"mod7/validation"
	"time"
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
