package main

import (
	"log"
	"net/http"
)

func main() {
	httpServeMux := http.ServeMux{}
	httpServeMux.Handle("/", http.FileServer(http.Dir(".")))
	srv := &http.Server{
		Addr:    ":8090",
		Handler: &httpServeMux,
	}
	log.Fatal(srv.ListenAndServe())

}
