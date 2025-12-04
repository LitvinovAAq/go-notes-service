package handlers

import (
	"errors"
	"fmt"
	"myproject/dto"
	"myproject/midleware"
	"myproject/models"
	"myproject/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrBadPathID = errors.New("invalid id in path")

func GetNote(s service.NoteService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// достаём userID из JWT
		userID, ok := midleware.GetUserID(ctx)
		if !ok || userID <= 0 {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":      "unauthorized",
				"request-id": rid,
			})
			return
		}

		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":      service.ErrInvalidID.Error(),
				"request_id": rid,
			})
			return
		}

		note, err := s.GetNote(ctx.Request.Context(), userID, id)
		if err != nil {
			respondWithError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, note)
	}
}

func GetAllNotes(s service.NoteService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := midleware.GetUserID(ctx)
		if !ok || userID <= 0 {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":      "unauthorized",
				"request-id": rid,
			})
			return
		}

		notes, err := s.GetAllNotes(ctx.Request.Context(), userID)
		if err != nil {
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
		userID, ok := midleware.GetUserID(ctx)
		if !ok || userID <= 0 {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":      "unauthorized",
				"request-id": rid,
			})
			return
		}

		var req dto.NoteRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":      "invalid JSON",
				"request-id": rid,
			})
			return
		}

		id, err := s.CreateNote(ctx.Request.Context(), userID, req.Title, req.Content)
		if err != nil {
			respondWithError(ctx, err)
			return
		}
		ctx.Header("Location", fmt.Sprintf("/notes/%d", id))
		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func DeleteNote(s service.NoteService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := midleware.GetUserID(ctx)
		if !ok || userID <= 0 {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":      "unauthorized",
				"request-id": rid,
			})
			return
		}

		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":      ErrBadPathID.Error(),
				"request-id": rid,
			})
			return
		}

		if err := s.DeleteNote(ctx.Request.Context(), userID, id); err != nil {
			respondWithError(ctx, err)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

func UpdateNote(s service.NoteService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := midleware.GetUserID(ctx)
		if !ok || userID <= 0 {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":      "unauthorized",
				"request-id": rid,
			})
			return
		}

		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":      ErrBadPathID.Error(),
				"request-id": rid,
			})
			return
		}

		var req dto.NoteUpdateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			rid := getRequestID(ctx)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":      "invalid JSON",
				"request-id": rid,
			})
			return
		}

		note, err := s.UpdateNote(ctx.Request.Context(), userID, id, req)
		if err != nil {
			respondWithError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NoteResponse{
			ID:      note.Id,
			Title:   note.Title,
			Content: note.Content,
		})
	}
}
