# iFlow 上下文：VANTUN 协议实现项目

## 项目概览

本项目旨在实现一个基于 QUIC 的高性能隧道协议 **VANTUN**，具备以下核心能力：

* **安全握手与会话协商**：通过 control stream 进行。
* **多类型逻辑流**：支持 interactive（交互）、bulk（批量）、telemetry（遥测）三种流类型。
* **FEC (前向纠错)**：增强数据传输的鲁棒性。
* **多路径 (Multipath)**：利用多条网络路径提升性能和可靠性。
* **上层拥塞控制 (Hybrid CC)**：结合底层 QUIC CC 和上层 Token-Bucket 限流。
* **可插拔混淆模块**：降低被误判为非法流量的风险。
* **最小可用客户端/服务端**：提供可通过命令行运行的 `client` 和 `server` 程序。

## 项目结构与关键文件

当前目录包含以下关键文件：

* `vantun-prd.md`: 项目需求文档，详细描述了项目目标、功能模块、里程碑和技术栈建议。
* `IFLOW.md`: 此文件，为 iFlow 提供项目上下文。

根据 `vantun-prd.md`，项目的技术栈建议如下：

* **语言**: Go
* **核心库**: `quic-go` 用于 QUIC 实现。
* **序列化**: `github.com/fxamacker/cbor` 用于 CBOR 编码。
* **FEC**: `github.com/klauspost/reedsolomon` 用于 Reed-Solomon 编码。
* **CLI**: `cobra/viper` 用于命令行接口和配置管理。

## 开发与运行

项目开发将遵循 `vantun-prd.md` 中定义的里程碑 (M1-M4) 进行迭代。

### 构建与运行

目前项目尚未包含源代码和构建脚本。根据 PRD，最终将提供以下 CLI 程序：

* **客户端**: `vantun client -c config.json`
* **服务端**: `vantun server -c config.json`

配置文件 `config.json` 将支持端口、认证 token、FEC 策略、混淆策略等设置。

### 开发约定

* **协议实现**: 遵循 PRD 中定义的协议基础、帧结构、流类型等。
* **功能模块**: 按 PRD 中的功能模块拆分（协议基础、数据转发、FEC、多路径、拥塞控制、混淆、遥测、CLI）进行开发。
* **代码质量**: 需满足验收标准，如吞吐率、延迟降低、稳定性等。