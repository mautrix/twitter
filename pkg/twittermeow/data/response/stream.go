package response

type StreamEvent struct {
	Topic   string             `json:"topic,omitempty"`
	Payload StreamEventPayload `json:"payload,omitempty"`
}

type StreamEventPayload struct {
	Config   *ConfigPayload   `json:"config,omitempty"`
	DmTyping *DmTypingPayload `json:"dm_typing,omitempty"`
	DmUpdate *DmUpdatePayload `json:"dm_update,omitempty"`
}

type ConfigPayload struct {
	SessionID             string `json:"session_id"`
	SubscriptionTTLMillis int    `json:"subscription_ttl_millis,omitempty"`
	HeartbeatMillis       int    `json:"heartbeat_millis,omitempty"`
}

type DmUpdatePayload struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
}

type DmTypingPayload = DmUpdatePayload
