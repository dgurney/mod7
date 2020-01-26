// Package validation handles the validation of user-provided OEM and CD keys.
package validation

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
	"fmt"
	"strconv"
)

// The boolean is only used for testing.
var valid = true

func validateCDKey(key string) error {
	site, err := strconv.ParseInt(key[0:3], 10, 0)
	if err != nil {
		valid = false
		return fmt.Errorf("the site number isn't a number")
	}
	main, err := strconv.ParseInt(key[4:11], 10, 0)
	if err != nil {
		valid = false
		return fmt.Errorf("the second segment isn't a number")
	}
	invalidSites := map[int64]int{333: 333, 444: 444, 555: 555, 666: 666, 777: 777, 888: 888, 999: 999}
	_, invalid := invalidSites[site]
	if invalid {
		valid = false
		fmt.Println("The site number is invalid: cannot be 333, 444, 555, 666, 777, 888, or 999.")
	}

	c := strconv.Itoa(int(main))
	secondSegmentInvalid := "The second segment is invalid:"
	if !checkdigitCheck(c) {
		valid = false
		fmt.Println(secondSegmentInvalid, "the last digit cannot be 0 or >= 8.")
	}
	sum := digitsum(main)
	if sum%7 != 0 {
		valid = false
		fmt.Printf("%s the digit sum (%d) must be divisible by 7.\n", secondSegmentInvalid, sum)
	}
	return nil
}

func validateECDKey(key string) error {
	_, err := strconv.ParseInt(key[0:4], 10, 0)
	if err != nil {
		valid = false
		return fmt.Errorf("the first segment isn't a number")
	}
	main, err := strconv.ParseInt(key[5:12], 10, 0)
	if err != nil {
		valid = false
		return fmt.Errorf("the second segment isn't a number")
	}

	// Error is safe to discard since we checked if it's a number before.
	last, _ := strconv.ParseInt(key[3:4], 10, 0)
	third, _ := strconv.ParseInt(key[2:3], 10, 0)

	if last != third+1 && last != third+2 {
		switch {
		case third == 8 && last != 9 && last != 0:
			valid = false
			fmt.Println("The first segment is invalid: The last digit must be 3rd digit + 1 or 2.")
		case third+1 >= 9 && last == 0 || third+2 >= 9 && last == 1:
			break
		default:
			valid = false
			fmt.Println("The first segment is invalid: The last digit must be 3rd digit + 1 or 2.")
		}
	}

	c := strconv.Itoa(int(main))
	secondSegmentInvalid := "The second segment is invalid:"
	if !checkdigitCheck(c) {
		valid = false
		fmt.Println(secondSegmentInvalid, "the last digit cannot be 0 or >= 8.")
	}
	sum := digitsum(main)
	if sum%7 != 0 {
		valid = false
		fmt.Printf("%s the digit sum (%d) must be divisible by 7.\n", secondSegmentInvalid, sum)
	}
	return nil
}

func validateOEM(key string) error {
	_, err := strconv.ParseInt(key[0:5], 10, 0)
	if err != nil {
		valid = false
		return fmt.Errorf("the first segment is not a number")
	}
	th, err := strconv.ParseInt(key[10:17], 10, 0)
	if err != nil {
		valid = false
		return fmt.Errorf("the third segment is not a number")
	}
	_, err = strconv.ParseInt(key[18:], 10, 0)
	if err != nil {
		valid = false
		return fmt.Errorf("the fourth segment is not a number")
	}
	julian, err := strconv.ParseInt(key[0:3], 10, 0)
	if julian == 0 || julian > 366 {
		valid = false
		fmt.Println("The date is invalid: valid date range 001-366.")
	}
	year := key[3:5]
	validYears := map[string]string{"95": "95", "96": "96", "97": "97", "98": "98", "99": "99", "00": "00", "01": "01", "02": "02", "03": "03"}
	_, validYear := validYears[year]
	if !validYear {
		valid = false
		fmt.Println("The year is invalid: cannot be less than 95 or above 03")
	}
	if key[6:9] != "OEM" {
		valid = false
		fmt.Println("The second segment is invalid: must be OEM.")
	}

	third := key[10:17]
	thirdSegmentInvalid := "The third segment is invalid:"
	if string(third[0]) != "0" {
		valid = false
		fmt.Println(thirdSegmentInvalid, "must begin with a 0.")
	}
	c := strconv.Itoa(int(th))
	if !checkdigitCheck(c) {
		valid = false
		fmt.Println(thirdSegmentInvalid, "last digit cannot be 0 or >= 8.")
	}
	sum := digitsum(th)
	if sum%7 != 0 {
		valid = false
		fmt.Printf("%s digit sum (%d) must be divisible by 7.\n", thirdSegmentInvalid, sum)
	}
	return nil
}

// ValidateKey validates the provided OEM or CD key. It should not be used when customized output is desired.
func ValidateKey(k string) bool {
	maybeValidMessage := "%s is valid if you get no further output.\n"
	unableToValidate := "Unable to validate key:"
	// Make sure the provided key has a chance of being valid.
	switch {
	case len(k) == 11 && k[3:4] == "-":
		fmt.Printf(maybeValidMessage, k)
		if err := validateCDKey(k); err != nil {
			valid = false
			fmt.Println(unableToValidate, err)
		}
	case len(k) == 12 && k[4:5] == "-":
		fmt.Printf(maybeValidMessage, k)
		if err := validateECDKey(k); err != nil {
			valid = false
			fmt.Println(unableToValidate, err)
		}
	case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
		fmt.Printf(maybeValidMessage, k)
		if err := validateOEM(k); err != nil {
			valid = false
			fmt.Println(unableToValidate, err)
		}
	default:
		valid = false
		fmt.Printf("%s doesn't even resemble a valid key.\n", k)
	}
	return valid
}
