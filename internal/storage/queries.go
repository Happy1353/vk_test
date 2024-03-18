package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"vk/internal/models"
)

func GetAllFilms() ([]models.Film, error) {
	const op = "postgres.GetAllFilms"
	query := "SELECT id, name, description, rating, release FROM films"
	rows, err := Storage.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	defer rows.Close()

	var films []models.Film

	for rows.Next() {
		var film models.Film
		err := rows.Scan(&film.ID, &film.Name, &film.Description, &film.Rating, &film.Release)
		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		films = append(films, film)
	}

	return films, nil
}

func GetAllActors() ([]models.Actor, error) {
	const op = "postgres.GetAllActors"
	query := "SELECT id, name, sex, birthday FROM actors"
	rows, err := Storage.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	defer rows.Close()

	var actors []models.Actor

	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.Birthday)
		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

func AddFilm(film models.CreateFilm) error {
	// Проверяем наличие актеров в базе данных перед добавлением фильма
	for _, actor := range film.Actors {
		actorID, err := getActorID(actor)
		if err != nil {
			return err
		}

		if actorID == 0 {
			return errors.New("actor not found: " + actor)
		}
	}

	// Добавление записи о фильме в таблицу films
	query := "INSERT INTO films (name, description, rating, release) VALUES ($1, $2, $3, $4) RETURNING id"
	var filmID int
	err := Storage.QueryRow(query, film.Name, film.Description, film.Rating, film.Release).Scan(&filmID)
	if err != nil {
		return fmt.Errorf("failed to add film: %w", err)
	}

	// Связывание актеров с добавленным фильмом
	for _, actor := range film.Actors {
		actorID, err := getActorID(actor)
		if err != nil {
			return err
		}

		// Связывание актера с фильмом в таблице film_actors
		query = "INSERT INTO film_actors (film_id, actor_id) VALUES ($1, $2)"
		_, err = Storage.Exec(query, filmID, actorID)
		if err != nil {
			return fmt.Errorf("failed to link actor with film: %w", err)
		}
	}

	return nil
}

func getActorID(actorName string) (int, error) {
	query := "SELECT id FROM actors WHERE name = $1"
	var actorID int
	err := Storage.QueryRow(query, actorName).Scan(&actorID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Актер не найден
			return 0, nil
		}
		// Произошла ошибка при выполнении запроса
		return 0, fmt.Errorf("failed to get actor ID: %w", err)
	}

	return actorID, nil
}

func DeleteFilm(id int) error {
	const op = "storage.DeleteFilm"
	query := "DELETE FROM films WHERE id = $1"

	_, err := Storage.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func FindFilm(id int) (models.Film, error) {
	const op = "storage.FindFilm"
	query := "SELECT id, name, description, rating, release FROM films WHERE id = $1"

	row, err := Storage.Query(query, id)
	if err != nil {
		return models.Film{}, fmt.Errorf("%s: %w", op, err)
	}
	defer row.Close()

	var film models.Film
	if row.Next() {
		err := row.Scan(&film.ID, &film.Name, &film.Description, &film.Rating, &film.Release)
		if err != nil {
			return models.Film{}, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		return models.Film{}, fmt.Errorf("%s: film not found", op)
	}

	return film, nil
}

func UpdateFilm(id int, updatedFilm models.Film) error {
	const op = "storage.UpdateFilm"
	query := "UPDATE films SET "
	var args []interface{}
	var count int = 1

	if updatedFilm.Name != "" {
		query += "name=$1, "
		args = append(args, updatedFilm.Name)
		count++
	}
	if updatedFilm.Description != "" {
		query += "description=$2, "
		args = append(args, updatedFilm.Description)
		count++
	}
	if updatedFilm.Description != "" {
		query += "rating=$3, "
		args = append(args, updatedFilm.Description)
		count++
	}
	if updatedFilm.Release != "" {
		query += "release=$4, "
		args = append(args, updatedFilm.Release)
		count++
	}

	if count == 1 {
		return errors.New("no fields to update")
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id=$"
	query += strconv.Itoa(count)
	args = append(args, id)

	_, err := Storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func FindActor(id int) (models.Actor, error) {
	const op = "storage.FindActor"
	query := "SELECT id, name, sex, birthday FROM actors WHERE id = $1"
	row, err := Storage.Query(query, id)
	if err != nil {
		return models.Actor{}, fmt.Errorf("%s: %w", op, err)
	}
	defer row.Close()

	var actor models.Actor
	if row.Next() {
		err := row.Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.Birthday)
		if err != nil {
			return models.Actor{}, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		return models.Actor{}, fmt.Errorf("%s: Actor not found", op)
	}

	return actor, nil
}

func AddActor(actor models.Actor) error {
	query := "INSERT INTO actors (name, sex, birthday) VALUES ($1, $2, $3) RETURNING id"
	row := Storage.QueryRow(query, actor.Name, actor.Sex, actor.Birthday)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to add actor: %w", err)
	}

	return nil
}

func DeleteActor(id int) error {
	const op = "storage.DeleteActor"
	query := "DELETE FROM actors WHERE id = $1"

	_, err := Storage.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func UpdateActor(id int, updatedFilm models.Actor) error {
	const op = "storage.UpdateActor"
	query := "UPDATE actors SET "
	var args []interface{}
	var count int = 1

	if updatedFilm.Name != "" {
		query += "name=$1, "
		args = append(args, updatedFilm.Name)
		count++
	}
	if updatedFilm.Sex != "" {
		query += "sex=$2, "
		args = append(args, updatedFilm.Sex)
		count++
	}
	if updatedFilm.Birthday != "" {
		query += "birthday=$3, "
		args = append(args, updatedFilm.Birthday)
		count++
	}

	if count == 1 {
		return errors.New("no fields to update")
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id=$"
	query += strconv.Itoa(count)
	args = append(args, id)

	_, err := Storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func GetActorsByFilmID(filmID int) ([]models.Actor, error) {
	const op = "postgres.GetActorsByFilmID"
	query := "SELECT a.id, a.name, a.sex, a.birthday FROM actors a JOIN film_actors fa ON a.id = fa.actor_id WHERE fa.film_id = $1"

	// Выполните запрос к базе данных для извлечения всех актеров фильма
	rows, err := Storage.Query(query, filmID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	// Создайте срез для хранения всех актеров
	var actors []models.Actor

	// Проитерируйтесь по результатам запроса и сканируйте их в структуры Actor
	for rows.Next() {
		var actor models.Actor
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.Birthday)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		actors = append(actors, actor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actors, nil
}
