package server

import (
	"net/http"
	"vk/internal/server/handlers"
)

func SetupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/api/v1/films", handlers.FilmsHandler)
	router.HandleFunc("/api/v1/film/", handlers.FilmHandler)
	router.HandleFunc("/api/v1/film", handlers.AddFilmHandler)
	router.HandleFunc("/api/v1/film_actors/", handlers.FindActorsFilm)

	router.HandleFunc("/api/v1/actor/", handlers.ActorHandler)
	router.HandleFunc("/api/v1/actor", handlers.AddActorHandler)
	router.HandleFunc("/api/v1/actors", handlers.ActorsHandler)

	return router
}
