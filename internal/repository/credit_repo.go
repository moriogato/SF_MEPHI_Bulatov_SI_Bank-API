package repository

import (
	"bank-api/internal/models"
	"database/sql"
)

type CreditRepository struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{db: db}
}

func (r *CreditRepository) Create(credit *models.Credit) error {
	query := `INSERT INTO credits (user_id, account_id, amount, interest_rate, term_months, monthly_payment, remaining_amount)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`
	return r.db.QueryRow(query, credit.UserID, credit.AccountID, credit.Amount,
		credit.InterestRate, credit.TermMonths, credit.MonthlyPayment, credit.RemainingAmount).
		Scan(&credit.ID, &credit.CreatedAt)
}

func (r *CreditRepository) CreatePaymentSchedule(schedule *models.PaymentSchedule) error {
	query := `INSERT INTO payment_schedules (credit_id, payment_number, payment_date, payment_amount, principal_amount, interest_amount)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRow(query, schedule.CreditID, schedule.PaymentNumber,
		schedule.PaymentDate, schedule.PaymentAmount, schedule.PrincipalAmount, schedule.InterestAmount).
		Scan(&schedule.ID)
}

func (r *CreditRepository) GetByUserID(userID int) ([]models.Credit, error) {
	rows, err := r.db.Query(`SELECT id, user_id, account_id, amount, interest_rate, term_months,
                              monthly_payment, remaining_amount, status, created_at, closed_at
                              FROM credits WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []models.Credit
	for rows.Next() {
		var c models.Credit
		var closedAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.UserID, &c.AccountID, &c.Amount, &c.InterestRate,
			&c.TermMonths, &c.MonthlyPayment, &c.RemainingAmount, &c.Status, &c.CreatedAt, &closedAt); err != nil {
			return nil, err
		}
		if closedAt.Valid {
			c.ClosedAt = &closedAt.Time
		}
		credits = append(credits, c)
	}
	return credits, nil
}

func (r *CreditRepository) GetSchedule(creditID int) ([]models.PaymentSchedule, error) {
	rows, err := r.db.Query(`SELECT id, credit_id, payment_number, payment_date, payment_amount,
                              principal_amount, interest_amount, status, paid_at, penalty_amount
                              FROM payment_schedules WHERE credit_id=$1 ORDER BY payment_number`, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		var paidAt sql.NullTime
		if err := rows.Scan(&s.ID, &s.CreditID, &s.PaymentNumber, &s.PaymentDate,
			&s.PaymentAmount, &s.PrincipalAmount, &s.InterestAmount, &s.Status, &paidAt, &s.PenaltyAmount); err != nil {
			return nil, err
		}
		if paidAt.Valid {
			s.PaidAt = &paidAt.Time
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *CreditRepository) UpdateStatus(creditID int, status string) error {
	_, err := r.db.Exec(`UPDATE credits SET status=$1 WHERE id=$2`, status, creditID)
	return err
}

func (r *CreditRepository) UpdateRemainingAmount(creditID int, amount float64) error {
	_, err := r.db.Exec(`UPDATE credits SET remaining_amount=remaining_amount-$1 WHERE id=$2`, amount, creditID)
	return err
}
