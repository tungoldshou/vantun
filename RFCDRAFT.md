# VANTUN RFC 文档草案

## 1. 简介

VANTUN 是一个基于 QUIC 的高性能隧道协议，旨在提供安全、可靠和高效的数据传输。它具有以下核心特性：

* **多路径支持**：利用多条网络路径提高吞吐量和可靠性。
* **前向纠错 (FEC)**：通过 Reed-Solomon 编码增强数据传输的鲁棒性。
* **自适应拥塞控制**：结合底层 QUIC 拥塞控制和上层令牌桶限流，并根据遥测数据动态调整参数。
* **流量混淆**：通过模仿 HTTP/3 流量格式降低被识别和阻断的风险。
* **多种流类型**：支持交互流、批量流和遥测流，以满足不同应用需求。

## 2. 协议架构

VANTUN 协议建立在 QUIC 传输协议之上，使用 TLS 1.3 进行安全握手。协议架构如下：

```
+--------------------------------------------------+
|                   Application                    |
+--------------------------------------------------+
|              VANTUN Session Layer                |
|  - Session Negotiation (Control Stream)          |
|  - Stream Management (Interactive/Bulk/Telemetry)|
|  - FEC Encoding/Decoding                        |
|  - Obfuscation                                  |
+--------------------------------------------------+
|              QUIC Transport Layer                |
|  - Connection Management                         |
|  - Stream Multiplexing                           |
|  - Congestion Control (BBR/Cubic)                |
|  - Loss Detection                                |
+--------------------------------------------------+
|              TLS 1.3 Security Layer              |
+--------------------------------------------------+
|               UDP Network Layer                  |
+--------------------------------------------------+
```

## 3. 会话建立

VANTUN 会话通过以下步骤建立：

1.  **QUIC 连接**：客户端和服务端通过 QUIC 建立安全连接。
2.  **控制流**：在 QUIC 连接上打开一个双向控制流（流类型 0）。
3.  **会话协商**：
    *   客户端发送 `SessionInit` 消息到控制流，包含协议版本、认证令牌（可选）和支持的特性列表。
    *   服务端验证 `SessionInit` 消息，然后发送 `SessionAccept` 消息作为响应，指示会话是否被接受以及服务端支持的特性。
4.  **流创建**：会话建立后，客户端和服务端可以创建交互流、批量流和遥测流。

## 4. 消息格式

### 4.1 控制消息

控制消息在控制流上交换，使用 CBOR 编码。

*   **SessionInit** (类型 0x01)：
    *   `version` (uint16)：协议版本。
    *   `token` (bytes)：可选的认证令牌。
    *   `supported_features` (array of strings)：客户端支持的特性列表。
*   **SessionAccept** (类型 0x02)：
    *   `accepted` (bool)：会话是否被接受。
    *   `reason` (string)：如果会话被拒绝，拒绝的原因。
    *   `server_features` (array of strings)：服务端支持的特性列表。

### 4.2 帧格式

数据在逻辑流上以帧的形式传输，采用 TLV（Type-Length-Value）格式：

*   **DATA** (类型 0x00)：携带应用数据。
*   **PADDING** (类型 0x01)：填充帧，用于混淆流量。
*   **TELEMETRY** (类型 0x02)：携带遥测数据。

## 5. 流类型

VANTUN 定义了三种逻辑流类型：

*   **交互流** (类型 1)：用于低延迟、高优先级的交互式数据，如 SSH 或游戏流量。
*   **批量流** (类型 2)：用于高吞吐量的批量数据传输，如文件下载。
*   **遥测流** (类型 3)：用于传输连接的遥测数据，如 RTT、丢包率和带宽估计。

## 6. FEC (前向纠错)

VANTUN 使用 Reed-Solomon 编码实现前向纠错。数据被分割成 K 个数据分片，然后生成 M 个冗余分片。只要接收到 K 个分片（无论是数据分片还是冗余分片），就可以恢复原始数据。

FEC 参数（K 和 M）可以根据网络条件动态调整。

## 7. 多路径

VANTUN 支持多路径传输，可以同时使用多条网络路径（例如，Wi-Fi 和蜂窝数据）。每条路径都有独立的 QUIC 连接，并周期性地探测其 RTT、丢包率和带宽。数据可以基于权重在多条路径上进行分流。

## 8. 拥塞控制

VANTUN 采用混合拥塞控制机制：

*   **底层控制**：利用 QUIC 内置的 BBR 或 Cubic 拥塞控制算法。
*   **上层控制**：实现令牌桶限流器，提供更精细的速率控制。
*   **动态调整**：控制器根据遥测流上报的网络状况（RTT、丢包率、带宽）动态调整令牌桶的速率和 FEC 的冗余度。

## 9. 流量混淆

为了降低被识别和阻断的风险，VANTUN 提供了流量混淆功能：

*   **HTTP/3 模式**：将数据帧封装成类似 HTTP/3 帧的格式。
*   **填充**：在数据流中添加随机填充帧，以隐藏真实的流量模式和数据长度。

## 10. 配置

VANTUN 支持通过命令行参数和 JSON 配置文件进行配置。

### 10.1 命令行参数

*   `-server`：以服务端模式运行。
*   `-addr`：监听地址（服务端）或连接地址（客户端）。
*   `-config`：JSON 配置文件路径。
*   `-log-level`：日志级别。
*   `-multipath`：启用多路径。
*   `-obfs`：启用流量混淆。
*   `-fec-data`：FEC 数据分片数。
*   `-fec-parity`：FEC 冗余分片数。

### 10.2 JSON 配置文件

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

## 11. 验收标准

*   在 50ms RTT, 0% loss 下，吞吐率 ≥ 90% 链路可用带宽。
*   在 200ms RTT, 5% loss 下，交互流延迟降低 ≥ 20%（相对无FEC）。
*   单连接稳定运行 ≥ 24h。
*   CLI client/server 可通过 JSON 配置独立运行。
*   输出日志可用于链路质量评估。