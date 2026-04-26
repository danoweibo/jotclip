package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/danoweibo/jotclip/api/internal/queue"
	"github.com/danoweibo/jotclip/api/internal/storage"
	"github.com/google/uuid"
)

type VideoHandler struct {
	r2    *storage.R2Client
	queue *queue.QueueClient
}

func NewVideoHandler(r2 *storage.R2Client, q *queue.QueueClient) *VideoHandler {
	return &VideoHandler{r2: r2, queue: q}
}

func (h *VideoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(500 << 20)

	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "No video file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "video/mp4"
	}

	key, err := h.r2.UploadVideo(r.Context(), file, header.Filename, contentType)
	if err != nil {
		http.Error(w, "Failed to upload video", http.StatusInternalServerError)
		return
	}

	videoID := uuid.New().String()

	// Enqueue transcoding job
	err = h.queue.EnqueueTranscode(r.Context(), queue.TranscodeVideoPayload{
		VideoKey:  key,
		VideoID:   videoID,
		ProjectID: "test-project",
	})
	if err != nil {
		http.Error(w, "Failed to enqueue transcode job", http.StatusInternalServerError)
		return
	}

	url, err := h.r2.GetPresignedURL(r.Context(), key)
	if err != nil {
		http.Error(w, "Failed to generate URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"video_id": videoID,
		"key":      key,
		"url":      url,
	})
}