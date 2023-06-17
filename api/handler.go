package api

import (
	"encoding/json"
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

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusNotFound, ApiResponse{
			Success: false,
			Error:   "Endpoint not found",
		})
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusMethodNotAllowed, ApiResponse{
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
			WriteJSON(res, http.StatusUnauthorized, ApiResponse{
				Success: false,
				Error:   "No token provided",
			})
			return
		}

		if tokenString != s.cfg.Auth {
			WriteJSON(res, http.StatusUnauthorized, ApiResponse{
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
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiResponse{Error: err.Error()})
		}
	}
}

func (s *APIServer) GetLeaderboards(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["which"]

	leaderboards := map[string]func() (*storage.PlayerDataResponse, error){
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

	return WriteJSON(w, http.StatusBadRequest, ApiResponse{
		Success: false,
		Error:   "Invalid Leaderboard",
	})
}

func (s *APIServer) InsertPlayer(w http.ResponseWriter, r *http.Request) error {
	createAccReq := new(models.Account)
	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		return err
	}

	if createAccReq.ID == 0 && createAccReq.Name == "" {
		return WriteJSON(w, 400, ApiResponse{
			Success: false,
			Error:   "ID cannot be 0",
		})
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
