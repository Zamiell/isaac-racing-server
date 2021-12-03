package server

// Used for single-player speedruns, e.g. R+7 S1
type HallOfFameEntry struct {
	Rank        int
	Racer       string
	ProfileName string
	Time        int
	Version     string
	Date        string
	Proof       string
	Site        string
}

// Used for online seasons
type HallOfFameEntryOnline struct {
	Rank              int
	Racer             string
	ProfileName       string
	AdjustedAverage   int
	UnadjustedAverage int
	ForfeitPenalty    int
	NumForfeits       int
}
