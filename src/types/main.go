package types

import (
	"time"
)

type SeasonId uint
type MatchId uint

type Season struct {
	Id        SeasonId   `json:"seasonId"`
	Year      uint       `json:"year"`
	MatchDays []MatchDay `json:"matchDays"`
}

type MatchDay struct {
	Id           MatchId   `json:"matchId"`
	Date         time.Time `json:"matchDate"`
	MatchResults []Match   `json:"matchResults"`
}

type Match struct {
	HomeTeam  string    `json:"homeTeam"`
	AwayTeam  string    `json:"awayTeam"`
	HomeScore uint      `json:"homeScore"`
	AwayScore uint      `json:"awayScore"`
	DateTime  time.Time `json:"dateTime"`
	Court     uint      `json:"court"`
}
