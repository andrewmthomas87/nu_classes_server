package main

import (
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/andrewmthomas87/northwestern"
	"github.com/andrewmthomas87/northwestern/database"
	"github.com/andrewmthomas87/northwestern/generated"
	"github.com/andrewmthomas87/northwestern/server/auth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"time"
)

type signInData struct {
	AccessCode string `json:"accessCode" binding:"required"`
}

func signInHandler(cookieName string, googlePeople *auth.GooglePeople, auth *auth.AuthToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data signInData
		err := c.BindJSON(&data)
		if err != nil {
			return
		}

		email, err := googlePeople.Me(data.AccessCode)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		tokenString, err := auth.TokenStringForUser(email)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		c.SetCookie(cookieName, tokenString, int(24*time.Hour/time.Second), "/", "", false, false)
	}
}

func authHandler(cookieName string, auth *auth.AuthToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(cookieName)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		email, err := auth.UserFromTokenString(tokenString)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		c.Set("email", email)
	}
}

func graphqlHandler(db *database.Database) gin.HandlerFunc {
	h := handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &northwestern.Resolver{Db: db}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	db, err := database.NewDatabase(viper.GetString("database.user"), viper.GetString("database.password"), viper.GetString("database.host"), viper.GetInt("database.port"), viper.GetString("database.database"))
	if err != nil {
		log.Fatal(err)
	}

	googlePeople := auth.NewGooglePeople()
	auth := auth.NewAuth(viper.GetString("auth.secret"), "HS256")

	router := gin.Default()

	router.POST("/sign-in", signInHandler(viper.GetString("auth.cookieName"), googlePeople, auth))

	authorized := router.Group("/")
	authorized.Use(authHandler(viper.GetString("auth.cookieName"), auth))

	authorized.POST("/query", graphqlHandler(db))
	authorized.GET("/", playgroundHandler())

	log.Fatal(router.Run())
}
