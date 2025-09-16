# Configuration Guide

## 🎯 Overview

This guide covers all configuration options for VANTUN, from basic setups to advanced optimizations.

## 📋 Configuration File Structure

VANTUN uses JSON format for configuration files. All configuration files follow the same basic structure:

```json
{
  "server": boolean,                    // Server or client mode
  "address": "host:port",               // Listen/connect address
  "log_level": "info",                  // Logging level
  "log_file": "/path/to/log",           // Log file path
  
  // Core Features
  "multipath": boolean,                 // Enable multipath
  "obfs": boolean,                      // Enable obfuscation
  "fec_data": 10,                       // FEC data shards
  "fec_parity": 3,                      // FEC parity shards
  
  // Performance Tuning
  "token_bucket_rate": 1000000,        // Rate limiting (bps)
  "token_bucket_capacity": 5000000,    // Bucket capacity (bits)
  
  // TLS Configuration (optional)
  "tls": {
    "cert": "/path/to/cert.pem",
    "key": "/path/to/key.pem"
  },
  
  // Advanced Options
  "performance": {
    "workers": 4,                      // Worker threads
    "buffer_size": 2097152,            // Buffer size (bytes)
    "read_buffer": 1048576,            // Read buffer (bytes)
    "write_buffer": 1048576            // Write buffer (bytes)
  },
  
  // Mobile Optimization
  "mobile": {
    "network_detection": true,         // Auto-detect network type
    "adaptive_fec": true,              // Adaptive FEC adjustment
    "power_optimization": true         // Battery optimization
  }
}
```

## 🚀 Basic Configuration

### Minimal Server Configuration
```json
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info"
}
```

### Minimal Client Configuration
```json
{
  "server": false,
  "address": "server.example.com:4242",
  "log_level": "info"
}
```

### Recommended Basic Configuration
```json
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3
}
```

## ⚙️ Core Configuration Options

### Server vs Client Mode
```json
{
  "server": true,    // Run as server
  "server": false    // Run as client
}
```

### Network Configuration
```json
{
  "address": "0.0.0.0:4242",        // Listen on all interfaces
  "address": "192.168.1.100:4242",  // Listen on specific IP
  "address": "server.example.com:4242"  // Connect to server
}
```

### Logging Configuration
```json
{
  "log_level": "debug",     // Detailed debugging information
  "log_level": "info",      // General information (default)
  "log_level": "warn",      // Warnings only
  "log_level": "error",     // Errors only
  "log_level": "fatal",     // Fatal errors only
  
  "log_file": "/var/log/vantun/server.log",  // Log to file
  "log_file": "",           // Log to stdout (default)
}
```

## 🔧 Feature Configuration

### Multipath Settings
```json
{
  "multipath": true,                    // Enable multipath
  "multipath": false,                   // Disable multipath
  
  "multipath_config": {
    "max_paths": 4,                     // Maximum paths to use
    "scheduler": "round_robin",         // Path scheduling algorithm
    "failover_timeout": 30,             // Failover timeout (seconds)
    "path_discovery": true              // Auto-discover paths
  }
}
```

### Obfuscation Settings
```json
{
  "obfs": true,                          // Enable obfuscation
  "obfs": false,                         // Disable obfuscation
  
  "obfuscation_config": {
    "mode": "http3",                    // Obfuscation mode
    "padding": true,                    // Enable padding
    "padding_ratio": 0.1,               // Padding ratio (0.0-1.0)
    "timing_jitter": 50,                // Timing jitter (ms)
    "entropy_target": 7.8               // Target entropy level
  }
}
```

### FEC Configuration
```json
{
  "fec_data": 10,                        // Number of data shards
  "fec_parity": 3,                       // Number of parity shards
  
  "fec_config": {
    "enabled": true,                     // Enable FEC
    "adaptive": true,                    // Adaptive FEC adjustment
    "loss_threshold": 0.05,              // Loss threshold (5%)
    "update_interval": 30,               // Update interval (seconds)
    "max_overhead": 0.5                  // Maximum overhead (50%)
  }
}
```

### Rate Limiting Configuration
```json
{
  "token_bucket_rate": 1000000,          // Rate limit (bits per second)
  "token_bucket_capacity": 5000000,      // Bucket capacity (bits)
  
  "rate_limiting_config": {
    "enabled": true,                     // Enable rate limiting
    "burst_allowance": 1.5,              // Burst multiplier
    "measurement_window": 60,            // Measurement window (seconds)
    "adaptive": true                     // Adaptive rate limiting
  }
}
```

## 🔒 TLS Configuration

### Self-Signed Certificate (Development)
```bash
# Generate self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/CN=vantun.local"
```

### TLS Configuration
```json
{
  "tls": {
    "enabled": true,                     // Enable TLS
    "cert": "/etc/vantun/server.crt",    // Certificate path
    "key": "/etc/vantun/server.key",     // Private key path
    "ca": "/etc/vantun/ca.crt",          // CA certificate path (optional)
    "verify_client": false,              // Verify client certificates
    "cipher_suites": [                   // Cipher suites (optional)
      "TLS_AES_128_GCM_SHA256",
      "TLS_AES_256_GCM_SHA384",
      "TLS_CHACHA20_POLY1305_SHA256"
    ],
    "min_version": "1.2",                // Minimum TLS version
    "max_version": "1.3"                 // Maximum TLS version
  }
}
```

### Let's Encrypt Integration
```json
{
  "tls": {
    "enabled": true,
    "cert": "/etc/letsencrypt/live/yourdomain.com/fullchain.pem",
    "key": "/etc/letsencrypt/live/yourdomain.com/privkey.pem",
    "auto_renew": true,
    "renew_before": 30                    // Renew 30 days before expiry
  }
}
```

## ⚡ Performance Optimization

### Basic Performance Settings
```json
{
  "performance": {
    "workers": 4,                        // Number of worker threads
    "buffer_size": 2097152,              // General buffer size (2MB)
    "read_buffer": 1048576,              // Read buffer (1MB)
    "write_buffer": 1048576,             // Write buffer (1MB)
    "socket_buffer": 4194304             // Socket buffer (4MB)
  }
}
```

### Advanced Performance Tuning
```json
{
  "performance": {
    "cpu_affinity": [0, 1, 2, 3],        // CPU affinity mask
    "numa_aware": true,                  // NUMA-aware scheduling
    "batch_size": 32,                    // Batch processing size
    "pipeline_depth": 16,                // Pipeline depth
    "memory_pool_size": 10000,           // Memory pool size
    "gc_threshold": 0.8                  // Garbage collection threshold
  }
}
```

### Network Optimization
```json
{
  "network": {
    "tcp_nodelay": true,                 // Disable Nagle's algorithm
    "tcp_quickack": true,                // Enable TCP quick acknowledgments
    "socket_buffer_size": 4194304,       // Socket buffer size (4MB)
    "mtu_discovery": true,               // Enable path MTU discovery
    "congestion_control": "bbr"          // Congestion control algorithm
  }
}
```

## 📱 Mobile Optimization

### Basic Mobile Settings
```json
{
  "mobile": {
    "enabled": true,                     // Enable mobile optimizations
    "network_detection": true,           // Auto-detect network type
    "adaptive_fec": true,                // Adaptive FEC for mobile
    "power_optimization": true           // Battery optimization
  }
}
```

### Advanced Mobile Configuration
```json
{
  "mobile": {
    "enabled": true,
    "network_detection": true,
    "adaptive_fec": true,
    "power_optimization": true,
    
    "wifi_optimization": {
      "enabled": true,                   // WiFi-specific optimizations
      "max_fec_overhead": 0.2,           // Max FEC overhead on WiFi
      "multipath_preferred": true        // Prefer multipath on WiFi
    },
    
    "cellular_optimization": {
      "enabled": true,                   // Cellular-specific optimizations
      "max_fec_overhead": 0.4,           // Higher FEC for cellular
      "aggressive_fec": true,             // Aggressive FEC on cellular
      "bandwidth_estimation": true        // Estimate cellular bandwidth
    },
    
    "roaming_optimization": {
      "enabled": true,                   // Roaming optimizations
      "detect_roaming": true,             // Auto-detect roaming
      "reduce_multipath": true,           // Reduce multipath when roaming
      "conservative_fec": true            // Conservative FEC when roaming
    }
  }
}
```

## 🔧 Advanced Configuration

### Logging Configuration
```json
{
  "logging": {
    "level": "info",                      // Log level
    "file": "/var/log/vantun/server.log", // Log file path
    "max_size": 100,                      // Max log file size (MB)
    "max_backups": 10,                    // Max backup files
    "max_age": 7,                         // Max age (days)
    "compress": true,                     // Compress old logs
    "json_format": false,                 // JSON log format
    "enable_metrics": true                // Enable metrics logging
  }
}
```

### Metrics Configuration
```json
{
  "metrics": {
    "enabled": true,                      // Enable metrics
    "port": 8080,                         // Metrics port
    "path": "/metrics",                   // Metrics path
    "namespace": "vantun",                // Metrics namespace
    "subsystem": "server",                // Metrics subsystem
    "labels": {                           // Default labels
      "instance": "server-1",
      "region": "us-east"
    }
  }
}
```

### Debug Configuration
```json
{
  "debug": {
    "enabled": false,                     // Enable debug mode
    "pprof_port": 6060,                   // pprof port
    "trace_enabled": false,               // Enable tracing
    "trace_endpoint": "http://localhost:14268/api/traces", // Trace endpoint
    "memory_profiling": true,             // Enable memory profiling
    "cpu_profiling": true,                // Enable CPU profiling
    "block_profiling": true               // Enable block profiling
  }
}
```

## 🎯 Environment-Specific Configurations

### High-Loss Network Configuration
```json
{
  "fec_data": 15,                        // More data shards
  "fec_parity": 5,                       // More parity shards
  "fec_config": {
    "adaptive": true,
    "loss_threshold": 0.08,              // Higher loss threshold
    "aggressive_mode": true               // Aggressive FEC mode
  },
  "performance": {
    "buffer_size": 4194304,              // Larger buffers
    "retry_attempts": 5                  // More retry attempts
  }
}
```

### Low-Latency Configuration
```json
{
  "performance": {
    "workers": 8,                        // More workers
    "buffer_size": 1048576,              # Smaller buffers for lower latency
    "batch_size": 16,                    # Smaller batch size
    "pipeline_depth": 8                  # Shallow pipeline
  },
  "network": {
    "tcp_nodelay": true,                 # Disable Nagle's algorithm
    "tcp_quickack": true,                # Quick acknowledgments
    "socket_buffer_size": 2097152        # Moderate buffer size
  }
}
```

### High-Throughput Configuration
```json
{
  "performance": {
    "workers": 16,                       # Many workers
    "buffer_size": 8388608,              # Large buffers
    "batch_size": 64,                    # Large batch size
    "pipeline_depth": 32                 # Deep pipeline
  },
  "network": {
    "socket_buffer_size": 16777216,      # Large socket buffers
    "mtu_discovery": true,               # MTU discovery
    "congestion_control": "bbr"          # BBR congestion control
  }
}
```

## 🔍 Configuration Validation

### Test Configuration
```bash
# Validate configuration file
vantun --config /etc/vantun/server.json --check

# Test with verbose output
vantun --config /etc/vantun/server.json --check --verbose
```

### Configuration Generator
```bash
# Generate sample configuration
vantun --generate-config server > server.json
vantun --generate-config client > client.json

# Generate with specific options
vantun --generate-config server --multipath --obfs > server-advanced.json
```

## 📊 Performance Monitoring

### Key Metrics to Monitor
- **Throughput**: Data transfer rate
- **Latency**: Round-trip time
- **Packet Loss**: Network loss rate
- **CPU Usage**: Processor utilization
- **Memory Usage**: Memory consumption
- **Connection Count**: Active connections

### Configuration for Monitoring
```json
{
  "metrics": {
    "enabled": true,
    "port": 8080,
    "path": "/metrics",
    "collect_interval": 10,
    "metrics": [
      "throughput",
      "latency",
      "packet_loss",
      "cpu_usage",
      "memory_usage",
      "connection_count"
    ]
  }
}
```

## 🚨 Security Configuration

### Basic Security
```json
{
  "tls": {
    "enabled": true,
    "cert": "/etc/vantun/server.crt",
    "key": "/etc/vantun/server.key"
  },
  "security": {
    "max_connections": 1000,             # Max concurrent connections
    "rate_limit": 100,                   # Connection rate limit
    "timeout": 300,                      # Connection timeout
    "max_packet_size": 65536             # Max packet size
  }
}
```

### Advanced Security
```json
{
  "security": {
    "enabled": true,
    "max_connections": 1000,
    "rate_limit": 100,
    "timeout": 300,
    "max_packet_size": 65536,
    "ip_whitelist": ["192.168.1.0/24"],  # IP whitelist
    "ip_blacklist": ["10.0.0.0/8"],      # IP blacklist
    "geo_blocking": {                    # Geographic blocking
      "enabled": false,
      "allowed_countries": ["US", "CA", "EU"],
      "blocked_countries": ["CN", "RU"]
    },
    "ddos_protection": {                 # DDoS protection
      "enabled": true,
      "threshold": 1000,                 # Packets per second
      "ban_time": 300                    # Ban duration (seconds)
    }
  }
}
```

---

## 🎉 Next Steps

After configuring VANTUN:

1. **[Test Configuration](Testing.md)** - Validate your setup
2. **[Monitor Performance](Monitoring.md)** - Set up monitoring
3. **[Optimize Performance](Performance-Optimization.md)** - Fine-tune settings
4. **[Secure Your Setup](Security.md)** - Implement security measures
5. **[Scale Deployment](Scaling.md)** - Scale for production

---

*For configuration examples and best practices, see [Configuration Examples](Configuration-Examples.md).*