package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "vk"
)

var Storage *sql.DB

func New() (*sql.DB, error) {
	const op = "storage.New"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	//db, err := sql.Open("postgres", storagePath)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	Storage = db

	// Создаем таблицу для фильмов
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS films (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		rating INTEGER, 
		release TEXT
	);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем таблицу для актеров
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS actors (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		sex TEXT,
		birthday TEXT
	);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем таблицу для связи фильмов с актерами
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS film_actors (
		film_id INTEGER,
		actor_id INTEGER,
		PRIMARY KEY (film_id, actor_id),
		FOREIGN KEY (film_id) REFERENCES films(id),
		FOREIGN KEY (actor_id) REFERENCES actors(id)
	);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
