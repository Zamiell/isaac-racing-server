package main

/*
	Imports
*/

import (
	"math/rand"
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

// From: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func getRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var allowedCharacters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	randomString := make([]rune, length)
	for i := range randomString {
		randomString[i] = allowedCharacters[rand.Intn(len(allowedCharacters))]
	}
	return string(randomString)
}

// From: https://stackoverflow.com/questions/24122821/go-golang-time-now-unixnano-convert-to-milliseconds
func makeTimestamp() int64 {
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
