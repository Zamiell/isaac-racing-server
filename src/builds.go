package server

var (
	// This must be kept in sync with the build exceptions in "isaac-tournament-bot".
	buildExceptions = [][]string{
		// #0 - n/a
		{},

		// -------------------
		// Treasure Room items
		// -------------------

		// #1 - Cricket's Head
		{},

		// #2 - Dr. Fetus
		{
			// Very annoying to use with the skeleton body.
			"Tainted Forgotten", // 35
		},

		// #3 - Epic Fetus
		{
			// The target keeps moving and you can't control it, making it impossible to target
			// enemies.
			"Tainted Lilith", // 32
		},

		// #4 - Ipecac
		{
			// The short-range brimstone causes self-damage.
			"Azazel", // 7

			// Can cause unavoidable damage if a clot is behind you or shoots at an obstacle near
			// you.
			"Tainted Eve", // 26
		},

		// #5 - Magic Mushroom
		{},

		// #6 - Mom's Knife
		{},

		// #7 - Polyphemus
		{},

		// #8 - Proptosis
		{},

		// #9 - Tech X
		{},

		// ----------------
		// Devil Room items
		// ----------------

		// #10 - Brimstone
		{
			// Gello fires the brim very slowly and auto-fire is not always accurate.
			"Tainted Lilith", // 32
		},

		// #11 - Maw of the Void
		{},

		// ----------------
		// Angel Room items
		// ----------------

		// #12 - Crown of Light
		{
			// Eve cannot use the razor with this start.
			"Eve", // 5

			// Crown is never full with the depleting hearts.
			"Tainted Magdalene", // 22

			// Crown can't be active with the clot mechanic.
			"Tainted Eve", // 26
		},

		// #13 - Godhead
		{
			// Small damage up for a tears down, resulting in a loss of DPS overall.
			"Azazel", // 7

			// Does nothing with the bone club.
			"The Forgotten", // 16

			// Does nothing with the bone club.
			"Tainted Forgotten", // 35
		},

		// #14 - Sacred Heart
		{},

		// #15 - Spirit Sword
		{
			// Annoying because the sword goes to Incubus.
			"Lilith", // 13

			// No synergy with the bone club.
			"The Forgotten", // 16

			// No synergy with the bone club.
			"Tainted Forgotten", // 35
		},

		// #16 - Revelation
		{},

		// ------
		// Builds
		// ------

		// #17 - 20/20 + The Inner Eye
		{},

		// #18 - C Section + Steven
		{},

		// #19 - Chocolate Milk + Steven
		{},

		// #20 - Cricket's Body + Steven
		{
			// The Brimstone beam inherits the split shots but they are not very good.
			"Azazel", // 7

			// The Brimstone beam inherits the split shots but they are not very good.
			"Tainted Azazel", // 28
		},

		// #21 - Dead Eye + Jesus Juice
		{
			// The Brimstone beam prevents it from working.
			"Azazel", // 7

			// It is hard to be accurate when shots come from Incubus.
			"Lilith", // 13

			// It is hard to be accurate with triple shot.
			"Keeper", // 14

			// The Brimstone beam prevents it from working.
			"Tainted Azazel", // 28

			// It is hard to be accurate with quad shot.
			"Tainted Keeper", // 33
		},

		// #22 - Death's Touch + Sad Onion
		{
			// The piercing shots do nothing for Azazel's Brimstone.
			"Azazel", // 7

			// The piercing shots do nothing for the bone club.
			"The Forgotten", // 16

			// The piercing shots do nothing for Tainted Azazel's Brimstone.
			"Tainted Azazel", // 28

			// The piercing shots do nothing for the bone club.
			"Tainted Forgotten", // 35
		},

		// #23 - Fire Mind + 13 Luck
		{
			// The synergy is only useful with a tear build.
			"Azazel", // 7

			// Luck does not apply to Incubus for some reason.
			"Lilith", // 13

			// The synergy is only useful with a tear build.
			"The Forgotten", // 16

			// The synergy is only useful with a tear build.
			"Tainted Azazel", // 28

			// Too dangerous to be synergistic.
			"Tainted Lost", // 31

			// The synergy is only useful with a tear build.
			"Tainted Forgotten", // 35
		},

		// #24 - Haemolacria + The Sad Onion
		{},

		// #25 - Incubus + Twisted Pair
		{},

		// #26 - Monstro's Lung + The Sad Onion
		{
			// Huge tears down, resulting in a loss of DPS overall.
			"Keeper", // 14

			// Tears down, worse than having no starter with the fetus.
			"Tainted Lilith", // 32

			// Huge tears down, resulting in a loss of DPS overall.
			"Tainted Keeper", // 33
		},

		// #27 - Technology + A Lump of Coal
		{},

		// #28 - Tech.5 + Jesus Juice
		{},

		// #29 - Pointy Rib + Eve's Mascara
		{},

		// #30 - Sawblade + Fate
		{
			// Very complicated to play orbitals with her because she can't protect herself from
			// losing the deal with soul hearts.
			"Bethany", // 18

			// Impossible to play orbitals with Tainted Eve's clots, they will disappear very
			// quickly.
			"Tainted Eve", // 26

			// With his health mechanic, it is too dangerous to use orbitals.
			"Tainted Lost", // 31
		},

		// #31 - Eye of the Occult + Loki's Horns + 15 Luck
		{
			// Homing brimstone is too powerful, resulting in a build with a low-skill requirement.
			"Azazel", // 7

			// It is only a damage up on the bone club.
			"The Forgotten", // 16

			// Homing brimstone is too powerful, resulting in a build with a low-skill requirement.
			"Tainted Azazel", // 28

			// It is only a damage up on the bone club.
			"Tainted Forgotten", // 35
		},

		// #32 - Distant Admiration + Friend Zone + Forever Alone + BFFS!
		{},
	}
)
