package main

import (
	"bufio"
	"mod7/oem"
	"mod7/tendigit"
	"mod7/validation"
	"os"
	"testing"
)

func TestOEM(t *testing.T) {
	ka := make([]string, 0)
	och := make(chan string)
	vch := make(chan bool)
	for i := 0; i < 500000; i++ {
		go oem.GenerateOEM(och)
		ka = append(ka, <-och)
	}
	for i := 0; i < len(ka); i++ {
		go validation.BatchValidate(ka[i], vch)
		if !<-vch {
			t.Errorf("Received invalid key %s!", ka[i])
		}

	}
}
func TestCD(t *testing.T) {
	ka := make([]string, 0)
	dch := make(chan string)
	vch := make(chan bool)
	for i := 0; i < 500000; i++ {
		go oem.GenerateOEM(dch)
		ka = append(ka, <-dch)
	}
	for i := 0; i < len(ka); i++ {
		go validation.BatchValidate(ka[i], vch)
		if !<-vch {
			t.Errorf("Received invalid key %s!", ka[i])
		}

	}
}

func BenchmarkOEM100(b *testing.B) {
	och := make(chan string)
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			go oem.GenerateOEM(och)
			<-och
		}
	}
}

func Benchmark10digit100(b *testing.B) {
	tch := make(chan string)
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			go tendigit.Generate10digit(tch)
			<-tch
		}
	}
}

func BenchmarkBatchValidate100(b *testing.B) {
	b.StopTimer()
	keyfile, err := os.Open("benchmark_files/benchmark_100.txt")
	if err != nil {
		b.Error(err)
	}
	defer keyfile.Close()
	var keys []string
	vch := make(chan bool)
	scanner := bufio.NewScanner(keyfile)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}
	kl := len(keys)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < kl; i++ {
			if keys[i] != "" {
				go validation.BatchValidate(keys[i], vch)
			}
		}
	}
}
