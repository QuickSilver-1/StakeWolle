-- Создание таблицы code для хранения реферальных кодов
CREATE TABLE code (
    code_id      SERIAL PRIMARY KEY,            -- Уникальный идентификатор кода
    code_string  VARCHAR(64) UNIQUE,            -- Строка реферального кода
    expires      TIMESTAMP                      -- Время истечения срока действия кода
);

-- Создание таблицы users для хранения данных пользователей
CREATE TABLE users (
    user_id      SERIAL PRIMARY KEY,            -- Уникальный идентификатор пользователя
    email        VARCHAR(50) UNIQUE NOT NULL,   -- Электронная почта пользователя
    password     VARCHAR(64) NOT NULL,          -- Хэш пароля пользователя
    ref_id       INT,                           -- Идентификатор реферера
    ref_code     INT,                           -- Идентификатор реферального кода, присвоенного пользователю
    FOREIGN KEY (ref_id) REFERENCES users (user_id) ON DELETE SET NULL,  -- Связь с реферером, при его удалении значение устанавливается в NULL
    FOREIGN KEY (ref_code) REFERENCES code (code_id) ON DELETE SET NULL  -- Связь с таблицей id кода, при удалении кода значение устанавливается в NULL
);
