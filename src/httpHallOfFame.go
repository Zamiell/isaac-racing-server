package main

import (
	"time"
	//	"math"
	//	"net/http"
	//	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
)

func httpHallOfFame(c *gin.Context) {
	// Local variables
	w := c.Writer

	type speedRun struct {
		Rank        int
		Racer       string
		ProfileName string
		Time        int
		CoOpBaby    bool
		Version     float32
		Date        time
		Proof       string
	}
	var season1r9 []speedRun
	var season1r14 []speedRun
	var season2r7 []speedRun
	var season3r7 []speedRun

	season1r9 = append(season1r9, speedRun{Rank: "1", Racer: "Dea1h", ProfileName: "Dea1h", Time: 5039, CoOpBaby: false, Version: "1.06.J99", Date: "2017-07-13T0800", Proof: "https://www.twitch.tv/videos/158833908"})
	season1r9 = append(season1r9, speedRun{Rank: "2", Racer: "Cyber_1", ProfileName: "Cyber_1", Time: 5462, CoOpBaby: false, Version: "0.5.19", Date: "2017-04-23T0800", Proof: "https://www.twitch.tv/videos/137587747"})
	season1r9 = append(season1r9, speedRun{Rank: "3", Racer: "CrafterLynx", ProfileName: "CrafterLynx", Time: 5620, CoOpBaby: false, Version: "0.10.2", Date: "2017-09-19T0800", Proof: "https://www.twitch.tv/videos/175962871"})
	season1r9 = append(season1r9, speedRun{Rank: "4", Racer: "Zamiel", ProfileName: "Zamiel", Time: 5629, CoOpBaby: false, Version: "1.06.J39", Date: "2017-04-12T0800", Proof: "https://www.twitch.tv/videos/135084542"})
	season1r9 = append(season1r9, speedRun{Rank: "5", Racer: "ReidMercury__", ProfileName: "ReidMercury__", Time: 5655, CoOpBaby: false, Version: "", Date: "2017-09-17T0800", Proof: "https://www.twitch.tv/videos/175491970"})
	season1r9 = append(season1r9, speedRun{Rank: "6", Racer: "ceehe", ProfileName: "ceehe", Time: 5776, CoOpBaby: false, Version: "", Date: "2017-10-18T0800", Proof: "https://www.twitch.tv/videos/180734439"})
	season1r9 = append(season1r9, speedRun{Rank: "7", Racer: "leo_ze_tron", ProfileName: "leo_ze_tron", Time: 5925, CoOpBaby: false, Version: "1.06.J104", Date: "2017-08-17T0800", Proof: "https://www.twitch.tv/videos/167785993"})
	season1r9 = append(season1r9, speedRun{Rank: "8", Racer: "SlashSP", ProfileName: "SlashSP", Time: 5938, CoOpBaby: false, Version: "", Date: "2017-09-22T0800", Proof: "https://www.twitch.tv/videos/176523741"})
	season1r9 = append(season1r9, speedRun{Rank: "9", Racer: "Shigan", ProfileName: "Shigan", Time: 5984, CoOpBaby: false, Version: "1.06.J75", Date: "2017-05-14T0800", Proof: "https://www.twitch.tv/videos/143486007"})
	season1r9 = append(season1r9, speedRun{Rank: "10", Racer: "thereisnofuture", ProfileName: "thereisnofuture", Time: 5999, CoOpBaby: false, Version: "1.06.J39", Date: "2017-04-14T0800", Proof: "https://www.twitch.tv/videos/135612266"})

	//season1r14 = append(season1r14, speedRun{Rank: "1", Racer: "Dea1h", ProfileName: "Dea1h", Time: 8909, CoOpBaby: false, Version: "1.06.J104", Date: "2017-09-24T0800", Proof: "https://www.twitch.tv/videos/177220345"})
	//season1r14 = append(season1r14, speedRun{Rank: "2", Racer: "", ProfileName: "", Time: 3931, CoOpBaby: false, Version: "1.06.J99", Date: "2017--T0800", Proof: ""})

	//season2r7 = append(season2r7, speedRun{Rank: "1", Racer: "Dea1h", ProfileName: "Dea1h", Time: 3673, CoOpBaby: false, Version: "1.06.J104", Date: "2017-07-20T0800", Proof: "https://www.twitch.tv/videos/160709749"})
	//season2r7 = append(season2r7, speedRun{Rank: "2", Racer: "", ProfileName: "", Time: 3931, CoOpBaby: false, Version: "1.06.J99", Date: "2017--T0800", Proof: ""})

	//season3r7 = append(season3r7, speedRun{Rank: "1", Racer: "Dea1h", ProfileName: "Dea1h", Time: 3931, CoOpBaby: false, Version: "1.06.J99", Date: "2017-11-20T0700", Proof: "https://www.twitch.tv/videos/158833908"})
	//season3r7 = append(season3r7, speedRun{Rank: "2", Racer: "", ProfileName: "", Time: 3931, CoOpBaby: false, Version: "1.06.J99", Date: "2017--T0800", Proof: ""})

	var allSeasonsResults [][]speedRun
	allSeasonsResults := make([][]season1r19, 2)

	data := TemplateData{
		Title:   "Hall Of Fame",
		HofData: allSeasonsResults,
	}
	httpServeTemplate(w, "halloffame", data)
}
