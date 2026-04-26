package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/danoweibo/jotclip/api/internal/storage"
)

type VideoHandler struct {
	r2 *storage.R2Client
}

func NewVideoHandler(r2 *storage.R2Client) *VideoHandler {
	return &VideoHandler{r2: r2}
}

func (h *VideoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// 500MB max
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

	// Get presigned URL for immediate playback
	url, err := h.r2.GetPresignedURL(r.Context(), key)
	if err != nil {
		http.Error(w, "Failed to generate URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"key": key,
		"url": url,
	})
}