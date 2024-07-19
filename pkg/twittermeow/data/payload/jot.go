package payload

import "encoding/json"

type JotLoggingCategory string

const (
	JotLoggingCategoryPerftown JotLoggingCategory = "perftown"
)

type JotLogPayload struct {
	Description string `json:"description,omitempty"`
	Product     string `json:"product,omitempty"`
	DurationMS  int64  `json:"duration_ms,omitempty"`
	EventValue  int64  `json:"event_value,omitempty"`
}

func (p *JotLogPayload) ToJSON() ([]byte, error) {
	val := []interface{}{p}
	return json.Marshal(&val)
}

type JotDebugLoggingCategory string

const (
	JotDebugLoggingCategoryClientEvent JotDebugLoggingCategory = "client_event"
)

type JotDebugLogPayload struct {
	Category                          JotDebugLoggingCategory `json:"_category_,omitempty"`
	FormatVersion                     int                     `json:"format_version,omitempty"`
	TriggeredOn                       int64                   `json:"triggered_on,omitempty"`
	Items                             []any                   `json:"items,omitempty"`
	EventNamespace                    EventNamespace          `json:"event_namespace,omitempty"`
	ClientEventSequenceStartTimestamp int64                   `json:"client_event_sequence_start_timestamp,omitempty"`
	ClientEventSequenceNumber         int                     `json:"client_event_sequence_number,omitempty"`
	ClientAppID                       string                  `json:"client_app_id,omitempty"`
}

type EventNamespace struct {
	Page   string `json:"page,omitempty"`
	Action string `json:"action,omitempty"`
	Client string `json:"client,omitempty"`
}
