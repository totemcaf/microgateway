package entity

import "time"

// Message is the unit of receive / transmit
type Message struct {
	Target      string
	Path        string
	Method      string
	ContentType string
	Payload     []byte
	// TODO Add headers
}

// MessageHistory  tracks the message and its metadata
type MessageHistory struct {
	ID         string
	Message    *Message
	Received   time.Time
	Deliveries []Delivery
}

// NewMessageHistory makes a new MessageHistory
func NewMessageHistory(id string, message *Message) *MessageHistory {
	return &MessageHistory{ID: id, Message: message}
}

// Delivery records each time the message was sent to the target
type Delivery struct {
	Sent               time.Time
	TargetResponseCode int
	TargetResponse     string
}

// Sent records a new instance of sending the message to the target
func (m *MessageHistory) Sent(t time.Time, responseCode int, response string) {
	m.Deliveries = append(m.Deliveries, Delivery{t, responseCode, response})
}
