package services

import (
	"encoding/json"
	"fmt"
	"hockeyAnalytics/internal/dto/nhl"
	"net/http"
	"net/url"
)

func FetchRealtimeStats(
	season string,
) ([]nhl.RealtimePlayerStats, error) {

	baseURL :=
		"https://api.nhle.com/stats/rest/en/skater/realtime"

	allPlayers :=
		[]nhl.RealtimePlayerStats{}

	start := 0

	for {

		params := url.Values{}

		params.Add(
			"isAggregate",
			"false",
		)

		params.Add(
			"isGame",
			"false",
		)

		params.Add(
			"start",
			fmt.Sprintf("%d", start),
		)

		params.Add(
			"limit",
			"100",
		)

		params.Add(
			"sort",
			`[{"property":"hits","direction":"DESC"}]`,
		)

		params.Add(
			"cayenneExp",
			"gameTypeId=2 and seasonId="+season,
		)

		fullURL :=
			baseURL + "?" + params.Encode()

		resp, err :=
			http.Get(fullURL)

		if err != nil {

			return nil, err
		}

		var result nhl.RealtimeResponse

		err = json.NewDecoder(
			resp.Body,
		).Decode(&result)

		resp.Body.Close()

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

		start += 100
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

		params := url.Values{}

		params.Add(
			"isAggregate",
			"false",
		)

		params.Add(
			"isGame",
			"false",
		)

		params.Add(
			"start",
			fmt.Sprintf("%d", start),
		)

		params.Add(
			"limit",
			"100",
		)

		params.Add(
			"sort",
			`[{"property":"points","direction":"DESC"}]`,
		)

		params.Add(
			"cayenneExp",
			"gameTypeId=2 and seasonId="+season,
		)

		fullURL :=
			baseURL + "?" + params.Encode()

		resp, err :=
			http.Get(fullURL)

		if err != nil {

			return nil, err
		}

		var result nhl.SummaryResponse

		err = json.NewDecoder(
			resp.Body,
		).Decode(&result)

		resp.Body.Close()

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

		start += 100
	}

	return allPlayers, nil
}
