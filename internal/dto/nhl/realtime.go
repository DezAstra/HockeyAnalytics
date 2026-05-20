package nhl

type RealtimeResponse struct {
	Data []RealtimePlayerStats `json:"data"`
}

type RealtimePlayerStats struct {
	PlayerID int `json:"playerId"`

	Hits int `json:"hits"`

	BlockedShots int `json:"blockedShots"`
}
