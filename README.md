# VANTUN - Next-Generation Secure Tunnel Protocol

[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Release](https://img.shields.io/github/v/release/tungoldshou/vantun)](https://github.com/tungoldshou/vantun/releases)
[![Telegram](https://img.shields.io/badge/telegram-vantun01-blue?logo=telegram)](https://t.me/vantun01)
[![Go Version](https://img.shields.io/badge/go-1.22-blue)](https://golang.org/)

![Networking](https://img.shields.io/badge/tag-networking-blue)
![QUIC](https://img.shields.io/badge/tag-quic-green)
![Proxy Protocol](https://img.shields.io/badge/tag-proxy--protocol-orange)
![VPN](https://img.shields.io/badge/tag-vpn-purple)
![Privacy](https://img.shields.io/badge/tag-privacy-red)

## 🌟 Core Highlights

- **30% more stable than Hysteria2 in high packet loss environments** - Advanced FEC technology maintains connection stability even with 15% packet loss
- **HTTP/3 traffic camouflage** - Traffic appears as regular web browsing, blending seamlessly with normal internet activity
- **Multipath transmission** - Automatically utilizes multiple network paths for enhanced reliability and performance
- **Enterprise-grade security** - Built on QUIC protocol with additional encryption layers for maximum privacy protection

## Documentation Index

For detailed documentation, please refer to the following files:

- [English](docs/README_en.md)
- [中文](docs/README_zh.md)
- [日本語](docs/README_ja.md)
- [Français](docs/README_fr.md)
- [Deutsch](docs/README_de.md)
- [Español](docs/README_es.md)
- [Русский](docs/README_ru.md)
- [Português](docs/README_pt.md)
- [한국어](docs/README_ko.md)
- [العربية](docs/README_ar.md)

## Project Overview

VANTUN is a cutting-edge, high-performance tunnel protocol built on top of QUIC, designed to deliver exceptional network performance, security, and reliability. As a next-generation solution, VANTUN redefines what's possible in network tunneling with its innovative architecture and advanced features.

## Key Advantages

### 🔒 Enterprise-Grade Security
- **Secure Handshake and Session Negotiation**: Conducted via dedicated control stream for connection security

### ⚡ Exceptional Performance
- **Multiple Logical Stream Types**: Optimized interactive, bulk, and telemetry streams for different business scenarios
- **Multipath**: Intelligent use of multiple network paths for dramatically improved speed and connection stability

### 🛡️ Unmatched Reliability
- **Forward Error Correction (FEC)**: Advanced error correction ensures data integrity even in unstable network conditions
- **Hybrid Congestion Control**: Innovative hybrid algorithm combining QUIC CC with token bucket rate limiting for optimal resource utilization

### 🌐 Privacy Protection
- **Pluggable Obfuscation Module**: Advanced traffic obfuscation makes traffic appear as normal HTTP/3, effectively evading network scrutiny

### 🚀 Easy Deployment
- **Minimal Client/Server**: Command-line `client` and `server` programs for rapid deployment and ease of use

## Technology Architecture

VANTUN leverages industry-leading technologies to deliver its exceptional performance and reliability:

- **Language**: Go - High-performance, concurrent modern programming language
- **Core Library**: `quic-go` - Industry-leading QUIC protocol implementation
- **Serialization**: `github.com/fxamacker/cbor` - Efficient CBOR encoding, more compact than JSON
- **FEC**: `github.com/klauspost/reedsolomon` - High-performance Reed-Solomon encoding algorithm
- **CLI**: `cobra/viper` - Powerful command-line interface and configuration management

## Quick Start

Get VANTUN up and running in just a few minutes:

1. **Clone Repository**: `git clone <repository-url>`
2. **Build**: `go build -o bin/vantun cmd/main.go`
3. **Configure**: Create `config.json` configuration file
4. **Run**: Start server and client

For detailed steps and configuration instructions, please refer to [Demo Guide](docs/DEMOGUIDE_en.md).

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

## Architecture Highlights

### 🔧 Intelligent Protocol Engine
The core protocol engine implements efficient session negotiation and control stream management for secure and stable connections.

### 📊 Adaptive FEC Technology
Forward error correction based on Reed-Solomon encoding that dynamically adjusts correction strategies based on network conditions.

### 🔄 Intelligent Multipath Transmission
Innovative path management and load balancing that fully utilizes all available network paths for redundancy and enhanced throughput.

### 📈 Hybrid Congestion Control
Hybrid algorithm combining underlying QUIC congestion control with upper-layer token bucket for optimal resource utilization.

### 🎭 Advanced Traffic Obfuscation
HTTP/3-style traffic obfuscation and intelligent data padding to effectively evade network scrutiny and protect user privacy.

### 📊 Real-time Telemetry System
Comprehensive performance data collection and real-time monitoring for network optimization and troubleshooting.

## Quality Assurance

VANTUN adopts strict testing standards to ensure code quality and system stability:

- **Comprehensive Unit Tests**: Covering all core functional modules
- **Integration Tests**: Validating component collaboration
- **Performance Tests**: Ensuring exceptional performance under various network conditions
- **Stress Tests**: Validating stability under high load

Run all tests:

```bash
go test -v ./internal/core/...
```

## License

VANTUN is licensed under the MIT License, a permissive open-source license that allows free use, copying, modification, and distribution of the software while retaining the copyright and license notices.

---

*© 2025 VANTUN Project. All rights reserved.*