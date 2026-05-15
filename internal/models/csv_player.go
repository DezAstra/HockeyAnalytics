package models

type CSVPlayer struct {
	Name           string
	Team           string
	Position       string
	GamesPlayed    int
	Goals          int
	Assists        int
	PlusMinus      int
	PenaltyMinutes int
	TimeOfIce      float64
}
