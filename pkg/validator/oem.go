package validator

import "strconv"

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
	if !checkdigitCheck(th) {
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
