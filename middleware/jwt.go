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
	Role string `json:"role"`
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

func GenerateJwt(username string, role string)(string, error){
	waktuKadaluwarsa := time.Now().Add(30*time.Minute)
	claims := &Claims{
		Username: username,
		Role: role,
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

func JwtAndAuthorization(roles... string)gin.HandlerFunc{
	return func (c *gin.Context){
		if c.Request.Method == "OPTIONS" {
			// Biarkan request OPTIONS melewati middleware JWT
			// Middleware CORS global harus menanganinya
			c.Next()
			return
		}
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
		c.Set("role", claims.Role)
		if len(roles) > 0 {
			valid := false
			for _, role := range roles{
				if claims.Role == role{
					valid = true
					break
				}
			}
			if !valid {
				c.JSON(http.StatusForbidden, gin.H{"message": "tidak memiliki akses untuk operasi ini"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
