package models

type HalloweenAccount struct {
	Payload  string `json:"payload"`
	RobloxID int64  `json:"robloxId"`
	Houses   int64  `json:"houses"`
	Candies  int64  `json:"candies"`
}

type GetHalloweenLB struct {
	HousesCount []struct {
		RobloxID   int64 `json:"robloxId"`
		HouseCount int64 `json:"value"`
	} `json:"houses"`
	CandiesCount []struct {
		RobloxID   int64 `json:"robloxId"`
		CandyCount int64 `json:"value"`
	} `json:"candies"`
}

func NewGetHalloweenLB() *GetHalloweenLB {
	return &GetHalloweenLB{
		HousesCount: []struct {
			RobloxID   int64 `json:"robloxId"`
			HouseCount int64 `json:"value"`
		}{},
		CandiesCount: []struct {
			RobloxID   int64 `json:"robloxId"`
			CandyCount int64 `json:"value"`
		}{},
	}
}
