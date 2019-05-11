package main

import (
	//	"math"
	//	"net/http"
	//	"strings"

	//	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	"github.com/gin-gonic/gin"
)

func httpHallOfFame(c *gin.Context) {
	// Local variables
	w := c.Writer

	var season1r9 []models.SpeedRun
	var season1r14 []models.SpeedRun
	var season2r7 []models.SpeedRun
	var season3r7 []models.SpeedRun
	var season4r7 []models.SpeedRun
	var season5r7 []models.SpeedRun
	//var season6r7 []models.SpeedRun

	// Season 1 R+9 Start
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        5039,
			Date:        "2017-07-13",
			Proof:       "https://www.twitch.tv/videos/158833908",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        2,
			Racer:       "Cyber_1",
			ProfileName: "Cyber_1",
			Time:        5462,
			Date:        "2017-04-23",
			Proof:       "https://www.twitch.tv/videos/137587747",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        3,
			Racer:       "CrafterLynx",
			ProfileName: "CrafterLynx",
			Time:        5620,
			Date:        "2017-09-19",
			Proof:       "https://www.twitch.tv/videos/175962871",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        4,
			Racer:       "Zamiel",
			ProfileName: "Zamiel",
			Time:        5629,
			Date:        "2017-04-12",
			Proof:       "https://www.twitch.tv/videos/135084542",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        5,
			Racer:       "ReidMercury__",
			ProfileName: "ReidMercury__",
			Time:        5655,
			Date:        "2017-09-17",
			Proof:       "https://www.twitch.tv/videos/175491970",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        6,
			Racer:       "ceehe",
			ProfileName: "ceehe",
			Time:        5776,
			Date:        "2017-10-18",
			Proof:       "https://www.twitch.tv/videos/180734439",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        7,
			Racer:       "leo_ze_tron",
			ProfileName: "leo_ze_tron",
			Time:        5925,
			Date:        "2017-08-17",
			Proof:       "https://www.twitch.tv/videos/167785993",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        8,
			Racer:       "SlashSP",
			ProfileName: "SlashSP",
			Time:        5938,
			Date:        "2017-09-22",
			Proof:       "https://www.twitch.tv/videos/176523741",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        9,
			Racer:       "Shigan",
			ProfileName: "Shigan",
			Time:        5984,
			Date:        "2017-05-14",
			Proof:       "https://www.twitch.tv/videos/143486007",
			Site:        "Twitch",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        10,
			Racer:       "thereisnofuture",
			ProfileName: "thereisnofuture",
			Time:        5999,
			Date:        "2017-04-14",
			Proof:       "https://www.twitch.tv/videos/135612266",
			Site:        "Twitch",
		},
	)

	// Season 1 R+14 Start
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        8909,
			Date:        "2017-09-24",
			Proof:       "https://www.twitch.tv/videos/177220345",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        2,
			Racer:       "CrafterLynx",
			ProfileName: "CrafterLynx",
			Time:        9061,
			Date:        "2017-09-09",
			Proof:       "https://www.twitch.tv/videos/175961290",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        3,
			Racer:       "ceehe",
			ProfileName: "ceehe",
			Time:        9782,
			Date:        "2017-10-08",
			Proof:       "https://www.twitch.tv/videos/213334990",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        4,
			Racer:       "yamayamadingdong",
			ProfileName: "yama",
			Time:        9837,
			Date:        "2017-09-15",
			Proof:       "https://www.twitch.tv/videos/174831927",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        5,
			Racer:       "SlashSP",
			ProfileName: "SlashSP",
			Time:        9960,
			Date:        "2017-09-22",
			Proof:       "https://www.twitch.tv/videos/176523011",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        6,
			Racer:       "Shigan",
			ProfileName: "Shigan",
			Time:        10188,
			Date:        "2017-05-14",
			Proof:       "https://www.twitch.tv/videos/143486007",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        7,
			Racer:       "ReidMercury__",
			ProfileName: "ReidMercury__",
			Time:        10431,
			Date:        "2017-08-18",
			Proof:       "https://www.twitch.tv/videos/167959080",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        8,
			Racer:       "Zamiel",
			ProfileName: "Zamiel",
			Time:        10616,
			Date:        "2017-04-12",
			Proof:       "https://www.twitch.tv/videos/135084542",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        9,
			Racer:       "MrPopoche1",
			ProfileName: "",
			Time:        10665,
			Date:        "2017-06-06",
			Proof:       "https://www.twitch.tv/videos/149949436",
			Site:        "Twitch",
		},
	)
	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        10,
			Racer:       "SergeBenamou31",
			ProfileName: "",
			Time:        11364,
			Date:        "2017-08-31",
			Proof:       "https://www.twitch.tv/videos/171259098",
			Site:        "Twitch",
		},
	)

	// Season 2 R+7 Start
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        3673,
			Date:        "2017-07-20",
			Proof:       "https://www.twitch.tv/videos/160709749",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        2,
			Racer:       "Shigan",
			ProfileName: "Shigan",
			Time:        4033,
			Date:        "2017-07-19",
			Proof:       "https://www.twitch.tv/videos/160356773",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        3,
			Racer:       "ceehe",
			ProfileName: "ceehe",
			Time:        4158,
			Date:        "2017-08-10",
			Proof:       "https://www.twitch.tv/videos/165878075",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        4,
			Racer:       "Zamiel",
			ProfileName: "Zamiel",
			Time:        4476,
			Date:        "2017-06-07",
			Proof:       "https://www.twitch.tv/videos/150106693",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        5,
			Racer:       "CrafterLynx",
			ProfileName: "CrafterLynx",
			Time:        4479,
			Date:        "2017-06-06",
			Proof:       "https://www.twitch.tv/videos/148982150",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        6,
			Racer:       "SlashSP",
			ProfileName: "SlashSP",
			Time:        4548,
			Date:        "2017-07-24",
			Proof:       "https://www.twitch.tv/videos/161482668",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        7,
			Racer:       "thereisnofuture",
			ProfileName: "thereisnofuture",
			Time:        4593,
			Date:        "2017-09-26",
			Proof:       "https://www.twitch.tv/videos/177693217",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        8,
			Racer:       "AdRyDN",
			ProfileName: "AdRyDN",
			Time:        4619,
			Date:        "2017-06-14",
			Proof:       "https://www.twitch.tv/videos/149545223",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        9,
			Racer:       "ReidMercury__",
			ProfileName: "ReidMercury__",
			Time:        4625,
			Date:        "2017-11-10",
			Proof:       "https://www.twitch.tv/videos/200127524",
			Site:        "Twitch",
		},
	)
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        10,
			Racer:       "thalen22",
			ProfileName: "thalen22",
			Time:        4628,
			Date:        "2017-07-31",
			Proof:       "https://www.twitch.tv/videos/163297005",
			Site:        "Twitch",
		},
	)

	// Season 3 R+7 Start
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        3797,
			Date:        "2017-12-04",
			Proof:       "https://www.twitch.tv/videos/206706029",
			Site:        "Twitch",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        2,
			Racer:       "ReidMercury__",
			ProfileName: "ReidMercury__",
			Time:        3844,
			Date:        "2017-11-28",
			Proof:       "https://www.twitch.tv/videos/205034914",
			Site:        "Twitch",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        3,
			Racer:       "bmzloop",
			ProfileName: "Loopy",
			Time:        4091,
			Date:        "2017-12-15",
			Proof:       "https://www.youtube.com/watch?v=g02cqQsrH1M",
			Site:        "YouTube",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        4,
			Racer:       "Shigan",
			ProfileName: "Shigan",
			Time:        4130,
			Date:        "2018-01-19",
			Proof:       "https://www.twitch.tv/videos/220225421",
			Site:        "Twitch",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        5,
			Racer:       "CrafterLynx",
			ProfileName: "CrafterLynx",
			Time:        4133,
			Date:        "2018-01-16",
			Proof:       "https://www.twitch.tv/videos/219308798",
			Site:        "Twitch",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        6,
			Racer:       "SapphireHX",
			ProfileName: "Sapphire",
			Time:        4249,
			Date:        "2018-02-06",
			Proof:       "https://www.youtube.com/watch?v=Zw2Ot5hjgZQ",
			Site:        "YouTube",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        7,
			Racer:       "Moucheron_Quipet",
			ProfileName: "MoucheronQuipet",
			Time:        4269,
			Date:        "2018-03-14",
			Proof:       "https://www.youtube.com/watch?v=gMT-caoKJE0",
			Site:        "YouTube",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        8,
			Racer:       "leo_ze_tron",
			ProfileName: "leo_ze_tron",
			Time:        4290,
			Date:        "2018-03-08",
			Proof:       "https://www.twitch.tv/videos/236448621",
			Site:        "Twitch",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        9,
			Racer:       "Zamiel",
			ProfileName: "Zamiel",
			Time:        4362,
			Date:        "2017-12-13",
			Proof:       "https://www.twitch.tv/videos/209120653",
			Site:        "Twitch",
		},
	)
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        10,
			Racer:       "ceehe",
			ProfileName: "ceehe",
			Time:        4448,
			Date:        "2017-12-13",
			Proof:       "https://www.twitch.tv/videos/209219042",
			Site:        "Twitch",
		},
	)

	// Season 4 R+7 Start
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Antizoubilamaka",
			ProfileName: "CRAZYEIGHTSFAN69",
			Time:        3406,
			Date:        "2018-09-07",
			Proof:       "https://www.twitch.tv/videos/307245073",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        2,
			Racer:       "Cyber_1",
			ProfileName: "Cyber_1",
			Time:        3475,
			Date:        "2018-05-11",
			Proof:       "https://www.twitch.tv/videos/260283073",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        3,
			Racer:       "leo_ze_tron",
			ProfileName: "leo_ze_tron",
			Time:        3522,
			Date:        "2018-07-11",
			Proof:       "https://www.twitch.tv/videos/283570282",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        4,
			Racer:       "thereisnofuture",
			ProfileName: "thereisnofuture",
			Time:        3543,
			Date:        "2018-08-31",
			Proof:       "https://www.twitch.tv/videos/304179755",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        5,
			Racer:       "Shigan",
			ProfileName: "Shigan",
			Time:        3617,
			Date:        "2018-10-29",
			Proof:       "https://www.twitch.tv/videos/328752825",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        6,
			Racer:       "ReidMercury__",
			ProfileName: "ReidMercury__",
			Time:        3622,
			Date:        "2018-05-07",
			Proof:       "https://www.twitch.tv/videos/258743172",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        7,
			Racer:       "ceehe",
			ProfileName: "ceehe",
			Time:        3727,
			Date:        "2018-04-22",
			Proof:       "https://www.twitch.tv/videos/252823258",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        8,
			Racer:       "CrafterLynx",
			ProfileName: "CrafterLynx",
			Time:        3776,
			Date:        "2018-09-29",
			Proof:       "https://www.twitch.tv/videos/316373160",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        9,
			Racer:       "bmzloop",
			ProfileName: "Loopy",
			Time:        3788,
			Date:        "2018-04-19",
			Proof:       "https://www.twitch.tv/videos/251906525",
			Site:        "Twitch",
		},
	)
	season4r7 = append(
		season4r7,
		models.SpeedRun{
			Rank:        10,
			Racer:       "thisguyisbarry",
			ProfileName: "thisguyisbarry",
			Time:        3814,
			Date:        "2018-05-22",
			Proof:       "https://www.twitch.tv/videos/264286364",
			Site:        "Twitch",
		},
	)

	// Season 5 R+7 Start
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        3773,
			Date:        "2018-12-11",
			Proof:       "https://www.twitch.tv/videos/347393032",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        2,
			Racer:       "Cyber_1",
			ProfileName: "Cyber_1",
			Time:        3806,
			Date:        "2019-01-11",
			Proof:       "https://www.twitch.tv/videos/365298470",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        3,
			Racer:       "CrafterLynx",
			ProfileName: "CrafterLynx",
			Time:        3857,
			Date:        "2019-01-23",
			Proof:       "https://www.twitch.tv/videos/368584877",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        4,
			Racer:       "thereisnofuture",
			ProfileName: "thereisnofuture",
			Time:        3871,
			Date:        "2018-12-05",
			Proof:       "https://www.twitch.tv/thereisnofuture/v/344976260",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        5,
			Racer:       "mgln",
			ProfileName: "mgln",
			Time:        3929,
			Date:        "2019-04-21",
			Proof:       "https://www.youtube.com/watch?v=LUjetqZ7O9I",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        6,
			Racer:       "sisuka",
			ProfileName: "sisuka",
			Time:        4172,
			Date:        "2018-12-23",
			Proof:       "https://www.youtube.com/watch?v=Rth9ITDWZ4o",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        7,
			Racer:       "Shigan",
			ProfileName: "Shigan",
			Time:        4228,
			Date:        "2019-03-12",
			Proof:       "https://www.twitch.tv/videos/394474921",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        8,
			Racer:       "Zamiel",
			ProfileName: "Zamiel",
			Time:        4231,
			Date:        "2018-11-22",
			Proof:       "https://www.twitch.tv/videos/339146750",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        9,
			Racer:       "Gamonymous",
			ProfileName: "Gamonymous",
			Time:        4293,
			Date:        "2019-01-23",
			Proof:       "https://www.twitch.tv/gamonymous__/v/368290597",
			Site:        "Twitch",
		},
	)
	season5r7 = append(
		season5r7,
		models.SpeedRun{
			Rank:        10,
			Racer:       "SlashSP",
			ProfileName: "SlashSP",
			Time:        4305,
			Date:        "2018-12-15",
			Proof:       "https://www.twitch.tv/videos/349688950",
			Site:        "Twitch",
		},
	)

	// Season 6 R+7 Start
	/*
		season6r7 = append(
			season6r7,
			models.SpeedRun{
				Rank:        2,
				Racer:       "",
				ProfileName: "",
				Time:        3844,
				Date:        "",
				Proof:       "",
				Site:        "Twitch",
			},
		)
	*/

	data := TemplateData{
		Title:      "Hall Of Fame",
		Season1r9:  season1r9,
		Season1r14: season1r14,
		Season2r7:  season2r7,
		Season3r7:  season3r7,
		Season4r7:  season4r7,
		Season5r7:  season5r7,
		//Season6r7:	season6r7,
	}
	httpServeTemplate(w, "halloffame", data)
}
