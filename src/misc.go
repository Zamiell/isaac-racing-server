package server

import (
	"math"
	"strconv"
	"time"
	"unicode"
)

func intInSlice(a int, slice []int) bool {
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

// getTimestamp returns the Epoch timestamp in milliseconds.
// From: https://stackoverflow.com/questions/24122821/go-golang-time-now-unixnano-convert-to-milliseconds
func getTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func stringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
