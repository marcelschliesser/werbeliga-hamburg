package types

import (
	"time"
)

type Season struct {
	Id        uint64     `json:"id"`
	Year      uint64     `json:"year"`
	MatchDays []MatchDay `json:"matchDays,omitempty"`
}

type MatchDay struct {
	Id           uint8         `json:"id"`
	Date         time.Time     `json:"date"`
	MatchResults []MatchResult `json:"matchResults,omitempty"`
}

type MatchResult struct {
	HomeTeam  string    `json:"homeTeam"`
	AwayTeam  string    `json:"awayTeam"`
	HomeScore uint64    `json:"homeScore"`
	AwayScore uint64    `json:"awayScore"`
	Time      time.Time `json:"time"`
	Court     uint64    `json:"court"`
}
