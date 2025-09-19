package model

type UltimaMensagem struct {
	Topic      string `json:"topic"`
	QoS        byte   `json:"qos"`
	Retained   bool   `json:"retained"`
	Timestamp  string `json:"timestamp"`
	PayloadRaw string `json:"payload_raw"`
}

type SimpleStatus struct {
	IsLocked    bool   `json:"isLocked"`
	IsConnected bool   `json:"isConnected"`
	LastUpdate  string `json:"lastUpdate"`
}

type HistoryItem struct {
	ID        string `json:"id"`
	Action    string `json:"action"`   // "locked" | "unlocked"
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`   // "app" | "keypad" | "auto"
}

type HistoryResponse struct {
	History []HistoryItem `json:"history"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
	User    *struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	} `json:"user,omitempty"`
}
