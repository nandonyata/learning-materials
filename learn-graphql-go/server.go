package main

import (
	"learn-graphql-go/graph"
	"learn-graphql-go/internal/auth"
	database "learn-graphql-go/internal/pkg/db/migrations/psql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = defaultPort
// 	}

// 	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

// 	srv.AddTransport(transport.Options{})
// 	srv.AddTransport(transport.GET{})
// 	srv.AddTransport(transport.POST{})

// 	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

// 	srv.Use(extension.Introspection{})
// 	srv.Use(extension.AutomaticPersistedQuery{
// 		Cache: lru.New[string](100),
// 	})

// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	http.Handle("/query", srv)

// 	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()

	// CORS config
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // or limit to "http://localhost:3000"
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))
	router.Use(auth.Middleware())

	database.InitDB()
	defer database.CloseDB()
	database.Migrate()
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/query").ServeHTTP(c.Writer, c.Request)
	})
	router.POST("/query", func(c *gin.Context) {
		server.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
