# gh-proxy-go

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

**gh-proxy-go** 是 [hunshcn/gh-proxy](https://github.com/hunshcn/gh-proxy) 的 Go 语言重写版。它是一个 GitHub 资源加速代理服务，支持 release、archive、raw 文件的代理下载以及 `git clone` 操作。部署为 Docker 容器或单文件二进制后，用户只需要在 GitHub 链接前加上代理地址即可加速访问。

### 功能

- 代理 GitHub release / archive / blob / raw 文件
- 代理 gist 文件
- 透传 `git clone`（Smart HTTP 协议）
- 支持 jsDelivr CDN 加速（可选，对 blob/raw 文件做 302 重定向）
- 白名单 / 黑名单 / 直通名单（按用户或仓库粒度控制）
- 流式传输，无内存上限，适合大文件
- 完整 CORS 支持
- 单文件二进制，无需运行时依赖

### 快速开始

```bash
# 下载编译好的 Linux 二进制
./dist/gh-proxy-go_linux_amd64

# 服务默认监听 http://0.0.0.0:11001
```

### 使用方式

将 `github.com` 替换为你的代理地址：

```bash
# 下载 Release
wget http://your-proxy:11001/cli/cli/releases/download/v2.67.0/gh_2.67.0_windows_amd64.msi

# git clone
git clone http://your-proxy:11001/cli/cli.git

# 源码压缩包
wget http://your-proxy:11001/cli/cli/archive/refs/heads/trunk.zip

# 单个文件
wget http://your-proxy:11001/cli/cli/blob/trunk/README.md
```

也支持完整 URL 形式：

```bash
# 两种写法效果一样
http://your-proxy:11001/https://github.com/user/repo
http://your-proxy:11001/user/repo
```

### 配置

通过环境变量配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `HOST` | `0.0.0.0` | 监听地址 |
| `PORT` | `11001` | 监听端口 |
| `ASSET_URL` | `https://hunshcn.github.io/gh-proxy` | 前端页面来源 |
| `JSDELIVR` | `false` | 启用 jsDelivr CDN（`1`/`true`/`yes`） |
| `SIZE_LIMIT` | `0`（不限） | 单文件大小上限（字节） |
| `WHITE_LIST` | 空 | 白名单，每行一个 user 或 user/repo |
| `BLACK_LIST` | 空 | 黑名单 |
| `PASS_LIST` | 空 | 直通名单（匹配的请求直接 302 到源站） |
| `DEBUG` | `false` | 调试模式 |

黑白名单示例：

```bash
export WHITE_LIST="vuejs/core
facebook/react
*/vite"

export BLACK_LIST="baduser/*"
```

### 编译构建

```bash
# 编译当前平台
go build -o dist/gh-proxy-go .

# 使用构建脚本（Windows PowerShell）
.\build.ps1 -Target linux

# 使用 Makefile（Linux / macOS）
make build-linux
make build-all
```

编译产物输出到 `dist/` 目录。

### Docker

```bash
docker build -t gh-proxy-go .
docker run -d --name ghproxy -p 11001:11001 gh-proxy-go
```

### 与原始 Python 版的差异

- 使用 Gin 框架替代 Flask，性能更高
- 单文件二进制（约 15MB），对比 Python 版 Docker 镜像约 940MB
- 环境变量配置替代 Python 源码内修改
- 流式传输依赖 Go 原生 `net/http`，更稳定
- 支持简写路径模式（自动补充 `github.com` 前缀）

### License

MIT
