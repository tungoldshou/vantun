# VANTUN Demo 指南

## 1. 编译

确保你已经安装了 Go 1.21 或更高版本。

```bash
# 克隆仓库
git clone <repository-url>
cd vantun

# 编译
go build -o bin/vantun cmd/main.go
```

## 2. 配置

VANTUN 支持通过命令行参数和 JSON 配置文件进行配置。

### 2.1 命令行参数

- `-server`: 以服务端模式运行。
- `-addr`: 监听地址（服务端）或连接地址（客户端），默认为 `localhost:4242`。
- `-config`: JSON 配置文件路径。
- `-log-level`: 日志级别。
- `-multipath`: 启用多路径。
- `-obfs`: 启用流量混淆。
- `-fec-data`: FEC 数据分片数，默认为 10。
- `-fec-parity`: FEC 冗余分片数，默认为 3。

### 2.2 JSON 配置文件

创建一个 `config.json` 文件：

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

## 3. 运行

### 3.1 启动服务端

```bash
# 使用命令行参数
./bin/vantun -server -addr :4242

# 或使用配置文件
./bin/vantun -config config.json -server
```

### 3.2 启动客户端

```bash
# 使用命令行参数
./bin/vantun -addr localhost:4242

# 或使用配置文件
./bin/vantun -config config.json
```

## 4. 验证

当客户端和服务端都启动后，你应该能看到以下输出：

**服务端:**
```
2025/09/14 04:00:00 Server running, waiting for streams...
2025/09/14 04:00:00 Accepted interactive stream, echoing data...
```

**客户端:**
```
2025/09/14 04:00:00 Client connected, opening interactive stream...
2025/09/14 04:00:00 Received echo: Hello from VANTUN client!
```

这表明客户端成功连接到服务端，打开一个交互流，发送了一条消息，并收到了服务端的回显。

## 5. 功能测试

你可以通过修改配置来测试不同的功能：

### 5.1 测试 FEC

修改 `config.json` 中的 `fec_data` 和 `fec_parity` 参数，然后重新启动客户端和服务端。

### 5.2 测试流量混淆

在客户端和服务端都启用 `-obfs` 参数或在配置文件中设置 `"obfs": true`。

### 5.3 测试多路径

多路径功能需要多个网络接口。你可以使用 Linux 的 `netem` 工具来模拟不同的网络条件。

## 6. 性能测试

你可以使用 `iperf` 或其他网络性能测试工具来评估 VANTUN 的性能。

## 7. 日志

VANTUN 会输出遥测数据和控制信息到标准输出，可以用于监控和调试。