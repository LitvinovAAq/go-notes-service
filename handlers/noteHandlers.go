package handlers

import (
	"errors"
	"fmt"
	"myproject/models"
	"myproject/service"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

var ErrBadPathID = errors.New("invalid id in path")

func GetNote(s service.NoteService) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err!=nil{
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrInvalidID.Error(), "request_id": rid})
			return 
		}
		note, err := s.GetNote(ctx.Request.Context(), id)
		if err != nil{
			respondWithError(ctx, err)
			return 
		}
		ctx.JSON(http.StatusOK, note)
	}
}

func GetAllNotes(s service.NoteService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		notes, err := s.GetAllNotes(ctx.Request.Context())
		if err != nil{
			respondWithError(ctx, err)
			return 
		}
		if notes == nil {
            notes = []models.Note{}
        }
		ctx.JSON(http.StatusOK, notes)
	}	
}

func CreateNote(s service.NoteService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req NoteRequest 
		if err := ctx.ShouldBindJSON(&req); err != nil{
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON", "request-id": rid})
			return
		}
		id, err := s.CreateNote(ctx.Request.Context(), req.Title, req.Content)
		if err!=nil{
			respondWithError(ctx, err)
			return 
		}
		ctx.Header("Location", fmt.Sprintf("/notes/%d", id))
		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func DeleteNote(s service.NoteService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil{
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrBadPathID.Error(), "request-id": rid})
			return 
		}
		if err := s.DeleteNote(ctx.Request.Context(), id); err != nil{
			respondWithError(ctx, err)
			return 
		}
		ctx.Status(http.StatusNoContent)
	}
}




