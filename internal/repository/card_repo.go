package repository

import (
	"database/sql"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

type CardDB struct {
	ID            int
	AccountID     int
	CardNumberEnc string
	CardHolder    string
	ExpiryDate    string
	CVVHash       string
	HMACSignature string
	CreatedAt     string
}

func (r *CardRepository) Create(card *CardDB) error {
	query := `INSERT INTO cards (account_id, card_number, card_holder, expiry_date, cvv_hash, hmac_signature)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	return r.db.QueryRow(query,
		card.AccountID,
		card.CardNumberEnc,
		card.CardHolder,
		card.ExpiryDate,
		card.CVVHash,
		card.HMACSignature,
	).Scan(&card.ID, &card.CreatedAt)
}

func (r *CardRepository) GetByAccountID(accountID int) ([]CardDB, error) {
	rows, err := r.db.Query(`SELECT id, account_id, card_number, card_holder, expiry_date, cvv_hash, hmac_signature, created_at
                              FROM cards WHERE account_id=$1`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []CardDB
	for rows.Next() {
		var c CardDB
		if err := rows.Scan(&c.ID, &c.AccountID, &c.CardNumberEnc, &c.CardHolder, &c.ExpiryDate, &c.CVVHash, &c.HMACSignature, &c.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}
