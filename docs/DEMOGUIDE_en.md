# VANTUN Demo Guide

## 1. Compilation

Make sure you have Go 1.21 or higher installed.

```bash
# Clone repository (if not already cloned)
# git clone <repository-url>
cd vantun

# Compilation
go build -o bin/vantun cmd/main.go
```

## 2. Configuration

VANTUN supports configuration through command-line arguments and JSON configuration files.

### 2.1 Command-Line Arguments

- `-server`: Run in server mode.
- `-addr`: Listen address (server) or connection address (client), default is `localhost:4242`.
- `-config`: JSON configuration file path.
- `-log-level`: Log level.
- `-multipath`: Enable multipath.
- `-obfs`: Enable traffic obfuscation.
- `-fec-data`: FEC data shards count, default is 10.
- `-fec-parity`: FEC parity shards count, default is 3.

### 2.2 JSON Configuration File

Create a `config.json` file:

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

## 3. Running

### 3.1 Starting the Server

```bash
# Using command-line arguments
./bin/vantun -server -addr :4242

# Or using configuration file
./bin/vantun -config config.json -server
```

### 3.2 Starting the Client

```bash
# Using command-line arguments
./bin/vantun -addr localhost:4242

# Or using configuration file
./bin/vantun -config config.json
```

## 4. Verification

When both client and server are started, you should see the following output:

**Server:**
```
2025/09/14 04:00:00 Server running, waiting for streams...
2025/09/14 04:00:00 Accepted interactive stream, echoing data...
```

**Client:**
```
2025/09/14 04:00:00 Client connected, opening interactive stream...
2025/09/14 04:00:00 Received echo: Hello from VANTUN client!
```

This indicates that the client successfully connected to the server, opened an interactive stream, sent a message, and received an echo from the server.

## 5. Feature Testing

You can test different features by modifying the configuration:

### 5.1 Testing FEC

Modify the `fec_data` and `fec_parity` parameters in `config.json`, then restart both client and server.

### 5.2 Testing Traffic Obfuscation

Enable the `-obfs` parameter on both client and server or set `"obfs": true` in the configuration file.

### 5.3 Testing Multipath

Multipath functionality requires multiple network interfaces. You can use Linux's `netem` tool to simulate different network conditions.

## 6. Performance Testing

You can use `iperf` or other network performance testing tools to evaluate VANTUN's performance.

## 7. Logging

VANTUN outputs telemetry data and control information to standard output, which can be used for monitoring and debugging.