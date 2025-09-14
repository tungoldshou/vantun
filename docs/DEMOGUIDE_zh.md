# VANTUN 演示与使用指南

本指南将帮助您快速上手VANTUN，从编译、配置到运行服务端和客户端。

## 1. 环境准备与编译

### 环境要求
- Go 1.21 或更高版本
- 支持的操作系统 (Linux, macOS, Windows)

### 编译步骤

```bash
# 克隆仓库 (如果尚未克隆)
# git clone <repository-url>
cd vantun

# 编译项目
go build -o bin/vantun cmd/main.go

# 验证编译结果
./bin/vantun -h
```

编译成功后，您将在`bin`目录下看到可执行文件`vantun`。

## 2. 配置详解

VANTUN 提供灵活的配置选项，支持通过命令行参数和 JSON 配置文件两种方式进行配置。

### 2.1 命令行参数详解

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-server` | 以服务端模式运行 | false (客户端模式) |
| `-addr` | 监听地址（服务端）或连接地址（客户端） | `localhost:4242` |
| `-config` | JSON 配置文件路径 | 无 |
| `-log-level` | 日志级别 (debug, info, warn, error) | info |
| `-multipath` | 启用多路径传输 | false |
| `-obfs` | 启用流量混淆 | false |
| `-fec-data` | FEC 数据分片数 | 10 |
| `-fec-parity` | FEC 冗余分片数 | 3 |

### 2.2 JSON 配置文件详解

创建一个 `config.json` 文件，包含以下配置项：

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

配置项说明：
- `server`: 运行模式 (true=服务端, false=客户端)
- `address`: 监听/连接地址
- `log_level`: 日志级别
- `multipath`: 是否启用多路径传输
- `obfs`: 是否启用流量混淆
- `fec_data`: FEC数据分片数
- `fec_parity`: FEC冗余分片数
- `token_bucket_rate`: 令牌桶速率 (字节/秒)
- `token_bucket_capacity`: 令牌桶容量 (字节)

你也可以参考 `test_configs/` 目录下的示例配置文件：
- `fec_client.json` 和 `fec_server.json`: FEC功能测试配置
- `obfs_client.json` 和 `obfs_server.json`: 流量混淆功能测试配置
- `multipath_client.json` 和 `multipath_server.json`: 多路径功能测试配置

## 3. 运行演示

### 3.1 启动服务端

打开一个终端窗口，执行以下命令启动VANTUN服务端：

```bash
# 使用命令行参数启动服务端
./bin/vantun -server -addr :4242

# 或使用配置文件启动服务端
./bin/vantun -config config.json -server
```

服务端启动后，您将看到类似以下的输出：
```
2025/09/14 04:00:00 Server running, waiting for streams...
```

### 3.2 启动客户端

打开另一个终端窗口，执行以下命令启动VANTUN客户端：

```bash
# 使用命令行参数启动客户端
./bin/vantun -addr localhost:4242

# 或使用配置文件启动客户端
./bin/vantun -config config.json
```

客户端启动后，您将看到类似以下的输出：
```
2025/09/14 04:00:00 Client connected, opening interactive stream...
2025/09/14 04:00:00 Received echo: Hello from VANTUN client!
```

这表明客户端已成功连接到服务端，并完成了首次数据交互。

## 4. 连接验证

当客户端和服务端都成功启动后，您将看到以下输出信息：

### 服务端输出:
```
2025/09/14 04:00:00 Server running, waiting for streams...
2025/09/14 04:00:00 Accepted interactive stream, echoing data...
```

### 客户端输出:
```
2025/09/14 04:00:00 Client connected, opening interactive stream...
2025/09/14 04:00:00 Received echo: Hello from VANTUN client!
```

这些输出表明：
1. 客户端成功连接到服务端
2. 客户端打开了一个交互流
3. 客户端发送了一条消息到服务端
4. 服务端接收并回显了该消息
5. 客户端成功接收到了回显消息

这验证了VANTUN隧道已正确建立并能正常传输数据。

## 5. 高级功能演示

VANTUN提供了多种高级功能，您可以按照以下步骤进行测试：

### 5.1 前向纠错(FEC)功能测试

FEC功能可以在网络不稳定时提供数据恢复能力。

```bash
# 使用FEC配置文件启动服务端
./bin/vantun -config test_configs/fec_server.json -server

# 使用FEC配置文件启动客户端
./bin/vantun -config test_configs/fec_client.json
```

在FEC配置中，您可以调整 `fec_data` 和 `fec_parity` 参数来控制纠错能力：
- `fec_data`: 数据分片数
- `fec_parity`: 冗余分片数

更高的冗余分片数可以提供更强的纠错能力，但会增加带宽开销。

### 5.2 流量混淆功能测试

流量混淆功能使VANTUN流量看起来像正常的HTTP/3流量，有效规避网络审查。

```bash
# 启用流量混淆的服务端
./bin/vantun -config test_configs/obfs_server.json -server

# 启用流量混淆的客户端
./bin/vantun -config test_configs/obfs_client.json
```

或者使用命令行参数：
```bash
# 使用命令行参数启用流量混淆
./bin/vantun -server -addr :4242 -obfs
./bin/vantun -addr localhost:4242 -obfs
```

### 5.3 多路径传输功能测试

多路径传输功能可以同时利用多个网络路径，提高传输速度和连接稳定性。

```bash
# 启用多路径传输的服务端
./bin/vantun -config test_configs/multipath_server.json -server

# 启用多路径传输的客户端
./bin/vantun -config test_configs/multipath_client.json
```

注意：多路径功能需要多个网络接口。在Linux系统中，您可以使用以下命令添加虚拟网络接口进行测试：

```bash
# 添加虚拟网络接口示例（需要root权限）
sudo ip link add name dummy0 type dummy
sudo ip addr add 192.168.100.1/24 dev dummy0
sudo ip link set dummy0 up
```

## 6. 性能测试

您可以使用 `iperf` 或其他网络性能测试工具来评估 VANTUN 的性能表现。

### 使用iperf进行性能测试

1. 首先启动VANTUN服务端：
```bash
./bin/vantun -server -addr :4242
```

2. 在另一个终端启动iperf服务端：
```bash
iperf3 -s
```

3. 通过VANTUN隧道连接到iperf服务端：
```bash
# 在客户端机器上
./bin/vantun -addr <server-ip>:4242
# 然后使用iperf客户端通过隧道测试
iperf3 -c localhost -p <tunneled-port>
```

### 性能优化建议

1. **调整FEC参数**：根据网络质量调整 `fec_data` 和 `fec_parity` 参数
2. **启用多路径传输**：在多网卡环境下启用多路径功能
3. **优化令牌桶参数**：根据带宽调整 `token_bucket_rate` 和 `token_bucket_capacity`

## 7. 日志与监控

VANTUN 提供详细的日志输出和遥测数据，帮助您监控和调试系统：

### 日志级别
- `debug`: 详细调试信息
- `info`: 一般信息（默认）
- `warn`: 警告信息
- `error`: 错误信息

### 遥测数据
VANTUN会输出以下遥测数据到标准输出：
- 网络延迟(RTT)
- 丢包率
- 带宽使用情况
- 拥塞窗口大小
- 在途字节数
- 传输速率

这些数据可用于性能分析和故障排查。