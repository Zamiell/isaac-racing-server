package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
)

// For holding the values of the "items.json" file
type JSONItem struct {
	Name        string `json:"name"`
	Damage      string `json:"dmg"`
	DamageX     string `json:"dmg_x"` // nolint:tagliatelle
	Health      string `json:"health"`
	SoulHearts  string `json:"soul_hearts"` // nolint:tagliatelle
	SinHearts   string `json:"sin_hearts"`  // nolint:tagliatelle
	Tears       string `json:"tears"`
	Delay       string `json:"delay"`
	DelayX      string `json:"delay_x"` // nolint:tagliatelle
	Speed       string `json:"speed"`
	ShotSpeed   string `json:"shot_speed"` // nolint:tagliatelle
	Height      string `json:"height"`
	Range       string `json:"range"`
	Luck        string `json:"luck"`
	Beelzebub   bool   `json:"beelzebub"`
	Bob         bool   `json:"bob"`
	Bookworm    bool   `json:"bookworm"`
	Conjoined   bool   `json:"conjoined"`
	Funguy      bool   `json:"funguy"`
	Guppy       bool   `json:"guppy"`
	Leviathan   bool   `json:"leviathan"`
	OhCrap      bool   `json:"ohcrap"`
	Seraphim    bool   `json:"seraphim"`
	SpiderBaby  bool   `json:"spiderbaby"`
	Spun        bool   `json:"spun"`
	Superbum    bool   `json:"superbum"`
	YesMother   bool   `json:"yesmother"`
	SpaceBar    bool   `json:"space"`
	HealthOnly  bool   `json:"health_only"`   // nolint:tagliatelle
	Intro       string `json:"introduced_in"` // nolint:tagliatelle
	Shown       bool   `json:"shown"`
	Summary     bool   `json:"in_summary"`     // nolint:tagliatelle
	SummaryName string `json:"summary_name"`   // nolint:tagliatelle
	SummaryCond string `json:"condition_name"` // nolint:tagliatelle
	Text        string `json:"text"`
}

// For holding the values of the "builds.json" file
type IsaacItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TournamentInfo struct {
	Name         string `json:"name"`
	ChallongeID  string `json:"challonge_id"` // nolint:tagliatelle
	ChallongeURL string `json:"challonge"`
	Date         string `json:"date"`
	Notability   string `json:"notability"`
	Organizer    string `json:"organizer"`
	Ruleset      string `json:"ruleset"`
	Description  string `json:"description"`
}

type NameSorter []os.FileInfo

func (a NameSorter) Len() int           { return len(a) }
func (a NameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool { return a[i].Name() > a[j].Name() }

const (
	projectName = "isaac-racing-server"
)

var (
	allItems       = make(map[string]*JSONItem)
	allItemNames   = make(map[int]string)
	allBuilds      = make([][]IsaacItem, 0)
	allTournaments = make([]TournamentInfo, 0)
)

func loadAllItems() {
	jsonFilePath := path.Join(projectPath, "public", "items.json")
	jsonFile, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		logger.Fatal("Failed to open \""+jsonFilePath+"\":", err)
	}

	// Create all the items
	if err := json.Unmarshal(jsonFile, &allItems); err != nil {
		logger.Fatal("Failed to unmarshal the items:", err)
		return
	}

	// Create a 2nd map of just item names
	for k, v := range allItems {
		itemID, _ := strconv.Atoi(k)
		allItemNames[itemID] = v.Name
	}
}

func loadAllBuilds() {
	// Open the JSON file and verify it was good
	jsonFilePath := path.Join(projectPath, "public", "builds.json")
	var jsonFile []byte
	if v, err := ioutil.ReadFile(jsonFilePath); err != nil {
		logger.Fatal("Failed to open \""+jsonFilePath+"\":", err)
	} else {
		jsonFile = v
	}

	// Create all the items
	if err := json.Unmarshal(jsonFile, &allBuilds); err != nil {
		logger.Fatal("Failed to unmarshal the builds:", err)
	}
}

func loadAllTournaments() {
	// Temporary var for each tournament
	var tournament TournamentInfo
	// Open the JSON files for tournaments and load them into TournamentInfo
	jsonFolderPath := path.Join(projectPath, "BoIR-trueskill/tournaments")
	fileList, err := ioutil.ReadDir(jsonFolderPath)
	if err != nil {
		logger.Error("Could not read the files in ", jsonFolderPath)
	}
	sort.Sort(NameSorter(fileList))
	for _, file := range fileList {
		// Create the full file path
		filePath := path.Join(jsonFolderPath, file.Name())
		var jsonFile []byte
		if v, err := ioutil.ReadFile(filePath); err != nil {
			// Fatal error if we cannot open a file
			logger.Fatal("Failed to open \""+filePath+"\":", err)
		} else {
			jsonFile = v
		}

		// Create all the tournament vars
		if err := json.Unmarshal(jsonFile, &tournament); err != nil {
			logger.Fatal("Failed to unmarshal the tournament:", err)
		}
		allTournaments = append(allTournaments, tournament)
	}
}

func getBuildName(startingBuildIndex int) string {
	startingBuild := allBuilds[startingBuildIndex]
	if len(startingBuild) == 1 {
		return startingBuild[0].Name
	}

	if len(startingBuild) == 2 {
		return startingBuild[0].Name + " + " + startingBuild[1].Name
	}

	return startingBuild[0].Name + " + more"
}

func getBuildID(startingBuildIndex int) int {
	startingBuild := allBuilds[startingBuildIndex]

	return startingBuild[0].ID
}

const RaceStatusOpen = "open"
const RaceStatusInProgress = "in progress"
const RaceFormatUnseeded = "unseeded"
const RaceFormatCustom = "custom"
const RaceFormatSeeded = "seeded"
const RaceFormatDiversity = "diversity"
const RaceGoalBeast = "The Beast"
const RaceGoalCustom = "custom"
const RacerStatusReady = "ready"
const RacerStatusRacing = "racing"
