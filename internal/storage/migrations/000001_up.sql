-- Создание таблицы для хранения информации о домах
CREATE TABLE houses (
    house_id SERIAL PRIMARY KEY,
    address VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    developer VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_address ON houses (address);

-- Создание таблицы для хранения информации о квартирах
CREATE TABLE flats (
    id SERIAL PRIMARY KEY,
    house_id INT NOT NULL REFERENCES houses(house_id) ON DELETE CASCADE,
    price INT NOT NULL,
    rooms INT NOT NULL,
    status VARCHAR(30) NOT NULL
);

-- Автозаполняем поле updated_at
CREATE FUNCTION updated_houses() RETURNS TRIGGER AS $$
BEGIN
  -- Проверяем наличие соответствия между house_id в таблицах flats и houses
    IF EXISTS (SELECT * FROM houses JOIN flats ON houses.house_id = flats.house_id) THEN
        UPDATE houses
        SET updated_at = NOW()
        WHERE houses.house_id IN (SELECT house_id FROM flats);
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
AFTER INSERT ON flats
FOR EACH ROW
EXECUTE FUNCTION updated_houses();

-- Создание таблицы для хранения информации о пользователях
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(40) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hash_pass VARCHAR(70) NOT NULL,
    role VARCHAR(15) NOT NULL
);

-- -- Создание таблицы для хранения избранных квартир пользователя
-- CREATE TABLE favorite_flats (
--     id SERIAL PRIMARY KEY,
--     user_id INT NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
--     house_id INT NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
--     flat_id INT NOT NULL REFERENCES flats(id) ON DELETE CASCADE
-- );
