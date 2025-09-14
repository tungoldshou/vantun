# VANTUN - High-Performance Tunnel Protocol Based on QUIC

[//]: # (VANTUN是一个基于QUIC的高性能隧道协议实现，具备以下核心功能：)

VANTUN is a high-performance tunnel protocol implementation based on QUIC with the following core features:

## 核心特性 / Key Features

1. **安全握手与会话协商**：通过control stream进行
   <br>**Secure Handshake and Session Negotiation**: Conducted via control stream
2. **多类型逻辑流**：支持interactive（交互）、bulk（批量）、telemetry（遥测）三种流类型
   <br>**Multiple Logical Stream Types**: Supports interactive, bulk, and telemetry stream types
3. **FEC (前向纠错)**：增强数据传输的鲁棒性
   <br>**Forward Error Correction (FEC)**: Enhances data transmission robustness
4. **多路径 (Multipath)**：利用多条网络路径提升性能和可靠性
   <br>**Multipath**: Utilizes multiple network paths to improve performance and reliability
5. **上层拥塞控制 (Hybrid CC)**：结合底层QUIC CC和上层Token-Bucket限流
   <br>**Hybrid Congestion Control**: Combines underlying QUIC CC with upper-layer token bucket rate limiting
6. **可插拔混淆模块**：降低被误判为非法流量的风险
   <br>**Pluggable Obfuscation Module**: Reduces the risk of being misidentified as illegal traffic
7. **最小可用客户端/服务端**：提供可通过命令行运行的`client`和`server`程序
   <br>**Minimal Client/Server**: Provides `client` and `server` programs that can be run via command line

## 技术栈 / Tech Stack

- **语言 / Language**: Go
- **核心库 / Core Library**: `quic-go` for QUIC implementation
- **序列化 / Serialization**: `github.com/fxamacker/cbor` for CBOR encoding
- **FEC**: `github.com/klauspost/reedsolomon` for Reed-Solomon encoding
- **CLI**: `cobra/viper` for command-line interface and configuration management

## 快速开始 / Quick Start

请参考 [DEMOGUIDE.md](DEMOGUIDE.md) 获取详细的构建、配置和运行说明。
<br>Please refer to [DEMOGUIDE.md](DEMOGUIDE.md) for detailed build, configuration, and running instructions.

## 项目结构 / Project Structure

```
vantun/
├── cmd/              # 命令行程序入口 / Command-line program entry
├── internal/
│   ├── cli/          # CLI配置管理 / CLI configuration management
│   └── core/         # 核心协议实现 / Core protocol implementation
├── go.mod            # Go模块定义 / Go module definition
└── README.md         # 项目说明文档 / Project documentation
```

## 核心模块 / Core Modules

### 1. 协议基础 / Protocol Foundation
- 会话协商和控制流
  <br>Session negotiation and control stream
- 数据流类型管理
  <br>Data stream type management

### 2. FEC模块 / FEC Module
- Reed-Solomon编码/解码
  <br>Reed-Solomon encoding/decoding
- 数据恢复能力
  <br>Data recovery capability

### 3. 多路径支持 / Multipath Support
- 路径管理
  <br>Path management
- 负载均衡策略
  <br>Load balancing strategies

### 4. 拥塞控制 / Congestion Control
- 令牌桶限流
  <br>Token bucket rate limiting
- 自适应调整
  <br>Adaptive adjustment

### 5. 混淆模块 / Obfuscation Module
- HTTP/3风格流量伪装
  <br>HTTP/3-style traffic obfuscation
- 数据填充
  <br>Data padding

### 6. 遥测系统 / Telemetry System
- 性能数据收集
  <br>Performance data collection
- 实时监控
  <br>Real-time monitoring

## 测试 / Testing

运行单元测试：
<br>Run unit tests:

```bash
go test -v ./internal/core/...
```

## 许可证 / License

[待定]
<br>[To be determined]