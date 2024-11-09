package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ilaumjd/random-learn/httpserver/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	db, _ := sql.Open("postgres", dbURL)

	_ = database.New(db)

	mux := http.NewServeMux()

	apiCfg := apiConfig{}

	// basic routing
	root := http.FileServer(http.Dir("."))
	fileServerHandler := http.StripPrefix("/app/", root)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(srv.ListenAndServe())
}
