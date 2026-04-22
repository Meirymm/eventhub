-- Users table
CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name    VARCHAR(50) NOT NULL DEFAULT '',
    last_name     VARCHAR(50) NOT NULL DEFAULT '',
    role          VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Venues table
CREATE TABLE IF NOT EXISTS venues (
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    address  TEXT,
    capacity INT DEFAULT 0
);

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Events table
CREATE TABLE IF NOT EXISTS events (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    venue_id    INT REFERENCES venues(id) ON DELETE SET NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    start_time  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tickets table
CREATE TABLE IF NOT EXISTS tickets (
    id         SERIAL PRIMARY KEY,
    event_id   INT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    price      NUMERIC(10,2) DEFAULT 0,
    qr_code    TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed data for demo
INSERT INTO venues (name, address, capacity) VALUES
    ('Алматы Арена', 'пр. Абая 44, Алматы', 10000),
    ('Конгресс-холл', 'ул. Сейфуллина 597, Алматы', 500)
ON CONFLICT DO NOTHING;

INSERT INTO categories (name) VALUES
    ('Концерт'), ('Конференция'), ('Выставка'), ('Спорт')
ON CONFLICT DO NOTHING;