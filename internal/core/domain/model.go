// Package domain holds the data model structs that are used between
// the different layers
package domain

// NoticationKind represents the existing kinds of notifications
type NoticationKind uint

const (
	// MessageKind is a chat message
	MessageKind = 0
	// ModeratedMessageKind is a moderated chat message
	ModeratedMessageKind = 1
	// PresenceJoinKind is a join or leave message
	PresenceJoinKind = 2
	// PresenceLeaveKind is a join or leave message
	PresenceLeaveKind = 3
	// ChatHistoryKind is a list of messages from the chat
	ChatHistoryKind = 4
	// PresenceListKind is a message with the room members presence
	PresenceListKind = 5
)

// Notication represents a general notification message
type Notication struct {
	UUID    string
	Payload interface{}
	To      string
	Kind    NoticationKind
}

// PresenceUpdateAction type of PresenceUpdates
type PresenceUpdateAction uint

const (
	// PresenceUpdateJoinAction is Presence Join Action
	PresenceUpdateJoinAction = 0

	// PresenceUpdateLeaveAction is Presence Leave Action
	PresenceUpdateLeaveAction = 1
)

// PresenceUpdate represents a presence update message
type PresenceUpdate struct {
	ID        string
	Action    PresenceUpdateAction
	Timestamp uint64
}

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
