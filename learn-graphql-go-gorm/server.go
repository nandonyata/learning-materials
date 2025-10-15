package main

import (
	database "learn-graphql-go-gorm/datalayer"
	"learn-graphql-go-gorm/datalayer/actions"
	"learn-graphql-go-gorm/graph"
	"learn-graphql-go-gorm/graph/resolvers"

	"learn-graphql-go-gorm/middlewares"
	"learn-graphql-go-gorm/services/product_service"
	"learn-graphql-go-gorm/services/user_service"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

func main() {
	// Database
	database.InitDB()
	defer database.CloseDB()
	database.Migrate()
	db := database.DB

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Actions
	userAction := actions.NewUserAction(db)
	productAction := actions.NewProductAction(db)

	// Services
	userService := user_service.NewService(userAction)
	productService := product_service.NewService(productAction)

	// GraphQL resolver
	resolver := &resolvers.Resolver{
		UserService:    userService,
		ProductService: productService,
	}

	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Router
	router := gin.Default()
	router.Use(middlewares.CORS)
	router.Use(middlewares.Middleware(userService))

	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/query").ServeHTTP(c.Writer, c.Request)
	})
	router.POST("/query", func(c *gin.Context) {
		server.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
