package repository

import (
	"bank-api/internal/models"
	"database/sql"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	query := `INSERT INTO transactions (from_account_id, to_account_id, amount, type, description)
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRow(query, tx.FromAccountID, tx.ToAccountID, tx.Amount, tx.Type, tx.Description).
		Scan(&tx.ID, &tx.CreatedAt)
}

func (r *TransactionRepository) GetByAccountID(accountID int) ([]models.Transaction, error) {
	rows, err := r.db.Query(`
        SELECT id, from_account_id, to_account_id, amount, type, description, created_at
        FROM transactions
        WHERE from_account_id = $1 OR to_account_id = $1
        ORDER BY created_at DESC
    `, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		var fromID, toID sql.NullInt64
		if err := rows.Scan(&tx.ID, &fromID, &toID, &tx.Amount, &tx.Type, &tx.Description, &tx.CreatedAt); err != nil {
			return nil, err
		}
		if fromID.Valid {
			id := int(fromID.Int64)
			tx.FromAccountID = &id
		}
		if toID.Valid {
			id := int(toID.Int64)
			tx.ToAccountID = &id
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// CreateTx создаёт транзакцию в рамках SQL-транзакции
func (r *TransactionRepository) CreateTx(tx *sql.Tx, t *models.Transaction) error {
	query := `INSERT INTO transactions (from_account_id, to_account_id, amount, type, description)
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return tx.QueryRow(query, t.FromAccountID, t.ToAccountID, t.Amount, t.Type, t.Description).
		Scan(&t.ID, &t.CreatedAt)
}
