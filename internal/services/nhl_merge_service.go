package services

import "hockeyAnalytics/internal/dto/nhl"

func MergePlayerStats(
	summary []nhl.PlayerSummary,
	realtime []nhl.RealtimePlayerStats,
) []nhl.CombinedPlayerStats {

	realtimeMap :=
		make(
			map[int]nhl.RealtimePlayerStats,
		)

	for _, player := range realtime {

		realtimeMap[player.PlayerID] =
			player
	}

	var combined []nhl.CombinedPlayerStats

	for _, player := range summary {

		rt :=
			realtimeMap[player.PlayerID]

		combined = append(
			combined,
			nhl.CombinedPlayerStats{
				PlayerSummary: player,

				Hits: rt.Hits,

				BlockedShots: rt.BlockedShots,
			},
		)
	}

	return combined
}
