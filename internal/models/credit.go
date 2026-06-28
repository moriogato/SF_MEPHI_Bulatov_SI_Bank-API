package models

import "time"

type Credit struct {
	ID              int        `json:"id"`
	UserID          int        `json:"user_id"`
	AccountID       int        `json:"account_id"`
	Amount          float64    `json:"amount"`
	InterestRate    float64    `json:"interest_rate"`
	TermMonths      int        `json:"term_months"`
	MonthlyPayment  float64    `json:"monthly_payment"`
	RemainingAmount float64    `json:"remaining_amount"`
	Status          string     `json:"status"` // active, paid, overdue
	CreatedAt       time.Time  `json:"created_at"`
	ClosedAt        *time.Time `json:"closed_at,omitempty"`
}

type PaymentSchedule struct {
	ID              int        `json:"id"`
	CreditID        int        `json:"credit_id"`
	PaymentNumber   int        `json:"payment_number"`
	PaymentDate     time.Time  `json:"payment_date"`
	PaymentAmount   float64    `json:"payment_amount"`
	PrincipalAmount float64    `json:"principal_amount"`
	InterestAmount  float64    `json:"interest_amount"`
	Status          string     `json:"status"` // pending, paid, overdue
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	PenaltyAmount   float64    `json:"penalty_amount"`
}
