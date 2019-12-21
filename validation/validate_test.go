package validation

import (
	"bufio"
	"os"
	"testing"
)

var validKeys = []string{
	"111-1111111",
	"000-0000007",
	"10000-OEM-0000007-00000",
	"32299-OEM-0840621-16752",
	"118-5688143",
}
var invalidKeys = []string{
	// Not even close to a valid key
	"1",
	"10000-OEM-0000007-1",
	// Invalid date
	"00099-OEM-0840621-16752",
	"36799-OEM-0840621-16752",
	// Invalid year
	"10094-OEM-0840621-16752",
	"10019-OEM-0840621-16752",
	// Invalid site
	"333-5688143",
	"444-5688143",
	"555-5688143",
	"666-5688143",
	"777-5688143",
	"888-5688143",
	"999-5688143",
	// Invalid check digit
	"10000-OEM-0140628-12345",
	"332-5683148",
	// Invalid third segment starting digit
	"10000-OEM-8040621-12345",
	// Invalid digit sum
	"10000-OEM-0000006-12345",
	"001-1234566",
	// Not a number
	"11a-1111111",
	"111-a111111",
	"1000a-OEM-0000007-11111",
	"10000-OEM-00000a7-11111",
	"10000-OEM-0000007-1111a",
	// Invalid second segment
	"10000-SEX-0000007-00000",
}

func TestSingleValidation(t *testing.T) {
	for _, key := range validKeys {
		v := ValidateKey(key)
		switch {
		default:
			t.Logf("%s is valid, as expected.", key)
		case !v:
			t.Errorf("Valid key %s did not pass validation!", key)
		}
	}
	for _, key := range invalidKeys {
		v := ValidateKey(key)
		switch {
		default:
			t.Logf("%s is not valid, as expected.", key)
		case v:
			t.Errorf("Invalid key %s passed validation!", key)
		}
	}
}

func TestBatchValidation(t *testing.T) {
	vch := make(chan bool)
	for _, key := range validKeys {
		go BatchValidate(key, vch)
		switch {
		default:
			t.Logf("%s is valid, as expected.", key)
		case !<-vch:
			t.Errorf("Valid key %s did not pass validation!", key)
		}
	}
	for _, key := range invalidKeys {
		go BatchValidate(key, vch)
		switch {
		default:
			t.Logf("%s is not valid, as expected.", key)
		case <-vch:
			t.Errorf("Invalid key %s passed validation!", key)
		}
	}
}

func BenchmarkBatchValidate100(b *testing.B) {
	b.StopTimer()
	keyfile, err := os.Open("../benchmark_files/benchmark_100.txt")
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
				go BatchValidate(keys[i], vch)
			}
		}
	}
}
