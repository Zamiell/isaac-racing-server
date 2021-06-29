package main

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/models"
)

var (
	achievementMap map[int][]string
)

func achievementsInit() {
	achievementMap = map[int][]string{
		// Complete x races
		1: {"This Wasn't So Bad", "Complete your first race."},
		2: {"I Think I Have the Hang of This", "Complete 10 races."},
		3: {"Intermediate Racist", "Complete 50 races."},
		4: {"Expert Racist", "Complete 100 races."},
		5: {"Orange Juice", "Complete 500 races."}, // Reference to Ou_J

		// Complete x races with every ruleset
		6: {"Dipping Your Toe in the Water", "Complete a race with every ruleset (unseeded, seeded, and diversity)."},
		7: {"Experimental Treatment", "Complete 10 races with every ruleset (unseeded, seeded, and diversity)."},
		8: {"Jack of All Trades", "Complete 100 races with every ruleset (unseeded, seeded, and diversity)."},

		// Get x average
		11: {"Average Joe", "Have an unseeded average time of 16:30 or less (with at least 50 races played and under a 25% forfeit rate)."},
		12: {"There's So Many Strats and So Little Time", "Have an unseeded average time of 16:00 or less (with at least 50 races played and under a 25% forfeit rate)."},
		13: {"There's No Skill in an RNG Based Game", "Have an unseeded average time of 15:30 or less (with at least 50 races played and under a 25% forfeit rate)."},
		14: {"Master of Consistency", "Have an unseeded average time of 15:00 or less (with at least 50 races played and under a 25% forfeit rate)."},

		// Every starting item
		21: {"Well Rounded - Unseeded", "Complete an unseeded race with every starting item."},
		22: {"Well Rounded - Seeded", "Complete a seeded race with every starting item."},
		23: {"Well Rounded - Practice", "Complete a practice mode run with every starting item."},

		// Streaks
		31: {"I'm Having a Good Day", "Finish 1st place in 3 races in a row with 5 people or more."},
		32: {"A Literal God, Tasteless", "Finish 1st place in 5 races in a row with 5 people or more."}, // Reference to a Starcraft 2 meme
		33: {"I'm Having a Bad Day", "Quit out of 3 races in a row with 5 people or more."},

		// Complete a race with x time
		41: {"Pretty Darn Fast", "Complete an unseeded race in less than 13 minutes."},
		42: {"Speeding Bullet", "Complete an unseeded race in less than 12 minutes."},
		43: {"Fast Like Sanic", "Complete an unseeded race in less than 11 minutes."},
		44: {"Dea1hly Fast", "Complete an unseeded race in less than 10 minutes."}, // Reference to Dea1h
		45: {"Fizzy Giraffe", "Complete an unseeded race in less than 9 minutes."}, // Reference to giraffeFizzoid

		// Complete a race with x time with y item, unseeded
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

		// Complete a race with x time with y item, seeded
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
		301: {"My Bombs Don't Hurt Me", "Complete a race where you had Dr. Fetus and Ipecac."},
		302: {"Green Lung Best Lung", "Complete a race where you had Monstro's Lung and Ipecac."},
		303: {"Sometimes You Should Take Tiny Planet", "Complete a race with Tiny Planet and either Brimstone or Technology."},
		304: {"Sometimes You Should Take Dunce Cap", "Complete a race with Epic Fetus and Dunce Cap."},
		305: {"This Build Takes Skill", "Complete a race with Ipecac and Toxic Shock."},
		306: {"Missile Yourself Through the Door", "Complete a race with Epic Fetus and Holy Mantle."},
		307: {"Sometimes You Should Stay on a Tears Build", "Complete a race with Cricket's Body and The Parasite."},
		308: {"It's Technically a DPS Up", "Complete a race with Ipecac and Cricket's Body."},
		309: {"Like It Wasn't Powerful Enough Already", "Complete a race with Mega Blast and Car Battery."},
		310: {"7 Shots Is Better Than One", "Complete a race with Mutant Spider and The Inner Eye."},
		311: {"Massive Multiplier", "Complete a race with Technology and A Lump of Coal."},

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
		405: {"Decisions, Decisions", "Complete a race where you \"stole\" an item from the Bosh Rush that had 2 or more \"starting\" items."},
		406: {"Worth It", "Complete a race where you took Experimental Treatment as the third piece of the Spun transformation."},
		407: {"Clutch Leviathan", "Complete a race where you used the Leviathan transformation to take a devil deal that you otherwise wouldn't have been able to."},
		408: {"Winners Don't Use Drugs", "Finish 2nd place in a race where you used a Tears Down pill."},
		410: {"Pretty Basic", "Complete a race without taking an item that gives damage up."},
		411: {"Last Resort at Critical Health", "Complete a race after having used the D4 or the D100 at least once."},
		412: {"Undeserved Win", "Finish 1st place in a race with at least 2 people after having activated a Guppy's Collar or a Broken Ankh."},
		413: {"Never Guppy", "Finish two races in a row where you had the Guppy transformation."},
		414: {"U Can't Touch This", "Finish a race without taking damage."},
		415: {"Maybe I Shouldn't Have Min-Maxed So Hard", "Finish a race with 12 hearts."},
		416: {"Better Late Than Never", "Complete a race where you found and used an Emperor card on The Chest."},
		417: {"Collaborative Victory", "Finish a race where you tied for 1st place."},
		418: {"Consolation Prize", "Complete a race where you didn't take any red heart damage and only received a Devil/Angel room on floor 2, floor 5, and floor 8."},
		419: {"Confident in Your Dodging Ability", "Complete a race where you took two devil deals that left you with only one heart remaining."},
		420: {"Not Very \"Difficult\"", "Complete a race where you spawned as Judas' Shadow with more than 2 hearts."}, // Reference to EladDifficult
		421: {"Bomb of Kings", "Complete a race where you killed Pin with 1 bomb and no other forms of damage."},       // Reference to the Battle of Kings showmatch series
		422: {"Curse of the Full Clear", "Complete a race where you entered every room of The Chest."},
	}

	// Delete every row in the database
	if err := db.Achievements.DeleteAll(); err != nil {
		logger.Fatal("Failed to delete all of the entries in the achievements table:", err)
	}

	// Add the achievement list to the database
	for id, achievement := range achievementMap {
		if err := db.Achievements.Insert(id, achievement[0], achievement[1]); err != nil {
			logger.Fatal("Failed to insert the achievements:", err)
		}
	}
	logger.Info("Added", len(achievementMap), "achievements.")
}

func achievementsCheck(userID int, username string) {
	// Get this racer's current achievements
	userAchievements, err := db.UserAchievements.GetAll(userID)
	if err != nil {
		logger.Error("Database error while getting the achievements for user "+strconv.Itoa(userID)+":", err)
		return
	}

	achievements1_8(userID, username, userAchievements)
	achievements11_14(userID, username, userAchievements)
}

// Achievement 1-8 (complete x races)
func achievements1_8(userID int, username string, userAchievements []int) {
	if intInSlice(1, userAchievements) &&
		intInSlice(2, userAchievements) &&
		intInSlice(3, userAchievements) &&
		intInSlice(4, userAchievements) &&
		intInSlice(5, userAchievements) &&
		intInSlice(6, userAchievements) &&
		intInSlice(7, userAchievements) &&
		intInSlice(8, userAchievements) {

		return
	}

	var finishedList []models.Race
	if v, err := db.RaceParticipants.GetFinishedRaces(userID); err != nil {
		logger.Error("Database error while getting the finished races for user "+strconv.Itoa(userID)+":", err)
		return
	} else {
		finishedList = v
	}

	// Achievement 1 - This Wasn't So Bad - Complete your first race.
	if !intInSlice(1, userAchievements) {
		if len(finishedList) >= 1 {
			achievementsGive(userID, username, 1)
		}
	}

	// Achievement 2 - I Think I Have the Hang of This - Complete 10 races.
	if !intInSlice(2, userAchievements) {
		if len(finishedList) >= 10 {
			achievementsGive(userID, username, 2)
		}
	}

	// Achievement 3 - Intermediate Racist - Complete 50 races.
	if !intInSlice(3, userAchievements) {
		if len(finishedList) >= 50 {
			achievementsGive(userID, username, 3)
		}
	}

	// Achievement 4 - Expert Racist - Complete 100 races.
	if !intInSlice(4, userAchievements) {
		if len(finishedList) >= 100 {
			achievementsGive(userID, username, 4)
		}
	}

	// Achievement 5 - Orange Juice - Complete 500 races.
	if !intInSlice(5, userAchievements) {
		if len(finishedList) >= 500 {
			achievementsGive(userID, username, 5)
		}
	}

	// Achievement 6-8 (complete x races with every ruleset)
	if !intInSlice(6, userAchievements) ||
		!intInSlice(7, userAchievements) ||
		!intInSlice(8, userAchievements) {

		// Count the number of races in each ruleset
		countUnseeded := 0
		countSeeded := 0
		countDiversity := 0
		for _, race := range finishedList {
			if race.Format == RaceFormatUnseeded {
				countUnseeded++
			} else if race.Format == RaceFormatSeeded {
				countSeeded++
			} else if race.Format == RaceFormatDiversity {
				countDiversity++
			}
		}

		// Achievement 6 - Dipping Your Toe in the Water - Complete a race with every ruleset (unseeded, seeded, and diversity).
		if !intInSlice(6, userAchievements) {
			if countUnseeded >= 1 && countSeeded >= 1 && countDiversity >= 1 {
				achievementsGive(userID, username, 6)
			}
		}

		// Achievement 7 - Experimental Treatment - Complete 10 races with every ruleset (unseeded, seeded, and diversity).
		if !intInSlice(7, userAchievements) {
			if countUnseeded >= 10 && countSeeded >= 10 && countDiversity >= 10 {
				achievementsGive(userID, username, 7)
			}
		}

		// Achievement 8 - Jack of All Trades - Complete 100 races with every ruleset (unseeded, seeded, and diversity).
		if !intInSlice(8, userAchievements) {
			if countUnseeded >= 100 && countSeeded >= 100 && countDiversity >= 100 {
				achievementsGive(userID, username, 8)
			}
		}
	}
}

// Achievement 11-14 (get x average)
func achievements11_14(userID int, username string, userAchievements []int) {
	if intInSlice(11, userAchievements) &&
		intInSlice(12, userAchievements) &&
		intInSlice(13, userAchievements) &&
		intInSlice(14, userAchievements) {

		return
	}

	// Get their average
	// TODO
}

func achievementsGive(userID int, username string, achievementID int) {
	// Give them the achievement in the database
	if err := db.UserAchievements.Insert(userID, achievementID); err != nil {
		logger.Error("Failed to give achievement #"+strconv.Itoa(achievementID)+" to user \""+username+"\":", err)
		return
	}

	// Send them a notification that they got this achievement
	// (they should not be offline, but check just in case they went offline immediately after finishing)
	s, ok := websocketSessions[username]
	if ok {
		type AchievementMessage struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		websocketEmit(s, "achievement", &AchievementMessage{
			achievementID,
			achievementMap[achievementID][0],
			achievementMap[achievementID][1],
		})
	}
}
