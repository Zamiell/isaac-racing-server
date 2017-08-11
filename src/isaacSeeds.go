package main

import (
	"math/rand"
	"time"
)

func isaacGetRandomSeed() string {
	// Get a random uint32
	rand.Seed(time.Now().UnixNano())
	seed := rand.Uint32()
	return isaacSeedToString(seed)
}

// This algorithm has been reverse engeineered from the game's binary by
// Killburn (and blcd, independently)
func isaacSeedToString(num uint32) string {
	chars := "ABCDEFGHJKLMNPQRSTWXYZ01234V6789"

	// Checksum
	// https://www.reddit.com/r/bindingofisaac/comments/2wvp6h/is_it_known_what_makes_a_seed_valid/csdppvx/
	var x byte
	x = 0
	tnum := num
	for tnum != 0 {
		x += byte(tnum)
		x += byte(x + (x >> 7))
		tnum >>= 5
	}
	num ^= 0x0FEF7FFD
	tnum = num<<8 | uint32(x)

	// Build the string
	s := ""
	for i := 0; i < 8; i++ {
		var charIndex int
		if i >= 0 && i <= 5 {
			charIndex = int(num >> uint(27-(i*5)) & 0x1F)
		} else if i == 6 {
			charIndex = int(tnum >> 5 & 0x1F)
		} else if i == 7 {
			charIndex = int(tnum & 0x1F)
		}
		s += string(chars[charIndex])
	}

	// Insert a space in the middle
	s = s[:4] + " " + s[4:]

	return s
}
