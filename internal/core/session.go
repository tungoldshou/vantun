package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/quic-go/quic-go"
)

// Session represents a VANTUN session over a QUIC connection.
type Session struct {
	conn quic.Connection
	// telemetryManager manages telemetry collection and reporting for this session.
	telemetryManager *TelemetryManager
}

// Config holds the configuration for a VANTUN session.
type Config struct {
	// Address is the server address to connect to or listen on.
	Address string
	// TLSConfig is the TLS configuration for the QUIC connection.
	TLSConfig *tls.Config
	// IsServer indicates if this session is for a server.
	IsServer bool
}

// NewSession creates a new VANTUN session based on the provided configuration.
// For a client, it establishes a connection to the server.
// For a server, it listens for incoming connections.
func NewSession(ctx context.Context, config *Config) (*Session, error) {
	if config.IsServer {
		return newServerSession(ctx, config)
	}
	return newClientSession(ctx, config)
}

func newClientSession(ctx context.Context, config *Config) (*Session, error) {
	conn, err := quic.DialAddr(ctx, config.Address, config.TLSConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	// Perform session negotiation handshake on the control stream.
	if err := performClientHandshake(ctx, conn); err != nil {
		conn.CloseWithError(0, "handshake failed")
		return nil, fmt.Errorf("handshake failed: %w", err)
	}

	Info("Client connected to %s", config.Address)
	
	// Create a session with telemetry manager
	session := &Session{conn: conn}
	
	// Open a telemetry stream for this session
	telemetryStream, err := session.OpenTelemetryStream(ctx)
	if err != nil {
		// Log the error but don't fail the session creation
		Warn("Failed to open telemetry stream: %v", err)
	} else {
		// Create and start telemetry manager
		session.telemetryManager = NewTelemetryManager(conn, telemetryStream, 1*time.Second)
		session.telemetryManager.Start()
	}

	return session, nil
}

func newServerSession(ctx context.Context, config *Config) (*Session, error) {
	listener, err := quic.ListenAddr(config.Address, config.TLSConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", config.Address, err)
	}

	// Accept connections in a loop to handle multiple clients
	// In a full implementation, this would typically be run in a separate goroutine
	go func() {
		for {
			conn, err := listener.Accept(ctx)
			if err != nil {
				// Check if the context was cancelled
				if ctx.Err() != nil {
					break
				}
				Error("Failed to accept connection: %v", err)
				continue
			}

			// Handle each connection in a separate goroutine
			go func(conn quic.Connection) {
				// Perform session negotiation handshake on the control stream.
				if err := performServerHandshake(ctx, conn); err != nil {
					conn.CloseWithError(0, "handshake failed")
					Error("Handshake failed: %v", err)
					return
				}

				Info("Server accepted connection from %s", conn.RemoteAddr().String())
				
				// Create a session for this connection
				session := &Session{conn: conn}
				
				// Accept telemetry stream
				telemetryStream, err := session.AcceptTelemetryStream(ctx)
				if err != nil {
					// Log the error but don't fail the session
					Warn("Failed to accept telemetry stream: %v", err)
				} else {
					// Create and start telemetry manager
					session.telemetryManager = NewTelemetryManager(conn, telemetryStream, 1*time.Second)
					session.telemetryManager.Start()
					defer session.telemetryManager.Stop()
				}
				
				// Accept multiple interactive streams in a loop
				for {
					stream, err := session.AcceptInteractiveStream(ctx)
					if err != nil {
						Error("Failed to accept interactive stream: %v", err)
						// If the session is closed, break out of the loop
						if ctx.Err() != nil {
							break
						}
						// Otherwise, continue accepting streams
						continue
					}
					
					// Handle each stream in a separate goroutine
					go func(s quic.Stream) {
						defer s.Close()
						Info("Accepted interactive stream, echoing data...")
						buf := make([]byte, 1024)
						for {
							n, err := s.Read(buf)
							if err != nil {
								Error("Read error: %v", err)
								break
							}
							if _, err := s.Write(buf[:n]); err != nil {
								Error("Write error: %v", err)
								break
							}
						}
						Info("Finished handling interactive stream")
					}(stream)
				}
			}(conn)
		}
	}()

	// Return a session with a nil conn for server mode
	// This is a workaround to allow the server to start
	return &Session{conn: nil}, nil
}

// performClientHandshake performs the client side of the session negotiation handshake.
func performClientHandshake(ctx context.Context, conn quic.Connection) error {
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return fmt.Errorf("failed to open control stream: %w", err)
	}
	defer stream.Close()

	// Send SessionInit message
	initPayload := &SessionInitPayload{
		Version:           1,
		Token:             nil, // TODO: Add token support
		SupportedFeatures: []string{}, // TODO: Add feature support
	}
	initData, err := EncodeSessionInit(initPayload)
	if err != nil {
		return fmt.Errorf("failed to encode SessionInit: %w", err)
	}

	msg := &Message{
		Type: SessionInit,
		Data: initData,
	}
	
	if err := WriteMessage(stream, msg); err != nil {
		return fmt.Errorf("failed to send SessionInit: %w", err)
	}

	// Receive SessionAccept message
	receivedMsg, err := ReadMessage(stream)
	if err != nil {
		return fmt.Errorf("failed to read SessionAccept: %w", err)
	}

	if receivedMsg.Type != SessionAccept {
		return fmt.Errorf("expected SessionAccept, got %d", receivedMsg.Type)
	}

	acceptPayload, err := DecodeSessionAccept(receivedMsg.Data)
	if err != nil {
		return fmt.Errorf("failed to decode SessionAccept payload: %w", err)
	}

	if !acceptPayload.Accepted {
		return fmt.Errorf("server rejected session: %s", acceptPayload.Reason)
	}

	Info("Session handshake completed successfully")
	return nil
}

// performServerHandshake performs the server side of the session negotiation handshake.
func performServerHandshake(ctx context.Context, conn quic.Connection) error {
	stream, err := conn.AcceptStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to accept control stream: %w", err)
	}
	defer stream.Close()

	// Receive SessionInit message
	receivedMsg, err := ReadMessage(stream)
	if err != nil {
		return fmt.Errorf("failed to read SessionInit: %w", err)
	}

	if receivedMsg.Type != SessionInit {
		return fmt.Errorf("expected SessionInit, got %d", receivedMsg.Type)
	}

	initPayload, err := DecodeSessionInit(receivedMsg.Data)
	if err != nil {
		return fmt.Errorf("failed to decode SessionInit payload: %w", err)
	}

	Info("Received SessionInit: Version=%d, Token=%v, Features=%v", initPayload.Version, initPayload.Token, initPayload.SupportedFeatures)

	// Send SessionAccept message
	acceptPayload := &SessionAcceptPayload{
		Accepted:       true,
		Reason:         "",
		ServerFeatures: []string{}, // TODO: Add feature support
	}
	acceptData, err := EncodeSessionAccept(acceptPayload)
	if err != nil {
		return fmt.Errorf("failed to encode SessionAccept: %w", err)
	}

	responseMsg := &Message{
		Type: SessionAccept,
		Data: acceptData,
	}
	
	if err := WriteMessage(stream, responseMsg); err != nil {
		return fmt.Errorf("failed to send SessionAccept: %w", err)
	}

	Info("Session handshake completed successfully")
	return nil
}

// Close closes the underlying QUIC connection and stops the telemetry manager.
func (s *Session) Close() error {
	// Stop the telemetry manager if it exists
	if s.telemetryManager != nil {
		s.telemetryManager.Stop()
	}
	
	// Close the connection if it exists
	if s.conn != nil {
		return s.conn.CloseWithError(0, "")
	}
	
	return nil
}

// Connection returns the underlying QUIC connection.
func (s *Session) Connection() quic.Connection {
	return s.conn
}

// TODO: Add methods for creating/accepting streams (interactive, bulk, telemetry).