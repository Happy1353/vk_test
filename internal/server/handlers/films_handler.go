package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"vk/internal/models"
	postgres "vk/internal/storage"
)

// @Summary Получить список всех фильмов
// @Description Получение списка всех фильмов
// @Tags films
// @Accept json
// @Produce json
// @Success 200 {array} models.Film
// @Failure 500 {string} string "Internal server error"
// @Router /films [get]
func FilmsHandler(w http.ResponseWriter, r *http.Request) {
	films, err := postgres.GetAllFilms()
	if err != nil {
		http.Error(w, "Failed to get films", http.StatusInternalServerError)
		return
	}

	if films == nil {
		films = []models.Film{}
	}

	filmsJSON, err := json.Marshal(films)
	if err != nil {
		http.Error(w, "Failed to marshal films", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(filmsJSON)
}

// @Summary Добавить новый фильм
// @Description Добавление нового фильма
// @Tags films
// @Accept json
// @Produce json
// @Param film body models.CreateFilm true "Новый фильм"
// @Success 201 {string} string "Film added successfully"
// @Failure 400 {string} string "Failed to parse request body"
// @Failure 500 {string} string "Failed to add film"
// @Router /film [post]
func AddFilmHandler(w http.ResponseWriter, r *http.Request) {
	var newFilm models.CreateFilm
	err := json.NewDecoder(r.Body).Decode(&newFilm)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	err = postgres.AddFilm(newFilm)
	if err != nil {
		http.Error(w, "Failed to add film", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func FilmHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		DeleteFilm(w, r)
	case http.MethodGet:
		FindFilm(w, r)
	case http.MethodPatch:
		UpdateFilm(w, r)
	}
}

// @Summary Получить информацию о фильме по ID
// @Description Получение информации о фильме по его идентификатору
// @Tags films
// @Accept json
// @Produce json
// @Param id path integer true "ID фильма"
// @Success 200 {object} models.Film
// @Failure 400 {string} string "Missing film ID or invalid film ID"
// @Failure 404 {string} string "Film not found"
// @Failure 500 {string} string "Internal server error"
// @Router /film/{id} [get]
func FindFilm(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id_str := parts[len(parts)-1]

	if id_str == "" {
		http.Error(w, "Missing film ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Invalid film ID", http.StatusBadRequest)
		return
	}

	film, err := postgres.FindFilm(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find film: %v", err), http.StatusInternalServerError)
		return
	}

	filmsJSON, err := json.Marshal(film)
	if err != nil {
		http.Error(w, "Failed to marshal films", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(filmsJSON)
}

// @Summary Удалить фильм
// @Description Удаление фильма по его идентификатору
// @Tags films
// @Accept json
// @Produce json
// @Param id path integer true "ID фильма"
// @Success 200 {string} string "Film deleted successfully"
// @Failure 400 {string} string "Missing film ID or invalid film ID"
// @Failure 500 {string} string "Failed to delete film"
// @Router /film/{id} [delete]
func DeleteFilm(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id_str := parts[len(parts)-1]

	if id_str == "" {
		http.Error(w, "Missing film ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Invalid film ID", http.StatusBadRequest)
		return
	}

	err = postgres.DeleteFilm(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete film: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Обновить информацию о фильме
// @Description Обновление информации о фильме по его идентификатору
// @Tags films
// @Accept json
// @Produce json
// @Param id path integer true "ID фильма"
// @Param film body models.Film true "Измененные данные фильма"
// @Success 200 {string} string "Film updated successfully"
// @Failure 400 {string} string "Invalid film ID or failed to decode request body"
// @Failure 500 {string} string "Failed to update film"
// @Router /film/{id} [patch]
func UpdateFilm(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id_str := parts[len(parts)-1]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Invalid film ID", http.StatusBadRequest)
		return
	}

	var updatedFilm models.Film
	err = json.NewDecoder(r.Body).Decode(&updatedFilm)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = postgres.UpdateFilm(id, updatedFilm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update film: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
