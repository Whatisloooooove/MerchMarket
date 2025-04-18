-- Создаем схему для приложения, если она не существует
CREATE SCHEMA IF NOT EXISTS merchshop;

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS merchshop.users (
    user_id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    coins INTEGER NOT NULL DEFAULT 1000 CHECK (coins >= 0)
);

-- Таблица товаров
CREATE TABLE IF NOT EXISTS merchshop.merch (
    merch_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL CHECK (stock >= 0)
);

-- Таблица транзакций между пользователями
CREATE TABLE IF NOT EXISTS merchshop.transactions (
    transaction_id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    amount INTEGER NOT NULL CHECK (amount > 0),
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (sender_id) REFERENCES users(user_id),
    FOREIGN KEY (receiver_id) REFERENCES users(user_id)
);

-- Таблица покупок товаров
CREATE TABLE IF NOT EXISTS merchshop.purchases (
    purchase_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    merch_id INTEGER NOT NULL,
    count INTEGER NOT NULL CHECK (count > 0),
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (merch_id) REFERENCES merch(merch_id)
);

-- Таблица истории изменения баланса монет
CREATE TABLE IF NOT EXISTS merchshop.coinhistory (
    change_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    coins_before INTEGER NOT NULL CHECK (coins_before >= 0),
    coins_after INTEGER NOT NULL CHECK (coins_after >= 0),
    change_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(user_id)
);