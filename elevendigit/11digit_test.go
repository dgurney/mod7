package elevendigit

import (
	"mod7/validation"
	"testing"
)

func TestECD(t *testing.T) {
	ka := make([]string, 0)
	dch := make(chan string)
	vch := make(chan bool)
	for i := 0; i < 500000; i++ {
		go Generate11digit(dch)
		ka = append(ka, <-dch)
	}
	for i := 0; i < len(ka); i++ {
		go validation.BatchValidate(ka[i], vch)
		if !<-vch {
			t.Errorf("Generated key %s is invalid!", ka[i])
		}

	}
}

func Benchmark11digit100(b *testing.B) {
	tch := make(chan string)
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			go Generate11digit(tch)
			<-tch
		}
	}
}
