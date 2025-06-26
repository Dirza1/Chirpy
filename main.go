package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.ServeMux{}
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	srv := &http.Server{
		Addr:    ":8090",
		Handler: &mux,
	}

	mux.HandleFunc("/healthz", healthz)

	log.Fatal(srv.ListenAndServe())

}

func healthz(writer http.ResponseWriter, request *http.Request) {
	text := []byte("OK")
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write(text)
}
