package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/MichalGul/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	fmt.Printf("Hallo Chirpy\n")
	const port = "8080"
	const filepathRoot = "."
	// Load env file
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	currentPlatform := os.Getenv("PLATFORM")
	dbConn, dbErr := sql.Open("postgres", dbURL)
	if dbErr != nil {
		log.Fatalf("Error connecting to database: %v", dbErr)
	}

	dbQueries := database.New(dbConn)
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       currentPlatform,
	}
	httpMultiplexer := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: httpMultiplexer,
	}

	httpMultiplexer.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))) // This will handle all paths with app prefix. Eg. app/assets
	httpMultiplexer.HandleFunc("GET /api/healthz", healthHandler)
	httpMultiplexer.HandleFunc("POST /api/validate_chirp", handleValidation)
	httpMultiplexer.HandleFunc("POST /api/users", apiConfig.handleUsers)

	httpMultiplexer.HandleFunc("GET /admin/metrics", apiConfig.returnServerHitsHandler)
	httpMultiplexer.HandleFunc("POST /admin/reset", apiConfig.resetServerHitsHandler)

	log.Printf("Serving on port 8080\n")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
