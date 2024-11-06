package middleware

import (
    "net/http"
    "strings"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/gin-gonic/gin"
)


var jwtKey = []byte(os.Getenv("SECRETKEY"))

// トークンを生成する関数
func GenerateJWT(userID string) (string, error) {
    // トークンの有効期限を設定
    expirationTime := time.Now().Add(5 * time.Minute)

    // トークンに含めるクレーム（ペイロード）を設定
    claims := &jwt.RegisteredClaims{
        Subject:   userID,
        ExpiresAt: jwt.NewNumericDate(expirationTime),
    }

    // トークンを生成
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
    // トークンを解析してクレームを取得
    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }

    return claims, nil
}


func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // リクエストヘッダーからトークンを取得
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // トークンを検証
        claims, err := ValidateJWT(tokenString)
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // クレーム（ユーザー情報など）をコンテキストに保存し、ハンドラーに渡す
        c.Set("userID", claims.Subject)
        c.Next()
    }
}
