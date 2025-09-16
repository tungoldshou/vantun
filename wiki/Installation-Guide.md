# Installation Guide

## 🎯 Overview

This guide provides comprehensive installation instructions for VANTUN across different platforms and deployment methods.

## 📋 Prerequisites

### System Requirements
- **Operating System**: Linux (Ubuntu/Debian/CentOS/RHEL/Alpine/Arch), macOS, Windows
- **Architecture**: x86_64, ARM64, ARMv7
- **Memory**: Minimum 64MB RAM, Recommended 256MB+
- **Storage**: 50MB free space
- **Network**: UDP port access (default: 4242)

### Network Requirements
- **Firewall**: Allow UDP traffic on configured ports
- **Root Access**: Required for system service installation
- **Internet**: For downloading binaries and updates

## 🚀 Installation Methods

### Method 1: One-Click Script (Recommended)

#### Linux/macOS
```bash
# Download and run installation script
curl -fsSL https://get.vantun.org | bash

# Or with wget
wget -qO- https://get.vantun.org | bash
```

#### Windows (PowerShell)
```powershell
# Download and run installation script
iwr -useb https://get.vantun.org | iex
```

**What the script does:**
- Detects your operating system and architecture
- Downloads the appropriate binary
- Installs VANTUN to `/usr/local/bin`
- Creates configuration directory (`/etc/vantun`)
- Sets up systemd services (Linux)
- Generates default configuration files

### Method 2: Docker Installation

#### Basic Docker
```bash
# Pull the latest image
docker pull tungoldshou/vantun:latest

# Run server
docker run -d \
  --name vantun-server \
  --restart unless-stopped \
  -p 4242:4242 \
  -v /etc/vantun:/etc/vantun \
  tungoldshou/vantun:latest \
  server --config /etc/vantun/server.json
```

#### Docker Compose
```bash
# Clone repository
git clone https://github.com/tungoldshou/vantun.git
cd vantun

# Start with Docker Compose
docker-compose up -d
```

### Method 3: Package Managers

#### Alpine Linux
```bash
# Add VANTUN repository
echo "https://vantun.org/alpine/v3.14/main" >> /etc/apk/repositories
wget -O /etc/apk/keys/vantun.rsa.pub https://vantun.org/alpine/vantun.rsa.pub

# Install VANTUN
apk update
apk add vantun
```

#### Arch Linux (AUR)
```bash
# Using yay
yay -S vantun

# Or manually
git clone https://aur.archlinux.org/vantun.git
cd vantun
makepkg -si
```

### Method 4: Manual Installation

#### Download Binary
```bash
# Get latest release version
VERSION=$(curl -s https://api.github.com/repos/tungoldshou/vantun/releases/latest | grep tag_name | cut -d'"' -f4)

# Download for your architecture
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
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Download binary
wget "https://github.com/tungoldshou/vantun/releases/download/${VERSION}/vantun-linux-${ARCH}"

# Make executable and install
chmod +x vantun-linux-${ARCH}
sudo mv vantun-linux-${ARCH} /usr/local/bin/vantun
```

#### Build from Source
```bash
# Prerequisites
sudo apt-get update
sudo apt-get install -y golang git make

# Clone repository
git clone https://github.com/tungoldshou/vantun.git
cd vantun

# Build
make build

# Install
sudo make install
```

## ⚙️ Post-Installation Setup

### 1. Create Configuration Directory
```bash
sudo mkdir -p /etc/vantun
sudo chmod 755 /etc/vantun
```

### 2. Generate Configuration Files

#### Server Configuration
```bash
sudo tee /etc/vantun/server.json > /dev/null << 'EOF'
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3,
  "token_bucket_rate": 1000000,
  "token_bucket_capacity": 5000000
}
EOF
```

#### Client Configuration
```bash
sudo tee /etc/vantun/client.json > /dev/null << 'EOF'
{
  "server": false,
  "address": "your-server.com:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "fec_data": 10,
  "fec_parity": 3,
  "local_addr": "127.0.0.1:1080",
  "socks5": true
}
EOF
```

### 3. System Service Setup (Linux)

#### Systemd Service (Ubuntu/Debian/CentOS)
```bash
# Create service file
sudo tee /etc/systemd/system/vantun-server.service > /dev/null << 'EOF'
[Unit]
Description=VANTUN Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/etc/vantun
ExecStart=/usr/local/bin/vantun server --config /etc/vantun/server.json
Restart=on-failure
RestartSec=5s
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable vantun-server
sudo systemctl start vantun-server
```

#### Launchd Service (macOS)
```bash
# Create service file
sudo tee /Library/LaunchDaemons/com.vantun.server.plist > /dev/null << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.vantun.server</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/vantun</string>
        <string>server</string>
        <string>--config</string>
        <string>/etc/vantun/server.json</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/vantun-server.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/vantun-server-error.log</string>
</dict>
</plist>
EOF

# Load and start service
sudo launchctl load /Library/LaunchDaemons/com.vantun.server.plist
sudo launchctl start com.vantun.server
```

## 🔍 Verification

### Check Installation
```bash
# Verify binary installation
which vantun
vantun --version

# Check configuration
ls -la /etc/vantun/
cat /etc/vantun/server.json
```

### Test Service Status
```bash
# Systemd (Linux)
sudo systemctl status vantun-server

# Launchd (macOS)
sudo launchctl list | grep vantun

# Check logs
sudo journalctl -u vantun-server -f  # Linux
sudo tail -f /var/log/vantun-server.log  # macOS
```

### Network Test
```bash
# Test if server is listening
sudo netstat -tulnp | grep 4242  # Linux
sudo lsof -i :4242  # macOS

# Test connectivity
nc -z -v localhost 4242
```

## 🛠️ Advanced Installation

### Multi-Instance Setup
```bash
# Create multiple configuration directories
sudo mkdir -p /etc/vantun/{instance1,instance2,instance3}

# Create separate configurations
for i in {1..3}; do
    sudo cp /etc/vantun/server.json /etc/vantun/instance$i/server.json
    sudo sed -i "s/4242/$((4242 + i))/g" /etc/vantun/instance$i/server.json
done

# Create separate services
for i in {1..3}; do
    sudo cp /etc/systemd/system/vantun-server.service /etc/systemd/system/vantun-server$i.service
    sudo sed -i "s/server.json/instance$i\/server.json/g" /etc/systemd/system/vantun-server$i.service
done

# Start all instances
for i in {1..3}; do
    sudo systemctl enable vantun-server$i
    sudo systemctl start vantun-server$i
done
```

### High-Availability Setup
```bash
# Install keepalived for failover
sudo apt-get install keepalived  # Ubuntu/Debian
sudo yum install keepalived      # CentOS/RHEL

# Configure keepalived
sudo tee /etc/keepalived/keepalived.conf > /dev/null << 'EOF'
vrrp_instance VI_1 {
    state MASTER
    interface eth0
    virtual_router_id 51
    priority 100
    advert_int 1
    
    authentication {
        auth_type PASS
        auth_pass vantun123
    }
    
    virtual_ipaddress {
        192.168.1.100/24
    }
}
EOF

# Start keepalived
sudo systemctl enable keepalived
sudo systemctl start keepalived
```

## 🔧 Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Find process using port 4242
sudo lsof -i :4242
sudo netstat -tulnp | grep 4242

# Kill the process
sudo kill -9 <PID>
```

#### Permission Denied
```bash
# Fix permissions
sudo chown -R root:root /etc/vantun
sudo chmod 755 /etc/vantun
sudo chmod 644 /etc/vantun/*.json

# Fix binary permissions
sudo chmod 755 /usr/local/bin/vantun
```

#### Service Won't Start
```bash
# Check logs
sudo journalctl -u vantun-server -n 50

# Test configuration
vantun --config /etc/vantun/server.json --check

# Verify network connectivity
ping -c 4 your-server.com
```

#### High CPU Usage
```bash
# Monitor CPU usage
top -p $(pgrep vantun)
htop

# Check for configuration issues
vantun --config /etc/vantun/server.json --log-level debug
```

#### Memory Issues
```bash
# Monitor memory usage
free -h
ps aux | grep vantun

# Check for memory leaks
valgrind --tool=memcheck vantun --config /etc/vantun/server.json
```

## 🔄 Updates

### Update VANTUN
```bash
# Stop service
sudo systemctl stop vantun-server

# Update using script
curl -fsSL https://get.vantun.org | bash

# Or update manually
wget https://github.com/tungoldshou/vantun/releases/latest/download/vantun-linux-amd64
sudo mv vantun-linux-amd64 /usr/local/bin/vantun
sudo chmod +x /usr/local/bin/vantun

# Restart service
sudo systemctl start vantun-server
```

### Rollback Updates
```bash
# Stop service
sudo systemctl stop vantun-server

# Restore previous version
sudo cp /usr/local/bin/vantun.backup /usr/local/bin/vantun

# Restart service
sudo systemctl start vantun-server
```

## 🗂️ Directory Structure

After installation:
```
/etc/vantun/                 # Configuration directory
├── server.json             # Server configuration
├── client.json             # Client configuration
├── server.crt              # TLS certificate (if applicable)
└── server.key              # TLS private key (if applicable)

/usr/local/bin/vantun       # Binary executable

/var/log/vantun/            # Log directory
├── server.log              # Server logs
└── client.log              # Client logs

/etc/systemd/system/        # Systemd services (Linux)
├── vantun-server.service   # Server service
└── vantun-client.service   # Client service
```

---

## 🎉 Next Steps

After successful installation:

1. **[Configure VANTUN](Configuration.md)** - Customize settings for your needs
2. **[Test the Installation](Testing.md)** - Verify everything works correctly
3. **[Secure Your Setup](Security.md)** - Implement security best practices
4. **[Monitor Performance](Monitoring.md)** - Set up monitoring and logging
5. **[Optimize Performance](Performance-Optimization.md)** - Fine-tune for your environment

---

*For additional help, visit our [Community Forum](Community.md) or [GitHub Issues](https://github.com/tungoldshou/vantun/issues).*