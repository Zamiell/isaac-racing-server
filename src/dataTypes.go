package main

/*
	Data types that are shared between different elements of the program
*/

// Used to track the current races in memory
type Race struct {
	ID              int
	Name            string
	Status          string /* open, starting, in progress, finished */
	Ruleset         Ruleset
	Captain         string
	SoundPlayed     bool
	DatetimeCreated int64
	DatetimeStarted int64
	Racers          map[string]*Racer
}
type Ruleset struct {
	Ranked        bool   `json:"ranked"`
	Solo          bool   `json:"solo"`
	Format        string `json:"format"`
	Character     string `json:"character"`
	Goal          string `json:"goal"`
	StartingBuild int    `json:"startingBuild"`
	Seed          string `json:"seed"`
}
type Racer struct {
	ID                   int
	Name                 string
	DatetimeJoined       int64
	Status               string
	Seed                 string
	FloorNum             int
	StageType            int
	DatetimeArrivedFloor int64
	Items                []*Item
	StartingItem         int /* Determined by seeing if room count is > 0 */
	Rooms                []*Room
	Place                int
	PlaceMid             int
	DatetimeFinished     int64
	RunTime              int64 /* in milliseconds */
	Comment              string
}
type Item struct {
	ID        int `json:"id"`
	FloorNum  int `json:"floorNum"`
	StageType int `json:"stageType"`
}
type Room struct {
	ID        string /* e.g. "5.999" */
	FloorNum  int
	StageType int
}
