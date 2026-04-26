package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/danoweibo/jotclip/api/internal/queue"
	"github.com/danoweibo/jotclip/api/internal/storage"
	"github.com/hibiken/asynq"
)

type TranscodeWorker struct {
	r2 *storage.R2Client
}

func NewTranscodeWorker(r2 *storage.R2Client) *TranscodeWorker {
	return &TranscodeWorker{r2: r2}
}

func (w *TranscodeWorker) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload queue.TranscodeVideoPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("🎬 Transcoding video: %s", payload.VideoKey)

	// Create temp directory
	tmpDir := filepath.Join(os.TempDir(), payload.VideoID)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Get presigned URL to download the video
	url, err := w.r2.GetPresignedURL(ctx, payload.VideoKey)
	if err != nil {
		return fmt.Errorf("failed to get presigned URL: %w", err)
	}

	// Download video to temp file
	inputPath := filepath.Join(tmpDir, "input.mp4")
	if err := downloadFile(url, inputPath); err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}

	// Transcode to HLS
	outputPath := filepath.Join(tmpDir, "index.m3u8")
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,
		"-codec:", "copy",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	log.Printf("✅ Transcoding complete for video: %s", payload.VideoID)
	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}