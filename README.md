# B站稍后再看邮件推送服务

自动获取 B站"稍后再看"列表，并通过 163 邮箱发送到 Gmail。

## 功能特性

- 获取 B站"稍后再看"视频列表
- 通过 163 邮箱 SMTP 发送邮件到 Gmail
- 支持立即执行和定时执行两种模式
- Cookie 过期自动发送邮件通知

## 环境要求

- Go 1.21+
- 163 邮箱账号（已开启 SMTP 服务）
- B站账号

## 配置

编辑 `config/config.yaml`：

```yaml
bilibili:
  sessdata: "your_sessdata"
  bili_jct: "your_bili_jct"
  buvid3: "your_buvid3"

email:
  smtp_host: "smtp.163.com"
  smtp_port: 25
  username: "your_email@163.com"
  password: "your_auth_code"
  from: "your_email@163.com"
  to: "your_gmail@gmail.com"
  alert_to: "your_gmail@gmail.com"

schedule: "0 9 * * *"
```

### 获取 B站 Cookie

1. 登录 [B站](https://www.bilibili.com)
2. 按 F12 打开开发者工具 → Application → Cookies
3. 复制 `SESSDATA`、`bili_jct`、`buvid3` 的值

### 获取 163 邮箱授权码

1. 登录 163 邮箱
2. 设置 → POP3/SMTP/IMAP
3. 开启 SMTP 服务
4. 获取授权码（不是登录密码）

## 使用方法

```bash
# 编译
go build -o bili-seelater .

# 立即执行一次
./bili-seelater run

# 启动定时推送服务（每天9点执行）
./bili-seelater serve
```

## Cron 表达式

`schedule` 字段支持标准 cron 格式：

| 表达式 | 说明 |
|--------|------|
| `0 9 * * *` | 每天早上9点 |
| `0 */6 * * *` | 每6小时 |
| `0 9,18 * * *` | 每天9点和18点 |

## 项目结构

```
bili-seelater/
├── cmd/
│   ├── root.go          # CLI 入口
│   ├── run.go           # run 命令
│   └── serve.go         # serve 命令
├── internal/
│   ├── bilibili/client.go
│   ├── email/sender.go
│   └── config/config.go
├── config/config.yaml
├── main.go
└── go.mod
```

## 依赖

- github.com/spf13/cobra - CLI 框架
- github.com/spf13/viper - 配置管理
- github.com/go-resty/resty/v2 - HTTP 客户端
- gopkg.in/gomail.v2 - 邮件发送
- github.com/robfig/cron/v3 - 定时任务
