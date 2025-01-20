package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcelschliesser/werbeliga-hamburg/types"
)

func ParseMatchResults(url string) ([]types.MatchResult, error) {
	formData := "season=24&match=607"

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

	// Find the table containing match results
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			// Skip header row
			if j == 0 {
				return
			}

			var match types.MatchResult

			// Extract data from columns
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

	return matches, nil
}
