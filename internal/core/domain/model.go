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
	// PresenceUpdateKind is a join or leave message
	PresenceUpdateKind = 2
	// ChatHistoryKind is a list of messages from the chat
	ChatHistoryKind = 3
	// PresenceListKind is a message with the room members presence
	PresenceListKind = 4
)

// Notication represents a general notification message
type Notication struct {
	Kind    NoticationKind
	To      string
	Payload []byte
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
