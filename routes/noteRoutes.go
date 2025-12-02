package routes

import (
	"myproject/handlers"
	"myproject/service"

	"github.com/gin-gonic/gin"
)

func RegisterNoteRoutes(r gin.IRouter, s service.NoteService) {

	r.GET("/notes", handlers.GetAllNotes(s))
	r.GET("/notes/:id", handlers.GetNote(s))
	r.POST("/notes", handlers.CreateNote(s))
	r.DELETE("/notes/:id", handlers.DeleteNote(s))
	r.PATCH("/notes/:id", handlers.UpdateNote(s))

}
