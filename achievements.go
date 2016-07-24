package main

func achievementsInit() {
	// Define all of the achievements
	achievementList := map[int][]string{
		// Complete X races
		1: {"This Wasn't So Bad", "Complete your first race."},
		2: {"I Think I Have the Hang of This", "Complete 10 races."},
		3: {"Intermediate Racist", "Complete 50 races."},
		4: {"Expert Racist", "Complete 100 races."},
		5: {"Seasoned Veteran", "Complete 500 races."},
		6: {"Orange Juice", "Complete 1000 races."},

		7: {"Dipping Your Toe in the Water", "Complete a race with every ruleset."},
		8: {"Experimental Treatment", "Complete 10 races with every ruleset."},
		9: {"Jack of All Trades", "Complete 100 races with every ruleset."},

		// Get X average
		10: {"Average Joe", "Have an unseeded average time of 16:30 or less (with at least 50 races played and under a 25% forfeit rate)."},
		11: {"There's So Many Strats and So Little Time", "Have an unseeded average time of 16:00 or less (with at least 50 races played and under a 25% forfeit rate)."},
		12: {"There's No Skill in an RNG Based Game", "Have an unseeded average time of 15:30 or less (with at least 50 races played and under a 25% forfeit rate)."},
		13: {"Master of Consistency", "Have an unseeded average time of 15:00 or less (with at least 50 races played and under a 25% forfeit rate)."},

		// Every starting item
		14: {"Well Rounded - Unseeded", "Complete an unseeded race with every starting item."},
		15: {"Well Rounded - Seeded", "Complete a seeded race with every starting item."},
		16: {"Well Rounded - Practice", "Complete a practice mode run with every starting item."},

		// Streaks
		17: {"I'm Having a Good Day", "Finish 1st place in 3 races in a row with 5 people or more."},
		18: {"A Literal God, Tasteless", "Finish 1st place in 5 races in a row with 5 people or more."},
		19: {"I'm Having a Bad Day", "Quit out of 3 races in a row with 5 people or more."},

		// Complete a race with X time
		20: {"Pretty Darn Fast", "Complete an unseeded race in less than 13 minutes."},
		21: {"Speeding Bullet", "Complete an unseeded race in less than 12 minutes."},
		22: {"Fast Like Sanic", "Complete an unseeded race in less than 11 minutes."},
		23: {"Dea1hly Fast", "Complete an unseeded race in less than 10 minutes."},
		24: {"Fizzy Giraffe", "Complete an unseeded race in less than 9 minutes."},

		// Complete a race with X time with X item, unseeded
		101: {"Unseeded Item Mastery - 20/20", "Complete an unseeded race in less than 12 minutes with a 20/20 start."},
		102: {"Unseeded Item Mastery - Chocolate Milk", "Complete an unseeded race in less than 12 minutes with a Chocolate Milk start."},
		103: {"Unseeded Item Mastery - Cricket's Body", "Complete an unseeded race in less than 12 minutes with a Cricket's Body start."},
		104: {"Unseeded Item Mastery - Cricket's Head", "Complete an unseeded race in less than 12 minutes with a Cricket's Head start."},
		105: {"Unseeded Item Mastery - Dead Eye", "Complete an unseeded race in less than 12 minutes with a Dead Eye start."},
		106: {"Unseeded Item Mastery - Death's Touch", "Complete an unseeded race in less than 12 minutes with a Death's Touch start."},
		107: {"Unseeded Item Mastery - Dr. Fetus", "Complete an unseeded race in less than 12 minutes with a Dr. Fetus start."},
		108: {"Unseeded Item Mastery - Epic Fetus", "Complete an unseeded race in less than 12 minutes with a Epic Fetus start."},
		109: {"Unseeded Item Mastery - Ipecac", "Complete an unseeded race in less than 12 minutes with a Ipecac start."},
		110: {"Unseeded Item Mastery - Judas' Shadow", "Complete an unseeded race in less than 12 minutes with a Judas' Shadow start."},
		111: {"Unseeded Item Mastery - Lil' Brimstone", "Complete an unseeded race in less than 12 minutes with a Lil' Brimstone start."},
		112: {"Unseeded Item Mastery - Magic Mushroom", "Complete an unseeded race in less than 12 minutes with a Magic Mushroom start."},
		113: {"Unseeded Item Mastery - Mom's Knife", "Complete an unseeded race in less than 12 minutes with a Mom's Knife start."},
		114: {"Unseeded Item Mastery - Monstro's Lung", "Complete an unseeded race in less than 12 minutes with a Monstro's Lung start."},
		115: {"Unseeded Item Mastery - Polyphemus", "Complete an unseeded race in less than 12 minutes with a Polyphemus start."},
		116: {"Unseeded Item Mastery - Proptosis", "Complete an unseeded race in less than 12 minutes with a Proptosis start."},
		117: {"Unseeded Item Mastery - Sacrificial Dagger", "Complete an unseeded race in less than 12 minutes with a Sacrificial Dagger start."},
		118: {"Unseeded Item Mastery - Tech.5", "Complete an unseeded race in less than 12 minutes with a Tech.5 start."},
		119: {"Unseeded Item Mastery - Tech X", "Complete an unseeded race in less than 12 minutes with a Tech X start."},

		// Complete a race with X time with X item, seeded
		201: {"Seeded Item Mastery - 20/20", "Complete a seeded race in less than 11 minutes with a 20/20 start."},
		202: {"Seeded Item Mastery - Chocolate Milk", "Complete a seeded race in less than 11 minutes with a Chocolate Milk start."},
		203: {"Seeded Item Mastery - Cricket's Body", "Complete a seeded race in less than 11 minutes with a Cricket's Body start."},
		204: {"Seeded Item Mastery - Cricket's Head", "Complete a seeded race in less than 11 minutes with a Cricket's Head start."},
		205: {"Seeded Item Mastery - Dead Eye", "Complete a seeded race in less than 11 minutes with a Dead Eye start."},
		206: {"Seeded Item Mastery - Death's Touch", "Complete a seeded race in less than 11 minutes with a Death's Touch start."},
		207: {"Seeded Item Mastery - Dr. Fetus", "Complete a seeded race in less than 11 minutes with a Dr. Fetus start."},
		208: {"Seeded Item Mastery - Epic Fetus", "Complete a seeded race in less than 11 minutes with a Epic Fetus start."},
		209: {"Seeded Item Mastery - Ipecac", "Complete a seeded race in less than 11 minutes with a Ipecac start."},
		210: {"Seeded Item Mastery - Judas' Shadow", "Complete a seeded race in less than 11 minutes with a Judas' Shadow start."},
		211: {"Seeded Item Mastery - Lil' Brimstone", "Complete a seeded race in less than 11 minutes with a Lil' Brimstone start."},
		212: {"Seeded Item Mastery - Magic Mushroom", "Complete a seeded race in less than 11 minutes with a Magic Mushroom start."},
		213: {"Seeded Item Mastery - Mom's Knife", "Complete a seeded race in less than 11 minutes with a Mom's Knife start."},
		214: {"Seeded Item Mastery - Monstro's Lung", "Complete a seeded race in less than 11 minutes with a Monstro's Lung start."},
		215: {"Seeded Item Mastery - Polyphemus", "Complete a seeded race in less than 11 minutes with a Polyphemus start."},
		216: {"Seeded Item Mastery - Proptosis", "Complete a seeded race in less than 11 minutes with a Proptosis start."},
		217: {"Seeded Item Mastery - Sacrificial Dagger", "Complete a seeded race in less than 11 minutes with a Sacrificial Dagger start."},
		218: {"Seeded Item Mastery - Tech.5", "Complete a seeded race in less than 11 minutes with a Tech.5 start."},
		219: {"Seeded Item Mastery - Tech X", "Complete a seeded race in less than 11 minutes with a Tech X start."},
		220: {"Seeded Item Mastery - Brimstone", "Complete a seeded race in less than 11 minutes with a Brimstone start."},
		221: {"Seeded Item Mastery - Incubus", "Complete a seeded race in less than 11 minutes with a Incubus start."},
		222: {"Seeded Item Mastery - Maw of the Void", "Complete a seeded race in less than 11 minutes with a Maw of the Void start."},
		223: {"Seeded Item Mastery - Crown of Light", "Complete a seeded race in less than 11 minutes with a Crown of Light start."},
		224: {"Seeded Item Mastery - Godhead", "Complete a seeded race in less than 11 minutes with a Godhead start."},
		225: {"Seeded Item Mastery - Sacred Heart", "Complete a seeded race in less than 11 minutes with a Sacred Heart start."},

		// Item synergies (2 items)
		301: {"Score, the Shop Has Mapping!", "Buy a piece of mapping from a shop after already having The Mind."},
		302: {"My Bombs Don't Hurt Me", "Complete a race where you had Dr. Fetus and Ipecac."},
		303: {"Green Lung Best Lung", "Complete a race where you had Monstro's Lung and Ipecac."},
		304: {"Sometimes You Should Take Tiny Planet", "Complete a race with Brimstone and Tiny Planet."},
		305: {"This Build Takes Skill", "Complete a race with Ipecac and Toxic Shock."},
		306: {"Propel Yourself Through the Door", "Complete a race with Epic Fetus and Holy Mantle."},
		307: {"Sometimes You Should Stay on a Tears Build", "Complete a race with Cricket's Body and The Parasite."},
		308: {"It's Technically a DPS Up", "Complete a race with Ipecac and Cricket's Body."},
		309: {"Like It Wasn't Powerful Enough Already", "Complete a race with Mega Blast and Car Battery."},
		355: {"7 Shots Is Better Than One", "Complete a race with Mutant Spider and The Inner Eye."},

		// Item synergies (3 items)
		351: {"Day of the Tentacles", "Complete a race with Monstro's Lung, Brimstone, and a homing item."},
		352: {"Ultimate Friends", "Complete a race with Lil' Brimstone, Incubus, and BFFS!"},
		353: {"It's Beautiful", "Complete a race with Epic Fetus, Brimstone, and Rubber Cement."},
		354: {"Shields Are Pretty Good", "Complete a race where you had Blood Rights, The Polaroid, and Scapular."},

		// Miscellaneous
		401: {"Last Man Standing", "Complete a race with at least 5 people where everyone else died or quit."},
		402: {"Marginal Expected Value", "Complete a race where you did not open any chests on The Chest."},
		403: {"Optimal Shoveling", "Complete a race where you immediately skipped floor 3, floor 5, and floor 7 using We Need to Go Deeper!"},
		404: {"Filthy Thief", "Complete a race where you \"stole\" an item from the Boss Rush."},
		405: {"Worth It", "Complete a race where you took Experimental Treatment as the third piece of the Spun transformation."},
		406: {"Clutch Leviathan", "Complete a race where you used the Leviathan transformation to take a devil deal that you otherwise wouldn't have been able to."},
		407: {"Winners Don't Use Drugs", "Finish 2nd place in a race where you used a Tears Down pill."},
		408: {"Pretty Basic", "Complete a race without taking an item that gives damage up."},
		409: {"Last Resort at Critical Health", "Complete a race after having used the D4 or the D100 at least once."},
		410: {"I Deserved This Win", "Finish 1st place in a race with at least 2 people after having procced a Guppy's Collar."},
		411: {"Never Guppy", "Finish two races in a row where you had the Guppy transformation."},
		412: {"U Can't Touch This", "Finish a race without taking damage."},
		413: {"Maybe I Shouldn't Have Min-Maxed So Hard", "Finish a race with 12 hearts."},
	}

	// Delete every row in the database
	if err := db.Achievements.DeleteAll(); err != nil {
		log.Fatal("Failed to delete all of the entries in the achievements table:", err)
	}

	// Add the achievement list to the database
	for id, achievement := range achievementList {
		if err := db.Achievements.Insert(id, achievement[0], achievement[1]); err != nil {
			log.Fatal("Failed to insert the achievements:", err)
		}
	}
}
