package main

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
