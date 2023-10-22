package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")

	if PORT == "" {
		log.Fatal("PORT is not found in Env")
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

	router.Mount("/v1", routerv1)

	log.Printf("Server running at Port: %v", PORT)
	serve := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	err := serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
