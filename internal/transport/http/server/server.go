package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gersastas/wallet-service/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
	stats      map[string]int
	mu         sync.Mutex
	wallet     map[string]*models.Wallet
	walletsMu  sync.Mutex
}

func New(address string) *Server {
	s := &Server{
		stats: make(map[string]int),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/wallets/create", s.createWallet)
	mux.HandleFunc("/wallets/get", s.getWallet)

	s.httpServer = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return s
}

func (s *Server) Run() error {
	err := s.httpServer.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

type WalletRequest struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type WalletResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) createWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		http.Error(w, "currency is required", http.StatusBadRequest)
		return
	}

	now := time.Now()

	walletID := uuid.New()

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "user_id is invalid", http.StatusBadRequest)
		return
	}

	wallet := &models.Wallet{
		ID:        walletID,
		UserID:    userUUID,
		Name:      req.Name,
		Balance:   0,
		Currency:  req.Currency,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	s.walletsMu.Lock()
	s.wallet[walletID.String()] = wallet
	s.walletsMu.Unlock()

	resp := WalletResponse{
		ID:        wallet.ID.String(),
		UserID:    wallet.UserID.String(),
		Name:      wallet.Name,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		CreatedAt: wallet.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.WithError(err).Error("failed to encode response")
	}
}

func (s *Server) getWallet(w http.ResponseWriter, r *http.Request) {

}
