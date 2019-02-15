package validation

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
