package repository

import (
	"bank-api/internal/models"
	"database/sql"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// GetDB возвращает соединение с БД для использования в транзакциях
func (r *AccountRepository) GetDB() *sql.DB {
	return r.db
}

// Create — создание счёта
func (r *AccountRepository) Create(account *models.Account) error {
	query := `INSERT INTO accounts (user_id, balance, currency) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, account.UserID, account.Balance, account.Currency).
		Scan(&account.ID, &account.CreatedAt)
}

// GetByUserID — получение всех счетов пользователя
func (r *AccountRepository) GetByUserID(userID int) ([]models.Account, error) {
	rows, err := r.db.Query(`SELECT id, user_id, balance, currency, created_at FROM accounts WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.Currency, &acc.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

// GetByID получает счёт по ID
func (r *AccountRepository) GetByID(accountID int) (*models.Account, error) {
	var acc models.Account
	query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE id=$1`
	err := r.db.QueryRow(query, accountID).Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.Currency, &acc.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &acc, err
}

// UpdateBalance обновляет баланс счёта
func (r *AccountRepository) UpdateBalance(accountID int, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err := r.db.Exec(query, amount, accountID)
	return err
}

// UpdateBalanceTx обновляет баланс в рамках транзакции
func (r *AccountRepository) UpdateBalanceTx(tx *sql.Tx, accountID int, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err := tx.Exec(query, amount, accountID)
	return err
}

// GetByIDTx получает счёт по ID в рамках транзакции
func (r *AccountRepository) GetByIDTx(tx *sql.Tx, accountID int) (*models.Account, error) {
	var acc models.Account
	query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE id=$1`
	err := tx.QueryRow(query, accountID).Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.Currency, &acc.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &acc, err
}
