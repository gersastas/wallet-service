package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gersastas/wallet-service/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
	wallet     map[string]*models.Wallet
	walletsMu  sync.Mutex
}

func New(address string) *Server {
	r := chi.NewRouter()

	s := &Server{
		wallet: make(map[string]*models.Wallet),
	}

	r.Post("/wallets", s.handleCreateWallet)
	r.Get("/wallets/{id}", s.handleGetWallet)

	s.httpServer = &http.Server{
		Addr:    address,
		Handler: r,
	}

	return s
}

func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
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

func (r *WalletRequest) Validate() error {
	if r.UserID == "" {
		return errors.New("user_id is required")
	}
	if _, err := uuid.Parse(r.UserID); err != nil {
		return errors.New("user_id is invalid")
	}
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Currency == "" {
		return errors.New("currency is required")
	}
	return nil
}

type WalletResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		s.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	walletID := uuid.New()
	userUUID, _ := uuid.Parse(req.UserID)

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

	s.sendJSON(w, resp, http.StatusCreated)
}

func (s *Server) handleGetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "id")
	if walletID == "" {
		s.sendError(w, "wallet_id is required", http.StatusBadRequest)
		return
	}

	s.walletsMu.Lock()
	wallet, exists := s.wallet[walletID]
	s.walletsMu.Unlock()

	if !exists {
		s.sendError(w, "wallet not found", http.StatusNotFound)
		return
	}

	if wallet.DeletedAt != nil {
		s.sendError(w, "wallet not found", http.StatusNotFound)
		return
	}

	resp := WalletResponse{
		ID:        wallet.ID.String(),
		UserID:    wallet.UserID.String(),
		Name:      wallet.Name,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		CreatedAt: wallet.CreatedAt,
	}

	s.sendJSON(w, resp, http.StatusOK)
}

func (s *Server) sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Error("failed to encode response")
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (s *Server) sendError(w http.ResponseWriter, message string, status int) {
	s.sendJSON(w, ErrorResponse{Error: message}, status)
}
