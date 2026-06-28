package models

import "time"

type Card struct {
	ID         int       `json:"id"`
	AccountID  int       `json:"account_id"`
	CardNumber string    `json:"card_number"` // только для ответа (расшифрованный)
	CardHolder string    `json:"card_holder"`
	ExpiryDate string    `json:"expiry_date"`
	CVV        string    `json:"cvv,omitempty"` // только при создании
	CreatedAt  time.Time `json:"created_at"`
}
