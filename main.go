package main

import (
	"crypto/rand"
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	setupDB()

	router := gin.Default()

	// Get store key
	storeKey := []byte(os.Getenv("STORE_KEY"))
	if len(storeKey) == 0 {
		log.Println("No store key found, generating one...")
		storeKey = make([]byte, 32)
		_, err := rand.Read(storeKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Setup sessions
	store := cookie.NewStore(storeKey)
	router.Use(sessions.Sessions("user", store))
	gob.Register(User{})

	// Setup templates
	router.LoadHTMLGlob("templates/*")

	// Auth group with middleware
	authed := router.Group("/")
	authed.Use(AuthMiddleware)

	// Routes
	router.GET("/", homePage)
	authed.GET("/secret", secret)

	// Login
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.POST("/login", loginUser)

	// Register
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	router.POST("/register", registerUser)

	// Logout
	router.GET("/logout", logoutUser)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Home page
func homePage(c *gin.Context) {
	// Check if user is logged in
	user, authed := AuthCheck(c)
	if !authed {
		c.HTML(http.StatusOK, "home.html", nil)
		return
	}

	// User is logged in
	c.HTML(http.StatusOK, "home.html", gin.H{
		"username": user.Username,
	})
}

// Secret page for logged in users only
func secret(c *gin.Context) {
	c.String(http.StatusOK, "This is a secret page for logged in users only")
}
