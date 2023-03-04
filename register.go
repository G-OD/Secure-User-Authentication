package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterForm struct {
	Username        string `form:"username"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
}

func registerUser(c *gin.Context) {
	// Get form data
	var form RegisterForm
	if err := c.ShouldBind(&form); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Invalid form data",
		})
		return
	}

	// Check if username passwords match
	if form.Password != form.ConfirmPassword {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Passwords do not match",
		})
		return
	}

	// Hash password with bcrypt for security. Cost of 12 is recommended
	password, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Error creating user",
		})
		return
	}

	// Create user
	if err := db.Create(&User{
		Username: form.Username,
		Password: string(password),
	}).Error; err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Error creating user",
		})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}
