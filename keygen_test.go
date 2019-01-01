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
			t.Errorf("Generated key %s is invalid!", ka[i])
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
			t.Errorf("Generated key %s is invalid!", ka[i])
		}

	}
}

func TestBatchValidation(t *testing.T) {
	vch := make(chan bool)
	validKeys := []struct {
		key string
	}{
		{"111-1111111"},
		{"000-0000007"},
		{"10000-OEM-0000007-11111"},
		{"32299-OEM-0840621-16752"},
		{"118-5688143"},
	}
	invalidKeys := []struct {
		key string
	}{
		// Not even close to a valid key
		{"1"},
		{"10000-OEM-0000007-1"},
		// Invalid date
		{"00099-OEM-0840621-16752"},
		{"36799-OEM-0840621-16752"},
		// Invalid year
		{"10094-OEM-0840621-16752"},
		{"10019-OEM-0840621-16752"},
		// Invalid site
		{"333-5688143"},
		{"444-5688143"},
		{"555-5688143"},
		{"666-5688143"},
		{"777-5688143"},
		{"888-5688143"},
		{"999-5688143"},
		// Invalid check digit
		{"10000-OEM-0140628-12345"},
		{"332-5683148"},
		// Invalid third segment starting digit
		{"10000-OEM-8040621-12345"},
		// Invalid digit sum
		{"10000-OEM-0000006-12345"},
		{"001-1234566"},
		// Not a number
		{"11a-1111111"},
		{"111-a111111"},
		{"1000a-OEM-0000007-11111"},
		{"10000-OEM-00000a7-11111"},
		{"10000-OEM-0000007-1111a"},
	}
	for _, kt := range validKeys {
		go validation.BatchValidate(kt.key, vch)
		switch {
		default:
			t.Logf("%s is valid, as expected.", kt.key)
		case !<-vch:
			t.Errorf("Valid key %s did not pass validation!", kt.key)
		}
	}
	for _, kt := range invalidKeys {
		go validation.BatchValidate(kt.key, vch)
		switch {
		default:
			t.Logf("%s is not valid, as expected.", kt.key)
		case <-vch:
			t.Errorf("Invalid key %s passed validation!", kt.key)
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
