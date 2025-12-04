CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    title   TEXT NOT NULL,
    content TEXT
);
