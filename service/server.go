package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/voice-of-colombo/service/database"
	"github.com/voice-of-colombo/service/graph"
	"github.com/voice-of-colombo/service/graph/generated"
	"github.com/voice-of-colombo/service/speech"
	"github.com/voice-of-colombo/service/speech/recording"
	"github.com/voice-of-colombo/service/user"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/api/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func setupRouter() {
	router := gin.Default()

	apiGroup := router.Group("/api")
	apiGroup.POST("/login", user.LoginUser)
	apiGroup.POST("/logout", user.LogoutUser)
	apiGroup.POST("/abcd", recording.GetData)

	authorized := router.Group("/api")
	authorized.Use(user.AuthorizeJWT())
	{
		apiGroup.GET("/playground", playgroundHandler())
		apiGroup.POST("/query", graphqlHandler())
		apiGroup.POST("/speech/evaluation/:speechId/upload", speech.SaveEvaluation)
		apiGroup.GET("/speech/evaluation/:speechId/file", speech.DownloadEvaluation)
	}

	router.Run("localhost:8080")
	// recording.GetData(nil)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	loadEnv()
	database.OpenDatabase()
	setupRouter()
}
