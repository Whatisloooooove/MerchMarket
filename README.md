# **Merch Market** 🛍️  

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue)](https://golang.org/)  
[![Coverage](https://img.shields.io/badge/Coverage-TODO%25-red)](https://github.com/your-username/merch-shop)  

Сервис для внутреннего магазина мерча компании, где сотрудники могут:  
✅ Покупать мерч компании за виртуальные монеты  
✅ Передавать монеты друг другу в качестве премий или участий в квестах  
✅ Просматривать историю транзакций и покупок

✨ **Каждый новый пользователь получает стартовые 1000 монет** ✨

---

## **🚀 Возможности**  

- **Регистрация и авторизация** (JWT)  
- **Покупка мерча** из каталога 
- **Перевод монет** между сотрудниками  
- **История операций**:  
  - Полученные/отправленные переводы  
  - Список купленных товаров  

---

## **🛠 Технологии**  

- **Язык**: [Go](https://go.dev/)
- **База данных**: PostgreSQL  
- **API**: RESTful, документация через Swagger
- **Аутентификация**: JWT  
- **Тестирование**: Unit, Integration, E2E  
- **Деплой**: Docker + Docker Compose  

---

## **📦 Запуск проекта**  

### **Требования**  
- Установленные [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/)  

### **Инструкция**  
1. Клонируйте репозиторий:  
   ```bash
   git clone https://github.com/Whatisloooooove/MerchMarket
   cd MerchMarket
   ```  
2. Запустите сервис:  
   ```bash
   docker-compose up --build
   ```  
3. Сервис будет доступен на:  
   - **API**: `http://localhost:8080`  

---

## **📚 API Документация**  

Основные эндпоинты:  

| Метод | Путь                       | Описание                         |
| ----- | -------------------------- | -------------------------------- |
| POST  | `/auth/register`           | Регистрация нового сотрудника    |
| POST  | `/auth/login`              | Авторизация (получение JWT)      |
| GET   | `/merch`                   | Список товаров                   |
| POST  | `/merch/buy`               | Покупка товара                   |
| POST  | `/coins/transfer`          | Перевод монет другому сотруднику |
| GET   | `/history/purchase`        | История покупок пользователя     |
| GET   | `/history/transfer`        | История покупок пользователя     |

Пример запроса:  
```bash
curl -X POST http://localhost:8080/auth/register \
  -H  "Content-Type: application/json" \
  -d '{"login":"user", "pass":"secret"}'
```
или вот так
```bash
curl -X GET http://localhost:8080/merch \
  -H "Authorization: <JWT_Token>"
```

Ответом на запрос вернётся json состоящий из:
error_code - код ошибки
message - сообщение об ошибке (успехе)
data - данные в формате задаваемом API

Например:
```
{
	error_code: 200,
	message: "", # пустой либо "Merch list loaded successfully"
	data: [
		{
 		name: "shoes",
 		price: 100,
 		stock: 20
 		},
 		...
	]
}
```
---

## **🧪 Тестирование**  

- **Api-тесты**:  
  ```bash
  go test  -v -count=1 ./test/api_test/...
  ```  
- **service-тесты**:  
  ```bash
  go test  -v -count=1 ./test/service_test/...
  ```  
- **storage-тесты**:  
  ```bash
  go test  -v -count=1 ./test/storage_test/...
  ```  

Покрытие кода: **TO-DO** (проверено `go cover`).   

## **🤝 Состав команды**
| Участник | Github | Роль| 
| ----- | ----------------- | -------------------------------- |
| Морочковский Владислав 🛌 | [Whatisloooooove](https://github.com/Whatisloooooove) | Team Lead and QA Engineer |
| Шалбай Алишер 🥷 | [reshile](https://github.com/reshile)     | API Developer |
| Накорнеева Юлия 🌞 | [nakorneeva](https://github.com/Yulia-Nakorneeva) | Database Engineer |
| Афиф Азиз 🧑‍💻 | [AzizAF1](https://github.com/AzizAF1) | Backend Core Developer |