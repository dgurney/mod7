package validation

import (
	"strconv"
)

func batchValidateCDKey(key string, v chan bool) {
	site, err := strconv.ParseInt(key[0:3], 10, 0)
	if err != nil {
		v <- false
		return
	}
	main, err := strconv.ParseInt(key[4:11], 10, 0)
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

func batchValidateECDKey(key string, v chan bool) {
	_, err := strconv.ParseInt(key[0:4], 10, 0)
	if err != nil {
		v <- false
		return
	}
	main, err := strconv.ParseInt(key[5:12], 10, 0)
	if err != nil {
		v <- false
		return
	}

	// Error is safe to discard since we checked if it's a number before.
	last, _ := strconv.ParseInt(key[3:4], 10, 0)
	third, _ := strconv.ParseInt(key[2:3], 10, 0)
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

func batchValidateOEMKey(key string, v chan bool) {
	_, err := strconv.ParseInt(key[0:5], 10, 0)
	if err != nil {
		v <- false
		return
	}
	th, err := strconv.ParseInt(key[10:17], 10, 0)
	if err != nil {
		v <- false
		return
	}
	_, err = strconv.ParseInt(key[18:], 10, 0)
	if err != nil {
		v <- false
		return
	}
	julian, err := strconv.ParseInt(key[0:3], 10, 0)
	if julian == 0 || julian > 366 {
		v <- false
		return
	}
	year := key[3:5]
	validYears := map[string]string{"95": "95", "96": "96", "97": "97", "98": "98", "99": "99", "00": "00", "01": "01", "02": "02", "03": "03"}
	_, valid := validYears[year]
	if !valid {
		v <- false
		return
	}
	if key[6:9] != "OEM" {
		v <- false
		return
	}
	third := key[10:17]
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

// BatchValidate is typically used to validate individual keys from an array.
func BatchValidate(k string, v chan bool) {
	// Determine key type
	switch {
	case len(k) == 12 && k[4:5] == "-":
		batchValidateECDKey(k, v)
	case len(k) == 11 && k[3:4] == "-":
		batchValidateCDKey(k, v)
	case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
		batchValidateOEMKey(k, v)
	default:
		v <- false
	}
}
