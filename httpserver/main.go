package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/ilaumjd/random-learn/httpserver/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	platform := os.Getenv("PLATFORM")

	dbURL := os.Getenv("DB_URL")
	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		platform: platform,
		db:       dbQueries,
	}

	mux := http.NewServeMux()

	// basic routing
	root := http.FileServer(http.Dir("."))
	fileServerHandler := http.StripPrefix("/app/", root)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("GET /api/healthz", handleHealthz)

	// sql routing
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlePostUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlePostChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
	platform       string
	db             *database.Queries
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"error": "` + msg + `"}`))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(payload)
	w.Write(response)
}
