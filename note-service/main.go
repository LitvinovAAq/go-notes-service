package main

import (
	"fmt"
	"myproject/db"
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
	srv := service.CreateNoteService(repo)

	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
    	c.JSON(200, gin.H{"status": "ok"})
	})

	r.Use(gin.Recovery())
	r.Use(midleware.RequestID())
	r.Use(midleware.RequestLogger())

	routes.RegisterNoteRoutes(r, srv)

	if err := database.Ping(); err != nil {
		panic("HE удалось подключиться к БД: " + err.Error())
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
