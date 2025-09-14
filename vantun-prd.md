# VANTUN 协议实现开发任务清单

## 🎯 项目目标

实现一个基于 QUIC 的高性能隧道协议，具备以下核心能力：

* 安全握手与会话协商（control stream）
* 多类型逻辑流（interactive、bulk、telemetry）
* FEC（前向纠错）增强
* 多路径（multipath）支持
* 上层拥塞控制（Hybrid CC）
* 可插拔混淆模块（仅限降低误判风险，不涉及非法用途）
* 提供最小可用 **client + server** 程序，命令行即可运行

---

## 📌 功能模块拆分

### 1. 协议基础

* [ ] 定义 SessionInit / SessionAccept 消息格式（CBOR）
* [ ] 实现 control stream（双向）
* [ ] 完成基础握手流程（基于 QUIC + TLS1.3）
* [ ] 定义逻辑流类型（interactive / bulk / telemetry）

### 2. 数据转发与 framing

* [ ] 定义帧结构（TLV：type, flags, length, payload）
* [ ] 实现 DATA / PADDING / TELEMETRY 基础帧
* [ ] 支持 SOCKS5/HTTP 代理模式（可选）

### 3. FEC 模块

* [ ] 集成 Reed-Solomon 库
* [ ] 在 bulk 流中添加冗余块生成/恢复逻辑
* [ ] 在 telemetry 中动态调整冗余比例

### 4. 多路径（Multipath）

* [ ] 探测多条路径的 RTT/loss/bandwidth
* [ ] 支持基于权重的分流（WRR）
* [ ] 在交互流中可选「复制包」机制

### 5. 拥塞控制（Hybrid CC）

* [ ] 底层：沿用 QUIC 自带（BBR / Cubic）
* [ ] 上层：实现 Token-Bucket 限流
* [ ] 控制器根据 telemetry 调整 send\_rate、FEC 强度、分流策略

### 6. 混淆（Obfuscation）

* [ ] 实现 real-http3 模式（模仿 HTTP/3 帧）
* [ ] 实现 padding/length shaping
* [ ] 控制面支持 obfs 协商与切换

### 7. Telemetry

* [ ] 实现周期性统计（RTT、loss、bw\_est）
* [ ] 在 telemetry stream 中传输
* [ ] 客户端/服务端本地日志（文本输出）

### 8. CLI 客户端 & 服务端

* [ ] client 程序：`vantun client -c config.json`
* [ ] server 程序：`vantun server -c config.json`
* [ ] config.json 支持：端口、认证 token、FEC 策略、obfs 策略
* [ ] 日志直接输出到 stdout，JSON 格式

---

## 🚀 里程碑（约 8-10 周）

**M1 – 最小可用版 (Week 1-3)**
✅ QUIC/TLS 基础握手
✅ control stream & session 协商
✅ interactive/bulk 数据流
✅ CLI client/server 最小跑通

**M2 – 增强稳定版 (Week 4-6)**
✅ FEC 集成（Reed-Solomon）
✅ Telemetry 基础（RTT、loss 上报）
✅ Token-bucket 控制器雏形

**M3 – 高级功能版 (Week 7-8)**
✅ Multipath（多出口探测 + 分流）
✅ Obfs (http3 + padding)
✅ Telemetry 动态调节（自适应 FEC + 速率）

**M4 – 验收版 (Week 9-10)**
✅ 完成 CLI 参数 & JSON 配置
✅ 测试（netem 场景、长连接 24h 稳定性）
✅ 交付 RFC 文档 + 可运行 demo

---

## 🔧 技术栈建议

* 语言：Go → quic-go 库
* 序列化：CBOR（github.com/fxamacker/cbor）
* FEC：Reed-Solomon（github.com/klauspost/reedsolomon）
* CLI：cobra/viper

---

## ✅ 验收标准

* 在 50ms RTT, 0% loss 下，吞吐率 ≥ 90% 链路可用带宽
* 在 200ms RTT, 5% loss 下，交互流延迟降低 ≥ 20%（相对无FEC）
* 单连接稳定运行 ≥ 24h
* CLI client/server 可通过 JSON 配置独立运行
* 输出日志可用于链路质量评估

---
