package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Dirza1/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("error opening database")
		os.Exit(1)
	}

	apiCfg := apiConfig{}
	apiCfg.Queries = database.New(db)
	apiCfg.PLATFORM = platform
	mux := http.ServeMux{}
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	srv := &http.Server{
		Addr:    ":8090",
		Handler: &mux,
	}

	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.reset)
	mux.HandleFunc("POST /api/chirps", chirps)
	mux.HandleFunc("POST /api/users", apiCfg.add_user)

	log.Fatal(srv.ListenAndServe())

}

func healthz(writer http.ResponseWriter, request *http.Request) {
	text := []byte("OK")
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write(text)
}

func chirps(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Chirp string `json:"body"`
		ID    string `json:"user_id"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, 400, "something went wrong")
	}
	validated_Chirp, err := validate_chirp(params.Chirp)
	if err != nil {
		respondWithError(writer, 400, "something went wrong")

	}

	fmt.Println(validated_Chirp)

}

func validate_chirp(chirp string) (string, error) {
	type returnValsTrue struct {
		NewChirp string `json:"cleaned_body"`
	}
	if len(chirp) > 140 {
		return "", errors.New("chirp to long")
	}
	cleanedChirp := checkForProfanity(chirp)
	respBody := returnValsTrue{
		NewChirp: cleanedChirp,
	}
	return respBody.NewChirp, nil
}

func (cfg *apiConfig) metrics(writer http.ResponseWriter, request *http.Request) {
	printValue := []byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load()))
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write(printValue)
}

func (cfg *apiConfig) reset(writer http.ResponseWriter, request *http.Request) {
	if cfg.PLATFORM != "dev" {
		respondWithError(writer, 403, "Forbidden")
		return
	}
	err := cfg.Queries.ResetUserDatabase(request.Context())
	if err != nil {
		respondWithError(writer, 400, "Issue during database reset")
		return
	}
	cfg.fileserverHits.Swap(0)
	writer.WriteHeader(200)
	writer.Write([]byte("Hits reset to 0 and database reset"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, request)
	})

}
func (cfg *apiConfig) add_user(writer http.ResponseWriter, request *http.Request) {
	type incomming struct {
		Email string `json:"email"`
	}
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(request.Body)
	inc := incomming{}
	err := decoder.Decode(&inc)
	if err != nil {
		respondWithError(writer, 500, "Somthing went wrong")
		return
	}

	DBuser, err := cfg.Queries.CreateUser(request.Context(), inc.Email)
	if err != nil {
		respondWithError(writer, 400, "something went wrong with creation of user")
		return
	}
	user := User{
		ID:        DBuser.ID,
		CreatedAt: DBuser.CreatedAt,
		UpdatedAt: DBuser.UpdatedAt,
		Email:     DBuser.Email,
	}
	respondWithJSON(writer, 201, user)

}

type apiConfig struct {
	fileserverHits atomic.Int32
	Queries        *database.Queries
	PLATFORM       string
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnValsFalse struct {
		Error string `json:"error"`
	}
	respBody := returnValsFalse{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func checkForProfanity(chrirp string) string {
	profaneSlice := []string{"kerfuffle", "sharbert", "fornax"}
	lowerdString := strings.ToLower(chrirp)
	splitString := strings.Split(lowerdString, " ")
	for i, word := range splitString {
		if slices.Contains(profaneSlice, word) {
			splitString[i] = "****"
		}
	}
	originalSplit := strings.Split(chrirp, " ")

	for i, word := range splitString {
		if word != "****" {
			splitString[i] = originalSplit[i]
		}
	}

	cleanedChirp := strings.Join(splitString, " ")

	return cleanedChirp
}
