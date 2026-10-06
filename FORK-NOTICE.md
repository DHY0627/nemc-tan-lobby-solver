# Fork Notice（本 fork 的说明）

本仓库是 **NeteaseBedrockGateway** 使用的一份 `nemc-tan-lobby-solver` 副本（fork），
用于给 [NeteaseBedrockGateway](https://github.com/DHY0627/NeteaseBedrockGateway) 提供
网易本地联机（TanLobby / NetherNet）协议实现。

## 来源

- 上游仓库：[`UCKETX/nemc-tan-lobby-solver`](https://github.com/UCKETX/nemc-tan-lobby-solver)（默认分支 `main`）
- 模块路径：`github.com/Happy2018new/nemc-tan-lobby-solver`（沿用上游 `go.mod` 的声明，便于以 `replace` 直接替换）
- 本 fork 的起点：上游 `619c4e0 protocol/login/dial.go: Minor changes`

## ⚠️ 许可证

**上游仓库没有 LICENSE 文件**（本 fork 也没有新增许可证），因此：

- 本 fork **不授予任何额外权利**；版权归原作者所有；
- 仅在本机自用（研究协议、连自己的服务器）通常不涉及分发条款；
- 若要**再分发**（包括把它 vendor 进公开仓库、发布二进制），请先联系原作者取得许可。

## 本 fork 相对上游的改动

共 14 个文件。按用途分组：

### 1. NetherNet 数据通道：不丢首包（NeteaseBedrockGateway 本次排错的核心修复）

| 文件 | 改动 |
|---|---|
| `core/nethernet/conn.go` | 新增 `bindChannelHandlers()`（幂等绑定收发处理器），`handleTransports()` 改为调用它 |
| `core/nethernet/listener.go` | 收到 `ReliableDataChannel` 后**立刻**绑定处理器，不再等两个通道都就绪 |
| `core/nethernet/dial.go` | 创建两个数据通道后各自调用一次 `bindChannelHandlers()` |

原因：底层 WebRTC 数据通道一建立就开始读取，而消息处理器原本要等两个通道协商完成才注册；
网易客户端把 DCEP 握手与第一个游戏包在同一毫秒发出，首包会被静默丢弃 → 客户端「连上了却一个字节都不发」。
详见 [NeteaseBedrockGateway/docs/troubleshooting.md](https://github.com/DHY0627/NeteaseBedrockGateway/blob/main/docs/troubleshooting.md)。

### 2. 不可靠数据通道支持

| 文件 | 改动 |
|---|---|
| `core/nethernet/message.go` | 新增 `handleUnreliableMessage()`（与可靠通道相同的分段帧格式） |

### 3. RakNet：cloudburst/Geyser Secure 模式的 4 字节安全 Cookie

| 文件 | 改动 |
|---|---|
| `core/raknet/dial.go` | Reply1 中 `Secure=true` 时保存 4 字节 Cookie，并在 Request2 里回传；否则服务器不响应 |
| `core/raknet/internal/message/open_connection_reply_1.go` | 解析 Secure 与 Cookie 字段 |
| `core/raknet/internal/message/open_connection_request_2.go` | 编码 Cookie 字段 |

### 4. SCTP：校验和改为 CRC32C(Castagnoli) + 小端

| 文件 | 改动 |
|---|---|
| `core/sctp/packet.go` | 校验和使用 CRC32C 而非 IEEE，并按小端写入（否则与网易客户端的 SCTP 实现校验和不一致） |

### 5. TanLobby 协议编解码补充

| 文件 | 改动 |
|---|---|
| `protocol/encoding/reader.go` / `writer.go` | 新增字段读写辅助 |
| `protocol/encoding/room_tips.go` | RoomTips 字段调整（ProtocolID / VersionString 等） |
| `protocol/packet/id.go` / `pool.go` | 包 ID 与包池补充 |
| `protocol/login/dial.go` | 拨号流程微调 |

## 构建

```bash
go build ./...        # 需要 Go 1.25+
```

作为依赖使用时，在上游项目里用 `replace` 指向本仓库即可：

```go
replace github.com/Happy2018new/nemc-tan-lobby-solver => github.com/DHY0627/nemc-tan-lobby-solver v0.0.0-<日期>-<提交>
```

（`NeteaseBedrockGateway` 目前把依赖以 `vendor/` 形式固定，因此 clone 后可直接编译。）
