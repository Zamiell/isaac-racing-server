package models

// SpeedRun struct contains data about a single runners run
type SpeedRun struct {
	Rank        int
	Racer       string
	ProfileName string
	Time        int
	Version     string
	Date        string
	Proof       string
	Site        string
}
