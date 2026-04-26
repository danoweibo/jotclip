package workers

import (
	"log"

	"github.com/danoweibo/jotclip/api/internal/queue"
	"github.com/danoweibo/jotclip/api/internal/storage"
	"github.com/hibiken/asynq"
)

func StartWorkerServer(redisAddr string, r2 *storage.R2Client) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypeTranscodeVideo, NewTranscodeWorker(r2).ProcessTask)

	log.Println("✅ Asynq worker server started")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Failed to start worker server: %v", err)
	}
}