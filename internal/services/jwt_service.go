package services

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"backend/internal/models"
)

func GenerateToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
		return "", errors.New("JWT secret not set")
	}
	expiration := time.Now().Add(24 * time.Hour)

	if user.Email == "" || user.Rol == "" || user.Nombre == "" {
		return "", errors.New("email, role and name must be provided")
	}

	claims := &models.JWTClaims{
		Username: user.Nombre,
		Email:    user.Email,
		Rol:      user.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
			Issuer:    "desarrollo-seguro-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	log.Println("Generated token for user:", user.Email)
	return tokenString, nil
}

func JwtAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing",
				"message": "Please provide a valid token",
			})
			c.Abort()
			return 
		}

		tokenparts := strings.Split(authHeader, " ")
		if len(tokenparts) != 2 || tokenparts[0] != "Bearer" {
			log.Println("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
				"message": "Please provide a valid token",
			})
			c.Abort()
			return 
		}

		token := tokenparts[1]
		claims, err := validateJWTToken(token)
		if err != nil {
			log.Println("Token validation error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
				"message": "Please provide a valid token",
			})
			c.Abort()
			return 
		}

		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("rol", claims.Rol)
		c.Set("claims", claims)
		c.Next()
	}
}

func validateJWTToken(tokenString string) (*models.JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*models.JWTClaims); ok && parsedToken.Valid {
		log.Println("Token valid for user:", claims.Email)
		return claims, nil
	}

	log.Println("Invalid token")
	return nil, errors.New("invalid token")
}
