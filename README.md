# CTM - Configuration Token Manager

**本项目完全由 Claude 完成**

一个用 Go 编写的命令行工具，用于管理多个配置和环境变量。

## 写在前面

本项目初衷是为了管理多个 claude api token，方便切换，现在写完看起来似乎可以应用到其他各种需要切换环境变量的情况。

另外一提，看到个平替的[免费使用 Claude Code 的方案](https://blyrin.cn/posts/cc-free/)，总之是用各平台免费提供的模型来代替 Claude，目前简单用了一下，简单任务还不错（得益于 Claude Code 的强大），但是相较于 Claude Sonnet 4.5 模型还是差点意思，可以考虑两种方案综合使用。

## 项目结构

```
ctm/
├── main.go              # 程序入口
├── cmd/                 # Cobra 命令定义
│   ├── root.go         # 根命令
│   ├── list.go         # list 命令
│   ├── use.go          # use 命令
│   ├── show.go         # show 命令
│   ├── config.go       # config 命令
│   └── init.go         # init 命令
├── logic/              # 业务逻辑
│   ├── types.go        # 数据结构定义
│   └── config.go       # 配置文件操作
├── config.json.example # 配置文件示例
└── go.mod
```

## 安装

```bash
go build -o ctm
```

## 配置文件

配置文件位于 `~/.ctm/config.json`，格式如下：

```json
{
  "providers": [
    {
      "name": "provider-name",
      "description": "provider description",
      "env": {
        "ANTHROPIC_BASE_URL": "https://example.com",
        "ANTHROPIC_API_KEY": "sk-..."
      }
    }
  ],
  "currentIndex": 0
}
```

## 使用方法

### 1. 列出所有配置

```bash
ctm list
```

输出示例：

```
Available providers:
* 0: gaccode
  1: v3 cm
  2: anyrouter
```

当前使用的配置会标记为绿色并带 `*` 前缀。

### 2. 切换配置

```bash
eval $(ctm use 1)
```

此命令会：

- 更新 `currentIndex` 为指定索引
- 输出所有环境变量的 export 命令
- 当前配置的环境变量设置为对应值
- 其他配置中的环境变量设置为空字符串

### 3. 查看配置详情

查看当前配置：

```bash
ctm show
```

查看指定索引的配置：

```bash
ctm show 0
```

### 4. 查看配置文件路径

```bash
ctm config
```

### 5. 生成 Shell 初始化脚本

```bash
ctm init
```

这会输出一个 shell 函数定义，你可以添加到 `~/.bashrc` 或 `~/.zshrc`：

```bash
ctm-use() {
    eval $(ctm use "$@")
}
```

然后你就可以直接使用：

```bash
ctm-use 1  # 自动切换并设置环境变量
```

### 6. 帮助信息

```bash
ctm help
```

## 特性

1. **自动环境变量管理**：切换配置时自动处理所有环境变量
2. **环境变量隔离**：非当前配置的环境变量会被设置为空，避免残留
3. **清晰的状态显示**：当前配置用绿色和 `*` 标记
4. **简洁的代码结构**：使用 Cobra 框架，代码组织清晰
5. **支持多种 Shell**：支持 bash、zsh、fish

## 注意事项

由于操作系统限制，子进程无法直接修改父进程的环境变量，因此需要使用 `eval $(ctm use <index>)` 的方式来应用环境变量。

建议使用 `ctm init` 生成的 shell 函数来简化操作。
