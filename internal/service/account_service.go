package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"errors"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

// CreateAccount — создание счёта для пользователя
func (s *AccountService) CreateAccount(userID int, currency string) (*models.Account, error) {
	if currency == "" {
		currency = "RUB"
	}
	account := &models.Account{
		UserID:   userID,
		Balance:  0,
		Currency: currency,
	}
	err := s.accountRepo.Create(account)
	return account, err
}

// GetUserAccounts — все счета пользователя
func (s *AccountService) GetUserAccounts(userID int) ([]models.Account, error) {
	return s.accountRepo.GetByUserID(userID)
}

// Deposit пополняет счёт
func (s *AccountService) Deposit(accountID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return s.accountRepo.UpdateBalance(accountID, amount)
}

// Withdraw снимает средства со счёта
func (s *AccountService) Withdraw(accountID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Проверяем баланс
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return errors.New("account not found")
	}
	if account.Balance < amount {
		return errors.New("insufficient funds")
	}

	return s.accountRepo.UpdateBalance(accountID, -amount)
}
