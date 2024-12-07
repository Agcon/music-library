package models

type UpdateSongRequest struct {
	Group       *string `json:"group,omitempty"`
	Title       *string `json:"title,omitempty"`
	ReleaseDate *string `json:"releaseDate,omitempty"`
	Text        *string `json:"text,omitempty"`
	FilePath    *string `json:"filePath,omitempty"`
}
