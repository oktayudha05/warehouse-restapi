package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct{
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey []byte
func init(){
	err := godotenv.Load(".env")
	if err != nil {
		panic("error load .env")
	}
	secretJWT := os.Getenv("JWT_KEY")
	jwtKey = []byte(secretJWT)
}

func GenerateJwt(username string)(string, error){
	waktuKadaluwarsa := time.Now().Add(30*time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(waktuKadaluwarsa),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func MiddlewareJwt()gin.HandlerFunc{
	return func (c *gin.Context){
		authHeader := c.GetHeader("Authorization")
		if authHeader == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"message": "tidak ada token"})
			c.Abort()
			return
		}
		tokenstring := string(authHeader)
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenstring, claims, func (token *jwt.Token)(interface{}, error){
			return jwtKey, nil
		})
		if err != nil || !token.Valid{
			c.JSON(http.StatusUnauthorized, gin.H{"message": "token tidak valid"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}