package nhl

type SummaryResponse struct {
	Data []PlayerSummary `json:"data"`
}

type PlayerSummary struct {
	PlayerID int `json:"playerId"`

	SkaterFullName string `json:"skaterFullName"`

	LastName string `json:"lastName"`

	TeamAbbrevs string `json:"teamAbbrevs"`

	PositionCode string `json:"positionCode"`

	GamesPlayed int `json:"gamesPlayed"`

	Goals int `json:"goals"`

	Assists int `json:"assists"`

	Points int `json:"points"`

	PlusMinus int `json:"plusMinus"`

	PenaltyMinutes int `json:"penaltyMinutes"`

	Shots int `json:"shots"`

	ShootingPct float64 `json:"shootingPct"`

	FaceoffWinPct float64 `json:"faceoffWinPct"`

	EvenStrengthGoals int `json:"evGoals"`

	PowerPlayGoals int `json:"ppGoals"`

	ShortHandedGoals int `json:"shGoals"`

	EvenStrengthPoints int `json:"evPoints"`

	PowerPlayPoints int `json:"ppPoints"`

	ShortHandedPoints int `json:"shPoints"`

	GameWinningGoals int `json:"gameWinningGoals"`

	TimeOnIcePerGame float64 `json:"timeOnIcePerGame"`

	SeasonID int `json:"seasonId"`
}
