package main

import (
	"context"
	"fmt"
	"myproject/cache"
	"myproject/db"
	"myproject/events"
	"myproject/midleware"
	"myproject/repository"
	"myproject/routes"
	"myproject/service"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.GetDB()
	if err != nil {
		panic(err)
	}
	defer database.Close()

	repo := repository.CreateNoteRepository(database)
	notesCache := cache.NewNotesCache()
	srv := service.CreateNoteService(repo, notesCache)

	// контекст для Kafka-consumer'а
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// запускаем consumer, который слушает user_registered и создаёт приветственные заметки
	if err := events.RunUserRegisteredConsumer(ctx, srv); err != nil {
		fmt.Println("failed to start Kafka consumer:", err)
	}

	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Use(gin.Recovery())
	r.Use(midleware.RequestID())
	r.Use(midleware.RequestLogger())

	routes.RegisterNoteRoutes(r, srv)

	if err := database.Ping(); err != nil {
		panic("Не удалось подключиться к БД: " + err.Error())
	}

	fmt.Println("Подключение к БД успешно!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
