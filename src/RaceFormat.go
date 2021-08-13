package server

type RaceFormat string

const (
	RaceFormatUnseeded  RaceFormat = "unseeded"
	RaceFormatSeeded    RaceFormat = "seeded"
	RaceFormatDiversity RaceFormat = "diversity"
	RaceFormatCustom    RaceFormat = "custom"
)
