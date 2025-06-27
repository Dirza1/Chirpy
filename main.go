package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	apiCfg := apiConfig{}
	mux := http.ServeMux{}
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	srv := &http.Server{
		Addr:    ":8090",
		Handler: &mux,
	}

	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /metrics", apiCfg.metrics)
	mux.HandleFunc("POST /reset", apiCfg.reset)

	log.Fatal(srv.ListenAndServe())

}

func healthz(writer http.ResponseWriter, request *http.Request) {
	text := []byte("OK")
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write(text)
}

func (cfg *apiConfig) metrics(writer http.ResponseWriter, request *http.Request) {
	printValue := []byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load()))
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write(printValue)
}

func (cfg *apiConfig) reset(writer http.ResponseWriter, request *http.Request) {
	cfg.fileserverHits.Swap(0)
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, request)
	})

}

type apiConfig struct {
	fileserverHits atomic.Int32
}
