-- +goose Up
CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,
  username VARCHAR(50) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  role VARCHAR(20) NOT NULL,
  password VARCHAR(255) NOT NULL
);

CREATE TABLE categories (
  kategori_id SERIAL PRIMARY KEY,
  nama_kategori VARCHAR(100) NOT NULL
);

CREATE TABLE tags (
  tag_id SERIAL PRIMARY KEY,
  nama_tag VARCHAR(100) NOT NULL
);

CREATE TABLE articles (
  artikel_id SERIAL PRIMARY KEY,
  judul VARCHAR(200) NOT NULL,
  konten TEXT NOT NULL,
  user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  tanggal_publikasi TIMESTAMPTZ
);

CREATE TABLE artikel_kategori (
  artikel_id INTEGER NOT NULL REFERENCES articles(artikel_id) ON DELETE CASCADE,
  kategori_id INTEGER NOT NULL REFERENCES categories(kategori_id) ON DELETE CASCADE,
  PRIMARY KEY (artikel_id, kategori_id)
);

CREATE TABLE artikel_tag (
  artikel_id INTEGER NOT NULL REFERENCES articles(artikel_id) ON DELETE CASCADE,
  tag_id INTEGER NOT NULL REFERENCES tags(tag_id) ON DELETE CASCADE,
  PRIMARY KEY (artikel_id, tag_id)
);

CREATE TABLE media (
  media_id SERIAL PRIMARY KEY,
  url TEXT NOT NULL,
  artikel_id INTEGER NOT NULL REFERENCES articles(artikel_id) ON DELETE CASCADE
);

CREATE TABLE comments (
  komentar_id SERIAL PRIMARY KEY,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  konten TEXT NOT NULL,
  user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  artikel_id INTEGER NOT NULL REFERENCES articles(artikel_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS media;
DROP TABLE IF EXISTS artikel_tag;
DROP TABLE IF EXISTS artikel_kategori;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;