package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"bank-api/internal/utils"
	"errors"
	"fmt"
	"time"
)

type CardService struct {
	cardRepo    *repository.CardRepository
	accountRepo *repository.AccountRepository
}

func NewCardService(cardRepo *repository.CardRepository, accountRepo *repository.AccountRepository) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
	}
}

// GenerateCard создаёт новую виртуальную карту для счёта
func (s *CardService) GenerateCard(accountID int, cardHolder string) (*models.Card, error) {
	// Проверяем, что счёт существует
	accounts, err := s.accountRepo.GetByUserID(accountID) // здесь нужна доработка
	if err != nil || len(accounts) == 0 {
		return nil, errors.New("account not found")
	}

	// Генерируем номер карты по алгоритму Луна
	cardNumber := utils.GenerateCardNumber()
	if !utils.IsValidLuhn(cardNumber) {
		return nil, errors.New("generated invalid card number")
	}

	// Генерируем CVV (3 цифры)
	cvv := fmt.Sprintf("%03d", time.Now().UnixNano()%1000)

	// Срок действия: +5 лет от текущей даты
	expiry := time.Now().AddDate(5, 0, 0).Format("01/06")

	// Шифруем номер карты
	encryptedNumber, err := utils.EncryptAES(cardNumber)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	// Хешируем CVV
	cvvHash, err := utils.HashCVV(cvv)
	if err != nil {
		return nil, fmt.Errorf("cvv hashing failed: %w", err)
	}

	// Вычисляем HMAC для целостности
	hmacSecret := []byte("hmac-secret-key-change-me") // В реальном проекте брать из .env
	hmacSignature := utils.ComputeHMAC(cardNumber, hmacSecret)

	// Сохраняем в БД
	cardDB := &repository.CardDB{
		AccountID:     accountID,
		CardNumberEnc: encryptedNumber,
		CardHolder:    cardHolder,
		ExpiryDate:    expiry,
		CVVHash:       cvvHash,
		HMACSignature: hmacSignature,
	}

	if err := s.cardRepo.Create(cardDB); err != nil {
		return nil, fmt.Errorf("failed to save card: %w", err)
	}

	// Возвращаем модель для клиента
	return &models.Card{
		ID:         cardDB.ID,
		AccountID:  cardDB.AccountID,
		CardNumber: cardNumber,
		CardHolder: cardDB.CardHolder,
		ExpiryDate: cardDB.ExpiryDate,
		CVV:        cvv,
		CreatedAt:  time.Now(), // можно парсить из БД
	}, nil
}

// GetCardsByAccount возвращает все карты счёта (с расшифровкой номера)
func (s *CardService) GetCardsByAccount(accountID int) ([]models.Card, error) {
	cardsDB, err := s.cardRepo.GetByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	var cards []models.Card
	hmacSecret := []byte("hmac-secret-key-change-me")

	for _, c := range cardsDB {
		// Расшифровываем номер карты
		decrypted, err := utils.DecryptAES(c.CardNumberEnc)
		if err != nil {
			continue // или вернуть ошибку
		}

		// Проверяем целостность через HMAC
		if !utils.VerifyHMAC(decrypted, c.HMACSignature, hmacSecret) {
			continue // или вернуть ошибку
		}

		cards = append(cards, models.Card{
			ID:         c.ID,
			AccountID:  c.AccountID,
			CardNumber: decrypted,
			CardHolder: c.CardHolder,
			ExpiryDate: c.ExpiryDate,
			CreatedAt:  time.Now(), // можно парсить
		})
	}
	return cards, nil
}
