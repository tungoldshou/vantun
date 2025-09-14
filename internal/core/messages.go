package core

// MessageType defines the type of a control message.
type MessageType uint8

const (
	// SessionInit is sent by the client to initiate a session.
	SessionInit MessageType = 0x01
	// SessionAccept is sent by the server to accept a session.
	SessionAccept MessageType = 0x02
	// StreamType is sent on a stream to identify its type.
	StreamType MessageType = 0x03
)

// Message represents a control message exchanged during session negotiation.
type Message struct {
	Type MessageType
	Data []byte // CBOR-encoded payload
}

// SessionInitPayload represents the payload for a SessionInit message.
type SessionInitPayload struct {
	// Version is the protocol version.
	Version uint16
	// Token is an optional authentication token.
	Token []byte
	// SupportedFeatures is a list of features the client supports.
	SupportedFeatures []string
}

// SessionAcceptPayload represents the payload for a SessionAccept message.
type SessionAcceptPayload struct {
	// Accepted indicates if the session was accepted.
	Accepted bool
	// Reason is an optional reason for rejection.
	Reason string
	// ServerFeatures is a list of features the server supports.
	ServerFeatures []string
}

// StreamTypePayload represents the payload for a StreamType message.
type StreamTypePayload struct {
	// Type is the type of the stream.
	Type uint8
}