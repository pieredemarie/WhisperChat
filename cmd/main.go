package main

import (
	"log"
	"net/http"
	"os"
	"whisperchat/internal/handler"
	"whisperchat/internal/repository/postgres"
	"whisperchat/internal/room"
	"whisperchat/internal/service"

	"github.com/joho/godotenv"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DB_URL")

	repo, err := postgres.NewPostgresRepo(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer repo.Close()

	manager := room.NewManager()
	svc := service.NewChatService(manager, repo)
	h := handler.NewHandler(svc)

	router := handler.NewRouter(h)

	addr := ":8080" // TODO: add to .env file
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(router)); err != nil {
		log.Fatal(err)
	}
}
