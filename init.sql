-- init.sql

CREATE TABLE IF NOT EXISTS mensajes (
    id SERIAL PRIMARY KEY,
    mensaje VARCHAR(255) NOT NULL
);
