# YoudaoTranslator Go

这是当前 TypeScript 翻译核心的独立 Go 移植版，输出与 Alfred Script Filter 兼容的 JSON。

## 支持范围

- Youdao 翻译 API
- Baidu 翻译 API
- CamelCase 查询预处理
- 中英文方向判断
- Youdao 基础释义、音标和网络释义
- Alfred `cmd` / `alt` modifier
- Copy Text 和 Quick Look URL
- 网络超时、HTTP 错误和响应格式错误提示

Go 版本与根目录现有 TypeScript 版本并行存在，不修改原实现。本目录暂不包含完整 Alfred `info.plist`、`.alfredworkflow` 安装包或发布脚本。

## 配置

程序读取以下环境变量：

| 变量 | 说明 |
|---|---|
| `key` | Youdao App Key 或 Baidu App ID |
| `secret` | API Secret |
| `platform` | 必须为 `Youdao` 或 `Baidu` |

程序使用最后一个命令行参数作为翻译内容。包含空格时请使用引号。

## 本地运行

从仓库根目录执行，以便 Alfred JSON 中的 `assets/...` 图标路径对应根目录资源：

```bash
go run ./go-translator/cmd/youdao-translator 'helloWorld'
```

完整示例：

```bash
key='your-key' \
secret='your-secret' \
platform='Youdao' \
go run ./go-translator/cmd/youdao-translator 'helloWorld'
```

Baidu：

```bash
key='your-appid' \
secret='your-secret' \
platform='Baidu' \
go run ./go-translator/cmd/youdao-translator '你好！'
```

stdout 只输出一个 Alfred Script Filter JSON 文档。配置或网络错误也会输出一个 `valid: false` 的 Alfred item。

## 测试

```bash
cd go-translator
go test ./...
go test -race ./...
go vet ./...
```

测试使用 fake HTTP transport，不调用真实翻译 API，也不需要有效密钥。

## 构建

本机：

```bash
cd go-translator
go build -trimpath -o ./youdao-translator ./cmd/youdao-translator
```

macOS Apple Silicon：

```bash
cd go-translator
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
  go build -trimpath -ldflags='-s -w' \
  -o ./youdao-translator-darwin-arm64 \
  ./cmd/youdao-translator
```

macOS Intel：

```bash
cd go-translator
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
  go build -trimpath -ldflags='-s -w' \
  -o ./youdao-translator-darwin-amd64 \
  ./cmd/youdao-translator
```

构建产物不应提交到 Git。

## 与原版本的兼容及修复

保留：

- `key`、`secret`、`platform` 环境变量
- 最后一个 argv 查询规则
- `Youdao`、`Baidu` 平台名
- 两家 API 的现有端点、参数和 MD5 协议
- `assets/translate.png` 和 `assets/translate-say.png`
- Alfred `cmd`、`alt`、`text.copy`、`quicklookurl`

修复：

- 缺配置或空查询不再崩溃
- HTTP 请求设置超时并校验状态码
- Provider 不保存请求级共享状态
- Baidu 回车参数使用译文
- 带标点或混合中文能够正确选择翻译方向
- Quick Look 查询进行 URL 编码
- 使用标准库 UTF-8 MD5
- 长标题按 Unicode rune 截断
- 错误 item 明确输出 `valid: false`
