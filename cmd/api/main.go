package main

import (
	"bank-api/internal/config"
	"bank-api/internal/handler"
	"bank-api/internal/middleware"
	"bank-api/internal/repository"
	"bank-api/internal/service"
	"bank-api/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Load()
	utils.InitJWT(cfg.JWTSecret)

	// 1. Подключение к БД
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logrus.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		logrus.Fatal("Database not reachable:", err)
	}

	// 2. Репозитории, сервисы, хендлеры
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo)
	accountHandler := handler.NewAccountHandler(accountService)

	cardRepo := repository.NewCardRepository(db)
	cardService := service.NewCardService(cardRepo, accountRepo)
	cardHandler := handler.NewCardHandler(cardService)

	// Email сервис
	emailService := service.NewEmailService(cfg)

	// Транзакции и переводы
	transactionRepo := repository.NewTransactionRepository(db)
	transferService := service.NewTransferService(accountRepo, transactionRepo, emailService)
	transferHandler := handler.NewTransferHandler(transferService)

	// CBR сервис
	cbrService := service.NewCBRService()
	cbrHandler := handler.NewCBRHandler(cbrService)

	//Кредит
	creditRepo := repository.NewCreditRepository(db)
	creditService := service.NewCreditService(creditRepo, accountRepo)
	creditHandler := handler.NewCreditHandler(creditService)

	// 3. Маршрутизатор
	r := mux.NewRouter()

	// 4. Публичные маршруты (без JWT)
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	// 5. Защищённые маршруты (с JWT)
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg))

	// Тестовый пинг
	protected.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "pong",
			"userID":  userID,
		})
	}).Methods("GET")

	// Счета
	protected.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts", accountHandler.GetUserAccounts).Methods("GET")
	protected.HandleFunc("/accounts/{id}/deposit", accountHandler.Deposit).Methods("POST")
	protected.HandleFunc("/accounts/{id}/withdraw", accountHandler.Withdraw).Methods("POST")

	// Карты
	protected.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	protected.HandleFunc("/cards", cardHandler.GetCards).Methods("GET")

	// Переводы
	protected.HandleFunc("/transfer", transferHandler.Transfer).Methods("POST")
	protected.HandleFunc("/transactions", transferHandler.GetTransactions).Methods("GET")

	// CBR
	protected.HandleFunc("/cbr/rate", cbrHandler.GetKeyRate).Methods("GET")

	//Кредит
	protected.HandleFunc("/credits", creditHandler.CreateCredit).Methods("POST")
	protected.HandleFunc("/credits", creditHandler.GetUserCredits).Methods("GET")
	protected.HandleFunc("/credits/{id}/schedule", creditHandler.GetSchedule).Methods("GET")

	// 6. Запуск сервера
	logrus.Info("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
