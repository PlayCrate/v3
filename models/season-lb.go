package models

type SeasonLBAccount struct {
	Payload     string `json:"payload"`
	RobloxID    int64  `json:"robloxId"`
	SeasonMain  int64  `json:"season_main"`
	SeasonEvent int64  `json:"season_event"`
}

type GetSeasonLB struct {
	SeasonMain []struct {
		RobloxID        int64 `json:"robloxId"`
		SeasonMainCount int64 `json:"value"`
	} `json:"season_main"`
	SeasonEvent []struct {
		RobloxID         int64 `json:"robloxId"`
		SeasonEventCount int64 `json:"value"`
	} `json:"season_event"`
}

func NewGetSeasonLB() *GetSeasonLB {
	return &GetSeasonLB{
		SeasonMain: []struct {
			RobloxID        int64 `json:"robloxId"`
			SeasonMainCount int64 `json:"value"`
		}{},
		SeasonEvent: []struct {
			RobloxID         int64 `json:"robloxId"`
			SeasonEventCount int64 `json:"value"`
		}{},
	}
}
