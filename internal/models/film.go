package models

type Film struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Rating      int    `json:"rating"`
	Release     string `json:"release"`
}

type CreateFilm struct {
	Film
	Actors []string `json:"actors"`
}
