package models

import "time"

type Account struct {
	ID            int       `json:"robloxId"`
	Name          string    `json:"robloxName"`
	Secrets       int       `json:"secrets,omitempty"`
	Eggs          int       `json:"eggs,omitempty"`
	Bubbles       int       `json:"bubbles,omitempty"`
	Power         int       `json:"power,omitempty"`
	Robux         int       `json:"robux,omitempty"`
	Playtime      int       `json:"playtime,omitempty"`
	LastSavedTime time.Time `json:"time_saved"`
}

func NewPlayer(ID int, Name string, Secrets int, Eggs int, Bubbles int, Power int, Robux int, Time int) *Account {
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
