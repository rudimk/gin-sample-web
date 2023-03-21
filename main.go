package main

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
)

// simulate some test accounts
var secrets = gin.H{
	"joe":    gin.H{"email": "joe@porter.run"},
	"ivan": gin.H{"email": "ivan@porter.run"},
}

func setupRouter() *gin.Engine {
  r := gin.Default()
  // Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"joe":    "bornintheusa",
		"ivan": "backintheussr",
	}))

  r.GET("/ping", func(c *gin.Context) {
    fmt.Println("REQUEST HEADERS ==> ")
    fmt.Println(c.Request.Header)
    fmt.Println("CLIENT IP IS", c.ClientIP())

    c.String(http.StatusOK, "pong")
  })
  
  r.GET("/ready", func(c *gin.Context) {
    fmt.Println("Received an unauthenticated healthcheck.")
    
    c.String(http.StatusOK, "check!")
  })
  
  authorized.GET("/readyz", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
      fmt.Println("Received an authenicated healthcheck.")
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

  return r
}

func main() {
  r := setupRouter()
  // Listen and Server in 0.0.0.0:8080
  r.Run(":8080")
}
