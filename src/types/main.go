package types

import (
	"time"
)

type MatchResult struct {
	HomeTeam  string
	AwayTeam  string
	HomeScore int8
	AwayScore int8
	DateTime  time.Time
	MatchDay  int8
}

type Season struct {
	Name     string
	Id       string
	Selected bool
}

type Match struct {
	Name     string
	Id       string
	Selected bool
}
