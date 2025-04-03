package response

type InitUploadMediaResponse struct {
	MediaID          int64  `json:"media_id,omitempty"`
	MediaIDString    string `json:"media_id_string,omitempty"`
	ExpiresAfterSecs int    `json:"expires_after_secs,omitempty"`
	MediaKey         string `json:"media_key,omitempty"`
}

type FinalizedUploadMediaResponse struct {
	MediaID          int64          `json:"media_id,omitempty"`
	MediaIDString    string         `json:"media_id_string,omitempty"`
	MediaKey         string         `json:"media_key,omitempty"`
	Size             int            `json:"size,omitempty"`
	ExpiresAfterSecs int            `json:"expires_after_secs,omitempty"`
	Image            Image          `json:"image,omitempty"`
	Video            Video          `json:"video,omitempty"`
	ProcessingInfo   ProcessingInfo `json:"processing_info,omitempty"`
}

type Image struct {
	ImageType string `json:"image_type,omitempty"`
	W         int    `json:"w,omitempty"`
	H         int    `json:"h,omitempty"`
}

type Video struct {
	VideoType string `json:"video_type,omitempty"`
}

type ProcessingState string

const (
	ProcessingStatePending    ProcessingState = "pending"
	ProcessingStateInProgress ProcessingState = "in_progress"
	ProcessingStateSucceeded  ProcessingState = "succeeded"
)

type ProcessingInfo struct {
	State           ProcessingState `json:"state,omitempty"`
	CheckAfterSecs  int             `json:"check_after_secs,omitempty"`
	ProgressPercent int             `json:"progress_percent,omitempty"`
}
