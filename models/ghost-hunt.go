package models

type GhostHuntAccount struct {
	Payload  string `json:"payload"`
	RobloxID int64  `json:"robloxId"`
}

type GhostHuntSerial struct {
	Serial int64 `json:"serial"`
}
