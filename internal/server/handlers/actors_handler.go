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

func ActorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		FindActor(w, r)
	case http.MethodDelete:
		DeleteActor(w, r)
	case http.MethodPatch:
		UpdateActor(w, r)
	}
}

// @Summary Получить информацию об актере по ID
// @Description Получение информации об актере по его идентификатору
// @Tags actors
// @Accept json
// @Produce json
// @Param id path integer true "ID актера"
// @Success 200 {object} models.Actor
// @Failure 400 {string} string "Missing actor ID or invalid actor ID"
// @Failure 404 {string} string "Actor not found"
// @Failure 500 {string} string "Internal server error"
// @Router /actor/{id} [get]
func FindActor(w http.ResponseWriter, r *http.Request) {
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

	film, err := postgres.FindActor(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find actor: %v", err), http.StatusInternalServerError)
		return
	}

	filmsJSON, err := json.Marshal(film)
	if err != nil {
		http.Error(w, "Failed to marshal actor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(filmsJSON)
}

// @Summary Добавить нового актера
// @Description Добавление нового актера в базу данных
// @Tags actors
// @Accept json
// @Produce json
// @Param actor body models.Actor true "Информация о новом актере"
// @Success 201 {string} string "Actor created"
// @Failure 400 {string} string "Failed to parse request body"
// @Failure 500 {string} string "Internal server error"
// @Router /actor [post]
func AddActorHandler(w http.ResponseWriter, r *http.Request) {
	var newActor models.Actor
	err := json.NewDecoder(r.Body).Decode(&newActor)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	err = postgres.AddActor(newActor)
	if err != nil {
		http.Error(w, "Failed to add actor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Получить список всех актеров
// @Description Получение списка всех актеров из базы данных
// @Tags actors
// @Accept json
// @Produce json
// @Success 200 {array} models.Actor
// @Failure 500 {string} string "Internal server error"
// @Router /actors [get]
func ActorsHandler(w http.ResponseWriter, r *http.Request) {
	actors, err := postgres.GetAllActors()
	if err != nil {
		http.Error(w, "Failed to get actors", http.StatusInternalServerError)
		return
	}

	if actors == nil {
		actors = []models.Actor{}
	}

	actorsJSON, err := json.Marshal(actors)
	if err != nil {
		http.Error(w, "Failed to marshal actors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(actorsJSON)
}

// @Summary Удалить актера по ID
// @Description Удаление актера из базы данных по его идентификатору
// @Tags actors
// @Accept json
// @Produce json
// @Param id path integer true "ID актера"
// @Success 200 {string} string "Actor deleted"
// @Failure 400 {string} string "Invalid actor ID"
// @Failure 500 {string} string "Internal server error"
// @Router /actor/{id} [delete]
func DeleteActor(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id_str := parts[len(parts)-1]

	if id_str == "" {
		http.Error(w, "Missing actor ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	err = postgres.DeleteActor(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete actor: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Обновить информацию об актере
// @Description Обновление информации об актере в базе данных по его идентификатору
// @Tags actors
// @Accept json
// @Produce json
// @Param id path integer true "ID актера"
// @Param actor body models.Actor true "Информация об обновленном актере"
// @Success 200 {string} string "Actor updated"
// @Failure 400 {string} string "Invalid actor ID or failed to decode request body"
// @Failure 500 {string} string "Internal server error"
// @Router /actor/{id} [patch]
func UpdateActor(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id_str := parts[len(parts)-1]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	var updatedActor models.Actor
	err = json.NewDecoder(r.Body).Decode(&updatedActor)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = postgres.UpdateActor(id, updatedActor)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update film: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Получить список всех фильмов, в которых участвовал актер
// @Description Получение списка всех фильмов, в которых участвовал актер, по его идентификатору
// @Tags film_actors
// @Accept json
// @Produce json
// @Param id path integer true "ID актера"
// @Success 200 {array} models.Film
// @Failure 400 {string} string "Invalid actor ID"
// @Failure 500 {string} string "Internal server error"
// @Router /film_actors/{id} [get]
func FindActorsFilm(w http.ResponseWriter, r *http.Request) {
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
	actors, err := postgres.GetActorsByFilmID(id)
	if err != nil {
		http.Error(w, "Failed to get actors", http.StatusInternalServerError)
		return
	}

	// Отправляем актеров в формате JSON
	actorsJSON, err := json.Marshal(actors)
	if err != nil {
		http.Error(w, "Failed to marshal actors", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(actorsJSON)

}
