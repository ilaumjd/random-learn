package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()

	apiCfg := apiConfig{}

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

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	_ = decoder.Decode(&params)

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}
	respBody := returnVals{
		Valid: true,
	}
	respondWithJSON(w, 200, respBody)
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
