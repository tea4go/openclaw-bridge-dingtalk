# 钉钉-OpenClaw Webhook 桥接服务 (Go + GoFrame)

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: [GoFrame](https://github.com/gogf/gf) v2
- **HTTP 客户端**: GoFrame gclient
- **配置管理**: GoFrame gcfg
- **日志**: GoFrame glog

## 项目结构

```
dingtalk-bridge-go/
├── main.go                 # 程序入口
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖锁定
├── config.yaml             # 配置文件
├── manifest/
│   └── config/
│       └── config.yaml     # 默认配置
├── internal/
│   ├── controller/         # HTTP 控制器
│   │   ├── dingtalk.go     # 钉钉 Webhook 处理
│   │   └── health.go       # 健康检查
│   ├── service/            # 业务逻辑层
│   │   ├── dingtalk/       # 钉钉服务
│   │   │   ├── crypto.go   # 签名验证
│   │   │   └── message.go  # 消息处理
│   │   └── openclaw/       # OpenClaw 服务
│   │       └── client.go   # API 客户端
│   └── model/              # 数据模型
│       ├── dingtalk.go     # 钉钉消息模型
│       └── openclaw.go     # OpenClaw 请求模型
├── utility/                # 工具函数
│   └── response/           # 响应封装
└── README.md               # 项目说明
```

## 快速开始

```bash
# 克隆/进入项目
cd dingtalk-bridge-go

# 下载依赖
go mod tidy

# 运行
go run main.go

# 编译
make build
```

## 配置说明

编辑 `config.yaml`:

```yaml
server:
  address: ":3000"
  
dingtalk:
  clientId: "dingbaw01rv4t8grihhq"
  clientSecret: "3wMpza07RCZFsE5Orvfh_tmTHvXNdDVuXgJuyE6CKS55FfGqzNqi3zel3tCnvpGL"
  appSecret: "3wMpza07RCZFsE5Orvfh_tmTHvXNdDVuXgJuyE6CKS55FfGqzNqi3zel3tCnvpGL"
  enableSignatureCheck: true

openclaw:
  gatewayUrl: "http://127.0.0.1:18789"
  hookToken: "dingtalk-bridge-token-2024"
  hookPath: "/hooks/dingtalk"
  sessionKey: "dingtalk:bridge:main"

webhook:
  path: "/webhook/dingtalk"
```

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /webhook/dingtalk | 钉钉 Webhook 接收 |
| GET | /health | 健康检查 |
| GET | / | 服务状态 |

## 特性

- ✅ GoFrame v2 框架，高性能、易维护
- ✅ 钉钉消息签名验证
- ✅ 会话管理
- ✅ 自动重试机制
- ✅ 结构化日志
- ✅ 优雅关闭

## 相关链接

- [GoFrame 文档](https://goframe.org/)
- [钉钉开放平台](https://open.dingtalk.com/)