package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

// ConnectionPool represents a pool of VANTUN connections
type ConnectionPool struct {
	// config holds the configuration for new connections
	config *Config
	// connections is a map of address to connection pool
	connections map[string]*connectionPool
	// mutex protects the connections map
	mutex sync.RWMutex
	// maxPoolSize is the maximum number of connections to keep in the pool
	maxPoolSize int
	// idleTimeout is the time after which idle connections are closed
	idleTimeout time.Duration
}

// connectionPool represents a pool of connections to a specific address
type connectionPool struct {
	// connections is a list of available connections
	connections []quic.Connection
	// mutex protects the connections list
	mutex sync.Mutex
	// address is the address of the remote endpoint
	address string
}

// NewConnectionPool creates a new ConnectionPool
func NewConnectionPool(config *Config, maxPoolSize int, idleTimeout time.Duration) *ConnectionPool {
	return &ConnectionPool{
		config:      config,
		connections: make(map[string]*connectionPool),
		maxPoolSize: maxPoolSize,
		idleTimeout: idleTimeout,
	}
}

// GetConnection gets a connection from the pool or creates a new one if none are available
func (cp *ConnectionPool) GetConnection(ctx context.Context, address string) (quic.Connection, error) {
	cp.mutex.RLock()
	pool, exists := cp.connections[address]
	cp.mutex.RUnlock()

	if !exists {
		cp.mutex.Lock()
		// Double-check if another goroutine created the pool
		pool, exists = cp.connections[address]
		if !exists {
			pool = &connectionPool{
				connections: make([]quic.Connection, 0),
				address:     address,
			}
			cp.connections[address] = pool
		}
		cp.mutex.Unlock()
	}

	// Try to get an existing connection from the pool
	conn := pool.getConnection()
	if conn != nil {
		return conn, nil
	}

	// Create a new connection
	return cp.createConnection(ctx, address)
}

// ReturnConnection returns a connection to the pool
func (cp *ConnectionPool) ReturnConnection(address string, conn quic.Connection) {
	cp.mutex.RLock()
	pool, exists := cp.connections[address]
	cp.mutex.RUnlock()

	if !exists {
		// If the pool doesn't exist, just close the connection
		conn.CloseWithError(0, "pool not found")
		return
	}

	// Try to add the connection to the pool
	if !pool.addConnection(conn, cp.maxPoolSize) {
		// If the pool is full, close the connection
		conn.CloseWithError(0, "pool full")
	}
}

// createConnection creates a new VANTUN connection
func (cp *ConnectionPool) createConnection(ctx context.Context, address string) (quic.Connection, error) {
	// Create a new config with the specific address
	config := &Config{
		Address:   address,
		TLSConfig: cp.config.TLSConfig,
		IsServer:  false, // Connection pool is for client connections
	}

	// Create a new session
	session, err := NewSession(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session.conn, nil
}

// Close closes all connections in the pool
func (cp *ConnectionPool) Close() error {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	// Close all connections in all pools
	for _, pool := range cp.connections {
		pool.mutex.Lock()
		for _, conn := range pool.connections {
			conn.CloseWithError(0, "pool closed")
		}
		pool.connections = nil
		pool.mutex.Unlock()
	}

	// Clear the connections map
	cp.connections = make(map[string]*connectionPool)

	return nil
}

// getConnection gets a connection from the pool
func (cp *connectionPool) getConnection() quic.Connection {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	// If there are no connections, return nil
	if len(cp.connections) == 0 {
		return nil
	}

	// Get the last connection from the list
	conn := cp.connections[len(cp.connections)-1]
	cp.connections = cp.connections[:len(cp.connections)-1]

	return conn
}

// addConnection adds a connection to the pool if there is space
func (cp *connectionPool) addConnection(conn quic.Connection, maxPoolSize int) bool {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	// If the pool is at maximum capacity, return false
	if len(cp.connections) >= maxPoolSize {
		return false
	}

	// Add the connection to the pool
	cp.connections = append(cp.connections, conn)
	return true
}

// ConnectionPoolSession wraps a ConnectionPool to provide session-like functionality
type ConnectionPoolSession struct {
	*ConnectionPool
}

// NewConnectionPoolSession creates a new ConnectionPoolSession
func NewConnectionPoolSession(config *Config, maxPoolSize int, idleTimeout time.Duration) *ConnectionPoolSession {
	return &ConnectionPoolSession{
		ConnectionPool: NewConnectionPool(config, maxPoolSize, idleTimeout),
	}
}

// OpenInteractiveStream opens a new interactive stream using a connection from the pool
func (cps *ConnectionPoolSession) OpenInteractiveStream(ctx context.Context, address string) (quic.Stream, error) {
	conn, err := cps.GetConnection(ctx, address)
	if err != nil {
		return nil, err
	}

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		// If we failed to open a stream, return the connection to the pool
		cps.ReturnConnection(address, conn)
		return nil, err
	}

	// Wrap the stream to return the connection to the pool when it's closed
	return &PooledStream{
		Stream:          stream,
		connection:      conn,
		address:         address,
		connectionPool:  cps.ConnectionPool,
		streamClosed:    make(chan struct{}),
	}, nil
}

// PooledStream wraps a QUIC stream to return the connection to the pool when closed
type PooledStream struct {
	quic.Stream
	connection     quic.Connection
	address        string
	connectionPool *ConnectionPool
	streamClosed   chan struct{}
	once           sync.Once
}

// Close closes the stream and returns the connection to the pool
func (ps *PooledStream) Close() error {
	// Ensure we only close the stream and return the connection once
	ps.once.Do(func() {
		// Close the underlying stream
		ps.Stream.Close()

		// Return the connection to the pool
		ps.connectionPool.ReturnConnection(ps.address, ps.connection)

		// Signal that the stream is closed
		close(ps.streamClosed)
	})

	return nil
}

// Closed returns a channel that is closed when the stream is closed
func (ps *PooledStream) Closed() <-chan struct{} {
	return ps.streamClosed
}