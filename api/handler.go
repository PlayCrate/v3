package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kattah7/v3/models"
	"github.com/kattah7/v3/storage"
	"github.com/redis/go-redis/v9"
)

type apiFunc func(http.ResponseWriter, *http.Request, *APIServer) error

type APIServer struct {
	ctx   context.Context
	store storage.Storage
	cfg   *models.Config
	rdb   *redis.Client
}

type ApiResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewAPIServer(ctx context.Context, cfg *models.Config, store storage.Storage, rdb *redis.Client) *APIServer {
	return &APIServer{
		ctx:   ctx,
		store: store,
		cfg:   cfg,
		rdb:   rdb,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = s.customHandler(route.HandlerFunc)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.WriteJSON(w, http.StatusOK, ApiResponse{
			Success: false,
			Error:   "Endpoint not found",
		})
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.WriteJSON(w, http.StatusOK, ApiResponse{
			Success: false,
			Error:   "Method not allowed",
		})
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         s.cfg.ListenAddress,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			return s.ctx
		},
	}

	log.Println("JSON API server running on port", s.cfg.ListenAddress)
	log.Fatal(srv.ListenAndServe())
}

func (s *APIServer) customHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: false,
				Error:   "No token provided",
			})
			return
		}

		if tokenString != s.cfg.Auth {
			s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: false,
				Error:   "Invalid token",
			})
			return
		}

		log.Println(r.Method, r.URL.Path)

		if err := f(w, r, s); err != nil {
			s.WriteJSON(w, http.StatusOK, ApiResponse{Error: err.Error()})
		}
	}
}

func (s *APIServer) WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func InsertPlayer(w http.ResponseWriter, r *http.Request, s *APIServer) error {
	createAccReq := new(models.Account)
	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		return err
	}

	if createAccReq.ID == 0 && createAccReq.Name == "" {
		return fmt.Errorf("ID cannot be 0")
	}

	account := models.NewPlayer(
		createAccReq.ID,
		createAccReq.Name,
		createAccReq.Secrets,
		createAccReq.Eggs,
		createAccReq.Bubbles,
		createAccReq.Power,
		createAccReq.Robux,
		createAccReq.Playtime,
	)

	if err := s.store.InsertAccounts(account); err != nil {
		return err
	}

	return s.WriteJSON(w, 200, ApiResponse{
		Success: true,
		Data:    account,
	})
}

func LeaderboardLookup(w http.ResponseWriter, r *http.Request, s *APIServer) error {
	if r.Method == "POST" {
		InsertAcc := new(models.AccountLookup)

		if err := json.NewDecoder(r.Body).Decode(InsertAcc); err != nil {
			return err
		}

		if InsertAcc.RobloxID == 0 {
			return fmt.Errorf("Missing robloxId")
		}

		acc, err := s.store.GetSpecificPlayer(InsertAcc.RobloxID)

		if err != nil {
			return err
		}

		return s.WriteJSON(w, http.StatusOK, ApiResponse{
			Success: true,
			Data:    acc,
		})
	}

	return fmt.Errorf("Invalid Method")
}

func PetsExistance(w http.ResponseWriter, r *http.Request, s *APIServer) error {
	if r.Method == "POST" {
		InsertAcc := new(models.PetsExistance)

		if err := json.NewDecoder(r.Body).Decode(InsertAcc); err != nil {
			return err
		}

		if InsertAcc.Payload == "" {
			return fmt.Errorf("Missing Payload")
		}

		switch InsertAcc.Payload {
		case "INSERT_PETS_EXISTANCE":
			if err := s.store.InsertPetsExistance(InsertAcc); err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Pets Inserted",
			})
		case "READ_PETS_EXISTANCE":
			cached, err := s.rdb.Get(context.Background(), "pets-exist").Result()
			if err != nil {
				return err
			}

			var cachedPets []*models.GetPetsExistance
			err = json.Unmarshal([]byte(cached), &cachedPets)
			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    cachedPets,
			})
		case "DELETE_PETS_EXISTANCE":
			if err := s.store.DeletePetsExistence(InsertAcc); err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Pets Removed",
			})
		default:
			return fmt.Errorf("Invalid Payload")
		}
	}

	return fmt.Errorf("Invalid Method")
}

func Auctions(w http.ResponseWriter, r *http.Request, s *APIServer) error {
	if r.Method == "POST" {
		Auction := new(models.AuctionAccount)
		if err := json.NewDecoder(r.Body).Decode(Auction); err != nil {
			return err
		}

		if Auction.Payload == "" {
			return fmt.Errorf("Invalid Payload")
		}

		switch Auction.Payload {
		case "LIST":
			if err := s.store.ListAuction(Auction); err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Inserted",
			})

		case "READ":
			auctions, err := s.store.GetAuctions()
			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    auctions,
			})
		case "DELETE":
			if err := s.store.RemoveAuction(Auction); err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Deleted",
			})
		case "PURCHASE":
			if err := s.store.PurchaseAuction(Auction); err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Purchased",
			})
		case "AUCTION_GET_CLAIMS":
			claims, err := s.store.GetAuctionClaims(Auction)

			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    claims,
			})
		case "AUCTION_CLAIM":
			if err := s.store.AuctionClaim(Auction); err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Claimed",
			})
		case "AUCTION_GET_LISTINGS":
			listing, err := s.store.GetAuctionListing(Auction)

			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    listing,
			})
		case "AUCTION_UNLIST":
			if err := s.store.AuctionUnlist(Auction); err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Unlisted",
			})
		}
	}

	return fmt.Errorf("Invalid Method")
}

func GetLeaderboards(w http.ResponseWriter, r *http.Request, s *APIServer) error {
	vars := mux.Vars(r)["which"]

	switch vars {
	case "eggs":
		{
			cached, err := s.rdb.Get(context.Background(), "eggs-lb").Result()
			if err != nil {
				return err
			}
			cachedEggs := new(any)
			err = json.Unmarshal([]byte(cached), &cachedEggs)

			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    cachedEggs,
			})
		}
	case "bubbles":
		{
			cached, err := s.rdb.Get(context.Background(), "bubbles-lb").Result()
			if err != nil {
				return err
			}
			cachedBubbles := new(any)
			err = json.Unmarshal([]byte(cached), &cachedBubbles)

			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    cachedBubbles,
			})
		}
	case "secrets":
		{
			cached, err := s.rdb.Get(context.Background(), "secrets-lb").Result()
			if err != nil {
				return err
			}

			cachedSecrets := new(any)
			err = json.Unmarshal([]byte(cached), &cachedSecrets)

			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    cachedSecrets,
			})
		}
	case "power":
		{
			cached, err := s.rdb.Get(context.Background(), "power-lb").Result()
			if err != nil {
				return err
			}

			cachedPower := new(any)
			err = json.Unmarshal([]byte(cached), &cachedPower)

			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    cachedPower,
			})
		}
	case "robux":
		{
			cached, err := s.rdb.Get(context.Background(), "robux-lb").Result()
			if err != nil {
				return err
			}

			cachedRobux := new(any)
			err = json.Unmarshal([]byte(cached), &cachedRobux)

			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    cachedRobux,
			})
		}
	case "playtime":
		{
			cached, err := s.rdb.Get(context.Background(), "playtime-lb").Result()
			if err != nil {
				return err
			}

			cachedPlaytime := new(any)
			err = json.Unmarshal([]byte(cached), &cachedPlaytime)
			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    cachedPlaytime,
			})
		}
	default:
		return fmt.Errorf("Invalid Leaderboard")
	}
}

func SeasonLB(w http.ResponseWriter, r *http.Request, s *APIServer) error {
	if r.Method == "POST" {
		SeasonLB := new(models.SeasonLBAccount)
		if err := json.NewDecoder(r.Body).Decode(SeasonLB); err != nil {
			return err
		}

		if SeasonLB.Payload == "" {
			return fmt.Errorf("Invalid Payload")
		}

		switch SeasonLB.Payload {
		case "INSERT_ACCOUNT":
			if err := s.store.InsertSeasonLB(SeasonLB); err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Successfully inserted into season leaderboard",
			})
		case "READ_LEADERBOARD":
			cached, err := s.rdb.Get(context.Background(), "season-lb").Result()
			if err != nil {
				return err
			}
			cachedSeason := new(any)
			err = json.Unmarshal([]byte(cached), &cachedSeason)

			if err != nil {
				return err
			}

			return s.WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    cachedSeason,
			})
		}
	}

	return fmt.Errorf("Invalid Method")
}
