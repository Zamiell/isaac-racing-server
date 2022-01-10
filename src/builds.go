package server

var (
	buildExceptions = [][]string{
		{}, // #0 - n/a

		// -------------------
		// Treasure Room items
		// -------------------

		// #1 - Cricket's Body
		// Cain - The range down causes Cain to barely be able to hit anything
		// Samson - ???
		// Azazel - The Brimstone beam inherits the split shots but they are not very good
		// Tainted Azazel - The Brimstone beam inherits the split shots but they are not very good
		{"Cain", "Samson", "Azazel", "Tainted Azazel"},

		{}, // #2 - Cricket's Head

		// #3 - Dead Eye
		// Azazel - The Brimstone beam prevents it from working
		// Lilith - It is hard to be accurate when shots come from Incubus
		// Keeper - It is hard to be accurate with triple shot
		// Tainted Azazel - The Brimstone beam prevents it from working
		// Tainted Keeper - It is hard to be accurate with triple shot
		{"Azazel", "Lilith", "Keeper", "Tainted Azazel", "Tainted Keeper"},

		// #4 - Death's Touch
		// The Forgotten - The piercing shots do nothing for the bone club
		{"The Forgotten"},

		// #5 - Dr. Fetus
		{},

		// #6 - Ipecac
		// Azazel - The short-range brimstone causes self-damage
		{"Azazel"},

		{}, // #7 - Magic Mushroom
		{}, // #8 - Mom's Knife
		{}, // #9 - Polyphemus
		{}, // #10 - Proptosis
		{}, // #11 - Tech.5
		{}, // #12 - Tech X
		{}, // #13 - C Section

		// ----------------
		// Devil Room items
		// ----------------

		{}, // #14 - Brimstone
		{}, // #15 - Maw of the Void

		// ----------------
		// Angel Room items
		// ----------------

		// #16 - Crown of Light
		// Eve - Eve cannot use the razor with this start
		{"Eve"},

		{}, // #17 - Sacred Heart
		{}, // #18 - Spirit Sword
		{}, // #19 - Revelation

		// -----------------
		// Secret Room items
		// -----------------

		{}, // #20 - Epic Fetus

		// ------------
		// Custom items
		// ------------

		{}, // #21 - Sawblade

		// ------
		// Builds
		// ------

		{}, // #22 - 20/20 + The Inner Eye
		{}, // #23 - Chocolate Milk + Steven
		{}, // #24 - Godhead + Cupid's Arrow
		{}, // #25 - Haemolacria + The Sad Onion

		// #26 - Incubus + Incubus
		// The Forgotten - The extra Incubus with bone clubs are not very helpful
		{"The Forgotten"},

		// #27 - Monstro's Lung + The Sad Onion
		// Keeper - The charge rate is too low to be very useful
		// Tainted Keeper - The charge rate is too low to be very useful
		{"Keeper", "Tainted Keeper"},

		{}, // #28 - Technology + A Lump of Coal
		{}, // #29 - Twisted Pair + Twisted Pair
		{}, // #30 - Pointy Rib + Eve's Mascara

		// #31 - Fire Mind + Mysterious Liquid + 13 Luck
		// Azazel - The synergy is only useful with a tear build
		// The Forgotten - The synergy is only useful with a tear build
		// Tainted Azazel - The synergy is only useful with a tear build
		{"Azazel", "The Forgotten", "Tainted Azazel"},

		// #32 - Eye of the Occult + Loki's Horns + 15 Luck
		// Azazel - The synergy is only useful with a tear build
		// The Forgotten - The synergy is only useful with a tear build
		// Tainted Azazel - The synergy is only useful with a tear build
		{"Azazel", "The Forgotten", "Tainted Azazel"},

		{}, // #33 - Distant Admiration + Friend Zone + Forever Alone + BFFS!
	}
)
