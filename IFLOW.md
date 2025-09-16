# VANTUN Project Context for iFlow CLI

## Project Overview

VANTUN is a next-generation, high-performance tunnel protocol built on top of QUIC. It's designed to deliver exceptional network performance, security, and reliability. The project is written in Go and leverages several key technologies:

- **Core Library**: `quic-go` for QUIC protocol implementation
- **Serialization**: `github.com/fxamacker/cbor` for efficient CBOR encoding
- **FEC**: `github.com/klauspost/reedsolomon` for Reed-Solomon encoding
- **CLI**: `cobra/viper` for command-line interface and configuration management

Key features include:
- Enterprise-grade security with secure handshake and session negotiation
- Exceptional performance with multiple logical stream types and multipath transmission
- Unmatched reliability with Forward Error Correction (FEC) and hybrid congestion control
- Privacy protection with pluggable obfuscation modules
- Easy deployment with minimal client/server programs

## Project Structure

```
vantun/
├── cmd/              # Command-line program entry
├── internal/
│   ├── cli/          # CLI configuration management
│   └── core/         # Core protocol implementation
├── docs/             # Documentation
├── go.mod            # Go module definition
└── README.md         # Project documentation
```

## Building and Running

### Prerequisites
- Go 1.22 or higher

### Compilation
```bash
go build -o bin/vantun cmd/main.go
```

### Configuration
VANTUN supports configuration through command-line arguments and JSON configuration files.

**Command-Line Arguments:**
- `-server`: Run in server mode
- `-addr`: Listen address (server) or connection address (client), default is `localhost:4242`
- `-config`: JSON configuration file path
- `-log-level`: Log level
- `-multipath`: Enable multipath
- `-obfs`: Enable traffic obfuscation
- `-fec-data`: FEC data shards count, default is 10
- `-fec-parity`: FEC parity shards count, default is 3

**JSON Configuration File Example:**
```json
{
  "server": false,
  "address": "localhost:4242",
  "log_level": "info",
  "multipath": false,
  "obfs": false,
  "fec_data": 10,
  "fec_parity": 3,
  "token_bucket_rate": 1000000,
  "token_bucket_capacity": 5000000
}
```

### Running

**Starting the Server:**
```bash
# Using command-line arguments
./bin/vantun -server -addr :4242

# Or using configuration file
./bin/vantun -config config.json -server
```

**Starting the Client:**
```bash
# Using command-line arguments
./bin/vantun -addr localhost:4242

# Or using configuration file
./bin/vantun -config config.json
```

## Development Conventions

### Code Organization
- Core protocol implementation is in `internal/core/`
- CLI configuration management is in `internal/cli/`
- All tests are colocated with their implementation files

### Testing
The project includes comprehensive unit tests, integration tests, performance tests, and stress tests.

Run all tests:
```bash
go test -v ./internal/core/...
```

### Key Components
1. **Session Management**: Handles QUIC connections and session negotiation in `session.go`
2. **Multipath Transmission**: Implemented in `multipath.go` for utilizing multiple network paths
3. **Forward Error Correction**: Adaptive FEC implementation in `adaptive_fec.go` and `fec.go`
4. **Traffic Obfuscation**: Obfuscation modules in `obfuscation.go` and HTTP/3 obfuscator in `http3_obfuscator.go`
5. **Telemetry**: Real-time performance monitoring in `telemetry.go`
6. **Congestion Control**: Hybrid approach combining QUIC CC with token bucket rate limiting in `token_bucket.go`