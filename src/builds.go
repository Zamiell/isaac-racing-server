package server

var (
	// This must be kept in sync with the build exceptions in "isaac-tournament-bot"
	buildExceptions = [][]string{
		{}, // #0 - n/a

		// -------------------
		// Treasure Room items
		// -------------------

		// #1 - Cricket's Body
		// Azazel - The Brimstone beam inherits the split shots but they are not very good
		// Tainted Azazel - The Brimstone beam inherits the split shots but they are not very good
		{"Azazel", "Tainted Azazel"},

		{}, // #2 - Cricket's Head

		// #3 - Dead Eye
		// Azazel - The Brimstone beam prevents it from working
		// Lilith - It is hard to be accurate when shots come from Incubus
		// Keeper - It is hard to be accurate with triple shot
		// Tainted Azazel - The Brimstone beam prevents it from working
		// Tainted Keeper - It is hard to be accurate with quad shot
		{"Azazel", "Lilith", "Keeper", "Tainted Azazel", "Tainted Keeper"},

		// #4 - Death's Touch
		// The Forgotten - The piercing shots do nothing for the bone club
		// Tainted Forgotten - The piercing shots do nothing for the bone club
		{"The Forgotten", "Tainted Forgotten"},

		// #5 - Dr. Fetus
		// Tainted Forgotten - Very annoying to use with the skeleton body
		{"Tainted Forgotten"},

		// #6 - Ipecac
		// Azazel - The short-range brimstone causes self-damage
		// Tainted Eve - Can cause unavoidable damage if a clot is behind you or shoots at an obstacle near you
		{"Azazel", "Tainted Eve"},

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

		// #14 - Brimstone
		// Tainted Lilith - Gello fires the brim very slowly and auto-fire is not always accurate
		{"Tainted Lilith"},

		{}, // #15 - Maw of the Void

		// ----------------
		// Angel Room items
		// ----------------

		// #16 - Crown of Light
		// Eve - Eve cannot use the razor with this start
		// Tainted Magdalene - Crown is never full with the depleting hearts
		// Tainted Eve - Crown can't be active with the clot mechanic
		{"Eve", "Tainted Magdalene", "Tainted Eve"},

		{}, // #17 - Sacred Heart

		// #18 - Spirit Sword
		// Lilith - Annoying because the sword goes to Incubus
		// The Forgotten - No synergy with the bone
		// Tainted Forgotten - No synergy with the bone
		{"Lilith", "The Forgotten", "Tainted Forgotten"},

		{}, // #19 - Revelation

		// -----------------
		// Secret Room items
		// -----------------

		// #20 - Epic Fetus
		// Tainted Lilith - The target keeps moving and you can't control it, making it impossible to target enemies
		{"Tainted Lilith"},

		// ------------
		// Custom items
		// ------------

		// #21 - Sawblade
		// Bethany - Very complicated to play orbitals with her because she can't protect herself from losing the deal with soul hearts
		// Tainted Eve - Impossible to play orbitals with Tainted Eve's clots, they will disappear very quickly
		// Tainted Lost - With his health mechanic, it is too dangerous to use orbitals
		{"Bethany", "Tainted Eve", "Tainted Lost"},

		// ------
		// Builds
		// ------

		{}, // #22 - 20/20 + The Inner Eye
		{}, // #23 - Chocolate Milk + Steven

		// #24 - Godhead + Cupid's Arrow
		// Azazel - Small damage up for a tears down, resulting in a loss of DPS overall
		// The Forgotten - Does nothing with the bone club
		// Tainted Forgotten - Does nothing with the bone club
		{"Azazel", "The Forgotten", "Tainted Forgotten"},

		{}, // #25 - Haemolacria + The Sad Onion
		{}, // #26 - Incubus + Incubus

		// #27 - Monstro's Lung + The Sad Onion
		// Keeper - Huge tears down, resulting in a loss of DPS overall
		// Tainted Keeper - Huge tears down, resulting in a loss of DPS overall
		{"Keeper", "Tainted Keeper"},

		{}, // #28 - Technology + A Lump of Coal
		{}, // #29 - Twisted Pair + Twisted Pair
		{}, // #30 - Pointy Rib + Eve's Mascara

		// #31 - Fire Mind + Mysterious Liquid + 13 Luck
		// Azazel - The synergy is only useful with a tear build
		// The Forgotten - The synergy is only useful with a tear build
		// Tainted Azazel - The synergy is only useful with a tear build
		// Tainted Lost - Too dangerous to be synergistic
		// Tainted Forgotten - The synergy is only useful with a tear build
		{"Azazel", "The Forgotten", "Tainted Azazel", "Tainted Lost", "Tainted Forgotten"},

		// #32 - Eye of the Occult + Loki's Horns + 15 Luck
		// Azazel - Homing brimstone is too powerful, resulting in a build with a low-skill requirement
		// The Forgotten - It is only a damage up on the bone club
		// Tainted Azazel - Homing brimstone is too powerful, resulting in a build with a low-skill requirement
		// Tainted Forgotten - It is only a damage up on the bone club
		{"Azazel", "The Forgotten", "Tainted Azazel", "Tainted Forgotten"},

		{}, // #33 - Distant Admiration + Friend Zone + Forever Alone + BFFS!
	}
)
