package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcelschliesser/werbeliga-hamburg/types"
)

const baseUrl string = "https://werbeliga.de/de/Spielplan,%20Tabelle%20&%20Torsch%C3%BCtzen"

type Crawler struct {
	httpClient http.Client
	baseUrl    string
}

func main() {
	c := NewCrawler(baseUrl, 10)

	// Starting Point (Season 2011)
	doc := c.FetchUrl(2, 1)

	seasons := ReturnSeasons(doc)

	for _, season := range seasons {
		matchDays := c.ReturnMatchDays(uint8(season.Id))
		season.MatchDays = matchDays
		log.Println(season.Year, len(season.MatchDays))
	}

	data, err := json.Marshal(seasons)
	if err != nil {
		log.Panicln(err)
	}
	os.WriteFile("data.json", data, 0644)

}

func NewCrawler(baseUrl string, timeoutSeconds int) *Crawler {
	return &Crawler{
		httpClient: http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		baseUrl: baseUrl,
	}
}

// ReturnSeasons return all Seasons with year and id
func ReturnSeasons(d *goquery.Document) []*types.Season {
	var seasons []*types.Season
	d.Find("select[id=season]").Find("option").Each(func(i int, s *goquery.Selection) {
		se := &types.Season{}

		yearString := strings.Split(strings.Split(s.Text(), "Saison")[1], "/")[0]
		year, err := strconv.ParseUint(yearString[1:], 10, 64)
		if err != nil {
			fmt.Printf("Failed to parse: %v\n", err)
			return
		}
		se.Year = year

		if id, ok := s.Attr("value"); ok {
			idunit, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				fmt.Printf("Failed to parse: %v\n", err)
				return
			}
			se.Id = idunit
		}
		seasons = append(seasons, se)
	})

	return seasons
}

func (c *Crawler) ReturnMatchDays(seasonId uint8) []types.MatchDay {
	var matchDays []types.MatchDay
	doc := c.FetchUrl(seasonId, 1)
	doc.Find("select[id=match]").Find("option").Each(func(i int, s *goquery.Selection) {
		m := types.MatchDay{}
		parseGameDate(&m, s.Text())

		if id, ok := s.Attr("value"); ok {
			idunit, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				fmt.Printf("Failed to parse: %v\n", err)
				return
			}
			m.Id = uint8(idunit)
		}
		matchDayDoc := c.FetchUrl(seasonId, m.Id)
		m.MatchResults = ReturnMatchResults(matchDayDoc)
		matchDays = append(matchDays, m)

	})

	return matchDays

}

func ReturnMatchResults(doc *goquery.Document) []types.MatchResult {
	var matches []types.MatchResult
	doc.Find("table").First().Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			rows := table.Find("tr")
			rowCount := rows.Length()

			// Skip header (first) and footer (last) rows
			if j == 0 || j == rowCount-1 {
				return
			}

			match := types.MatchResult{}

			cols := row.Find("td")
			if cols.Length() >= 4 {
				timeStr := strings.TrimSpace(cols.Eq(1).Text())
				matchStr := strings.TrimSpace(cols.Eq(2).Text())
				resultStr := strings.TrimSpace(cols.Eq(3).Text())

				// Parse teams
				teams := strings.Split(matchStr, ":")
				if len(teams) == 2 {
					match.HomeTeam = strings.TrimSpace(teams[0])
					match.AwayTeam = strings.TrimSpace(teams[1])
				}

				// Parse scores if available
				scores := strings.Split(resultStr, ":")
				if len(scores) == 2 {
					// Only parse if it's not "- : -"
					if scores[0] != "-" && scores[1] != "-" {
						fmt.Sscanf(scores[0], "%d", &match.HomeScore)
						fmt.Sscanf(scores[1], "%d", &match.AwayScore)
					}
				}

				// Parse date and time
				if t, err := time.Parse("15:04", timeStr); err == nil {
					match.Time = t
				}

				matches = append(matches, match)
			}
		})
	})
	return matches
}

func (c *Crawler) FetchUrl(season, match uint8) *goquery.Document {
	v := url.Values{
		"season": []string{fmt.Sprintf("%d", season)},
		"match":  []string{fmt.Sprintf("%d", match)},
	}
	req, err := http.NewRequest("POST", c.baseUrl, strings.NewReader(v.Encode()))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Errorf("failed to parse HTML: %v", err)
		return nil
	}
	return doc
}

func parseGameDate(g *types.MatchDay, s string) error {
	// Extract date part by splitting on "-" and trimming spaces
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid date format")
	}

	datePart := strings.TrimSpace(parts[1])
	t, err := time.Parse("02.01.2006", datePart)
	if err != nil {
		return err
	}
	g.Date = t
	return nil
}
