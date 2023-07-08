package models

import (
	"database/sql"
	"time"
)

type PlayerDataResponse struct {
	F2P    []*Account `json:"f2p,omitempty"`
	NonF2P []*Account `json:"nof2p,omitempty"`
	Other  []*Account `json:"other,omitempty"`
}

type Account struct {
	ID            int64     `json:"robloxId"`
	Name          string    `json:"robloxName"`
	Secrets       int64     `json:"secrets,omitempty"`
	Eggs          int64     `json:"eggs,omitempty"`
	Bubbles       int64     `json:"bubbles,omitempty"`
	Power         int64     `json:"power,omitempty"`
	Robux         int64     `json:"robux,omitempty"`
	Playtime      int64     `json:"playtime,omitempty"`
	LastSavedTime time.Time `json:"time_saved"`
}

type AccountLookup struct {
	RobloxID   int64  `json:"robloxId"`
	RobloxName string `json:"robloxName"`

	Secrets  int64 `json:"secrets"`
	Eggs     int64 `json:"eggs"`
	Bubbles  int64 `json:"bubbles"`
	Power    int64 `json:"power"`
	Playtime int64 `json:"playtime"`
	Robux    int64 `json:"robux"`

	SecretsRank  int64 `json:"secretsRank"`
	EggsRank     int64 `json:"eggsRank"`
	BubblesRank  int64 `json:"bubblesRank"`
	PowerRank    int64 `json:"powerRank"`
	PlaytimeRank int64 `json:"playtimeRank"`
	RobuxRank    int64 `json:"robuxRank"`

	F2PSecretsRank sql.NullInt64 `json:"freeToPlaySecretsRank"`
	F2PEggsRank    sql.NullInt64 `json:"freeToPlayEggsRank"`
	F2PBubblesRank sql.NullInt64 `json:"freeToPlayBubblesRank"`
	F2PPowerRank   sql.NullInt64 `json:"freeToPlayPowerRank"`
}

func NewPlayer(ID int64, Name string, Secrets int64, Eggs int64, Bubbles int64, Power int64, Robux int64, Time int64) *Account {
	return &Account{
		ID:            ID,
		Name:          Name,
		Secrets:       Secrets,
		Eggs:          Eggs,
		Bubbles:       Bubbles,
		Power:         Power,
		Robux:         Robux,
		Playtime:      Time,
		LastSavedTime: time.Now().UTC(),
	}
}
