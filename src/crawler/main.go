package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcelschliesser/werbeliga-hamburg/types"
)

func parseGameDate(g *types.Match, s string) error {
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

func ParseMatchResults(url string) ([]types.MatchResult, error) {
	formData := "season=23&match=560"

	req, err := http.NewRequest("POST", url, strings.NewReader(formData))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var matches []types.MatchResult

	// Get Match Results
	doc.Find("table").First().Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {

			rows := table.Find("tr")
			rowCount := rows.Length()

			// Skip header (first) and footer (last) rows
			if j == 0 || j == rowCount-1 {
				return
			}

			var match types.MatchResult

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
				// Note: You'll need to combine the date from the header with the time
				// This is a simplified example
				if t, err := time.Parse("15:04", timeStr); err == nil {
					match.DateTime = t
				}

				matches = append(matches, match)
			}
		})
	})

	// Get Season
	doc.Find("select[id=season]").Find("option").Each(func(i int, s *goquery.Selection) {
		se := types.Season{}
		se.Name = s.Text()
		if id, ok := s.Attr("value"); ok {
			se.Id = id
		}
		if _, ok := s.Attr("selected"); ok {
			se.Selected = ok
		}
		fmt.Println(i, se)
	})

	// Get Match
	doc.Find("select[id=match]").Find("option").Each(func(i int, s *goquery.Selection) {
		m := types.Match{}
		parseGameDate(&m, s.Text())

		if id, ok := s.Attr("value"); ok {
			m.Id = id
		}

		if _, ok := s.Attr("selected"); ok {
			m.Selected = ok
		}
		fmt.Println(i, m)

	})

	return matches, nil

}
