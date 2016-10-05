package main

/*
 *  Imports
 */

import (
	"unicode"
)

/*
 *  Miscellaneous functions
 */

func intInSlice(a int, slice []int) bool {
	for _, b := range slice {
		if b == a {
			return true
		}
	}
	return false
}

// From: http://stackoverflow.com/questions/31961882/how-to-check-if-there-is-a-special-character-in-string-or-if-a-character-is-a-sp
func hasSymbol(str string) bool {
	for _, letter := range str {
		if unicode.IsSymbol(letter) {
			return true
		}
	}
	return false
}
