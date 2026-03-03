package database

import (
	"database/sql"

	"github.com/gersastas/wallet-service/internal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, name, balance, currency, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		query,
		wallet.ID,
		wallet.UserID,
		wallet.Name,
		wallet.Balance,
		wallet.Currency,
		wallet.CreatedAt,
		wallet.UpdatedAt,
		wallet.DeletedAt,
	)

	return err
}

func (r *WalletRepository) GetByID(id uuid.UUID) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, name, balance, currency, created_at, updated_at, deleted_at
		FROM wallets
		WHERE id = $1 AND deleted_at IS NULL
	`

	wallet := &models.Wallet{}
	err := r.db.QueryRow(query, id).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Name,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
		&wallet.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return wallet, nil
}
