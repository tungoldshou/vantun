# VANTUN Documentation

Welcome to the VANTUN documentation! This site contains comprehensive guides and technical documentation for the next-generation secure tunnel protocol.

## 🚀 Quick Navigation

### 📚 Main Documentation
- **[🏠 Home](index.html)** - Interactive documentation homepage
- **[📖 Detailed README](https://github.com/tungoldshou/vantun/blob/main/README_DETAILED.md)** - Complete project overview with benchmarks
- **[🐳 Docker Deployment](DOCKER_DEPLOYMENT.md)** - Container deployment strategies
- **[🔬 Technical Deep Dive](TECHNICAL_DEEP_DIVE.md)** - Architecture and implementation details

### 🛠️ Installation & Deployment
- **One-Click Script**: `curl -fsSL https://get.vantun.org | bash`
- **Docker**: `docker run -d --name vantun-server -p 4242:4242 tungoldshou/vantun:latest server`
- **Manual**: Download from [GitHub Releases](https://github.com/tungoldshou/vantun/releases)

### 📊 Performance & Benchmarking
- **Benchmark Script**: `./scripts/benchmark.sh`
- **Interactive Charts**: Generated performance visualizations
- **Protocol Comparison**: Detailed analysis vs Hysteria2, V2Ray, WireGuard

## 🎯 Key Features

<div align="center">

| Feature | VANTUN | Hysteria2 | V2Ray | WireGuard |
|---------|--------|-----------|-------|-----------|
| **FEC Support** | ✅ Adaptive | ❌ None | ❌ None | ❌ None |
| **Multipath** | ✅ Intelligent | ❌ Single | ❌ Single | ❌ Single |
| **Obfuscation** | ✅ HTTP/3 | ✅ Brutal | ✅ Various | ❌ None |
| **Performance (15% loss)** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ | ⭐ |
| **Mobile Optimization** | ✅ Excellent | ✅ Good | ⚠️ Fair | ❌ Poor |

</div>

## 📈 Performance Highlights

### Throughput Comparison (100Mbps baseline)
```
Network Condition | VANTUN | Hysteria2 | V2Ray | WireGuard
------------------|--------|-----------|-------|----------
0% Loss           |   98   |    94     |  89   |    85
5% Loss           |   95   |    76     |  62   |    48
10% Loss          |   89   |    58     |  44   |    36
15% Loss          |   82   |    41     |  31   |    25
20% Loss          |   58   |    18     |  15   |    12
```

## 🚀 Quick Start

### Option 1: One-Click Script (Recommended)
```bash
curl -fsSL https://get.vantun.org | bash
```

### Option 2: Docker
```bash
# Server
docker run -d --name vantun-server -p 4242:4242 tungoldshou/vantun:latest server

# Client
docker run -d --name vantun-client -p 1080:1080 tungoldshou/vantun:latest client
```

### Option 3: Docker Compose
```bash
docker-compose up -d
```

## 📚 Documentation Structure

```
vantun/
├── README_DETAILED.md          # Complete project documentation
├── docs/                       # GitHub Pages documentation
│   ├── index.html             # Interactive documentation site
│   ├── DOCKER_DEPLOYMENT.md   # Docker deployment guide
│   ├── TECHNICAL_DEEP_DIVE.md # Technical architecture details
│   └── ...                    # Additional documentation
├── wiki/                      # GitHub Wiki source files
│   ├── Home.md                # Wiki homepage
│   ├── Quick-Start.md         # Quick start guide
│   ├── Protocol-Comparison.md # Performance comparisons
│   ├── Technical-Deep-Dive.md # Architecture details
│   ├── Benchmarking.md        # Testing and benchmarking
│   └── ...                    # Additional Wiki pages
├── scripts/                   # Utility scripts
│   ├── install.sh            # One-click installation
│   ├── benchmark.sh          # Performance benchmarking
│   └── docker-entrypoint.sh  # Docker entry point
└── WIKI_SETUP_GUIDE.md       # GitHub Wiki setup instructions
```

## 🔧 Configuration Example

### Basic Server Configuration
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

### Basic Client Configuration
```json
{
  "server": false,
  "address": "your-server.com:4242",
  "log_level": "info",
  "multipath": true,
  "obfs": true,
  "local_addr": "127.0.0.1:1080",
  "socks5": true
}
```

## 📊 Benchmarking

Run comprehensive performance tests:
```bash
./scripts/benchmark.sh
```

This will generate:
- Throughput comparison charts
- Latency performance analysis
- Connection stability metrics
- Interactive HTML reports

## 🌐 Community & Support

- **💬 Telegram**: [@vantun01](https://t.me/vantun01)
- **🐛 Issues**: [GitHub Issues](https://github.com/tungoldshou/vantun/issues)
- **📧 Email**: support@vantun.org
- **📖 Wiki**: [GitHub Wiki](https://github.com/tungoldshou/vantun/wiki)

## 🔗 Quick Links

- **[GitHub Repository](https://github.com/tungoldshou/vantun)** - Source code and releases
- **[Interactive Documentation](index.html)** - Beautiful HTML documentation
- **[Docker Hub](https://hub.docker.com/r/tungoldshou/vantun)** - Container images
- **[Performance Charts](index.html#benchmarking)** - Interactive performance visualizations

---

*This documentation is automatically deployed via GitHub Actions. For the latest updates, please check the [GitHub repository](https://github.com/tungoldshou/vantun).*