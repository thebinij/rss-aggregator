package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thebinij/rss-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT is not found in Env")
	}

	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("DB_URL is not found in Env")
	}

	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("Failed to connect to Database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	routerv1 := chi.NewRouter()
	routerv1.Get("/healthz", handlerReadiness)
	routerv1.Get("/err", handleErr)
	routerv1.Post("/users", apiCfg.handlerCreateUser)
	routerv1.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	routerv1.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	routerv1.Get("/feeds", apiCfg.handlerGetFeeds)
	routerv1.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	routerv1.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	routerv1.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", routerv1)

	log.Printf("Server running at Port: %v", PORT)
	serve := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
