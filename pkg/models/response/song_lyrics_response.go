package models

type SongLyricsResponse struct {
	Verses      []string `json:"verses"`
	TotalPages  int      `json:"totalPages"`
	CurrentPage int      `json:"currentPage"`
}
