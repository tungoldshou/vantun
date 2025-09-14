package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"math/big"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vantun/internal/cli"
	"vantun/internal/core"
	"github.com/quic-go/quic-go"
)

var (
	isServer     = flag.Bool("server", false, "Run as server")
	addr         = flag.String("addr", "localhost:4242", "Address to listen on (server) or connect to (client)")
	configFile   = flag.String("config", "", "Path to JSON configuration file")
	logLevel     = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	multipath    = flag.Bool("multipath", false, "Enable multipath")
	obfs         = flag.Bool("obfs", false, "Enable obfuscation")
	fecDataShards   = flag.Int("fec-data", 10, "Number of FEC data shards")
	fecParityShards = flag.Int("fec-parity", 3, "Number of FEC parity shards")
)

func main() {
	flag.Parse()

	var configManager *cli.ConfigManager
	var config *cli.Config

	// Load configuration from JSON file if provided
	if *configFile != "" {
		configManager = cli.NewConfigManager(*configFile)
		if err := configManager.StartHotReload(); err != nil {
			core.Error("Failed to start config manager: %v", err)
			os.Exit(1)
		}
		defer configManager.StopHotReload()
		
		// Get the initial config
		config = configManager.GetConfig()
	} else {
		// Use command-line flags or defaults
		config = &cli.Config{
			Server:              *isServer,
			Address:             *addr,
			LogLevel:            *logLevel,
			Multipath:           *multipath,
			Obfs:                *obfs,
			FECData:             *fecDataShards,
			FECParity:           *fecParityShards,
			TokenBucketRate:     1000000,   // Default 1 MB/s
			TokenBucketCapacity: 5000000,  // Default 5 MB capacity
		}
	}

	// Set up logging
	core.InitLogger(config.LogLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		core.Info("Shutting down...")
		cancel()
	}()

	// Create TLS config
	tlsConfig := generateTLSConfig()

	// Create core configuration
	currentConfig := configManager.GetConfig()
	
	coreConfig := &core.Config{
		Address:   currentConfig.Address,
		TLSConfig: tlsConfig,
		IsServer:  currentConfig.Server,
	}

	// Create token bucket
	tokenBucket := core.NewTokenBucket(currentConfig.TokenBucketRate, currentConfig.TokenBucketCapacity)
	
	// Create adaptive FEC
	adaptiveFEC, err := core.NewAdaptiveFEC(currentConfig.FECData, currentConfig.FECParity, 1, 10)
	if err != nil {
		core.Error("Failed to create adaptive FEC: %v", err)
		os.Exit(1)
	}
	
	// Create obfuscator if enabled
	obfuscator := core.NewObfuscator(core.ObfuscatorConfig{
		Enabled: currentConfig.Obfs,
	})
	
	// Variables for sessions
	var session *core.Session
	var obfsSession *core.ObfuscatorSession
	var multipathSession *core.MultipathSession
	
	// Create session
	if currentConfig.Multipath {
		// Create multipath session
		core.Info("Creating multipath session")
		// Create a token bucket controller for multipath
		tokenBucket := core.NewTokenBucket(currentConfig.TokenBucketRate, currentConfig.TokenBucketCapacity)
		adaptiveFEC, err := core.NewAdaptiveFEC(currentConfig.FECData, currentConfig.FECParity, 1, 10)
		if err != nil {
			core.Error("Failed to create adaptive FEC: %v", err)
			os.Exit(1)
		}
		
		// Create a token bucket controller (will be updated with connection later)
		controller := core.NewTokenBucketController(tokenBucket, adaptiveFEC, nil)
		controller.Start()
		defer controller.Stop()
		
		multipathSession = core.NewMultipathSession(coreConfig, controller, adaptiveFEC)
		// Add the primary path
		if err := multipathSession.AddPath(ctx, currentConfig.Address); err != nil {
			core.Error("Failed to add path to multipath session: %v", err)
			os.Exit(1)
		}
		// TODO: Add additional paths for true multipath support
	} else {
		// Create regular session
		session, err = core.NewSession(ctx, coreConfig)
		if err != nil {
			core.Error("Failed to create session: %v", err)
			os.Exit(1)
		}
		
		// If obfuscation is enabled, wrap the session
		if currentConfig.Obfs {
			obfsSession = core.NewObfuscatorSession(session, obfuscator)
		}
	}
	
	// Create token bucket controller
	// Note: For multipath, we would need to pass the connection from the primary path
	var controller *core.TokenBucketController
	if currentConfig.Multipath {
		// TODO: Implement controller for multipath
		core.Info("Token bucket controller not implemented for multipath yet")
	} else if !currentConfig.Server {
		// Only start token bucket controller for client mode
		// Server mode handles telemetry within the session
		controller = core.NewTokenBucketController(tokenBucket, adaptiveFEC, session.Connection())
		// Note: The telemetry stream is now handled within the session itself
		controller.Start()
		defer controller.Stop()
	} else {
		core.Info("Token bucket controller not started for server mode")
	}

	if currentConfig.Multipath {
		defer multipathSession.Close()
	} else {
		defer session.Close()
	}

	if currentConfig.Multipath && currentConfig.Server {
		core.Info("Multipath server running, waiting for connections...")
		// For a multipath server, we would need to implement a different logic
		// For now, we'll just keep the main goroutine alive
		<-ctx.Done()
	} else if currentConfig.Multipath && !currentConfig.Server {
		core.Info("Multipath client connected, opening interactive stream...")
		// For demo, open one interactive stream and send/receive data
		// Note: This is a simplified implementation that only uses one path
		// A real multipath implementation would be more complex
		stream, err := multipathSession.OpenStream(ctx)
		if err != nil {
			core.Error("Failed to open interactive stream: %v", err)
			os.Exit(1)
		}
		defer stream.Close()

		// Send a message
		message := "Hello from VANTUN multipath client!"
		if _, err := stream.Write([]byte(message)); err != nil {
			core.Error("Failed to send message: %v", err)
			os.Exit(1)
		}

		// Read the echo
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			core.Error("Failed to read echo: %v", err)
			os.Exit(1)
		}
		core.Info("Received echo: %s", string(buf[:n]))
	} else if currentConfig.Obfs && !currentConfig.Server {
		core.Info("Obfuscated client connected, opening interactive stream...")
		// For demo, open one interactive stream and send/receive data
		stream, err := obfsSession.OpenInteractiveStream(ctx)
		if err != nil {
			core.Error("Failed to open interactive stream: %v", err)
			os.Exit(1)
		}
		defer stream.Close()

		// Send a message
		message := "Hello from VANTUN obfuscated client!"
		if _, err := stream.Write([]byte(message)); err != nil {
			core.Error("Failed to send message: %v", err)
			os.Exit(1)
		}

		// Read the echo
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			core.Error("Failed to read echo: %v", err)
			os.Exit(1)
		}
		core.Info("Received echo: %s", string(buf[:n]))
	} else if !currentConfig.Server {
		core.Info("Client connected, opening interactive stream...")
		// For demo, open one interactive stream and send/receive data
		var stream quic.Stream
		var err error
		
		if currentConfig.Obfs {
			stream, err = obfsSession.OpenInteractiveStream(ctx)
		} else {
			stream, err = session.OpenInteractiveStream(ctx)
		}
		
		if err != nil {
			core.Error("Failed to open interactive stream: %v", err)
			os.Exit(1)
		}
		defer stream.Close()

		// Send a message
		message := "Hello from VANTUN client!"
		if _, err := stream.Write([]byte(message)); err != nil {
			core.Error("Failed to send message: %v", err)
			os.Exit(1)
		}

		// Read the echo
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			core.Error("Failed to read echo: %v", err)
			os.Exit(1)
		}
		core.Info("Received echo: %s", string(buf[:n]))
	} else {
		// Server mode
		core.Info("Server running, waiting for connections...")
		// For server mode, we just wait for shutdown signal
		<-ctx.Done()
	}

	// Wait for a bit to see logs
	time.Sleep(1 * time.Second)
}

// generateTLSConfig generates a simple self-signed TLS config for testing.
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		// Add localhost to the list of valid subjects to avoid certificate verification errors.
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"vantun"},
		// For testing purposes, we skip certificate verification.
		// In production, this should be removed.
		InsecureSkipVerify: true,
	}
}