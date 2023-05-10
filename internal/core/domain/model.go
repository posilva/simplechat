// Package domain holds the data model structs that are used between
// the different layers
package domain

// Message represents a message sent from a source to a destination
type Message struct {
	Payload string
	From    string
	To      string
}

// ModeratedMessage represents a message that was moderated and can be
// stored and shared with other players
type ModeratedMessage struct {
	Message
	ID              string
	FilteredPayload string
	Level           uint
}
