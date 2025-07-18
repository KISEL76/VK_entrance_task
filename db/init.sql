CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  login VARCHAR(32) NOT NULL UNIQUE CHECK (char_length(login) >= 3),
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE goods (
  id SERIAL PRIMARY KEY,
  title VARCHAR(100) NOT NULL CHECK (char_length(title) >= 5),
  description VARCHAR(500),
  image_url TEXT,
  price NUMERIC NOT NULL CHECK (price >= 0),
  author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT NOW()
);