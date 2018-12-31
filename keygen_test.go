package main

import (
	"bufio"
	"mod7/oem"
	"mod7/tendigit"
	"mod7/validation"
	"os"
	"testing"
)

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
