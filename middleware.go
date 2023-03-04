package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthCheck(c *gin.Context) (User, bool) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		return User{}, false
	}

	// Check if user exists in database in case they were deleted or the password was changed
	dbUser := User{}
	if err := db.Where(&User{
		Username: user.(User).Username,
		Password: user.(User).Password,
	}).First(&dbUser).Error; err != nil {
		return User{}, false
	}

	return user.(User), true
}

func AuthMiddleware(c *gin.Context) {
	user, authed := AuthCheck(c)
	if !authed {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Set("user", user)
	c.Next()
}
