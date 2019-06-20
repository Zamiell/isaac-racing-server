package main

import (
	"strconv"
	"time"
	"unicode"
)

/*
	Miscellaneous functions
*/

func intInSlice(a int, slice []int) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

func stringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

func isAlphaNumericUnderscore(str string) bool {
	isValid := true
	for _, character := range str {
		if unicode.IsLetter(character) {
			continue
		}

		if character >= '0' && character <= '9' {
			continue
		}

		if character == '_' {
			continue
		}

		isValid = false
		break
	}

	return isValid
}

// From: https://stackoverflow.com/questions/24122821/go-golang-time-now-unixnano-convert-to-milliseconds
func getTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func getOrdinal(n int) string {
	s := []string{"th", "st", "nd", "rd"}
	v := n % 100
	test := (v - 20) % 10
	var ord string
	if test >= 0 && test <= 3 {
		ord = s[test]
	} else if v >= 0 && v <= 3 {
		ord = s[v]
	} else {
		ord = s[0]
	}
	return strconv.Itoa(n) + ord
}
