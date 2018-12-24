package main

import (
	"flag"
	"fmt"
	"mod7/oem"
	"mod7/tendigit"
	"mod7/validation"
	"time"
)

const version = "1.1.0"

func main() {
	b := flag.Bool("b", false, "Generate both keys.")
	o := flag.Bool("o", false, "Generate an OEM key.")
	d := flag.Bool("d", false, "Generate a 10-digit key (aka CD Key).")
	r := flag.Int("r", 1, "Generate n keys.")
	t := flag.Bool("t", false, "Show how long the generation took.")
	v := flag.String("v", "", "Validate a CD or OEM key")
	bv := flag.String("bv", "", "Batch validate a key file. The key file should be a plain text file (with a .txt extension) with 1 key per line.")
	ver := flag.Bool("ver", false, "Show version number and exit")
	flag.Parse()
	if *r < 1 {
		*r = 1
	}
	var started time.Time
	if *t {
		started = time.Now()
	}
	CDKeych := make(chan string)
	OEMKeych := make(chan string)
	if *ver {
		fmt.Printf("mod7 v%s by Daniel Gurney\n", version)
		return
	}
	if len(*bv) > 0 {
		validation.ValidateKey("batchCheck", *bv)
		return
	}
	if len(*v) > 0 {
		validation.ValidateKey(*v, "")
		return
	}
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
