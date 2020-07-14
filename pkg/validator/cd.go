package validator

import "strconv"

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
	if !checkdigitCheck(main) {
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
