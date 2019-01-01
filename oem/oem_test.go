package oem

import (
	"mod7/validation"
	"testing"
)

func TestOEM(t *testing.T) {
	ka := make([]string, 0)
	och := make(chan string)
	vch := make(chan bool)
	for i := 0; i < 500000; i++ {
		go GenerateOEM(och)
		ka = append(ka, <-och)
	}
	for i := 0; i < len(ka); i++ {
		go validation.BatchValidate(ka[i], vch)
		if !<-vch {
			t.Errorf("Generated key %s is invalid!", ka[i])
		}

	}
}

func BenchmarkOEM100(b *testing.B) {
	och := make(chan string)
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			go GenerateOEM(och)
			<-och
		}
	}
}
