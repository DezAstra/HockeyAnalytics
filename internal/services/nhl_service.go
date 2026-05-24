package services

import (
	"encoding/json"
	"fmt"
	"hockeyAnalytics/internal/dto/nhl"
	"net/http"
	"net/url"
	"time"
)

const (
	nhlPageLimit   = 100
	nhlMaxAttempts = 3
)

var nhlHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

func buildNHLStatsURL(
	baseURL string,
	season string,
	start int,
	sort string,
) string {
	params := url.Values{}

	params.Set(
		"isAggregate",
		"false",
	)

	params.Set(
		"isGame",
		"false",
	)

	params.Set(
		"start",
		fmt.Sprintf("%d", start),
	)

	params.Set(
		"limit",
		fmt.Sprintf("%d", nhlPageLimit),
	)

	params.Set(
		"sort",
		sort,
	)

	params.Set(
		"cayenneExp",
		"gameTypeId=2 and seasonId="+season,
	)

	return baseURL + "?" + params.Encode()
}

func isRetryableNHLStatus(
	statusCode int,
) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode >= http.StatusInternalServerError
}

func fetchNHLPage[T any](
	fullURL string,
	result *T,
) error {
	var lastErr error

	for attempt := 1; attempt <= nhlMaxAttempts; attempt++ {
		req, err := http.NewRequest(
			http.MethodGet,
			fullURL,
			nil,
		)

		if err != nil {
			return err
		}

		req.Header.Set(
			"User-Agent",
			"hockeyAnalytics/1.0 (+https://localhost)",
		)

		req.Header.Set(
			"Accept",
			"application/json",
		)

		resp, err := nhlHTTPClient.Do(req)

		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		if isRetryableNHLStatus(resp.StatusCode) {
			lastErr = fmt.Errorf(
				"NHL API returned retryable status %d for %s",
				resp.StatusCode,
				fullURL,
			)

			resp.Body.Close()
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		if resp.StatusCode < http.StatusOK ||
			resp.StatusCode >= http.StatusMultipleChoices {
			resp.Body.Close()
			return fmt.Errorf(
				"NHL API returned status %d for %s",
				resp.StatusCode,
				fullURL,
			)
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		resp.Body.Close()

		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		return nil
	}

	return lastErr
}

func FetchRealtimeStats(
	season string,
) ([]nhl.RealtimePlayerStats, error) {
	baseURL :=
		"https://api.nhle.com/stats/rest/en/skater/realtime"

	allPlayers :=
		[]nhl.RealtimePlayerStats{}

	start := 0

	for {
		fullURL := buildNHLStatsURL(
			baseURL,
			season,
			start,
			`[{"property":"hits","direction":"DESC"}]`,
		)

		var result nhl.RealtimeResponse

		err := fetchNHLPage(
			fullURL,
			&result,
		)

		if err != nil {
			return nil, err
		}

		if len(result.Data) == 0 {
			break
		}

		allPlayers = append(
			allPlayers,
			result.Data...,
		)

		start += nhlPageLimit
	}

	return allPlayers, nil
}

func FetchSeasonSummary(
	season string,
) ([]nhl.PlayerSummary, error) {
	baseURL :=
		"https://api.nhle.com/stats/rest/en/skater/summary"

	allPlayers :=
		[]nhl.PlayerSummary{}

	start := 0

	for {
		fullURL := buildNHLStatsURL(
			baseURL,
			season,
			start,
			`[{"property":"points","direction":"DESC"}]`,
		)

		var result nhl.SummaryResponse

		err := fetchNHLPage(
			fullURL,
			&result,
		)

		if err != nil {
			return nil, err
		}

		if len(result.Data) == 0 {
			break
		}

		allPlayers = append(
			allPlayers,
			result.Data...,
		)

		start += nhlPageLimit
	}

	return allPlayers, nil
}
