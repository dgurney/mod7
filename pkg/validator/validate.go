package validator

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
	"strconv"
)

func (c cd) validate(v chan bool) {
	site, err := strconv.ParseInt(c.key[0:3], 10, 0)
	if err != nil {
		v <- false
		return
	}
	main, err := strconv.ParseInt(c.key[4:11], 10, 0)
	if err != nil {
		v <- false
		return
	}
	invalidSites := map[int64]int{333: 333, 444: 444, 555: 555, 666: 666, 777: 777, 888: 888, 999: 999}
	_, invalid := invalidSites[site]
	if invalid {
		v <- false
		return
	}
	check := strconv.Itoa(int(main))
	if !checkdigitCheck(check) {
		v <- false
		return
	}
	sum := digitsum(main)
	if sum%7 != 0 {
		v <- false
		return
	}
	v <- true
}

func (e elevencd) validate(v chan bool) {
	_, err := strconv.ParseInt(e.key[0:4], 10, 0)
	if err != nil {
		v <- false
		return
	}
	main, err := strconv.ParseInt(e.key[5:12], 10, 0)
	if err != nil {
		v <- false
		return
	}

	// Error is safe to discard since we checked if it's a number before.
	last, _ := strconv.ParseInt(e.key[3:4], 10, 0)
	third, _ := strconv.ParseInt(e.key[2:3], 10, 0)
	if last != third+1 && last != third+2 {
		switch {
		case third == 8 && last != 9 && last != 0:
			v <- false
			return
		case third+1 >= 9 && last == 0 || third+2 >= 9 && last == 1:
			break
		default:
			v <- false
			return
		}
	}

	c := strconv.Itoa(int(main))
	if !checkdigitCheck(c) {
		v <- false
		return
	}
	sum := digitsum(main)
	if sum%7 != 0 {
		v <- false
		return
	}
	v <- true
}

func (o oem) validate(v chan bool) {
	_, err := strconv.ParseInt(o.key[0:5], 10, 0)
	if err != nil {
		v <- false
		return
	}
	th, err := strconv.ParseInt(o.key[10:17], 10, 0)
	if err != nil {
		v <- false
		return
	}
	_, err = strconv.ParseInt(o.key[18:], 10, 0)
	if err != nil {
		v <- false
		return
	}
	julian, err := strconv.ParseInt(o.key[0:3], 10, 0)
	if julian == 0 || julian > 366 {
		v <- false
		return
	}
	year := o.key[3:5]
	validYears := map[string]string{"95": "95", "96": "96", "97": "97", "98": "98", "99": "99", "00": "00", "01": "01", "02": "02", "03": "03"}
	_, valid := validYears[year]
	if !valid {
		v <- false
		return
	}
	if o.key[6:9] != "OEM" {
		v <- false
		return
	}
	third := o.key[10:17]
	if string(third[0]) != "0" {
		v <- false
		return
	}
	c := strconv.Itoa(int(th))
	if !checkdigitCheck(c) {
		v <- false
		return
	}
	sum := digitsum(th)
	if sum%7 != 0 {
		v <- false
		return
	}
	v <- true

}

// Validate is used to validate a provided key.
func Validate(k string, v chan bool) {
	// Determine key type
	switch {
	case len(k) == 12 && k[4:5] == "-":
		ecd := elevencd{
			key: k,
		}
		ecd.validate(v)
	case len(k) == 11 && k[3:4] == "-":
		cd := cd{
			key: k,
		}
		cd.validate(v)
	case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
		oem := oem{
			key: k,
		}
		oem.validate(v)
	default:
		v <- false
	}
}
