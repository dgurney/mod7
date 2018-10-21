package main

import (
	"flag"
	"fmt"
	"mod7/oem"
	"mod7/tendigit"
)

func main() {
	b := flag.Bool("b", false, "Generate both keys")
	o := flag.Bool("o", false, "Generate an OEM key")
	d := flag.Bool("d", false, "Generate a 10-digit key (aka CD Key)")
	r := flag.Int("r", 1, "Repeat n times. I can't figure out why you'd ever need this, but whatever...")
	flag.Parse()
	if *r < 1 {
		*r = 1
	}
	for i := 0; i < *r; i++ {
		switch {
		case *d:
			tendigit.Generate10digit()
		case *o:
			oem.GenerateOEM()
		case *b:
			oem.GenerateOEM()
			tendigit.Generate10digit()
		default:
			fmt.Println("You must specify what you want to generate! Usage:")
			flag.PrintDefaults()
		}
	}
}
