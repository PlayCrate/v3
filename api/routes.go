package api

// Route is the model for the router setup
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc apiFunc
}

// Routes are the main setup for our Router
type Routes []Route

var routes = Routes{
	Route{"InsertLeaderboards", "POST", "/leaderboard", InsertPlayer},
	Route{"GetLeaderboards", "GET", "/leaderboard/{which}", GetLeaderboards},
	Route{"LeaderboardLookup", "POST", "/lb-lookup", LeaderboardLookup},

	Route{"Auction", "POST", "/auction", Auctions},

	Route{"SeasonLB", "POST", "/season-lb", SeasonLB},
	Route{"HalloweenLB", "POST", "/halloween-lb", HalloweenLB},
	Route{"GhostHunt", "POST", "/ghost-hunt", GhostHunt},

	Route{"PetsExistance", "POST", "/pets-exist", PetsExistance},
}
