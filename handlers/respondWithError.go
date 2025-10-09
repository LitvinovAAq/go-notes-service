package handlers

import (
	"errors"
	"myproject/internal/logger"
	"myproject/midleware"
	"myproject/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getRequestID(ctx *gin.Context) string{
	if v, ok := ctx.Get(midleware.RequestIDKey); ok{
		rid := v.(string)
		return rid
	}
	return ""
}

func respondWithError(c *gin.Context, err error) {
	rid := getRequestID(c)
	prefix := ""
	if rid != "" {
		prefix = "request_id=" + rid + " "
	}

	switch {
	case errors.Is(err, service.ErrInvalidID):
		logger.Errorf("%sinvalid_id: %v", prefix, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrInvalidID.Error(), "request_id": rid})
		return

	case errors.Is(err, service.ErrNoteNotFound):
		logger.Errorf("%snot_found: %v", prefix, err)
		c.JSON(http.StatusNotFound, gin.H{"error": service.ErrNoteNotFound.Error(), "request_id": rid})
		return

	case errors.Is(err, service.ErrTitleRequired):
		logger.Errorf("%stitle_required: %v", prefix, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTitleRequired.Error(), "request_id": rid})
		return

	case errors.Is(err, service.ErrTitleTooLong):
		logger.Errorf("%stitle_too_long: %v", prefix, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTitleTooLong.Error(), "request_id": rid})
		return

	case errors.Is(err, service.ErrContentTooLong):
		logger.Errorf("%scontent_too_long: %v", prefix, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": service.ErrContentTooLong.Error(), "request_id": rid})
		return
	}

	// log.Printf("%sinternal_error: %v", prefix, err)
	logger.Errorf("%sinternal_error: %v", prefix, err)

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
		"request_id": rid,
	})
}