package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"razi/config"
	"razi/models"
	"razi/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	// Validate incoming JSON

	// if err := c.ShouldBindJSON(&user); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// Validating incoming Form data
	user.Username = c.PostForm("username")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	//encryting password
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// Insert into DB
	var userID int
	query := `INSERT INTO users (username, email, password) 
          OUTPUT inserted.id 
          VALUES (@username, @email, @password)`

	err = config.DB.QueryRow(query, sql.Named("username", user.Username), sql.Named("email", user.Email),
		sql.Named("password", string(hashedpass))).Scan(&userID)

	// Query Row returns the inserted id to "userID"

	if err != nil {
		log.Println("DB error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func LoginUser(c *gin.Context) {
	var input models.User

	// Bind JSON input to struct
	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	// 	return
	// }

	input.Username = c.PostForm("username")
	input.Password = c.PostForm("password")

	// Fetch user from database by username
	var user models.User
	err := config.DB.QueryRow(`SELECT id, username, password FROM users WHERE username = @p1`, input.Username).
		Scan(&user.ID, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
		"token": token,
	})
}
