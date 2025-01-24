package types

import (
	"time"
)

type Season struct {
	Id        uint       `json:"id"`
	Year      uint       `json:"year"`
	MatchDays []MatchDay `json:"matchDays,omitempty"`
}

type MatchDay struct {
	Id           uint          `json:"id"`
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
