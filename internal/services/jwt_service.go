package services

import (
	"errors"
	"log"
	"net/http"
	"os"
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

	if user.Email == "" || user.Role == "" || user.Nombre == "" {
		return "", errors.New("email, role and name must be provided")
	}

	claims := &models.JWTClaims{
		Username: user.Nombre,
		Email:    user.Email,
		Role:      user.Role,
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
		var token string
		if cookieToken, err := c.Cookie("JWTtoken"); err == nil && cookieToken != "" {
			token = cookieToken
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization token not found",
				"message": "Please provide a valid token by login or refresh your token",
			})
			c.Abort()
			return
		}

		claims, err := validateJWTToken(token)
		if err != nil {
			if err.Error() == "token has expired" {
				log.Println("Token has expired")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Token has expired",
					"message": "Please refresh your token or login again",
				})
				c.Abort()
				return
			}
			log.Println("Token validation error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "Please provide a valid token",
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
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
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*models.JWTClaims); ok && parsedToken.Valid {
		log.Println("Token valid for user:", claims.Email)
		return claims, nil
	}

	log.Println("Invalid token")
	return nil, errors.New("invalid token")
}

func ValidateUserAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if cookieToken, err := c.Cookie("JWTtoken"); err == nil && cookieToken != "" {
			log.Println("Cookie token found:", cookieToken)
			token = cookieToken
		} else {
			log.Println("No token found in cookies")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization token not found",
				"message": "Please provide a valid token by login or refresh your token",
			})
			c.Abort()
			return
		}

		role, err := getRoleFromToken(token)
		if err != nil {
			log.Println("Error getting role from token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "Please provide a valid token",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "You do not have permission to access this resource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getRoleFromToken(tokenString string) (string, error) {
	claims, err := validateJWTToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.Role, nil
}
