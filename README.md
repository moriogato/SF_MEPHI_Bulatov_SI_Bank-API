# Bank API  
REST API для банковского сервиса на Go с поддержкой счетов, карт, переводов, кредитов и аналитики.

## Содержание  

О проекте  
Функционал  
Технологии  
Установка и запуск  
Переменные окружения  
Структура проекта  
API Эндпоинты  
База данных  
Безопасность  
Интеграции  
Тестирование  

## О проекте  

Bank API - это бэкенд-сервис для банковских операций, разработанный на Go. Проект реализует чистую архитектуру с разделением на слои: модели, репозитории, сервисы и обработчики запросов. Сервис поддерживает полный цикл банковских операций с защитой данных и интеграцией с внешними сервисами.  
Проект разработан в рамках учебного задания и демонстрирует следующие подходы:  
* Чистая архитектура  
* JWT-аутентификация  
* Шифрование чувствительных данных  
* SQL-транзакции  
* Работа с внешними API (SOAP, SMTP)  

## Функционал  

### Управление пользователями  

* Регистрация с проверкой уникальности email и username  
* Аутентификация с выдачей JWT-токена (срок действия 24 часа)  

### Банковские счета  

* Создание счетов в рублях  
* Пополнение баланса  
* Снятие средств  
* Просмотр всех счетов пользователя  
* Переводы между счетами (с SQL-транзакциями)  

### Виртуальные карты  
* Генерация номера карты по алгоритму Луна  
* Шифрование номера карты (AES-256)  
* Хеширование CVV (bcrypt)  
* HMAC для проверки целостности  
* Просмотр карт с расшифровкой для владельца  
* Срок действия карты: +5 лет  

### Кредитные операции  

* Оформление кредита с аннуитетными платежами  
* Расчет ежемесячного платежа  
* Генерация графика платежей  
* Просмотр активных кредитов  
* Получение графика платежей по кредиту  
* Автоматический расчет процентов   

### Аналитика  

* Статистика доходов/расходов за месяц  
* Кредитная нагрузка (общая сумма активных кредитов)  

### Интеграции  

* Центральный банк РФ: получение ключевой ставки через SOAP API  
* SMTP: отправка email-уведомлений о платежах  

## Технологии  

|  Компонент            | 	Технология       |     	Версия    |
|-----------------------|--------------------|----------------|
|  Язык	|Go	  | 1.23|
|  База данных	|  PostgreSQL  |	18  |
|Роутинг	|  gorilla/mux	| - |
|Аутентификация	| golang-jwt/jwt |	v5 |
|Шифрование |	bcrypt, AES-256, HMAC-SHA256 |	- |
|Логирование |	logrus |	- |
|Парсинг XML |	beevik/etree |	- |
|SMTP |	go-mail/mail |	v2 |
|Драйвер БД |	lib/pq |	- |  

## Установка и запуск  

### Требования  

* Go 1.23 или выше  
* PostgreSQL 17 или выше  
* Git  
* Docker (опционально)  

### Запуск через Docker  

1 Клонируйте репозиторий:  
```bash
git clone https://github.com/moriogato/SF_MEPHI_Bulatov_SI_Bank-API.git  
cd bank-api  
```
2 Создайте файл .env на основе .env.example и настройте его:  
```bash
cp .env.example .env  
```
3 Запустите контейнеры:  
```bash
docker-compose up -d  
```
4 Проверьте работу:  
```bash
curl http://localhost:8080/ping  
```
Приложение будет доступно по адресу: http://localhost:8080  

### Локальный запуск (без Docker)  

1 Установите PostgreSQL и создайте базу данных:  
```sql
CREATE DATABASE bankdb;  
```
2 Выполните миграции:  
```bash
psql -U postgres -d bankdb -f migrations/001_init.sql  
```
3 Скопируйте и настройте .env:  
```bash
cp .env.example .env  
```
4 Установите зависимости:  
```bash
go mod download  
```
5 Запустите приложение:  
```bash
go run cmd/api/main.go  
```
## Переменные окружения  

Создайте файл .env в корне проекта со следующими переменными:  
|  Переменная | Описание |	Пример значения |
|-------------|----------|------------------|
| DB_HOST |	Хост PostgreSQL |	localhost |
| DB_PORT |	Порт PostgreSQL |	5432 |
| DB_USER	| Пользователь PostgreSQL |	postgres |
| DB_PASSWORD	| Пароль PostgreSQL	| your_password |
| DB_NAME |	Имя базы данных	| bankdb |
| JWT_SECRET |	Секретный ключ для JWT	| your-secret-key |
| SMTP_HOST |	SMTP-сервер	| smtp.yandex.ru
| SMTP_PORT |	SMTP-порт	| 465
| SMTP_USER	| Email отправителя	| your-email@yandex.ru
| SMTP_PASSWORD	| Пароль приложения | your-app-password


## Структура проекта  

```text
bank-api/  
├── cmd/  
│   └── api/  
│       └── main.go                 # Точка входа  
├── internal/  
│   ├── config/                     # Конфигурация  
│   │   └── config.go  
│   ├── handler/                    # HTTP-обработчики  
│   │   ├── account_handler.go  
│   │   ├── analytics_handler.go  
│   │   ├── auth_handler.go  
│   │   ├── card_handler.go  
│   │   ├── cbr_handler.go  
│   │   ├── credit_handler.go  
│   │   └── transfer_handler.go  
│   ├── middleware/                 # Промежуточные слои  
│   │   └── auth.go  
│   ├── models/                     # Модели данных  
│   │   ├── account.go  
│   │   ├── card.go  
│   │   ├── credit.go  
│   │   ├── transaction.go  
│   │   └── user.go  
│   ├── repository/                 # Работа с БД  
│   │   ├── account_repo.go  
│   │   ├── card_repo.go  
│   │   ├── credit_repo.go  
│   │   ├── transaction_repo.go  
│   │   └── user_repo.go  
│   ├── service/                    # Бизнес-логика  
│   │   ├── account_service.go  
│   │   ├── analytics_service.go  
│   │   ├── auth_service.go  
│   │   ├── card_service.go  
│   │   ├── cbr_service.go  
│   │   ├── credit_service.go  
│   │   ├── email_service.go  
│   │   └── transfer_service.go  
│   └── utils/                      # Вспомогательные функции  
│       ├── crypto.go  
│       ├── jwt_utils.go  
│       └── luhn.go  
├── migrations/  
│   └── 001_init.sql  
├── .env.example  
├── .gitignore  
├── Dockerfile  
├── docker-compose.yml  
├── go.mod  
├── go.sum  
└── README.md
```

## API Эндпоинты

### Публичные эндпоинты (без JWT)  

| Метод |	Путь |	Описание |
|-------|------|-----------|
| POST |	/register	| Регистрация нового пользователя|
| POST |	/login | Аутентификация, получение JWT-токена|

### Защищённые эндпоинты (требуют JWT)  

Проверка  
| Метод |	Путь |	Описание |
|-------|------|-----------|
| GET |	/ping |	Проверка работоспособности токена|  

Счета  
| Метод |	Путь |	Описание |
|-------|------|-----------|
| POST |	/accounts |	Создание банковского счёта |
| GET |	/accounts	| Получение всех счетов пользователя |
| POST |	/accounts/{id}/deposit |	Пополнение счёта |
| POST |	/accounts/{id}/withdraw |	Снятие со счёта |

Карты  
| Метод |	Путь |	Описание |
|-------|------|-----------|
| POST	| /cards	| Генерация виртуальной карты |
| GET	| /cards?account_id={id}	| Получение карт счёта |

Переводы   
| Метод |	Путь |	Описание |
|-------|------|-----------|
| POST |	/transfer |	Перевод между счетами |
| GET |	/transactions?account_id={id} |	История операций по счёту |

Кредиты  
| Метод |	Путь |	Описание |
|-------|------|-----------|
| POST |	/credits	| Оформление кредита |
| GET |	/credits	| Получение всех кредитов пользователя |
| GET |	/credits/{id}/schedule	| График платежей по кредиту |

Интеграции  
| Метод |	Путь |	Описание |
|-------|------|-----------|
| GET |	/cbr/rate	| Получение ключевой ставки ЦБ РФ  |

Аналитика  
| Метод |	Путь |	Описание |
|-------|------|-----------|
| GET	| /analytics/monthly?account_id={id}	| Статистика за месяц | 
| GET |	/analytics/credit-load	| Кредитная нагрузка |

### Примеры запросов  
Регистрация:
```bash
curl -X POST http://localhost:8080/register \  
  -H "Content-Type: application/json" \  
  -d '{"username":"ivan","email":"ivan@test.com","password":"12345678"}'  
```
Логин:  
```bash
curl -X POST http://localhost:8080/login \  
  -H "Content-Type: application/json" \  
  -d '{"email":"ivan@test.com","password":"12345678"}'  
```
Создание счёта:  
```bash
curl -X POST http://localhost:8080/accounts \  
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \  
  -H "Content-Type: application/json" \  
  -d '{"currency":"RUB"}'  
```
Перевод:  
```bash
curl -X POST http://localhost:8080/transfer \  
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \  
  -H "Content-Type: application/json" \  
  -d '{"from_account_id":1,"to_account_id":2,"amount":100.50,"description":"Перевод"}'
```
## База данных  

### Схема  

users
| Поле  |	Тип  |	Описание |
|-------|------|-----------|
| id |	SERIAL |	Первичный ключ |
| username |	VARCHAR(50) |	Уникальное имя пользователя |
| email |	VARCHAR(100) |	Уникальный email |
| password_hash |	TEXT	| Хеш пароля (bcrypt) |
| created_at |	TIMESTAMP |	Дата создания |

accounts  
| Поле  |	Тип  |	Описание |
|-------|------|-----------|
| id	| SERIAL| 	Первичный ключ |
| user_id |	INT |	Внешний ключ на users(id) |
| balance |	DECIMAL(15,2)	| Текущий баланс |
| currency |	VARCHAR(3)	| Валюта (RUB) |
| created_at | TIMESTAMP |	Дата создания |

cards  
| Поле  |	Тип  |	Описание |
|-------|------|-----------|
| id	| SERIAL |	Первичный ключ |
| account_id	| INT |	Внешний ключ на accounts(id) |
| card_number	| BYTEA	| Зашифрованный номер карты |
| card_holder	| VARCHAR(100) |	Владелец карты |
| expiry_date |	VARCHAR(5)	| Срок действия (MM/YY) |
| cvv_hash	| TEXT	| Хеш CVV (bcrypt) |
| hmac_signature	| TEXT	| HMAC для проверки целостности|
| created_at	| TIMESTAMP	| Дата создания| 

transactions  
| Поле  |	Тип  |	Описание |
|-------|------|-----------|
| id	| SERIAL |	Первичный ключ |
| from_account_id	| INT	| Внешний ключ на accounts(id) |
| to_account_id |	INT |	Внешний ключ на accounts(id) |
| amount |	DECIMAL(15,2)	| Сумма операции |
| type	| VARCHAR(20)	| Тип операции |
| description	| TEXT |	Описание |
| created_at	| TIMESTAMP |	Дата операции |

credits  
| Поле  |	Тип  |	Описание |
|-------|------|-----------|
| id	| SERIAL	| Первичный ключ |
| user_id	| INT |	Внешний ключ на users(id) |
| account_id	| INT	| Внешний ключ на accounts(id) |
| amount	| DECIMAL(15,2) | Сумма кредита|
| interest_rate	| DECIMAL(5,2) |	Процентная ставка|
| term_months	| INT	| Срок в месяцах |
| monthly_payment	| DECIMAL(15,2) |	Ежемесячный платеж|
| remaining_amount |	DECIMAL(15,2) |	Остаток задолженности |
| status |	VARCHAR(20)	| Статус (active/paid/overdue) |
| created_at |	TIMESTAMP	| Дата оформления |
| closed_at	| TIMESTAMP	| Дата закрытия |

payment_schedules
| Поле  |	Тип  |	Описание |
|-------|------|-----------|
| id	| SERIAL	| Первичный ключ |
| credit_id	| INT	| Внешний ключ на credits(id) |
| payment_number	| INT	| Номер платежа |
| payment_date	| DATE	| Дата платежа |
| payment_amount	| DECIMAL(15,2)	| Сумма платежа |
| principal_amount	| DECIMAL(15,2)	| Основной долг |
| interest_amount	| DECIMAL(15,2)	| Проценты |
| status	| VARCHAR(20)	| Статус (pending/paid/overdue) |
| paid_at	| TIMESTAMP	| Дата оплаты |
| penalty_amount |	DECIMAL(15,2)	| Сумма штрафа |

## Безопасность  
### Шифрование и хеширование
| Данные |	Метод  |	Описание |
|--------|---------|-----------|
| Пароли пользователей | bcrypt	| Хеширование при регистрации |
| Номер карты	| AES-256-GCM	| Шифрование с мастер-ключом |
| CVV	| bcrypt	| Хеширование, невозможность расшифровки |
| Целостность данных	| HMAC-SHA256	| Проверка целостности карт |

### Аутентификация   
* JWT-токены с подписью HS256  
* Срок действия токена: 24 часа  
* Проверка через middleware для защищённых маршрутов  
* Извлечение user_id из контекста запроса  

### Защита данных   
* Параметризованные SQL-запросы (защита от SQL-инъекций)  
* Валидация входных данных  
* Проверка прав доступа к счетам и картам  
* Все секреты хранятся в переменных окружения  

## Интеграции  

### Центральный банк РФ  

Сервис получает актуальную ключевую ставку через SOAP API ЦБ РФ.  
* Эндпоинт: https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx  
* Метод: KeyRate  
* Параметры:  
   fromDate - дата начала периода (YYYY-MM-DD)  
   ToDate - дата окончания периода (YYYY-MM-DD)  
* Ответ: XML с текущей ставкой рефинансирования  
* Маржа банка: +5% к ключевой ставке  

### SMTP-уведомления  
Автоматическая отправка email-уведомлений о платежах.   
Поддерживаемые провайдеры:
* Yandex Mail (рекомендуется)
* Mailgun  
* SendGrid  
* Любой SMTP-сервер   
Формат письма: HTML    
События для уведомлений:   
* Переводы между счетами  
* Пополнение баланса  

## Тестирование  

### Тестирование API через PowerShell   

Регистрация:  
```powershell
$body = @{username="test"; email="test@test.com"; password="12345678"} | ConvertTo-Json  
Invoke-RestMethod -Uri "http://localhost:8080/register" -Method Post -Body $body -ContentType "application/json"  
```
Логин:  
```powershell
$body = @{email="test@test.com"; password="12345678"} | ConvertTo-Json  
$response = Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body $body -ContentType "application/json"
$token = $response.token  
```
Создание счёта:   
```powershell
$headers = @{Authorization = "Bearer $token"}  
$body = @{currency="RUB"} | ConvertTo-Json  
Invoke-RestMethod -Uri "http://localhost:8080/accounts" -Method Post -Body $body -ContentType "application/json" -Headers $headers  
```
Пополнение счёта:  
```powershell
$body = @{amount = 1000} | ConvertTo-Json  
Invoke-RestMethod -Uri "http://localhost:8080/accounts/1/deposit" -Method Post -Body $body -ContentType "application/json" -Headers $headers
```
Перевод:  
```powershell
$body = @{from_account_id=1; to_account_id=2; amount=100.00; description="Test transfer"} | ConvertTo-Json  
Invoke-RestMethod -Uri "http://localhost:8080/transfer" -Method Post -Body $body -ContentType "application/json" -Headers $headers
```
### Тестирование через cURL (Linux/macOS)

Регистрация:  
```bash
curl -X POST http://localhost:8080/register \  
  -H "Content-Type: application/json" \  
  -d '{"username":"test","email":"test@test.com","password":"12345678"}'
```
Логин:  
```bash
TOKEN=$(curl -X POST http://localhost:8080/login \  
  -H "Content-Type: application/json" \  
  -d '{"email":"test@test.com","password":"12345678"}' \  
  | jq -r '.token')
```
Создание счёта:  
```bash
curl -X POST http://localhost:8080/accounts \  
  -H "Authorization: Bearer $TOKEN" \  
  -H "Content-Type: application/json" \  
  -d '{"currency":"RUB"}'
```
Пополнение счёта:  
```bash
curl -X POST http://localhost:8080/accounts/1/deposit \  
  -H "Authorization: Bearer $TOKEN" \  
  -H "Content-Type: application/json" \  
  -d '{"amount":1000}'
```
