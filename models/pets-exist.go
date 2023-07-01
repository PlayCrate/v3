package models

type PetsExistance struct {
	Payload  string           `json:"payload"`
	RobloxID int64            `json:"robloxId"`
	Pets     []map[string]int `json:"pets"`
}

type GetPetsExistance struct {
	PetID    string `json:"petId"`
	PetCount int    `json:"petCount"`
}
