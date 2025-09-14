# VANTUN - 基于QUIC的高性能隧道协议

VANTUN是一个基于QUIC的高性能隧道协议实现，具备以下核心功能：

## 核心特性

1. **安全握手与会话协商**：通过control stream进行
2. **多类型逻辑流**：支持interactive（交互）、bulk（批量）、telemetry（遥测）三种流类型
3. **FEC (前向纠错)**：增强数据传输的鲁棒性
4. **多路径 (Multipath)**：利用多条网络路径提升性能和可靠性
5. **上层拥塞控制 (Hybrid CC)**：结合底层QUIC CC和上层Token-Bucket限流
6. **可插拔混淆模块**：降低被误判为非法流量的风险
7. **最小可用客户端/服务端**：提供可通过命令行运行的`client`和`server`程序

## 技术栈

- **语言**: Go
- **核心库**: `quic-go`用于QUIC实现
- **序列化**: `github.com/fxamacker/cbor`用于CBOR编码
- **FEC**: `github.com/klauspost/reedsolomon`用于Reed-Solomon编码
- **CLI**: `cobra/viper`用于命令行接口和配置管理

## 快速开始

### 构建项目

```bash
go build -o bin/vantun cmd/main.go
```

### 运行服务端

```bash
./bin/vantun -server -addr :4242
```

### 运行客户端

```bash
./bin/vantun -addr localhost:4242
```

## 配置选项

可以通过命令行参数或配置文件来配置VANTUN：

- `-server`: 运行服务端模式
- `-addr`: 监听地址（服务端）或连接地址（客户端）
- `-config`: JSON配置文件路径
- `-log-level`: 日志级别（debug, info, warn, error）
- `-multipath`: 启用多路径
- `-obfs`: 启用混淆
- `-fec-data`: FEC数据分片数
- `-fec-parity`: FEC奇偶分片数

## 项目结构

```
vantun/
├── cmd/              # 命令行程序入口
├── internal/
│   ├── cli/          # CLI配置管理
│   └── core/         # 核心协议实现
├── go.mod            # Go模块定义
└── README.md         # 项目说明文档
```

## 核心模块

### 1. 协议基础
- 会话协商和控制流
- 数据流类型管理

### 2. FEC模块
- Reed-Solomon编码/解码
- 数据恢复能力

### 3. 多路径支持
- 路径管理
- 负载均衡策略

### 4. 拥塞控制
- 令牌桶限流
- 自适应调整

### 5. 混淆模块
- HTTP/3风格流量伪装
- 数据填充

### 6. 遥测系统
- 性能数据收集
- 实时监控

## 测试

运行单元测试：

```bash
go test -v ./internal/core/...
```

## 许可证

[待定]