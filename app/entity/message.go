package entity

import "time"

// Message ar the units of receive / transmit
type Message struct {
	ID          string
	Target      string
	Path        string
	Method      string
	ContentType string
	Payload     []byte
	Received    time.Time
	Deliveries  []Delivery
}

// Delivery records each time the message was sent to the target
type Delivery struct {
	Sent           time.Time
	TargetResponse int
}

// Sent records a new instance of sending the message to the target
func (m *Message) Sent(t time.Time, response int) {
	m.Deliveries = append(m.Deliveries, Delivery{t, response})
}
