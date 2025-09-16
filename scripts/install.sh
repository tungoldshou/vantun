#!/bin/bash

# VANTUN One-Click Installation Script
# Supports: Ubuntu/Debian, CentOS/RHEL, Alpine, Arch Linux
# Version: 1.0.0

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
VANTUN_VERSION="latest"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/vantun"
SERVICE_DIR="/etc/systemd/system"
GITHUB_REPO="tungoldshou/vantun"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root"
        exit 1
    fi
}

# Detect OS
detect_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
    elif type lsb_release >/dev/null 2>&1; then
        OS=$(lsb_release -si | tr '[:upper:]' '[:lower:]')
        VER=$(lsb_release -sr)
    elif [[ -f /etc/lsb-release ]]; then
        . /etc/lsb-release
        OS=$DISTRIB_ID
        VER=$DISTRIB_RELEASE
    else
        OS=$(uname -s)
        VER=$(uname -r)
    fi
    
    log_info "Detected OS: $OS $VER"
}

# Check system requirements
check_requirements() {
    log_info "Checking system requirements..."
    
    # Check architecture
    ARCH=$(uname -m)
    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l|armhf)
            ARCH="armv7"
            ;;
        *)
            log_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    # Check available memory (minimum 64MB)
    MEMORY=$(free -m | awk 'NR==2{print $2}')
    if [[ $MEMORY -lt 64 ]]; then
        log_warning "Low memory detected: ${MEMORY}MB. VANTUN requires at least 64MB"
    fi
    
    # Check disk space (minimum 50MB)
    DISK=$(df / | awk 'NR==2{print $4}')
    if [[ $DISK -lt 51200 ]]; then
        log_error "Insufficient disk space: ${DISK}KB. At least 50MB required"
        exit 1
    fi
    
    log_success "System requirements check passed"
}

# Install dependencies
install_dependencies() {
    log_info "Installing dependencies..."
    
    case $OS in
        ubuntu|debian)
            apt-get update -qq
            apt-get install -y curl wget ca-certificates
            ;;
        centos|rhel|fedora)
            yum install -y curl wget ca-certificates
            ;;
        alpine)
            apk add --no-cache curl wget ca-certificates
            ;;
        arch)
            pacman -Sy --noconfirm curl wget ca-certificates
            ;;
        *)
            log_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    
    log_success "Dependencies installed"
}

# Download VANTUN
download_vantun() {
    log_info "Downloading VANTUN..."
    
    # Get latest release if version is "latest"
    if [[ $VANTUN_VERSION == "latest" ]]; then
        VANTUN_VERSION=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    
    DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$VANTUN_VERSION/vantun-linux-$ARCH"
    TEMP_FILE=$(mktemp)
    
    log_info "Downloading from: $DOWNLOAD_URL"
    
    if curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
        chmod +x "$TEMP_FILE"
        log_success "Download completed"
    else
        log_error "Download failed"
        exit 1
    fi
}

# Install VANTUN
install_vantun() {
    log_info "Installing VANTUN..."
    
    # Create directories
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$CONFIG_DIR/scripts"
    
    # Move binary
    mv "$TEMP_FILE" "$INSTALL_DIR/vantun"
    chmod +x "$INSTALL_DIR/vantun"
    
    # Create symlink
    ln -sf "$INSTALL_DIR/vantun" /usr/local/bin/vantun
    
    log_success "VANTUN installed to $INSTALL_DIR/vantun"
}

# Create configuration
create_config() {
    log_info "Creating configuration..."
    
    # Server configuration
    cat > "$CONFIG_DIR/server.json" << EOF
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3,
  "token_bucket_rate": 1000000,
  "token_bucket_capacity": 5000000,
  "tls": {
    "cert": "$CONFIG_DIR/server.crt",
    "key": "$CONFIG_DIR/server.key"
  }
}
EOF

    # Client configuration
    cat > "$CONFIG_DIR/client.json" << EOF
{
  "server": false,
  "address": "your-server.com:4242",
  "log_level": "info",
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

    # Generate self-signed certificate for testing
    if command -v openssl >/dev/null 2>&1; then
        log_info "Generating self-signed certificate..."
        openssl req -x509 -newkey rsa:4096 -keyout "$CONFIG_DIR/server.key" -out "$CONFIG_DIR/server.crt" -days 365 -nodes -subj "/C=US/ST=State/L=City/O=VANTUN/CN=vantun.local" 2>/dev/null
        chmod 600 "$CONFIG_DIR/server.key"
        log_success "Self-signed certificate generated"
    else
        log_warning "OpenSSL not found. Please provide your own certificates."
    fi
    
    log_success "Configuration files created in $CONFIG_DIR"
}

# Create systemd service
create_service() {
    log_info "Creating systemd service..."
    
    cat > "$SERVICE_DIR/vantun-server.service" << EOF
[Unit]
Description=VANTUN Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$CONFIG_DIR
ExecStart=$INSTALL_DIR/vantun --config $CONFIG_DIR/server.json
Restart=on-failure
RestartSec=5s
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOF

    cat > "$SERVICE_DIR/vantun-client.service" << EOF
[Unit]
Description=VANTUN Client
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$CONFIG_DIR
ExecStart=$INSTALL_DIR/vantun --config $CONFIG_DIR/client.json
Restart=on-failure
RestartSec=5s
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOF

    # Reload systemd
    systemctl daemon-reload
    
    log_success "Systemd services created"
}

# Create management scripts
create_scripts() {
    log_info "Creating management scripts..."
    
    # Start script
    cat > "$CONFIG_DIR/scripts/start.sh" << 'EOF'
#!/bin/bash
# VANTUN Start Script

case "$1" in
    server)
        systemctl start vantun-server
        echo "VANTUN server started"
        ;;
    client)
        systemctl start vantun-client
        echo "VANTUN client started"
        ;;
    *)
        echo "Usage: $0 {server|client}"
        exit 1
        ;;
esac
EOF

    # Stop script
    cat > "$CONFIG_DIR/scripts/stop.sh" << 'EOF'
#!/bin/bash
# VANTUN Stop Script

case "$1" in
    server)
        systemctl stop vantun-server
        echo "VANTUN server stopped"
        ;;
    client)
        systemctl stop vantun-client
        echo "VANTUN client stopped"
        ;;
    *)
        echo "Usage: $0 {server|client}"
        exit 1
        ;;
esac
EOF

    # Status script
    cat > "$CONFIG_DIR/scripts/status.sh" << 'EOF'
#!/bin/bash
# VANTUN Status Script

case "$1" in
    server)
        systemctl status vantun-server
        ;;
    client)
        systemctl status vantun-client
        ;;
    *)
        echo "Usage: $0 {server|client}"
        exit 1
        ;;
esac
EOF

    # Log script
    cat > "$CONFIG_DIR/scripts/logs.sh" << 'EOF'
#!/bin/bash
# VANTUN Logs Script

case "$1" in
    server)
        journalctl -u vantun-server -f
        ;;
    client)
        journalctl -u vantun-client -f
        ;;
    *)
        echo "Usage: $0 {server|client}"
        exit 1
        ;;
esac
EOF

    chmod +x "$CONFIG_DIR/scripts/"*.sh
    
    log_success "Management scripts created"
}

# Create uninstall script
create_uninstall() {
    log_info "Creating uninstall script..."
    
    cat > "$CONFIG_DIR/scripts/uninstall.sh" << 'EOF'
#!/bin/bash
# VANTUN Uninstall Script

read -p "Are you sure you want to uninstall VANTUN? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Uninstallation cancelled."
    exit 1
fi

echo "Stopping VANTUN services..."
systemctl stop vantun-server vantun-client 2>/dev/null || true
systemctl disable vantun-server vantun-client 2>/dev/null || true

echo "Removing files..."
rm -f /etc/systemd/system/vantun-*.service
rm -f /usr/local/bin/vantun
rm -rf /etc/vantun

echo "Reloading systemd..."
systemctl daemon-reload

echo "VANTUN has been successfully uninstalled."
EOF

    chmod +x "$CONFIG_DIR/scripts/uninstall.sh"
    
    log_success "Uninstall script created"
}

# Print usage information
print_usage() {
    echo
    log_success "VANTUN installation completed!"
    echo
    echo "Usage:"
    echo "  Server: vantun --config /etc/vantun/server.json"
    echo "  Client: vantun --config /etc/vantun/client.json"
    echo
    echo "Management scripts:"
    echo "  Start:   /etc/vantun/scripts/start.sh {server|client}"
    echo "  Stop:    /etc/vantun/scripts/stop.sh {server|client}"
    echo "  Status:  /etc/vantun/scripts/status.sh {server|client}"
    echo "  Logs:    /etc/vantun/scripts/logs.sh {server|client}"
    echo "  Uninstall: /etc/vantun/scripts/uninstall.sh"
    echo
    echo "Systemd services:"
    echo "  systemctl {start|stop|restart|status} vantun-server"
    echo "  systemctl {start|stop|restart|status} vantun-client"
    echo
    echo "Configuration files:"
    echo "  Server: /etc/vantun/server.json"
    echo "  Client: /etc/vantun/client.json"
    echo
    echo "Remember to:"
    echo "1. Edit /etc/vantun/client.json and change 'your-server.com' to your actual server address"
    echo "2. Configure firewall rules if necessary"
    echo "3. Generate proper TLS certificates for production use"
    echo
}

# Main installation function
main() {
    echo "========================================"
    echo "    VANTUN One-Click Installer"
    echo "========================================"
    echo
    
    check_root
    detect_os
    check_requirements
    install_dependencies
    download_vantun
    install_vantun
    create_config
    create_service
    create_scripts
    create_uninstall
    print_usage
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "VANTUN One-Click Installer"
        echo "Usage: $0 [options]"
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --version VER  Install specific version (default: latest)"
        exit 0
        ;;
    --version)
        VANTUN_VERSION="$2"
        shift 2
        ;;
esac

# Run main function
main "$@"