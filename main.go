package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)

	appLog := log.Default()
	appLog.Println(fmt.Sprintf("GOMAXPROCS: %d", runtime.GOMAXPROCS(-1)))

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(NewMetricsMiddleware())

	router.GET("/prometheus", gin.WrapH(promhttp.Handler()))
	router.POST("/jwt/generate", func(ctx *gin.Context) {
		claims := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "test",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

		tokenString, err := token.SignedString(privateKey)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	})

	router.GET("/jwt/validate", func(ctx *gin.Context) {
		h := ctx.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(h, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return privateKey.Public(), nil
		})
		if err != nil {
			log.Println(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.JSON(http.StatusOK, gin.H{
				"result": true,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"result": false,
			})
		}
	})

	router.Run()
}
