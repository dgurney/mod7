package main

import (
	"fmt"
	"mod7/oem"
	"mod7/tendigit"
	"time"
)

// Benchmark generates 6000000 keys and shows the elapsed time. It's meant to be much more understandable and user-accessible than "make bench"
func Benchmark() {
	och := make(chan string)
	dch := make(chan string)
	started := time.Now()
	count := 0
	for i := 0; i < 3000000; i++ {
		count++
		go oem.GenerateOEM(och)
		go tendigit.Generate10digit(dch)
		<-och
		<-dch
	}

	fmt.Printf("Took %s to generate %d keys.\n", time.Since(started).Round(time.Millisecond), count*2)
	return
}
