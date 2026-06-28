package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"errors"
	"math"
	"time"
)

type CreditService struct {
	creditRepo  *repository.CreditRepository
	accountRepo *repository.AccountRepository
}

func NewCreditService(creditRepo *repository.CreditRepository, accountRepo *repository.AccountRepository) *CreditService {
	return &CreditService{
		creditRepo:  creditRepo,
		accountRepo: accountRepo,
	}
}

// CalculateAnnuityPayment рассчитывает аннуитетный платёж
func (s *CreditService) CalculateAnnuityPayment(amount float64, annualRate float64, months int) float64 {
	monthlyRate := annualRate / 100 / 12
	if monthlyRate == 0 {
		return amount / float64(months)
	}
	factor := math.Pow(1+monthlyRate, float64(months))
	return amount * monthlyRate * factor / (factor - 1)
}

// CreateCredit оформляет кредит
func (s *CreditService) CreateCredit(userID, accountID int, amount float64, termMonths int) (*models.Credit, []models.PaymentSchedule, error) {
	// Получаем ключевую ставку ЦБ (в реальном проекте)
	// Пока используем фиксированную ставку
	annualRate := 21.0 // % (позже заменим на реальную ставку)

	// Проверяем, что счёт принадлежит пользователю
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, nil, err
	}
	if account == nil || account.UserID != userID {
		return nil, nil, errors.New("account not found or not owned by user")
	}

	// Рассчитываем аннуитетный платёж
	monthlyPayment := s.CalculateAnnuityPayment(amount, annualRate, termMonths)

	// Создаём кредит
	credit := &models.Credit{
		UserID:          userID,
		AccountID:       accountID,
		Amount:          amount,
		InterestRate:    annualRate,
		TermMonths:      termMonths,
		MonthlyPayment:  monthlyPayment,
		RemainingAmount: amount,
		Status:          "active",
	}

	if err := s.creditRepo.Create(credit); err != nil {
		return nil, nil, err
	}

	// Генерируем график платежей
	schedules := s.generateSchedule(credit.ID, amount, annualRate, termMonths, monthlyPayment)

	// Сохраняем график
	for _, sch := range schedules {
		if err := s.creditRepo.CreatePaymentSchedule(&sch); err != nil {
			return nil, nil, err
		}
	}

	return credit, schedules, nil
}

// generateSchedule создаёт график платежей
func (s *CreditService) generateSchedule(creditID int, amount, annualRate float64, months int, monthlyPayment float64) []models.PaymentSchedule {
	monthlyRate := annualRate / 100 / 12
	remaining := amount
	var schedules []models.PaymentSchedule

	for i := 1; i <= months; i++ {
		interest := remaining * monthlyRate
		principal := monthlyPayment - interest
		if principal > remaining {
			principal = remaining
			interest = monthlyPayment - principal
		}
		remaining -= principal

		schedules = append(schedules, models.PaymentSchedule{
			CreditID:        creditID,
			PaymentNumber:   i,
			PaymentDate:     time.Now().AddDate(0, i, 0),
			PaymentAmount:   monthlyPayment,
			PrincipalAmount: principal,
			InterestAmount:  interest,
			Status:          "pending",
			PenaltyAmount:   0,
		})

		if remaining <= 0 {
			break
		}
	}
	return schedules
}

// GetUserCredits возвращает кредиты пользователя
func (s *CreditService) GetUserCredits(userID int) ([]models.Credit, error) {
	return s.creditRepo.GetByUserID(userID)
}

// GetSchedule возвращает график платежей по кредиту
func (s *CreditService) GetSchedule(creditID int) ([]models.PaymentSchedule, error) {
	return s.creditRepo.GetSchedule(creditID)
}
