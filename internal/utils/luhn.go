package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateCardNumber генерирует валидный номер карты (16 цифр, алгоритм Луна)
func GenerateCardNumber() string {
	// Генерируем 15 случайных цифр
	digits := make([]int, 16)

	// Первая цифра: 4 (Visa) — можно менять на 5 (MasterCard)
	digits[0] = 4

	// Заполняем позиции 1..14 случайными цифрами (индексы 1-14)
	for i := 1; i < 15; i++ {
		digits[i] = rand.Intn(10)
	}

	// Вычисляем контрольную сумму для первых 15 цифр
	// Алгоритм Луна: идём справа налево, удваиваем каждую вторую цифру
	sum := 0
	for i := 14; i >= 0; i-- {
		digit := digits[i]
		// Позиция считается с конца: (15 - i) — нечётная = удваиваем
		// В индексах: 14, 12, 10, ... — удваиваем (потому что 15-i нечётное)
		if (15-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	// Контрольная цифра: (10 - sum%10) % 10
	checkDigit := (10 - sum%10) % 10
	digits[15] = checkDigit

	// Собираем строку
	result := ""
	for _, d := range digits {
		result += strconv.Itoa(d)
	}

	return result
}

// IsValidLuhn проверяет номер карты по алгоритму Луна
func IsValidLuhn(cardNumber string) bool {
	if len(cardNumber) != 16 {
		return false
	}

	sum := 0
	// Идём справа налево
	for i := 15; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')
		// Удваиваем каждую вторую цифру, начиная с конца
		if (16-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
