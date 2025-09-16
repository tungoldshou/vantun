#!/bin/sh
# VANTUN Docker Entrypoint Script

set -e

# Configuration
CONFIG_DIR="/etc/vantun"
LOG_DIR="/var/log/vantun"
DEFAULT_CONFIG="$CONFIG_DIR/config.json"

# Create config directory if it doesn't exist
mkdir -p "$CONFIG_DIR" "$LOG_DIR"

# Generate default configuration if none exists
if [ ! -f "$DEFAULT_CONFIG" ]; then
    echo "Generating default configuration..."
    
    case "$1" in
        "server")
            cat > "$DEFAULT_CONFIG" << 'EOF'
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info",
  "log_file": "/var/log/vantun/server.log",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3,
  "token_bucket_rate": 1000000,
  "token_bucket_capacity": 5000000
}
EOF
            ;;
        "client")
            cat > "$DEFAULT_CONFIG" << 'EOF'
{
  "server": false,
  "address": "server.example.com:4242",
  "log_level": "info",
  "log_file": "/var/log/vantun/client.log",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3,
  "token_bucket_rate": 1000000,
  "token_bucket_capacity": 5000000,
  "local_addr": "127.0.0.1:1080",
  "socks5": true
}
EOF
            ;;
        *)
            cat > "$DEFAULT_CONFIG" << 'EOF'
{
  "server": false,
  "address": "localhost:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3
}
EOF
            ;;
    esac
fi

# Set proper permissions
chmod 644 "$DEFAULT_CONFIG"

# Handle different command modes
case "$1" in
    "server")
        echo "Starting VANTUN server..."
        shift
        exec vantun server --config "$DEFAULT_CONFIG" "$@"
        ;;
    "client")
        echo "Starting VANTUN client..."
        shift
        exec vantun client --config "$DEFAULT_CONFIG" "$@"
        ;;
    "version"|"-v"|"--version")
        exec vantun --version
        ;;
    "help"|"-h"|"--help")
        exec vantun --help
        ;;
    *)
        # Pass through all arguments to vantun binary
        exec vantun "$@"
        ;;
esac