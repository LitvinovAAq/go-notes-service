package handlers

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"

    "user-service/auth"
    "user-service/service"
)


func RegisterUser(s service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		user, err := s.RegisterUser(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			switch err {
			case service.ErrEmailRequired,
				service.ErrEmailInvalid,
				service.ErrPasswordRequired,
				service.ErrPasswordTooShort,
				service.ErrEmailAlreadyTaken:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			}
		}

		c.JSON(http.StatusCreated, RegisterResponse{
			ID:    user.Id,
			Email: user.Email,
		})
	}
}

func LoginUser(s service.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
            return
        }

        user, err := s.LoginUser(c.Request.Context(), req.Email, req.Password)
        if err != nil {
            if errors.Is(err, service.ErrInvalidCredentials) {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            return
        }

        token, err := auth.GenerateToken(user.Id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
            return
        }

        c.JSON(http.StatusOK, LoginResponse{Token: token})
    }
}

