package auth

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

func secret() []byte {
    s := os.Getenv("JWT_SECRET")
    if s == "" {
        // только для локальной разработки
        s = "dev-secret"
    }
    return []byte(s)
}

func GenerateToken(userID int) (string, error) {
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret())
}

func ParseToken(tokenStr string) (int, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return secret(), nil
    })
    if err != nil {
        return 0, ErrInvalidToken
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return 0, ErrInvalidToken
    }

    return claims.UserID, nil
}
