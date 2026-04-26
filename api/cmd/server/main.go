package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	enginegrpc "github.com/danoweibo/jotclip/api/internal/grpc"
	"github.com/danoweibo/jotclip/api/internal/handlers"
	"github.com/danoweibo/jotclip/api/internal/queue"
	"github.com/danoweibo/jotclip/api/internal/queue/workers"
	"github.com/danoweibo/jotclip/api/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()

	// PostgreSQL
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("✅ PostgreSQL connected")

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}
	fmt.Println("✅ Redis connected")

	

	// gRPC Engine client
	engine := enginegrpc.NewEngineClient("localhost:50051")
	defer engine.Close()
	fmt.Println("✅ Engine gRPC client connected")

	// Router
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"jotclip-api"}`))
	})

	// Test gRPC endpoint
	r.Get("/test-engine", func(w http.ResponseWriter, r *http.Request) {
		resp, err := engine.AnalyzeVideo("test-123", "http://example.com/video.mp4", "", "en")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Queue client
	queueClient := queue.NewQueueClient("localhost:6379")
	defer queueClient.Close()
	fmt.Println("✅ Queue client initialized")

	

	// inside main(), after Redis setup:
	r2 := storage.NewR2Client()
	fmt.Println("✅ R2 storage client initialized")

	videoHandler := handlers.NewVideoHandler(r2, queueClient)
	r.Post("/videos/upload", videoHandler.Upload)

	// Start worker server in background
	go workers.StartWorkerServer("localhost:6379", r2)

	port := os.Getenv("PORT")
	fmt.Printf("🚀 Jotclip API running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}