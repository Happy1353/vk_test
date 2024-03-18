package models

type Actor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
}
