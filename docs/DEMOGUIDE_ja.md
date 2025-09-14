# VANTUN デモと使用ガイド

このガイドは、コンパイル、設定からサーバーとクライアントの実行まで、VANTUNをすぐに使い始めるのに役立ちます。

## 1. 環境準備とコンパイル

### 環境要件
- Go 1.21 またはそれ以上
- サポートされているオペレーティングシステム (Linux, macOS, Windows)

### コンパイル手順

```bash
# リポジトリのクローン (まだクローンしていない場合)
# git clone <repository-url>
cd vantun

# プロジェクトのコンパイル
go build -o bin/vantun cmd/main.go

# コンパイル結果の確認
./bin/vantun -h
```

コンパイルが成功すると、`bin`ディレクトリに実行ファイル`vantun`が表示されます。

## 2. 設定の詳細

VANTUNは、コマンドライン引数とJSON設定ファイルの2つの方法で設定できる柔軟な設定オプションを提供します。

### 2.1 コマンドライン引数の詳細

| 引数 | 説明 | デフォルト値 |
|------|------|--------|
| `-server` | サーバーモードで実行 | false (クライアントモード) |
| `-addr` | リッスンアドレス（サーバー）または接続アドレス（クライアント） | `localhost:4242` |
| `-config` | JSON設定ファイルのパス | なし |
| `-log-level` | ログレベル (debug, info, warn, error) | info |
| `-multipath` | マルチパス伝送を有効化 | false |
| `-obfs` | トラフィック難読化を有効化 | false |
| `-fec-data` | FECデータフラグメント数 | 10 |
| `-fec-parity` | FEC冗長フラグメント数 | 3 |

### 2.2 JSON設定ファイルの詳細

以下の設定項目を含む `config.json` ファイルを作成します：

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

設定項目の説明：
- `server`: 実行モード (true=サーバー, false=クライアント)
- `address`: リッスン/接続アドレス
- `log_level`: ログレベル
- `multipath`: マルチパス伝送を有効化するかどうか
- `obfs`: トラフィック難読化を有効化するかどうか
- `fec_data`: FECデータフラグメント数
- `fec_parity`: FEC冗長フラグメント数
- `token_bucket_rate`: トークンバケットレート (バイト/秒)
- `token_bucket_capacity`: トークンバケット容量 (バイト)

`test_configs/`ディレクトリのサンプル設定ファイルも参照できます：
- `fec_client.json` と `fec_server.json`: FEC機能テスト設定
- `obfs_client.json` と `obfs_server.json`: トラフィック難読化機能テスト設定
- `multipath_client.json` と `multipath_server.json`: マルチパス機能テスト設定

## 3. デモの実行

### 3.1 サーバーの起動

ターミナルウィンドウを開き、以下のコマンドを実行してVANTUNサーバーを起動します：

```bash
# コマンドライン引数でサーバーを起動
./bin/vantun -server -addr :4242

# または設定ファイルでサーバーを起動
./bin/vantun -config config.json -server
```

サーバーが起動すると、以下のような出力が表示されます：
```
2025/09/14 04:00:00 Server running, waiting for streams...
```

### 3.2 クライアントの起動

別のターミナルウィンドウを開き、以下のコマンドを実行してVANTUNクライアントを起動します：

```bash
# コマンドライン引数でクライアントを起動
./bin/vantun -addr localhost:4242

# または設定ファイルでクライアントを起動
./bin/vantun -config config.json
```

クライアントが起動すると、以下のような出力が表示されます：
```
2025/09/14 04:00:00 Client connected, opening interactive stream...
2025/09/14 04:00:00 Received echo: Hello from VANTUN client!
```

これは、クライアントがサーバーに正常に接続し、最初のデータ交換を完了したことを示しています。

## 4. 接続の検証

クライアントとサーバーが正常に起動すると、以下の出力情報が表示されます：

### サーバー出力:
```
2025/09/14 04:00:00 Server running, waiting for streams...
2025/09/14 04:00:00 Accepted interactive stream, echoing data...
```

### クライアント出力:
```
2025/09/14 04:00:00 Client connected, opening interactive stream...
2025/09/14 04:00:00 Received echo: Hello from VANTUN client!
```

これらの出力は以下のことを示しています：
1. クライアントがサーバーに正常に接続
2. クライアントがインタラクティブストリームを開いた
3. クライアントがサーバーにメッセージを送信
4. サーバーがメッセージを受信してエコー
5. クライアントがエコーメッセージを正常に受信

これにより、VANTUNトンネルが正しく確立され、正常にデータを転送できることを検証します。

## 5. 高度な機能デモ

VANTUNは複数の高度な機能を提供しており、以下の手順でテストできます：

### 5.1 前方誤り訂正(FEC)機能テスト

FEC機能は、ネットワークが不安定な場合にデータ復元能力を提供します。

```bash
# FEC設定ファイルでサーバーを起動
./bin/vantun -config test_configs/fec_server.json -server

# FEC設定ファイルでクライアントを起動
./bin/vantun -config test_configs/fec_client.json
```

FEC設定では、`fec_data`と`fec_parity`パラメータを調整して誤り訂正能力を制御できます：
- `fec_data`: データフラグメント数
- `fec_parity`: 冗長フラグメント数

より高い冗長フラグメント数はより強い誤り訂正能力を提供しますが、帯域幅のオーバーヘッドが増加します。

### 5.2 トラフィック難読化機能テスト

トラフィック難読化機能により、VANTUNトラフィックが通常のHTTP/3トラフィックのように見え、ネットワーク検閲を効果的に回避できます。

```bash
# トラフィック難読化を有効にしたサーバー
./bin/vantun -config test_configs/obfs_server.json -server

# トラフィック難読化を有効にしたクライアント
./bin/vantun -config test_configs/obfs_client.json
```

またはコマンドライン引数を使用：
```bash
# コマンドライン引数でトラフィック難読化を有効化
./bin/vantun -server -addr :4242 -obfs
./bin/vantun -addr localhost:4242 -obfs
```

### 5.3 マルチパス伝送機能テスト

マルチパス伝送機能により、複数のネットワークパスを同時に利用し、転送速度と接続の安定性を向上させることができます。

```bash
# マルチパス伝送を有効にしたサーバー
./bin/vantun -config test_configs/multipath_server.json -server

# マルチパス伝送を有効にしたクライアント
./bin/vantun -config test_configs/multipath_client.json
```

注意：マルチパス機能には複数のネットワークインターフェースが必要です。Linuxシステムでは、以下のコマンドで仮想ネットワークインターフェースを追加してテストできます：

```bash
# 仮想ネットワークインターフェース追加の例（root権限が必要）
sudo ip link add name dummy0 type dummy
sudo ip addr add 192.168.100.1/24 dev dummy0
sudo ip link set dummy0 up
```

## 6. パフォーマンステスト

`iperf`または他のネットワークパフォーマンステストツールを使用して、VANTUNのパフォーマンスを評価できます。

### iperfを使用したパフォーマンステスト

1. まずVANTUNサーバーを起動：
```bash
./bin/vantun -server -addr :4242
```

2. 別のターミナルでiperfサーバーを起動：
```bash
iperf3 -s
```

3. VANTUNトンネルを介してiperfサーバーに接続：
```bash
# クライアントマシンで
./bin/vantun -addr <server-ip>:4242
# 次にiperfクライアントを使用してトンネルを介してテスト
iperf3 -c localhost -p <tunneled-port>
```

### パフォーマンス最適化の推奨事項

1. **FECパラメータの調整**：ネットワーク品質に応じて `fec_data` と `fec_parity` パラメータを調整
2. **マルチパス伝送の有効化**：マルチNIC環境でマルチパス機能を有効化
3. **トークンバケットパラメータの最適化**：帯域幅に応じて `token_bucket_rate` と `token_bucket_capacity` を調整

## 7. ログと監視

VANTUNは詳細なログ出力とテレメトリデータを提供し、システムの監視とデバッグを支援します：

### ログレベル
- `debug`: 詳細なデバッグ情報
- `info`: 一般的な情報（デフォルト）
- `warn`: 警告情報
- `error`: エラー情報

### テレメトリデータ
VANTUNは標準出力に以下のテレメトリデータを出力します：
- ネットワーク遅延(RTT)
- パケット損失率
- 帯域幅使用状況
- 輻輳ウィンドウサイズ
- 送信中のバイト数
- 転送レート

これらのデータは、パフォーマンス分析とトラブルシューティングに使用できます。