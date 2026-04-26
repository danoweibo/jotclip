package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type QueueClient struct {
	client *asynq.Client
}

func NewQueueClient(redisAddr string) *QueueClient {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &QueueClient{client: client}
}

func (q *QueueClient) Close() {
	q.client.Close()
}

func (q *QueueClient) EnqueueTranscode(ctx context.Context, payload TranscodeVideoPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeTranscodeVideo, data)
	_, err = q.client.EnqueueContext(ctx, task)
	return err
}

func (q *QueueClient) EnqueueAnalyze(ctx context.Context, payload AnalyzeVideoPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeAnalyzeVideo, data)
	_, err = q.client.EnqueueContext(ctx, task)
	return err
}