package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sowjanyac93/bestWallet/modules/initializers"
	"github.com/sowjanyac93/bestWallet/modules/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
		Address   string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to encrypt Password",
		})
		return
	}

	user := models.User{FirstName: body.FirstName, LastName: body.LastName, Email: body.Email, Password: string(hash)}
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

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read input data",
		})
		return
	}

	// Check whether email and password exists
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
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
	// Generate new JWT token with expired timer(negative time)
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
	var user models.User
	c.BindJSON(&user)
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Error": "Failed to create user",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"message":   "User created successfully",
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
		c.JSON(http.StatusOK, &users)
	}
}

func DeleteUser(c *gin.Context) {
	var user models.User
	userId := c.Param("id")
	result := initializers.DB.Where("id = ?", userId).Delete(&user)
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
