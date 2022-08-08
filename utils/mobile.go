package utils

import "strings"

func IsDigit(r rune) bool {
	return r == '0' || r == '1' || r == '2' || r == '3' || r == '4' || r == '5' || r == '6' || r == '7' || r == '8' || r == '9'
}

func ParseMobile(mobile string) string {
	parsedMobile := ""

	// remove all `non digit` characters
	for _, c := range mobile {
		if IsDigit(c) {
			parsedMobile = parsedMobile + string(c)
		}
	}

	// remove (one or more) `0` from the beginning
	parsedMobile = strings.TrimLeft(parsedMobile, "0")

	// remove `91` if number of digits is 12
	if len(parsedMobile) == 12 {
		parsedMobile = strings.TrimPrefix(parsedMobile, "91")
	}

	// mobile number has to be `10` digits long
	if len(parsedMobile) != 10 {
		return ""
	}

	// mobile number should start with `6, 7, 8, or 9`
	firstDigit := parsedMobile[0]
	allowedFirstDigits := []rune{'9', '8', '7', '6', '5'}
	isValid := false
	for _, allowedFirstDigit := range allowedFirstDigits {
		if rune(firstDigit) == allowedFirstDigit {
			isValid = true
			break
		}
	}
	if !isValid {
		return ""
	}

	return parsedMobile
}

func IsMobileValid(mobile string) bool {
	return ParseMobile(mobile) != ""
}
