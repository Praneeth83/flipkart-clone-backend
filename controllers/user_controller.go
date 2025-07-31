package controllers

import (
	"net/http"
	"test1/config"
	"test1/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	var existingUser models.User
	result := config.DB.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "User already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
	}
	user.Password = string(hashedPassword)

	config.DB.Create(&user)
	return c.JSON(http.StatusCreated, echo.Map{"message": "Signup successful", "user": user})
}

func Login(c echo.Context) error {
	var inputUser models.User
	if err := c.Bind(&inputUser); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	var dbUser models.User
	result := config.DB.Where("email = ?", inputUser.Email).First(&dbUser)
	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Login successful", "user": dbUser})
}
