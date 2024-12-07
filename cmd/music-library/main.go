package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"music-library/databases"
	_ "music-library/docs"
	"music-library/pkg/handlers"
	"music-library/pkg/logging"
	"music-library/pkg/models"
	"os"
)

// @title       Music Library API
// @version     1.0
// @description This is a RESTful API for a music library.
// @host        localhost:8084
// @BasePath    /
func main() {
	logging.Init()
	if err := databases.Connect(); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := models.MigrateSong(databases.DB); err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	r := gin.Default()

	r.GET("/songs", handlers.GetSongs)
	r.GET("/song/:id", handlers.GetSongText)
	r.POST("/song", handlers.AddSong)
	r.PUT("/song/:id", handlers.UpdateSong)
	r.DELETE("/song/:id", handlers.DeleteSong)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("SERVER_PORT")
	logging.Log.Info("Запуск сервера...")
	err := r.Run(port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
