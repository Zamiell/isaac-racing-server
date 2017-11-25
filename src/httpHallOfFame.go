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

	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        5039,
			Version:     "1.06.J99",
			Date:        "2017-07-13",
			Proof:       "https://www.twitch.tv/videos/158833908",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        2,
			Racer:       "Cyber_1",
			ProfileName: "Cyber_1",
			Time:        5462,
			Version:     "0.5.19",
			Date:        "2017-04-23",
			Proof:       "https://www.twitch.tv/videos/137587747",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        3,
			Racer:       "CrafterLynx",
			ProfileName: "CrafterLynx",
			Time:        5620, Version: "0.10.2",
			Date:  "2017-09-19",
			Proof: "https://www.twitch.tv/videos/175962871",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        4,
			Racer:       "Zamiel",
			ProfileName: "Zamiel",
			Time:        5629,
			Version:     "1.06.J39",
			Date:        "2017-04-12",
			Proof:       "https://www.twitch.tv/videos/135084542",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        5,
			Racer:       "ReidMercury__",
			ProfileName: "ReidMercury__",
			Time:        5655,
			Version:     "",
			Date:        "2017-09-17",
			Proof:       "https://www.twitch.tv/videos/175491970",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        6,
			Racer:       "ceehe",
			ProfileName: "ceehe",
			Time:        5776,
			Version:     "",
			Date:        "2017-10-18",
			Proof:       "https://www.twitch.tv/videos/180734439",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        7,
			Racer:       "leo_ze_tron",
			ProfileName: "leo_ze_tron",
			Time:        5925,
			Version:     "1.06.J104",
			Date:        "2017-08-17",
			Proof:       "https://www.twitch.tv/videos/167785993",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        8,
			Racer:       "SlashSP",
			ProfileName: "SlashSP",
			Time:        5938,
			Version:     "",
			Date:        "2017-09-22",
			Proof:       "https://www.twitch.tv/videos/176523741",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        9,
			Racer:       "Shigan",
			ProfileName: "Shigan",
			Time:        5984,
			Version:     "1.06.J75",
			Date:        "2017-05-14",
			Proof:       "https://www.twitch.tv/videos/143486007",
		},
	)
	season1r9 = append(
		season1r9,
		models.SpeedRun{
			Rank:        10,
			Racer:       "thereisnofuture",
			ProfileName: "thereisnofuture",
			Time:        5999,
			Version:     "1.06.J39",
			Date:        "2017-04-14",
			Proof:       "https://www.twitch.tv/videos/135612266",
		},
	)

	season1r14 = append(
		season1r14,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        8909,
			Version:     "1.06.J104",
			Date:        "2017-09-24",
			Proof:       "https://www.twitch.tv/videos/177220345",
		},
	)
	/*
		season1r14 = append(
			season1r14,
			models.SpeedRun{
				Rank:        "2",
				Racer:       "",
				ProfileName: "",
				Time:        3931,
				Version:     "1.06.J99",
				Date:        "2017--",
				Proof:       "",
			},
		)
	*/
	season2r7 = append(
		season2r7,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        3673,
			Version:     "1.06.J104",
			Date:        "2017-07-20",
			Proof:       "https://www.twitch.tv/videos/160709749",
		},
	)
	/*
		season2r7 = append(
			season2r7,
			models.SpeedRun{Rank: 2,
				Racer:       "",
				ProfileName: "",
				Time:        3931,
				Version:     "1.06.J99",
				Date:        "2017--",
				Proof:       "",
			},
		)
	*/
	season3r7 = append(
		season3r7,
		models.SpeedRun{
			Rank:        1,
			Racer:       "Dea1h",
			ProfileName: "Dea1h",
			Time:        3931,
			Version:     "1.06.J99",
			Date:        "2017-11-20T0700",
			Proof:       "https://www.twitch.tv/videos/158833908",
		},
	)
	/*
		season3r7 = append(
			season3r7,
			models.SpeedRun{
				Rank:        "2",
				Racer:       "",
				ProfileName: "",
				Time:        3931,
				Version:     "1.06.J99",
				Date:        "2017--",
				Proof:       "",
			},
		)
	*/

	data := TemplateData{
		Title:     "Hall Of Fame",
		Season1r9: season1r9,
	}
	httpServeTemplate(w, "halloffame", data)
}
