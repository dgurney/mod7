package validation

import "strconv"

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

// BatchValidate validates an array of keys. Compared to the regular validation function, the output is terse for easier grepping.
func BatchValidate(k string, v chan bool) {
	// Determine key type
	switch {
	case len(k) == 11 && k[3:4] == "-":
		go batchValidateCDKey(k, v)
	case len(k) == 23 && k[5:6] == "-" && k[9:10] == "-" && k[17:18] == "-" && len(k[18:]) == 5:
		go batchValidateOEMKey(k, v)
	default:
		v <- false
	}
}
