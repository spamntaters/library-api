-- +goose Up
CREATE TABLE authors (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	bio TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE books (
	id SERIAL PRIMARY KEY,
	title VARCHAR(500) NOT NULL,
	author_id INT NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
	isbn VARCHAR(13),
	isbn13 VARCHAR(13),
	published_date DATE,
	page_count INT,
	description TEXT,
	cover_url TEXT,
	owned BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE series (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	description TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE series_books (
	series_id INT NOT NULL REFERENCES series(id) ON DELETE CASCADE,
	book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
	position INT NOT NULL,
	PRIMARY KEY (series_id, book_id)
);

CREATE TABLE tags (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE book_tags (
	book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
	tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
	PRIMARY KEY (book_id, tag_id)
);

CREATE INDEX idx_books_author_id ON books(author_id);
CREATE INDEX idx_books_owned ON books(owned);
CREATE INDEX idx_books_isbn ON books(isbn);
CREATE INDEX idx_books_isbn13 ON books(isbn13);

-- +goose Down
DROP TABLE IF EXISTS book_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS series_books;
DROP TABLE IF EXISTS series;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS authors;
