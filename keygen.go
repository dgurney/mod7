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
		if len(keys) == 0 {
			fmt.Println("The key file is empty.")
			return
		}
		for i := 0; i < len(keys); i++ {
			go validation.BatchValidate(keys[i], vch)
			switch {
			default:
				fmt.Printf("%s is invalid\n", keys[i])
			case <-vch:
				fmt.Printf("%s is valid\n", keys[i])
			}
		}
		return
	}
	if len(*v) > 0 {
		validation.ValidateKey(*v)
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
