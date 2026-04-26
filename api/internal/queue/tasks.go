package queue

const (
	TypeTranscodeVideo = "video:transcode"
	TypeAnalyzeVideo   = "video:analyze"
)

type TranscodeVideoPayload struct {
	VideoKey  string `json:"video_key"`
	VideoID   string `json:"video_id"`
	ProjectID string `json:"project_id"`
}

type AnalyzeVideoPayload struct {
	VideoKey   string `json:"video_key"`
	VideoID    string `json:"video_id"`
	ProjectID  string `json:"project_id"`
	ScriptText string `json:"script_text"`
	Language   string `json:"language"`
}