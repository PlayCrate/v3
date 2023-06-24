package models

import "time"

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
