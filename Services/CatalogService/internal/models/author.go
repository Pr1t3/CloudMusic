package models

type Author struct {
	Id     int    `json:"id"`
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
}
