package models

type AddSongRequest struct {
	Group string `json:"group" binding:"required"`
	Title string `json:"title" binding:"required"`
}
