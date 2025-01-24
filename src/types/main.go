package types

import (
	"time"
)

type SeasonId uint
type MatchId uint

type Season struct {
	Id        SeasonId   `json:"id"`
	Year      uint       `json:"year"`
	MatchDays []MatchDay `json:"matchDays,omitempty"`
}

type MatchDay struct {
	Id           MatchId       `json:"id"`
	Date         time.Time     `json:"date"`
	MatchResults []MatchResult `json:"matchResults,omitempty"`
}

type MatchResult struct {
	HomeTeam  string    `json:"homeTeam"`
	AwayTeam  string    `json:"awayTeam"`
	HomeScore uint      `json:"homeScore"`
	AwayScore uint      `json:"awayScore"`
	Time      time.Time `json:"time"`
	Court     uint      `json:"court"`
}
