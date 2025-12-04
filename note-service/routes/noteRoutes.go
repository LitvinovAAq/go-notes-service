package routes

import (
	"myproject/handlers"
	"myproject/midleware"
	"myproject/service"

	"github.com/gin-gonic/gin"
)

func RegisterNoteRoutes(r gin.IRouter, s service.NoteService) {
	// Группа маршрутов, которые требуют авторизации
	auth := r.Group("/")
	auth.Use(midleware.AuthMiddleware())

	auth.GET("/notes", handlers.GetAllNotes(s))
	auth.GET("/notes/:id", handlers.GetNote(s))
	auth.POST("/notes", handlers.CreateNote(s))
	auth.DELETE("/notes/:id", handlers.DeleteNote(s))
	auth.PATCH("/notes/:id", handlers.UpdateNote(s))
}
