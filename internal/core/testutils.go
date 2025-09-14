package core

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

// MockQUICStream is a mock implementation of quic.Stream for testing
type MockQUICStream struct {
	readData  []byte
	writeData []byte
	streamID  quic.StreamID
}

func (m *MockQUICStream) Read(p []byte) (n int, err error) {
	if len(m.readData) == 0 {
		return 0, nil
	}
	n = copy(p, m.readData)
	m.readData = m.readData[n:]
	return n, nil
}

func (m *MockQUICStream) Write(p []byte) (n int, err error) {
	m.writeData = append(m.writeData, p...)
	return len(p), nil
}

func (m *MockQUICStream) Close() error {
	return nil
}

func (m *MockQUICStream) CancelRead(quic.StreamErrorCode) {
}

func (m *MockQUICStream) CancelWrite(quic.StreamErrorCode) {
}

func (m *MockQUICStream) LocalAddr() net.Addr {
	return &net.UDPAddr{}
}

func (m *MockQUICStream) RemoteAddr() net.Addr {
	return &net.UDPAddr{}
}

func (m *MockQUICStream) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockQUICStream) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockQUICStream) SetWriteDeadline(t time.Time) error {
	return nil
}

func (m *MockQUICStream) StreamID() quic.StreamID {
	return m.streamID
}

func (m *MockQUICStream) Context() context.Context {
	return context.Background()
}

func (m *MockQUICStream) SendStream() quic.SendStream {
	return m
}

func (m *MockQUICStream) ReceiveStream() quic.ReceiveStream {
	return m
}

// MockQUICConnection is a mock implementation of quic.Connection for testing
type MockQUICConnection struct {
	addr     string
	rtt      time.Duration
	loss     float64
	streams  []quic.Stream
	streamID quic.StreamID
	closed   bool
	mutex    sync.RWMutex
}

func (m *MockQUICConnection) AcceptStream(ctx context.Context) (quic.Stream, error) {
	m.mutex.RLock()
	closed := m.closed
	m.mutex.RUnlock()
	
	if closed {
		return nil, fmt.Errorf("connection closed")
	}
	
	if len(m.streams) == 0 {
		// Block until context is done
		<-ctx.Done()
		return nil, ctx.Err()
	}
	stream := m.streams[0]
	m.streams = m.streams[1:]
	return stream, nil
}

func (m *MockQUICConnection) AcceptUniStream(ctx context.Context) (quic.ReceiveStream, error) {
	m.mutex.RLock()
	closed := m.closed
	m.mutex.RUnlock()
	
	if closed {
		return nil, fmt.Errorf("connection closed")
	}
	
	return nil, nil
}

func (m *MockQUICConnection) OpenStreamSync(ctx context.Context) (quic.Stream, error) {
	m.mutex.RLock()
	closed := m.closed
	m.mutex.RUnlock()
	
	if closed {
		return nil, fmt.Errorf("connection closed")
	}
	
	stream := &MockQUICStream{streamID: m.streamID}
	m.streamID++
	m.streams = append(m.streams, stream)
	return stream, nil
}

func (m *MockQUICConnection) OpenStream() (quic.Stream, error) {
	m.mutex.RLock()
	closed := m.closed
	m.mutex.RUnlock()
	
	if closed {
		return nil, fmt.Errorf("connection closed")
	}
	
	stream := &MockQUICStream{streamID: m.streamID}
	m.streamID++
	m.streams = append(m.streams, stream)
	return stream, nil
}

func (m *MockQUICConnection) OpenUniStreamSync(ctx context.Context) (quic.SendStream, error) {
	m.mutex.RLock()
	closed := m.closed
	m.mutex.RUnlock()
	
	if closed {
		return nil, fmt.Errorf("connection closed")
	}
	
	return nil, nil
}

func (m *MockQUICConnection) OpenUniStream() (quic.SendStream, error) {
	m.mutex.RLock()
	closed := m.closed
	m.mutex.RUnlock()
	
	if closed {
		return nil, fmt.Errorf("connection closed")
	}
	
	return nil, nil
}

func (m *MockQUICConnection) CloseWithError(quic.ApplicationErrorCode, string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.closed = true
	return nil
}

func (m *MockQUICConnection) Context() context.Context {
	return context.Background()
}

func (m *MockQUICConnection) ConnectionState() quic.ConnectionState {
	return quic.ConnectionState{}
}

func (m *MockQUICConnection) SendDatagram([]byte) error {
	m.mutex.RLock()
	closed := m.closed
	m.mutex.RUnlock()
	
	if closed {
		return fmt.Errorf("connection closed")
	}
	
	return nil
}

func (m *MockQUICConnection) ReceiveDatagram(context.Context) ([]byte, error) {
	m.mutex.RLock()
	closed := m.closed
	m.mutex.RUnlock()
	
	if closed {
		return nil, fmt.Errorf("connection closed")
	}
	
	return nil, nil
}

func (m *MockQUICConnection) RemoteAddr() net.Addr {
	return &net.UDPAddr{}
}

func (m *MockQUICConnection) LocalAddr() net.Addr {
	return &net.UDPAddr{}
}

// SetStreamID sets the next stream ID to be used for new streams
func (m *MockQUICConnection) SetStreamID(id quic.StreamID) {
	m.streamID = id
}

// AddStream adds a stream to the connection's stream list
func (m *MockQUICConnection) AddStream(stream quic.Stream) {
	m.streams = append(m.streams, stream)
}

// createTestServer creates a test VANTUN server for testing purposes
func createTestServer() (string, *Session, error) {
	// Generate a self-signed certificate for testing
	cert, key, err := generateTestCert()
	if err != nil {
		return "", nil, err
	}

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return "", nil, err
	}

	// Create TLS config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"vantun"},
	}

	// Listen on a random port
	listener, err := quic.ListenAddr("localhost:0", tlsConfig, nil)
	if err != nil {
		return "", nil, err
	}

	// Get the actual address
	addr := listener.Addr().String()

	// Handle connections in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	
	// Channel to signal when the server is ready
	ready := make(chan struct{})
	
	go func() {
		defer listener.Close()
		defer cancel()
		
		// Signal that the server is ready to accept connections
		close(ready)
		
		// Accept connections and handle handshake
		for {
			conn, err := listener.Accept(ctx)
			if err != nil {
				return
			}
			
			// Handle each connection in a separate goroutine
			go func(c quic.Connection) {
				defer c.CloseWithError(0, "test completed")
				
				// Perform session negotiation handshake on the control stream.
				if err := performServerHandshake(ctx, c); err != nil {
					Error("Handshake failed: %v", err)
					return
				}
				
				// Accept streams and echo data
				for {
					stream, err := c.AcceptStream(ctx)
					if err != nil {
						return
					}
					
					// Handle each stream in a separate goroutine
					go func(s quic.Stream) {
						defer s.Close()
						Info("Handling new stream %v", s)
						
						// First, read and handle the stream type message
						Info("About to read stream type message")
						msg, err := ReadMessage(s)
						if err != nil {
							Info("Failed to read stream type message: %v", err)
							return
						}
						Info("Successfully read message: Type=%d, Data length=%d", msg.Type, len(msg.Data))
						
						// If it's a stream type message, just acknowledge it by echoing it back
						if msg.Type == StreamType {
							Info("Received stream type message, echoing it back")
							// Echo the stream type message back
							if err := WriteMessage(s, msg); err != nil {
								Error("Failed to write stream type echo: %v", err)
								return
							}
							Info("Successfully echoed stream type message")
						} else {
							// For non-stream type messages, write it back and continue
							Info("Received non-stream type message, echoing data")
							if _, err := s.Write(msg.Data); err != nil {
								Error("Failed to write echo: %v", err)
								return
							}
						}
						
						// Continue reading and echoing data
						Info("Entering data echoing loop")
						buf := make([]byte, 4096)
						for {
							Info("About to read data")
							n, err := s.Read(buf)
							Info("Read returned: n=%d, err=%v", n, err)
							if err != nil {
								Info("Read error: %v", err)
								return
							}
							// Log the data being received
							Info("Received %d bytes: %v", n, buf[:n])
							// Log the data being echoed
							Info("Echoing %d bytes: %v", n, buf[:n])
							if _, err := s.Write(buf[:n]); err != nil {
								Error("Failed to write echo: %v", err)
								return
							}
						}
					}(stream)
				}
			}(conn)
		}
	}()

	// Wait for the server to be ready
	<-ready

	// Return the address and nil session since the session is not actually used
	return addr, nil, nil
}

// generateTestCert generates a self-signed certificate for testing
func generateTestCert() ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	return certPEM, keyPEM, nil
}