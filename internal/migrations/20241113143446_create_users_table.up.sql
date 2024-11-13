CREATE TABLE users (
    user_id     SERIAL PRIMARY KEY,
    email       VARCHAR(30) UNIQUE,
    password    VARCHAR(30),
    ref_id      INT,
    ref_code    VARCHAR(30),
    FOREIGN KEY (ref_id) REFERENCES users (user_id)
);