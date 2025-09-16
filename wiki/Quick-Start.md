# Quick Start Guide

Get VANTUN up and running in under 5 minutes with this quick start guide.

## 🚀 One-Command Installation

### Option 1: One-Click Script (Recommended)
```bash
# Install VANTUN with one command
curl -fsSL https://get.vantun.org | bash

# Or with wget
wget -qO- https://get.vantun.org | bash
```

### Option 2: Docker (Easiest)
```bash
# Run server
docker run -d --name vantun-server -p 4242:4242 tungoldshou/vantun:latest server

# Run client
docker run -d --name vantun-client -p 1080:1080 tungoldshou/vantun:latest client -addr SERVER_IP:4242
```

### Option 3: Manual Installation
```bash
# Download latest release
wget https://github.com/tungoldshou/vantun/releases/latest/download/vantun-linux-amd64

# Make executable
chmod +x vantun-linux-amd64

# Move to PATH
sudo mv vantun-linux-amd64 /usr/local/bin/vantun
```

## ⚙️ Basic Configuration

### Server Configuration
Create `/etc/vantun/server.json`:
```json
{
  "server": true,
  "address": "0.0.0.0:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true
}
```

### Client Configuration
Create `/etc/vantun/client.json`:
```json
{
  "server": false,
  "address": "YOUR_SERVER_IP:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "local_addr": "127.0.0.1:1080",
  "socks5": true
}
```

## 🎯 Start VANTUN

### Start Server
```bash
# Using systemd (if installed via script)
sudo systemctl start vantun-server

# Or run directly
vantun --config /etc/vantun/server.json
```

### Start Client
```bash
# Using systemd (if installed via script)
sudo systemctl start vantun-client

# Or run directly
vantun --config /etc/vantun/client.json
```

## ✅ Verify Installation

### Check Server Status
```bash
# Check if server is running
curl -v telnet://YOUR_SERVER_IP:4242

# Check logs
journalctl -u vantun-server -f
```

### Test Client Connection
```bash
# Test SOCKS5 proxy
curl -x socks5://127.0.0.1:1080 https://ipinfo.io

# Check logs
journalctl -u vantun-client -f
```

## 🌐 Use Your Tunnel

### Configure Applications
Set SOCKS5 proxy to `127.0.0.1:1080` in your applications:

#### Browser (Firefox)
1. Settings → Network Settings → Manual proxy configuration
2. SOCKS Host: `127.0.0.1` Port: `1080`
3. Check "SOCKS5" and "Proxy DNS when using SOCKS v5"

#### Browser (Chrome)
```bash
# Run Chrome with proxy
google-chrome --proxy-server="socks5://127.0.0.1:1080"
```

#### System-wide (Linux)
```bash
# Set proxy environment variables
export SOCKS_PROXY="socks5://127.0.0.1:1080"
export ALL_PROXY="socks5://127.0.0.1:1080"
```

## 📊 Performance Check

### Quick Speed Test
```bash
# Test without proxy
speedtest-cli

# Test with proxy
speedtest-cli --proxy socks5://127.0.0.1:1080
```

### Check Connection Info
```bash
# Check your IP through the tunnel
curl -x socks5://127.0.0.1:1080 https://ipinfo.io

# Check tunnel statistics
curl http://localhost:8080/metrics  # If metrics enabled
```

## 🔧 Next Steps

### 1. Optimize Performance
- [Performance Optimization Guide](Performance-Optimization.md)
- [Advanced Configuration](Advanced-Configuration.md)

### 2. Enhance Security
- [Security Features](Security-Features.md)
- [TLS Configuration](TLS-Configuration.md)

### 3. Scale Deployment
- [Docker Deployment](Docker-Deployment.md)
- [Cloud Deployment](Cloud-Deployment.md)
- [Enterprise Setup](Enterprise-Deployment.md)

## 🆘 Troubleshooting

### Common Issues

#### Connection Failed
```bash
# Check if VANTUN is running
sudo systemctl status vantun-client

# Check logs for errors
journalctl -u vantun-client -n 50

# Verify server is reachable
nc -v YOUR_SERVER_IP 4242
```

#### Slow Performance
```bash
# Check network conditions
ping -c 10 YOUR_SERVER_IP

# Try different FEC settings
# Edit /etc/vantun/client.json and adjust fec_data/fec_parity
```

#### High CPU Usage
```bash
# Check VANTUN process
top -p $(pgrep vantun)

# Try reducing workers or disabling multipath
# Edit configuration and restart
```

### Get Help
- 📧 Email: support@vantun.org
- 💬 Telegram: [@vantun01](https://t.me/vantun01)
- 🐛 Issues: [GitHub Issues](https://github.com/tungoldshou/vantun/issues)

## 📚 Further Reading

- [Technical Deep Dive](Technical-Deep-Dive.md) - Understand how VANTUN works
- [Protocol Comparison](Protocol-Comparison.md) - Compare with other solutions
- [Benchmarking](Benchmarking.md) - Performance testing guide

---

**🎉 Congratulations!** You now have VANTUN running. For advanced configuration and optimization, continue to our [Configuration Guide](Configuration.md).