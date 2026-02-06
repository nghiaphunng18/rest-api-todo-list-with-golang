package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/internal/config"
	"todo-api/internal/database"
	"todo-api/internal/handlers"
	appMiddleware "todo-api/internal/middleware"
	"todo-api/internal/repository"
	"todo-api/internal/services"
	"todo-api/internal/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load configuration: %v", err)
	}
	
	// connect to database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	log.Println("Connected to database success")

	// load config jwt
	jwtConfig := config.LoadJWTConfig()
	utils.SetJWTConfig(jwtConfig)

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	todoRepo := repository.NewTodoRepository(db) 
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	// config route
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public Routes
	r.Group(func(r chi.Router) {
        r.Get("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Todo App API v1")
        })

        // authentication endpoints
        r.Post("/register", authHandler.Register)
        r.Post("/login", authHandler.Login)
    })

	// Protected Routes
	r.Group(func(r chi.Router) {
        r.Use(appMiddleware.AuthMiddleware(cfg.JWTSecret)) 

		r.Route("/todos", func(r chi.Router) {
            r.Get("/", todoHandler.GetTodos)
            r.Post("/", todoHandler.CreateTodo)
			r.Get("/{id}", todoHandler.GetTodoByID)
			r.Put("/{id}", todoHandler.UpdateTodo)
			r.Delete("/{id}", todoHandler.DeleteTodo)
        })
    })

	// run server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Starting server on port %s\n", cfg.ServerPort)

	log.Fatal(http.ListenAndServe(serverAddr, r))
}
