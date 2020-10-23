package main

/*
   Copyright (C) 2020 Daniel Gurney
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	g "github.com/dgurney/mod7/v4/pkg/generator"
	"github.com/dgurney/mod7/v4/pkg/validator"
)

// Used if mod7 is not built using the makefile.
const version = "4.0.0"

// git describe --tags --dirty
var gitVersion string

func getVersion() string {
	if len(gitVersion) == 0 {
		return version
	}
	return gitVersion
}

func main() {
	all := flag.Bool("a", false, "Generate all kinds of keys.")
	bench := flag.Int("bench", 0, "Benchmark generation and validation of N keys.")
	batchvalidate := flag.String("bv", "", "Batch validate a key file. The key file should be a plain text file (with a .txt extension) with 1 key per line.")
	cd := flag.Bool("d", false, "Generate a 10-digit key (aka CD Key).")
	elevencd := flag.Bool("e", false, "Generate an 11-digit CD key.")
	oem := flag.Bool("o", false, "Generate an OEM key.")
	repeat := flag.Int("r", 1, "Generate n keys.")
	t := flag.Bool("t", false, "Show how long the generation or batch validation took.")
	validate := flag.String("v", "", "Validate a CD or OEM key.")
	ver := flag.Bool("ver", false, "Show version information and exit")
	flag.Parse()
	if *repeat < 1 {
		*repeat = 1
	}

	var started time.Time
	if *t {
		started = time.Now()
	}

	if *ver {
		fmt.Printf("mod7 v%s by Daniel Gurney\n", getVersion())
		return
	}

	if *bench != 0 {
		fmt.Println("Running key generation benchmark...")
		k := generationBenchmark(*bench)
		fmt.Println("Running key validator benchmark...")
		validationBenchmark(k)
		return
	}

	if len(*batchvalidate) > 0 {
		if filepath.Ext(*batchvalidate) != ".txt" {
			fmt.Println("The key file must be a plain text file with a .txt extension. Tricking this check will not do anything interesting, so don't bother.")
			return
		}
		keyfile, err := os.Open(*batchvalidate)
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
				go validator.Validate(keys[i], vch)
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

	if len(*validate) > 0 {
		vch := make(chan bool)
		go validator.Validate(*validate, vch)
		switch {
		default:
			fmt.Printf("%s is not valid.\n", *validate)
		case <-vch:
			fmt.Printf("%s is valid.\n", *validate)
		}
		return
	}

	CDKeych := make(chan string, runtime.NumCPU())
	eCDKeych := make(chan string, runtime.NumCPU())
	OEMKeych := make(chan string, runtime.NumCPU())
	if !*all && !*cd && !*elevencd && !*oem {
		fmt.Println("You must specify what you want to do! Usage:")
		flag.PrintDefaults()
		return
	}
	if *elevencd && *oem && *cd {
		*elevencd, *oem, *cd = false, false, false
		*all = true
	}
	// a and key type are mutually exclusive
	if *elevencd && *all || *oem && *all || *cd && *all {
		*all = false
	}
	oemkey := g.OEM{}
	ecdkey := g.ElevenCD{}
	cdkey := g.CD{}
	for i := 0; i < *repeat; i++ {
		if *all {
			go oemkey.Generate(OEMKeych)
			go cdkey.Generate(CDKeych)
			go ecdkey.Generate(eCDKeych)
			fmt.Println(<-OEMKeych)
			fmt.Println(<-CDKeych)
			fmt.Println(<-eCDKeych)
		}
		if *elevencd {
			go ecdkey.Generate(eCDKeych)
			fmt.Println(<-eCDKeych)
		}
		if *cd {
			go cdkey.Generate(CDKeych)
			fmt.Println(<-CDKeych)
		}
		if *oem {
			go oemkey.Generate(OEMKeych)
			fmt.Println(<-OEMKeych)
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
		mult := 0
		switch {
		default:
			switch {
			case *repeat > 1:
				fmt.Printf("Took %s to generate %d keys.\n", ended, *repeat)
				return
			case *repeat == 1:
				fmt.Printf("Took %s to generate %d key.\n", ended, *repeat)
				return
			}
		case *elevencd && *oem || *elevencd && *cd || *oem && *cd:
			mult = 2
		case *all:
			mult = 3
		}
		fmt.Printf("Took %s to generate %d keys.\n", ended, *repeat*mult)
	}
}
