package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"

    "user-service/db"
    "user-service/events"
    "user-service/handlers"
    "user-service/repository"
    "user-service/service"
)

func main() {
    r := gin.Default()

    database, err := db.GetDB()
    if err != nil {
        log.Fatalf("failed to connect to db: %v", err)
    }
    log.Println("connected to users-db")

    kafkaWriter := events.NewUserRegisteredWriter()
    defer kafkaWriter.Close()

    userRepo := repository.NewUserRepository(database)
    userSvc := service.NewUserService(userRepo, kafkaWriter)

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "user-service ok"})
    })

    r.POST("/users/register", handlers.RegisterUser(userSvc))
    r.POST("/auth/login", handlers.LoginUser(userSvc)) // ← добавили

    if err := r.Run(":8082"); err != nil {
        log.Fatalf("failed to run user-service: %v", err)
    }
}
