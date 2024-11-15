CREATE TABLE code (
    code_id      SERIAL PRIMARY KEY,
    code_string  VARCHAR(64) UNIQUE,
    expires TIMESTAMP
);

CREATE TABLE users (
    user_id     SERIAL PRIMARY KEY,
    email       VARCHAR(50) UNIQUE NOT NULL,
    password    VARCHAR(64) NOT NULL,
    ref_id      INT,
    ref_code    INT,
    FOREIGN KEY (ref_id) REFERENCES users (user_id) ON DELETE SET NULL,
    FOREIGN KEY (ref_code) REFERENCES code (code_id) ON DELETE SET NULL
);