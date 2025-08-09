package controllers

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"razi/config"
	"razi/middleware"
	"strconv"
	"testing"

	"razi/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	os.Setenv("JWT_Secret", "testsecret")

	config.InitializeDB()
	router := gin.Default()
	auth := router.Group("/password")
	auth.Use(middleware.JWTAuthMiddleware())

	auth.POST("/create", StorePassword)
	auth.DELETE("/delete/:id", DeletePassword)
	auth.GET("/all", GetAllPasswords)
	auth.PUT("/update/:id", UpdatePassword)
	return router
}

// 🔹 Get latest password ID for specific user_id
func GetID(t *testing.T, userID int) int {
	fmt.Printf("\n[TEST] Fetching latest password ID for user_id=%d\n", userID)
	var id int
	query := "SELECT TOP 1 id FROM passwords WHERE userid=@userid ORDER BY id DESC"

	row := config.DB.QueryRow(query, sql.Named("userid", userID))
	err := row.Scan(&id)
	if err != nil {
		fmt.Printf("[ERROR] Failed to fetch latest password ID: %v\n", err)
		t.Fatalf("Failed to fetch latest password ID: %v", err)
	}
	fmt.Printf("[TEST] Found password ID: %d for user_id=%d\n", id, userID)
	return id
}

func createTestUser(t *testing.T) int {
	fmt.Println("[TEST] Creating test user...")
	var userID int
	err := config.DB.QueryRow(`
		INSERT INTO users (username, email, password)
		OUTPUT INSERTED.id
		VALUES (@p1, @p2, @p3)
	`, "testuser", "Testemail@gmail.com", "hashedpassword").Scan(&userID)

	if err != nil {
		t.Fatalf("[ERROR] Failed to create test user: %v", err)
	}
	fmt.Printf("[TEST] Created test user with ID=%d\n", userID)
	return userID
}

func deleteTestUser(t *testing.T, userID int) {
	fmt.Printf("[TEST] Cleaning up test user ID=%d\n", userID)
	_, err := config.DB.Exec("DELETE FROM users WHERE id=@userid", sql.Named("userid", userID))
	if err != nil {
		t.Fatalf("[ERROR] Failed to delete test user ID=%d: %v", userID, err)
	}
	fmt.Printf("[TEST] Deleted test user ID=%d\n", userID)
}

func TestPasswordCRUDIntegration(t *testing.T) {
	router := setupTestRouter()

	// 1️⃣ Create user
	testUserID := createTestUser(t)

	// 2️⃣ Generate token
	token, _ := utils.GenerateToken(testUserID, "testuser")
	fmt.Printf("[TEST] Generated JWT Token for user_id=%d: %s\n", testUserID, token)

	// 3️⃣ CREATE password
	fmt.Println("[TEST] Creating password entry...")
	createBody := bytes.NewBufferString(`platform_name=TestPlatform&account_email=test@example.com&password=pass123`)
	req, _ := http.NewRequest("POST", "/password/create", createBody)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("[TEST] Create response: %d %s\n", w.Code, w.Body.String())
	assert.Equal(t, http.StatusCreated, w.Code, "Create should return 201")

	// 4️⃣ GET ID from DB
	passwordID := GetID(t, testUserID)

	// 5️⃣ UPDATE password
	fmt.Println("[TEST] Updating password entry...")
	updateBody := bytes.NewBufferString(`password=newpass456`)
	req3, _ := http.NewRequest("PUT", "/password/update/"+strconv.Itoa(passwordID), updateBody) // ✅ no colon
	req3.Header.Set("Authorization", "Bearer "+token)
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	fmt.Printf("[TEST] Update response: %d %s\n", w3.Code, w3.Body.String())
	assert.Equal(t, http.StatusCreated, w3.Code, "Update should return 201")

	// 6️⃣ DELETE password
	fmt.Println("[TEST] Deleting password entry...")
	req4, _ := http.NewRequest("DELETE", "/password/delete/"+strconv.Itoa(passwordID), nil) // ✅ no colon
	req4.Header.Set("Authorization", "Bearer "+token)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	fmt.Printf("[TEST] Delete response: %d %s\n", w4.Code, w4.Body.String())
	assert.Equal(t, http.StatusOK, w4.Code, "Delete should return 200")

	// 7️⃣ CLEANUP test user
	deleteTestUser(t, testUserID)
}
