package server

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/*
	Diversity ruleset definitions
*/

var goldenTrinketModifier = 32768

var validDiversityActiveItems = [...]int{
	// Rebirth items
	33, 34, 35, 36, 37, 38, 39, 40, 41, 42,
	44, 45, 47, 49, 56, 58, 65, 66, 77, 78, // Book of Belial Birthright (59) is banned
	83, 84, 85, 86, 93, 97, 102, 105, 107, 111,
	123, 124, 126, 127, 130, 133, 135, 136, 137, 145,
	146, 147, 158, 160, 164, 166, 171, 175, 177, 181,
	186, 192, 282, 285, 286, 287, 288, 289, 290, 291, // D100 (283) and D4 (284) are banned
	292, 293, 294, 295, 296, 297, 298, 323, 324, 325,
	326, 338,

	// Afterbirth items
	347, 348, 349, 351, 352, 357, 382, 383, 386, 396,
	406, 419, 421, 422, 427, 434, 437, 439, 441, // Broken Glass Cannon (474) is skipped

	// Afterbirth+ items
	475, 476, 477, 478, 479, 480, 481, 482, 483, 484,
	485, 486, 487, 488, 490, 504, 507, 510, // D Infinity (489) is banned

	// Booster Pack items
	512, 515, 516, 521, 522, 523, 527, 536, 545,

	// Repentance items
	263, 555, 556, 557, 577, 578, 580, 582, 585, 604, // Mom's Shovel (552) is banned
	605, 609, 611, 623, 625, 628, 631, 635, 638, 639, // Genesis (622) and R Key (636) are banned
	640, 642, 650, 653, 655, 685, 687, 704, 705, 706, // Esau Jr (703) is banned
	709, 710, 711, 712, 713, 719, 720, 722, 723, 728, // Recall (714) and Hold (715) are banned
	729,
}

var validDiversityPassiveItems = [...]int{
	// Rebirth items
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	// <3 (15), Raw Liver (16), Lunch (22), Dinner (23), Dessert (24), Breakfast (25),
	// and Rotten Meat (26) are banned
	11, 12, 13, 14, 17, 18, 19, 20, 21, 27,
	// Mom's Underwear (29), Moms Heels (30) and Moms Lipstick (31) are banned
	28, 32, 46, 48, 50, 51, 52, 53, 54, 55, 57,
	60, 62, 63, 64, 67, 68, 69, 70, 71, 72,
	73, 74, 75, 76, 79, 80, 81, 82, 87, 88,
	89, 90, 91, 94, 95, 96, 98, 99, 100, 101, // Super Bandage (92) is banned
	103, 104, 106, 108, 109, 110, 112, 113, 114, 115,
	116, 117, 118, 119, 120, 121, 122, 125, 128, 129,
	131, 132, 134, 138, 139, 140, 141, 142, 143, 144,
	148, 149, 150, 151, 152, 153, 154, 155, 156, 157,
	159, 161, 162, 163, 165, 167, 168, 169, 170, 172,
	173, 174, 178, 179, 180, 182, 183, 184, 185, 187, // Stem Cells (176) is banned
	188, 189, 190, 191, 193, 195, 196, 197, 198, 199, // Magic 8 Ball (194) is banned
	200, 201, 202, 203, 204, 205, 206, 207, 208, 209,
	210, 211, 212, 213, 214, 215, 216, 217, 218, 219,
	220, 221, 222, 223, 224, 225, 227, 228, 229, 230, // Black Lotus (226) is banned
	// Key Piece #1 (238) and Key Piece #2 (239) are banned
	231, 232, 233, 234, 236, 237, 240, 241, 242, 243,
	244, 245, 246, 247, 248, 249, 250, 251, 252, 254, // Magic Scab (253) is banned
	255, 256, 257, 259, 260, 261, 262, 264, 265, 266, // Missing No. (258) is banned
	267, 268, 269, 270, 271, 272, 273, 274, 275, 276,
	277, 278, 279, 280, 281, 299, 300, 301, 302, 303,
	304, 305, 306, 307, 308, 309, 310, 311, 312, 313,
	314, 315, 316, 317, 318, 319, 320, 321, 322, 327,
	// The Body (334) and Safety Pin (339) are banned
	328, 329, 330, 331, 332, 333, 335, 336, 337, 340,
	341, 342, 343, 345, // Match Book (344) and A Snack (346) are banned

	// Afterbirth items
	350, 353, 354, 356, 358, 359, 360, 361, 362, 363, // Mom's Pearls (355) is banned
	364, 365, 366, 367, 368, 369, 370, 371, 372, 373,
	374, 375, 376, 377, 378, 379, 380, 381, 384, 385,
	387, 388, 389, 390, 391, 392, 393, 394, 395, 397,
	398, 399, 400, 401, 402, 403, 404, 405, 407, 408,
	409, 410, 411, 412, 413, 414, 415, 416, 417, 418,
	420, 423, 424, 425, 426, 429, 430, 431, 432, 433, // PJs (428) is banned
	435, 436, 438, 440,

	// Afterbirth+ items
	442, 443, 444, 445, 446, 447, 448, 449, 450, 451,
	// Dad's Lost Coin (455) and Midnight Snack (456) are banned
	452, 453, 454, 457, 458, 459, 460, 461, 462, 463,
	464, 465, 466, 467, 468, 469, 470, 471, 472, 473,
	491, 492, 493, 494, 495, 496, 497, 498, 499, 500,
	501, 502, 503, 505, 506, 508, 509,

	// Booster Pack items
	511, 513, 514, 517, 518, 519, 520, 524, 525, 526,
	528, 529, 530, 531, 532, 533, 534, 535, 537, 538,
	// Marrow (541) is banned
	539, 540, 542, 543, 544, 546, 547, 548, 549,

	// Repentance items
	// Broken Shovel Top (550) and Broken Shovel Bottom (551) are banned
	553, 554, 558, 559, 560, 561, 562, 563, 564, 565,
	566, 567, 568, 569, 570, 571, 572, 573, 574, 575,
	576, 579, 581, 583, 584, 586, 588, 589, 591, 592,
	593, 594, 595, 596, 597, 598, 599, 600, 601, 602,
	603, 606, 607, 608, 610, 612, 614, 615, 616, 617,
	// Knife Piece #1 (626), Knife Piece #2 (627), and Dogma (633) are banned
	618, 619, 621, 624, 629, 632, 634, 637, 641, 643,
	// Damocles (Passive) (656) is banned
	644, 645, 646, 647, 649, 651, 652, 654, 657, 658,
	// Tropicamide (659) and Dad's Note (668) are banned
	660, 661, 663, 664, 665, 667, 669, 670, 671, 672,
	673, 674, 675, 676, 677, 678, 679, 680, 681, 682,
	683, 684, 686, 688, 689, 690, 691, 692, 693, 694,
	// Supper (707) is banned
	695, 696, 697, 698, 699, 700, 701, 702, 708, 716,
	717, 724, 725, 726, 727, 730, 731, 732,

	// Custom items
	1003, // Sawblade
}

var validDiversityTrinkets = [...]int{
	// Rebirth trinkets
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 48, 49, 50, 51, // 47 is Polaroid (obsolete)
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,

	// Afterbirth trinkets
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81,
	82, 83, 84, 86, 87, 88, 89, 90, // Karma (85) is banned

	// Afterbirth+ trinkets
	91, 92, 93, 94, 95, 96, 97, 98, 99, 100,
	101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
	111, 112, 113, 114, 115, 116, 117, 118, 119,

	// Booster pack trinkets
	120, 121, 122, 123, 124, 125, 126, 127, 128,

	// Repentance trinkets
	129, 130, 131, 132, 133, 134, 135, 136, 137,
	138, 139, 140, 141, 142, 143, 144, 145, 146,
	147, 148, 149, 150, 151, 152, 153, 155, 156, // Dice Bag (154) is banned
	157, 158, 159, 160, 161, 162, 163, 164, 165,
	166, 167, 168, 169, 170, 171, 172, 173, 174,
	175, 176, 177, 178, 179, 180, 181, 182, 183,
	184, 185, 186, 187, 188, 189,
}

var taintedLostItemsBanned = []int{
	9, 10, 11, 13, 20, 36, 40, 45, 53, 60, 62, 72,
	73, 78, 81, 82, 83, 96, 98, 108, 112, 115, 117,
	119, 126, 129, 133, 135, 138, 142, 146, 148, 156,
	157, 159, 161, 162, 172, 173, 179, 180, 184, 185,
	186, 193, 202, 204, 205, 207, 210, 211, 212, 214,
	218, 219, 223, 225, 227, 232, 236, 242, 243, 262,
	264, 265, 274, 276, 278, 279, 281, 282, 290, 292,
	296, 298, 299, 301, 302, 303, 311, 312, 313, 314,
	321, 326, 332, 335, 337, 354, 363, 371, 375, 387,
	404, 409, 412, 413, 423, 425, 436, 442, 448, 449,
	452, 457, 482, 486, 487, 493, 501, 508, 509, 513,
	514, 525, 535, 538, 539, 542, 543, 546, 549, 554,
	560, 568, 569, 571, 581, 610, 612, 614, 615, 629,
	639, 645, 652, 655, 658, 664, 671, 672, 674, 675,
	676, 677, 688, 690, 692, 693, 694, 695, 702, 713,
	724,
}

var characterItemBlacklist = map[string][]int{
	"Isaac":             {534},           // Schoolbag
	"Magdalene":         {534},           // Schoolbag
	"Cain":              {46},            // Lucky foot
	"Judas":             {534},           // Schoolbag
	"Blue Baby":         {534},           // Schoolbag
	"Eve":               {117, 122, 534}, // Dead Bird, Whore of Babylon, Schoolbag
	"Samson":            {157},           // Bloody Lust
	"Lazarus":           {214},           // Anemic
	"The Lost":          {313, 534},      // Holy Mantle, Schoolbag
	"Lilith":            {412, 534},      // Cambion Conception, Schoolbag
	"Keeper":            {230, 534, 672}, // Abaddon, Schoolbag, A Pound of Flesh
	"Apollyon":          {441, 534, 625}, // Mega Blast, Schoolbag, Mega Mush
	"Bethany":           {230, 584},      // Abaddon, Book of Virtues
	"Jacob & Esau":      {534},           // Schoolbag
	"Tainted Isaac":     {619},           // Birthright
	"Tainted Magdalene": {205, 534},      // Sharp Plug, Schoolbag
	"Tainted Judas":     {534},           // Schoolbag
	"Tainted Blue Baby": {534},           // Schoolbag
	"Tainted Eve":       {534, 713},      // Schoolbag, Sumptorium
	"Tainted Azazel":    {726},           // Hemoptysis
	"Tainted Lazarus":   {534},           // Schoolbag
	"Tainted Eden":      {619},           // Birthright
	"Tainted Lost":      taintedLostItemsBanned,
	"Tainted Keeper":    {230, 672}, // Abaddon, A Pound of Flesh
	"Tainted Apollyon":  {534},      // Schoolbag
	"Tainted Bethany":   {534},      // Schoolbag
	"Tainted Jacob":     {534, 722}, // Schoolbag, Anima Sola
}

/*
	Diversity helper functions
*/

func diversityGetSeed(ruleset Ruleset) string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Get 1 random unique active item
	var items []int
	item := validDiversityActiveItems[rand.Intn(len(validDiversityActiveItems))] // nolint: gosec
	items = append(items, item)

	// Get 3 random unique passive items
	for i := 1; i <= 3; i++ {
		for {
			// Initialize the PRNG and get a random element from the slice
			// (if we don't do this, it will use a seed of 1)
			randomIndex := rand.Intn(len(validDiversityPassiveItems)) // nolint: gosec
			item := validDiversityPassiveItems[randomIndex]

			if itemIsBannedOnThisCharacter(item, ruleset.Character) {
				continue
			}

			// Ensure this item is unique
			if intInSlice(item, items) {
				continue
			}

			items = append(items, item)
			break
		}
	}

	// Get 1 random trinket
	randomIndex := rand.Intn(len(validDiversityTrinkets)) // nolint: gosec
	trinket := validDiversityTrinkets[randomIndex]
	// The server has a 10% chance to make the trinket golden
	if rand.Intn(10) == 0 { // nolint: gosec
		trinket += goldenTrinketModifier
	}
	items = append(items, trinket)

	// The "seed" value is used to communicate the 5 random diversity items to the client
	seed := ""
	for _, item := range items {
		seed += strconv.Itoa(item) + ","
	}
	seed = strings.TrimSuffix(seed, ",") // Remove the trailing comma

	return seed
}

func itemIsBannedOnThisCharacter(item int, character string) bool {
	itemBlacklist, ok := characterItemBlacklist[character]
	if !ok {
		// There are no banned items for this character
		return false
	}

	return intInSlice(item, itemBlacklist)
}
