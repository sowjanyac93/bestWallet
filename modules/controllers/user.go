package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sowjanyac93/bestWallet/modules/initializers"
	"github.com/sowjanyac93/bestWallet/modules/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body models.User

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	if !createUserCommon(c, &body) {
		return
	}

	user := models.User{ID: uuid.New(), Username: body.Username, FirstName: body.FirstName, LastName: body.LastName, Email: body.Email, Password: body.Password, Address: body.Address}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to Create User",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account Created Successfully",
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	// Check whether email and password exist
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID.String() == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid Email or Password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid Email or Password",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	jwtTokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", jwtTokenStr, 7200, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"jwt":     jwtTokenStr,
	})
}

func Logout(c *gin.Context) {
	// Generate new JWT token with expired timer (negative time)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "",
		"exp": time.Now().Add(-time.Minute * 10).Unix(),
	})

	jwtTokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", jwtTokenStr, -600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged Out Successfully",
	})
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "Logged in",
	})
}

func CreateUser(c *gin.Context) {
	var body models.User
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	if !createUserCommon(c, &body) {
		return
	}

	// Use MarshalJSON to convert User to JSON
	userJSON, err := body.MarshalJSON()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to serialize User to JSON",
		})
		return
	}

	user := models.User{ID: uuid.New(), Username: body.Username, FirstName: body.FirstName, LastName: body.LastName, Email: body.Email, Password: body.Password, Address: body.Address}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to create user",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"user":    json.RawMessage(userJSON),
			"message": "User created successfully",
		})
	}
}

func GetUsers(c *gin.Context) {
	users := []models.User{}
	result := initializers.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to fetch users",
		})
	} else {
		// Use Marshal to convert Users to JSON
		var usersJSONArray []json.RawMessage
		for _, user := range users {
			userJSON, err := user.MarshalJSON()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"Error": "Failed to serialize User to JSON",
				})
				return
			}
			usersJSONArray = append(usersJSONArray, json.RawMessage(userJSON))
		}

		c.JSON(http.StatusOK, gin.H{
			"users":   usersJSONArray,
			"message": "Users fetched successfully",
		})
	}
}

func DeleteUser(c *gin.Context) {
	// Get the userID from the URL parameter
	userID := c.Param("id")

	// Get the authenticated user from the context
	authUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized access",
		})
		return
	}

	// Check if the authenticated user is trying to delete their own account
	if authUser.(models.User).ID.String() == userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Cannot delete your own account",
		})
		return
	}

	// Delete the user with the specified userID
	result := initializers.DB.Where("id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to delete user",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "User deleted successfully",
		})
	}
}

// checkEmailExists checks if the given email already exists in the database
func checkEmailExists(c *gin.Context, email string) bool {
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		// Email already exists
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Email already in use. Please choose a different email.",
		})
		return true
	}
	return false
}

// hashPassword hashes the given password using bcrypt
func hashPassword(c *gin.Context, password string) (string, bool) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to encrypt Password",
		})
		return "", false
	}
	return string(hash), true
}

// createUserCommon is a common function for creating a user with shared functionality
func createUserCommon(c *gin.Context, user *models.User) bool {
	if checkEmailExists(c, user.Email) {
		return false
	}

	hashedPassword, success := hashPassword(c, user.Password)
	if !success {
		return false
	}

	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return true
}
