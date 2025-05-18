

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    photo VARCHAR(255), -- stores URL or path to photo
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE IF NOT EXISTS projects (
--     id SERIAL PRIMARY KEY,
--     name TEXT NOT NULL,
--     owner_id INT REFERENCES users (id),
--     created_at TIMESTAMP DEFAULT NOW()
-- );