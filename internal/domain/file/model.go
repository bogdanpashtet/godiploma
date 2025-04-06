package file

import (
	"encoding/json"
	"time"
)

type SessionEvent struct {
	Timestamp  time.Time
	SessionID  string
	UserID     int64
	EventType  string
	Payload    map[string]interface{} // Распарсенный JSON
	RawPayload json.RawMessage        `json:"-"` // Можно хранить и сырой JSON для экономии
	LineNumber int                    `json:"-"` // Для отладки
}
