// Package validation handles the validation of user-provided OEM and CD keys.
package validation

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func checkdigitCheck(c string) bool {
	// Check digit cannot be 0 or >= 8.
	if c[len(c)-1:] == "0" || c[len(c)-1:] >= "8" {
		return false
	}
	return true
}

func digitsum(num int64) int64 {
	var s int64
	for num != 0 {
		digit := num % 10
		s += digit
		num /= 10
	}
	return s
}

func validateCDKey(key string) error {
	site, err := strconv.ParseInt(key[0:3], 10, 0)
	if err != nil {
		return fmt.Errorf("the site number isn't a number")
	}
	main, err := strconv.ParseInt(key[4:11], 10, 0)
	if err != nil {
		return fmt.Errorf("the second segment isn't a number")
	}
	invalidSites := map[int64]int{333: 333, 444: 444, 555: 555, 666: 666, 777: 777, 888: 888, 999: 999}
	_, invalid := invalidSites[site]
	if invalid {
		fmt.Println("The site number is invalid: cannot be 333, 444, 555, 666, 777, 888, or 999.")
	}

	c := strconv.Itoa(int(main))
	secondSegmentInvalid := "The second segment is invalid:"
	if !checkdigitCheck(c) {
		fmt.Println(secondSegmentInvalid, "the last digit cannot be 0 or >= 8.")
	}
	sum := digitsum(main)
	if sum%7 != 0 {
		fmt.Printf("%s the digit sum (%d) must be divisible by 7.\n", secondSegmentInvalid, sum)
	}
	return nil
}

func validateOEM(key string) error {
	_, err := strconv.ParseInt(key[0:5], 10, 0)
	if err != nil {
		return fmt.Errorf("the first segment is not a number")
	}
	th, err := strconv.ParseInt(key[10:17], 10, 0)
	if err != nil {
		return fmt.Errorf("the third segment is not a number")
	}
	_, err = strconv.ParseInt(key[18:], 10, 0)
	if err != nil {
		return fmt.Errorf("the fourth segment is not a number")
	}
	julian, err := strconv.ParseInt(key[0:3], 10, 0)
	if julian == 0 || julian > 366 {
		fmt.Println("The date is invalid: valid date range 001-366.")
	}
	year := key[3:5]
	validYears := map[string]string{"95": "95", "96": "96", "97": "97", "98": "98", "99": "99", "00": "00", "01": "01", "02": "02", "03": "03"}
	_, valid := validYears[year]
	if !valid {
		fmt.Println("The year is invalid: cannot be less than 95 or above 03")
	}

	third := key[10:17]
	thirdSegmentInvalid := "The third segment is invalid:"
	if string(third[0]) != "0" {
		fmt.Println(thirdSegmentInvalid, "must begin with a 0.")
	}
	c := strconv.Itoa(int(th))
	if !checkdigitCheck(c) {
		fmt.Println(thirdSegmentInvalid, "last digit cannot be 0 or >= 8.")
	}
	sum := digitsum(th)
	if sum%7 != 0 {
		fmt.Printf("%s digit sum (%d) must be divisible by 7.\n", thirdSegmentInvalid, sum)
	}
	return nil
}

// ValidateKey validates the provided OEM or CD key.
func ValidateKey(k, kf string) {
	batchCheck := false
	if k == "batchCheck" {
		batchCheck = true
	}
	maybeValidMessage := "%s is valid if you get no further output.\n"
	unableToValidate := "Unable to validate key:"
	switch {
	case !batchCheck:
		// Make sure the provided key has a chance of being valid.
		switch {
		case len(k) == 11 && k[3:4] == "-":
			fmt.Printf(maybeValidMessage, k)
			if err := validateCDKey(k); err != nil {
				fmt.Println(unableToValidate, err)
			}
		case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
			fmt.Printf(maybeValidMessage, k)
			if err := validateOEM(k); err != nil {
				fmt.Println(unableToValidate, err)
			}
		default:
			fmt.Printf("%s doesn't even resemble a valid key.\n", k)
		}
	case batchCheck:
		keyfile, err := os.Open(kf)
		if err != nil {
			fmt.Println("Unable to open key file:", err)
			return
		}
		kfStat, err := os.Stat(kf)
		if err != nil {
			fmt.Println("Unable to stat key file:", err)
			return
		}
		if filepath.Ext(kfStat.Name()) != ".txt" {
			fmt.Println("This file doesn't have a .txt extension. Nothing interesting will happen if you pass anything other than a plain text file with the required characteristics.")
			return
		}
		defer keyfile.Close()
		var keys []string

		scanner := bufio.NewScanner(keyfile)
		for scanner.Scan() {
			keys = append(keys, scanner.Text())
		}
		// Make sure the provided keys have a chance of being valid.
		for i := 0; i < len(keys); i++ {
			switch {
			case len(keys[i]) == 11 && keys[i][3:4] == "-":
				fmt.Printf(maybeValidMessage, keys[i])
				if err := validateCDKey(keys[i]); err != nil {
					fmt.Println(unableToValidate, err)
				}
			case len(keys[i]) == 23 && keys[i][5:6] == "-" && keys[i][9:10] == "-" && keys[i][17:18] == "-" && len(keys[i][18:]) == 5:
				fmt.Printf(maybeValidMessage, keys[i])
				if err := validateOEM(keys[i]); err != nil {
					fmt.Println(unableToValidate, err)
				}
			default:
				fmt.Printf("%s doesn't even resemble a valid key.\n", keys[i])
			}
		}
	}

}
