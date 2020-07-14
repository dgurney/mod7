package validator

import "strconv"

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
