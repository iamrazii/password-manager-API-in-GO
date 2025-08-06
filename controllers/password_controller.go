package controllers

import (
	"database/sql"
	"net/http"
	"razi/config"
	"razi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StorePassword(c *gin.Context) {

	// Initial verification of user
	userid, found := c.Get("user_id")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var input models.Password
	input.Platform_Name = c.PostForm("platform_name")
	input.Account_Email = c.PostForm("account_email")
	input.Password = c.PostForm("password")

	// Optional: Basic validation
	if input.Platform_Name == "" || input.Account_Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}
	query := `INSERT INTO passwords (userid, password, platform_name, account_email)
          VALUES (@p1, @p2, @p3, @p4)`

	_, err := config.DB.Exec(query, userid, input.Password, input.Platform_Name, input.Account_Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save password"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Password saved successfully"})
}

func UpdatePassword(c *gin.Context) {
	userid, found := c.Get("user_id")
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}

	passwordidstr := c.Param("id")
	passwordID, err := strconv.Atoi(passwordidstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password ID"})
		return
	}

	var newpassword models.UpdatePassword

	newpassword.Password = c.PostForm("password")
	query := `
UPDATE passwords 
SET password = @password 
WHERE id = @passwordID AND user_id = @userID
`
	_, err = config.DB.Exec(query,
		sql.Named("password", newpassword.Password),
		sql.Named("passwordID", passwordID),
		sql.Named("userID", userid),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Password changed successfully"})
}

func GetAllPasswords(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	query := `SELECT  platform_name,account_email ,created_at, password FROM passwords WHERE user_id = @userID`

	// returns *sql.rows (returned rows from query)
	rows, err := config.DB.Query(query, sql.Named("userID", userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch passwords"})
		return
	}

	defer rows.Close() // clean up

	// Iterate over rows and build the result slice
	var passwords []models.Password
	for rows.Next() {
		var p models.Password
		err := rows.Scan(&p.Id, &p.Platform_Name, &p.Account_Email, &p.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading password data"})
			return
		}
		passwords = append(passwords, p)
	}

	c.JSON(http.StatusOK, gin.H{"passwords": passwords})
}

func DeletePassword(c *gin.Context) {

	userID, found := c.Get("user_id") // Verify User
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to access"})
		return
	}

	passwordIDstr := c.Param("id")
	passwordID, err := strconv.Atoi(passwordIDstr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password format"})
		return
	}

	query := `Delete from passwords where id =@id and userid=@userid`
	result, err := config.DB.Exec(query, sql.Named("id", passwordID), sql.Named("userid", userID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete the password"})
	}

	affectedrows, _ := result.RowsAffected()
	if affectedrows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No such password found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password deleted successfully"})
}
