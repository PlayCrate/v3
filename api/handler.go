package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kattah7/v3/models"
	"github.com/kattah7/v3/storage"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listenAddr string
	store      storage.Storage
	cfg        *models.Config
}

type ApiResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewAPIServer(cfg *models.Config, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: cfg.ListenAddress,
		store:      store,
		cfg:        cfg,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/leaderboard", Authorization(makeHTTPHandleFunc(s.InsertPlayer), s)).Methods(http.MethodPost)
	router.HandleFunc("/leaderboard/{which}", Authorization(makeHTTPHandleFunc(s.GetLeaderboards), s)).Methods(http.MethodGet)
	router.HandleFunc("/auction", Authorization(makeHTTPHandleFunc(s.Auctions), s))
	router.HandleFunc("/pets-exist", Authorization(makeHTTPHandleFunc(s.PetsExistance), s))
	router.HandleFunc("/lb-lookup", Authorization(makeHTTPHandleFunc(s.LeaderboardLookup), s))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, ApiResponse{
			Success: false,
			Error:   "Endpoint not found",
		})
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, ApiResponse{
			Success: false,
			Error:   "Method not allowed",
		})
	})

	log.Println("JSON API server running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func Authorization(handlerFunc http.HandlerFunc, s *APIServer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		tokenString := req.Header.Get("Authorization")
		if tokenString == "" {
			WriteJSON(res, http.StatusOK, ApiResponse{
				Success: false,
				Error:   "No token provided",
			})
			return
		}

		if tokenString != s.cfg.Auth {
			WriteJSON(res, http.StatusOK, ApiResponse{
				Success: false,
				Error:   "Invalid token",
			})
			return
		}

		handlerFunc(res, req)
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)

		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusOK, ApiResponse{Error: err.Error()})
		}
	}
}

func (s *APIServer) LeaderboardLookup(w http.ResponseWriter, r *http.Request) error {
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

		return WriteJSON(w, http.StatusOK, ApiResponse{
			Success: true,
			Data:    acc,
		})
	}

	return fmt.Errorf("Invalid Method")
}

func (s *APIServer) PetsExistance(w http.ResponseWriter, r *http.Request) error {
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

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Pets Inserted",
			})
		case "READ_PETS_EXISTANCE":
			pets, err := s.store.GetPetsExistance()

			if err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    pets,
			})
		case "DELETE_PETS_EXISTANCE":
			if err := s.store.DeletePetsExistence(InsertAcc); err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Pets Removed",
			})
		default:
			return fmt.Errorf("Invalid Payload")
		}
	}

	return fmt.Errorf("Invalid Method")
}

func (s *APIServer) Auctions(w http.ResponseWriter, r *http.Request) error {
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

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Inserted",
			})

		case "READ":
			auctions, err := s.store.GetAuctions()
			if err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    auctions,
			})
		case "DELETE":
			if err := s.store.RemoveAuction(Auction); err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Deleted",
			})
		case "PURCHASE":
			if err := s.store.PurchaseAuction(Auction); err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Purchased",
			})
		case "AUCTION_GET_CLAIMS":
			claims, err := s.store.GetAuctionClaims(Auction)

			if err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    claims,
			})
		case "AUCTION_CLAIM":
			if err := s.store.AuctionClaim(Auction); err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Claimed",
			})
		case "AUCTION_GET_LISTINGS":
			listing, err := s.store.GetAuctionListing(Auction)

			if err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    listing,
			})
		case "AUCTION_UNLIST":
			if err := s.store.AuctionUnlist(Auction); err != nil {
				return err
			}

			return WriteJSON(w, http.StatusOK, ApiResponse{
				Success: true,
				Data:    "Auction Unlisted",
			})
			// case "AUCTION_EXPIRE_LIST":
			// 	expires, err := s.store.AuctionExpireList(Auction)

			// 	if err != nil {
			// 		return err
			// 	}

			// 	return WriteJSON(w, http.StatusOK, ApiResponse{
			// 		Success: true,
			// 		Data:    expires,
			// 	})
			// case "AUCTION_EXPIRE_CLAIM":
			// 	if err := s.store.AuctionExpireClaim(Auction); err != nil {
			// 		return err
			// 	}

			// 	return WriteJSON(w, http.StatusOK, ApiResponse{
			// 		Success: true,
			// 		Data:    "Auction Expire Claimed",
			// 	})
		}
	}

	return fmt.Errorf("Invalid Method")
}

func (s *APIServer) GetLeaderboards(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["which"]

	leaderboards := map[string]func() (*models.PlayerDataResponse, error){
		"eggs":     s.store.GetEggs,
		"bubbles":  s.store.GetBubbles,
		"secrets":  s.store.GetSecrets,
		"power":    s.store.GetPower,
		"robux":    s.store.GetRobux,
		"playtime": s.store.GetPlaytime,
	}

	if leaderboardFunc, ok := leaderboards[vars]; ok {
		data, err := leaderboardFunc()
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, ApiResponse{
			Success: true,
			Data:    data,
		})
	}

	return fmt.Errorf("Invalid Leaderboard")
}

func (s *APIServer) InsertPlayer(w http.ResponseWriter, r *http.Request) error {
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

	return WriteJSON(w, 200, ApiResponse{
		Success: true,
		Data:    account,
	})
}
