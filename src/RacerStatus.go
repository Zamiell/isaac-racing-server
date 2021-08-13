package server

type RacerStatus string

const (
	RacerStatusReady        RacerStatus = "ready"
	RacerStatusRacing       RacerStatus = "racing"
	RacerStatusQuit         RacerStatus = "quit"
	RacerStatusDisqualified RacerStatus = "disqualified"
)
