package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func loginUser(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Invalid form data",
		})
		return
	}

	// Get user from database
	var dbUser User
	if err := db.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Set session
	session := sessions.Default(c)
	session.Set("user", dbUser)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func logoutUser(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/")
}
