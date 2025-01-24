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

type Crawler struct {
	httpClient http.Client
	baseUrl    string
}

func main() {

	c := NewCrawler(os.Getenv("URL"), 10)

	seasons := c.fetchAllSeasons()

	c.fetchAllMatchIds(&seasons)

	// TODO: parallelize requests
	for i := range seasons {
		for j := range seasons[i].MatchDays {
			doc := c.FetchUrl(uint(seasons[i].Id), uint(seasons[i].MatchDays[j].Id))
			res := ReturnMatchResults(doc)
			seasons[i].MatchDays[j].MatchResults = res
		}
	}

	data, err := json.Marshal(seasons)
	if err != nil {
		log.Panicln(err)
	}
	os.WriteFile("data.json", data, 0644)
	log.Println(len(seasons))
}

// fetchAllMatchIds will fetch all MatchIds to given SeasonIds
func (c *Crawler) fetchAllMatchIds(seasons *[]types.Season) {

	for i := range *seasons {
		season := &(*seasons)[i] // TODO: Understand this "Hack" :-D
		doc := c.FetchUrl(uint(season.Id), 1)
		doc.Find("select[id=match]").Find("option").Each(func(i int, s *goquery.Selection) {
			var m types.MatchDay
			if id, ok := s.Attr("value"); ok {
				iduint, err := strconv.ParseUint(id, 10, 16)
				if err != nil {
					log.Fatalln(err.Error())
				}
				m.Id = types.MatchId(iduint)
				season.MatchDays = append(season.MatchDays, m)
			}
		})
	}
}

// fetchAllSeasonIds is the starting point and fetch all current seasons
func (c *Crawler) fetchAllSeasons() []types.Season {

	var firstSeasonId uint = 2
	var seasons []types.Season

	doc := c.FetchUrl(firstSeasonId, 1)

	doc.Find("select[id=season]").Find("option").Each(func(i int, s *goquery.Selection) {

		var season types.Season

		if id, ok := s.Attr("value"); ok {
			iduint, err := strconv.ParseUint(id, 10, 16)
			if err != nil {
				log.Fatalln(err.Error())
			}
			season.Id = types.SeasonId(iduint)
			season.Year = yearFromString(s.Text())
			seasons = append(seasons, season)

		}
	})

	return seasons
}

// NewCrawler initialize a crawler with the werbeliga.de baseUrl
// and a timeout per request of 10 seconds
func NewCrawler(baseUrl string, timeoutSeconds int) *Crawler {
	return &Crawler{
		httpClient: http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		baseUrl: baseUrl,
	}
}

func (c *Crawler) ReturnMatchDays(seasonId uint) []types.MatchDay {
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
			m.Id = types.MatchId(idunit)
		}
		matchDayDoc := c.FetchUrl(seasonId, uint(m.Id))
		m.MatchResults = ReturnMatchResults(matchDayDoc)
		matchDays = append(matchDays, m)

	})

	return matchDays

}

func ReturnMatchResults(doc *goquery.Document) []types.Match {
	var matches []types.Match
	doc.Find("table").First().Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			rows := table.Find("tr")
			rowCount := rows.Length()

			// Skip header (first) and footer (last) rows
			if j == 0 || j == rowCount-1 {
				return
			}

			match := types.Match{}

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
					match.DateTime = t
				}

				matches = append(matches, match)
			}
		})
	})
	return matches
}

func (c *Crawler) FetchUrl(season, match uint) *goquery.Document {
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
		log.Fatalln(err.Error())
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
	log.Println(t)
	return nil
}

func yearFromString(yearString string) uint {
	s := strings.Split(strings.Split(yearString, "Saison")[1], "/")[0]
	year, err := strconv.ParseUint(s[1:], 10, 64)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return uint(year)
}
