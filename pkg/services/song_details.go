package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SongDetails struct {
	Text     string `json:"text"`
	FilePath string `json:"filePath"`
}

func FetchSongDetails(group, song string) (*SongDetails, error) {
	api := fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", group, song)

	response, err := http.Get(api)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернуло статус %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	var details SongDetails
	if err := json.Unmarshal(body, &details); err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	return &details, nil
}
