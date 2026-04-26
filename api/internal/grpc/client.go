package engine

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EngineClient struct {
	client EngineServiceClient
	conn   *grpc.ClientConn
}

func NewEngineClient(addr string) *EngineClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to engine: %v", err)
	}

	return &EngineClient{
		client: NewEngineServiceClient(conn),
		conn:   conn,
	}
}

func (e *EngineClient) Close() {
	e.conn.Close()
}

func (e *EngineClient) AnalyzeVideo(videoID, videoURL, scriptText, language string) (*AnalyzeVideoResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return e.client.AnalyzeVideo(ctx, &AnalyzeVideoRequest{
		VideoId:    videoID,
		VideoUrl:   videoURL,
		ScriptText: scriptText,
		Language:   language,
	})
}