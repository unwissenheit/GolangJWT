package main

import (
    "net/http"

	"github.com/unwissenheit/GolangJWT/middleware"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // トークンを生成するエンドポイント
    r.POST("/login", func(c *gin.Context) {
        userID := c.PostForm("userID")
        token, err := middleware.GenerateJWT(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"token": token})
    })

    // JWT認証が必要なエンドポイント
    r.GET("/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
        userID := c.MustGet("userID").(string)
        c.JSON(http.StatusOK, gin.H{"message": "Hello, " + userID})
    })

    r.Run(":8000")
}
