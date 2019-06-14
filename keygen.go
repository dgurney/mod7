package main

import (
	"bufio"
	"flag"
	"fmt"
	"mod7/oem"
	"mod7/tendigit"
	"mod7/validation"
	"os"
	"path/filepath"
	"time"
)

// Used if mod7 is not built using the makefile.
const version = "1.2.9"

// git describe --tags --dirty
var gitVersion string

func getVersion() string {
	switch {
	default:
		return gitVersion
	case len(gitVersion) == 0:
		return version
	}
}

func main() {
	b := flag.Bool("b", false, "Generate both keys.")
	bench := flag.Bool("bench", false, "Benchmark the generation time of 6000000 keys.")
	o := flag.Bool("o", false, "Generate an OEM key.")
	d := flag.Bool("d", false, "Generate a 10-digit key (aka CD Key).")
	r := flag.Int("r", 1, "Generate n keys.")
	t := flag.Bool("t", false, "Show how long the generation or batch validation took.")
	v := flag.String("v", "", "Validate a CD or OEM key")
	a := flag.Bool("a", false, "Generate and print every valid key of the specified type (-o, -d, or -b). Note that OEM will take forever.")
	total := flag.Bool("total", false, "Print the total amount of valid keys")
	bv := flag.String("bv", "", "Batch validate a key file. The key file should be a plain text file (with a .txt extension) with 1 key per line.")
	ver := flag.Bool("ver", false, "Show version information and exit")
	flag.Parse()
	if *r < 1 {
		*r = 1
	}
	var started time.Time
	if *t {
		started = time.Now()
	}

	if *ver {
		fmt.Printf("mod7 v%s by Daniel Gurney\n", getVersion())
		return
	}

	if *bench {
		Benchmark()
		return
	}

	if len(*bv) > 0 {
		if filepath.Ext(*bv) != ".txt" {
			fmt.Println("The key file must be a plain text file with a .txt extension. Tricking this check will not do anything interesting, so don't bother.")
			return
		}
		keyfile, err := os.Open(*bv)
		if err != nil {
			fmt.Println("Unable to open key file:", err)
			return
		}
		defer keyfile.Close()
		var keys []string
		vch := make(chan bool)
		scanner := bufio.NewScanner(keyfile)
		for scanner.Scan() {
			keys = append(keys, scanner.Text())
		}
		kl := len(keys)
		if kl == 0 {
			fmt.Println("The key file is empty.")
			return
		}
		for i := 0; i < kl; i++ {
			if keys[i] != "" {
				go validation.BatchValidate(keys[i], vch)
				switch {
				default:
					fmt.Printf("%s is invalid\n", keys[i])
				case <-vch:
					fmt.Printf("%s is valid\n", keys[i])
				}
			}
		}

		if *t {
			var ended time.Duration
			switch {
			case time.Since(started).Round(time.Second) > 1:
				ended = time.Since(started).Round(time.Millisecond)
			default:
				ended = time.Since(started).Round(time.Microsecond)
			}
			if ended < 1 {
				// Oh Windows...
				fmt.Println("Could not display elapsed time correctly :(")
				return
			}
			switch {
			case len(keys) > 1:
				fmt.Printf("Took %s to validate %d keys.\n", ended, kl)
			default:
				fmt.Printf("Took %s to validate %d key.\n", ended, kl)
			}
		}
		return
	}

	if len(*v) > 0 {
		validation.ValidateKey(*v)
		return
	}

	if *total {
		fmt.Println("Calculating the total amount of valid CD and OEM keys...")
		validcd := make(chan int)
		validoem := make(chan int)
		go total10digit(validcd)
		go totaloem(validoem)
		fmt.Println("Amount of valid CD keys:", <-validcd)
		fmt.Println("Amount of valid OEM keys:", <-validoem)
		return
	}

	if *a {
		switch {
		default:
			fmt.Println("You must supply either -d, -o, or -b.")
			return
		case *d:
			cdCh := make(chan bool)
			generateAllCD(cdCh)
		case *o:
			oemCh := make(chan bool)
			generateAllOEM(oemCh)
		case *b:
			cdCh := make(chan bool)
			oemCh := make(chan bool)
			go generateAllCD(cdCh)
			go generateAllOEM(oemCh)
			<-cdCh
			<-oemCh

		}
		if *t {
			ended := time.Since(started).Round(time.Millisecond)
			fmt.Printf("Took %s to generate all keys.\n", ended)
		}
		return
	}

	CDKeych := make(chan string)
	OEMKeych := make(chan string)
	for i := 0; i < *r; i++ {
		switch {
		case *d:
			go tendigit.Generate10digit(CDKeych)
			fmt.Println(<-CDKeych)
		case *o:
			go oem.GenerateOEM(OEMKeych)
			fmt.Println(<-OEMKeych)
		case *b:
			go oem.GenerateOEM(OEMKeych)
			go tendigit.Generate10digit(CDKeych)
			fmt.Println(<-CDKeych)
			fmt.Println(<-OEMKeych)
		default:
			fmt.Println("You must specify what you want to do! Usage:")
			flag.PrintDefaults()
			return
		}
	}
	if *t {
		var ended time.Duration
		switch {
		case time.Since(started).Round(time.Second) > 1:
			ended = time.Since(started).Round(time.Millisecond)
		default:
			ended = time.Since(started).Round(time.Microsecond)
		}
		if ended < 1 {
			// Oh Windows...
			fmt.Println("Could not display elapsed time correctly :(")
			return
		}
		switch {
		default:
			switch {
			case *r > 1:
				fmt.Printf("Took %s to generate %d keys.\n", ended, *r)
			case *r == 1:
				fmt.Printf("Took %s to generate %d key.\n", ended, *r)
			}
		case *b:
			fmt.Printf("Took %s to generate %d keys.\n", ended, *r*2)
		}
	}
}
