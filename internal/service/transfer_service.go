package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"errors"
	"fmt"
)

type TransferService struct {
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
	emailService    *EmailService
}

func NewTransferService(
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
	emailService *EmailService,
) *TransferService {
	return &TransferService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		emailService:    emailService,
	}
}

// Transfer выполняет перевод между счетами (с SQL-транзакцией)
func (s *TransferService) Transfer(fromAccountID, toAccountID int, amount float64, description string) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if fromAccountID == toAccountID {
		return nil, errors.New("cannot transfer to the same account")
	}

	// Начинаем SQL-транзакцию
	tx, err := s.accountRepo.GetDB().Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Важно: откат при ошибке
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Получаем счёт отправителя
	fromAccount, err := s.accountRepo.GetByIDTx(tx, fromAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get from account: %w", err)
	}
	if fromAccount == nil {
		err = errors.New("from account not found")
		return nil, err
	}

	// Проверяем баланс
	if fromAccount.Balance < amount {
		err = errors.New("insufficient funds")
		return nil, err
	}

	// Получаем счёт получателя
	toAccount, err := s.accountRepo.GetByIDTx(tx, toAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get to account: %w", err)
	}
	if toAccount == nil {
		err = errors.New("to account not found")
		return nil, err
	}

	// Обновляем балансы (в рамках одной SQL-транзакции)
	if err = s.accountRepo.UpdateBalanceTx(tx, fromAccountID, -amount); err != nil {
		return nil, fmt.Errorf("failed to update from account: %w", err)
	}
	if err = s.accountRepo.UpdateBalanceTx(tx, toAccountID, amount); err != nil {
		return nil, fmt.Errorf("failed to update to account: %w", err)
	}

	// Сохраняем транзакцию
	transfer := &models.Transaction{
		FromAccountID: &fromAccountID,
		ToAccountID:   &toAccountID,
		Amount:        amount,
		Type:          "transfer",
		Description:   description,
	}
	if err = s.transactionRepo.CreateTx(tx, transfer); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Сбрасываем err, чтобы defer не делал Rollback
	err = nil

	// Отправляем email-уведомление (если настроено)
	if s.emailService != nil {
		// В реальном проекте здесь нужно получить email пользователя из БД
		// Пока отправляем заглушку
		_ = s.emailService.SendPaymentNotification("user@example.com", amount, description)
	}

	return transfer, nil
}

// GetAccountTransactions возвращает историю операций по счёту
func (s *TransferService) GetAccountTransactions(accountID int) ([]models.Transaction, error) {
	return s.transactionRepo.GetByAccountID(accountID)
}
