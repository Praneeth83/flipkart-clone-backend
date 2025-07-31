package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"test1/config"
	"test1/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if user.Role != "buyer" && user.Role != "seller" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Role must be 'buyer' or 'seller'"})
	}

	var existing models.User
	if err := config.DB.Where("email = ? AND role = ?", user.Email, user.Role).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "User already exists"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Password hash error"})
	}
	user.Password = string(hashed)

	config.DB.Create(&user)

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Signup successful",
		"user": echo.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Login(c echo.Context) error {
	var input models.User
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if input.Role != "buyer" && input.Role != "seller" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Role must be 'buyer' or 'seller'"})
	}

	var dbUser models.User
	if err := config.DB.Where("email = ? AND role = ?", input.Email, input.Role).First(&dbUser).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	claims := jwt.MapClaims{
		"user_id": dbUser.ID,
		"email":   dbUser.Email,
		"role":    dbUser.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   signedToken,
		"user": echo.Map{
			"id":    dbUser.ID,
			"name":  dbUser.Name,
			"email": dbUser.Email,
			"role":  dbUser.Role,
		},
	})
}

func AutoLogin(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing token"})
	}

	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or expired token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	role := claims["role"].(string)

	var user models.User
	if err := config.DB.Where("email = ? AND role = ?", email, role).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": echo.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
