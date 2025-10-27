package main

import (
	database "learn-graphql-go-gorm/datalayer"
	"learn-graphql-go-gorm/datalayer/actions"
	"learn-graphql-go-gorm/graph"
	"learn-graphql-go-gorm/graph/directives"
	"learn-graphql-go-gorm/graph/resolvers"

	"learn-graphql-go-gorm/middlewares"
	"learn-graphql-go-gorm/services/product_service"
	"learn-graphql-go-gorm/services/user_service"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	// Create GraphQL config with directive
	config := graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			Auth: directives.Auth,
		},
	}

	// Create GraphQL server with WebSocket support
	srv := handler.New(graph.NewExecutableSchema(config))

	// Add transports
	srv.AddTransport(transport.POST{})    // For queries and mutations
	srv.AddTransport(transport.Websocket{ // For subscriptions
		KeepAlivePingInterval: 10, // Ping every 10 seconds to keep connection alive
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development
				// In production, check against your frontend URL
				return true
			},
		},
	})

	// Add extensions
	srv.Use(extension.Introspection{})

	// Router
	router := gin.Default()
	router.Use(middlewares.CORS)
	router.Use(middlewares.Middleware(userService))

	// Playground
	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/query").ServeHTTP(c.Writer, c.Request)
	})

	// GraphQL endpoint - handles queries, mutations (POST) and subscriptions (GET for WebSocket upgrade)
	router.POST("/query", gin.WrapH(srv))
	router.GET("/query", gin.WrapH(srv)) // NEW: For WebSocket connections

	log.Printf("🚀 Server ready at http://localhost:%s/", port)
	log.Printf("📡 Subscriptions ready at ws://localhost:%s/query", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
