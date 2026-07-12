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
│   ├── unset.go        # unset 命令
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

当前使用的配置会标记为绿色并带 `*` 前缀。

### 2. 初始化 Shell 配置

```bash
ctm init
```

`ctm init` 会根据当前 shell 自动写入：

- zsh：`~/.zshrc`
- bash：`~/.bashrc`

它会写入 PATH，并生成 `ctm-use` 和 `ctm-unset` 两个 shell 函数。

写入内容示例：

```bash
# >>> ctm init >>>
export PATH='/path/to/ctm-dir':"$PATH"

ctm-use() {
    local ctm_exports
    local ctm_unsets
    ctm_exports=$(ctm use "$@") || return
    ctm_unsets=$(ctm unset) || return
    eval "$ctm_unsets"
    eval "$ctm_exports"
    ctm show
}

ctm-unset() {
    local ctm_unsets
    ctm_unsets=$(ctm unset "$@") || return
    eval "$ctm_unsets"
}
# <<< ctm init <<<
```

重复执行 `ctm init` 时，会更新 `# >>> ctm init >>>` 和 `# <<< ctm init <<<` 之间的配置块，不会重复追加。

执行后重新打开 shell，或手动 source 对应配置文件：

```bash
source ~/.zshrc  # zsh
source ~/.bashrc # bash
```

### 3. 切换配置

推荐使用 `ctm init` 生成的 `ctm-use` 函数：

```bash
ctm-use 1
```

也可以手动执行：

```bash
ctm_exports=$(ctm use 1) && \
ctm_unsets=$(ctm unset) && \
eval "$ctm_unsets" && \
eval "$ctm_exports" && \
ctm show
```

### 4. 清理 CTM 管理的环境变量

执行过 `ctm init` 并重新打开 shell，或 source 对应 rc 文件后，推荐使用：

```bash
ctm-unset
```

如果还没有初始化 shell，也可以手动执行：

```bash
eval "$(ctm unset)"
```

`ctm unset` 会读取所有 provider 中出现过的环境变量名，并输出对应的 unset 命令：

```bash
unset ANTHROPIC_API_KEY
unset ANTHROPIC_BASE_URL
```

### 5. 查看配置详情

查看当前配置：

```bash
ctm show
```

查看指定索引的配置：

```bash
ctm show 0
```

### 6. 查看配置文件路径

```bash
ctm config
```

### 7. 帮助信息

```bash
ctm help
```

## 特性

1. **自动环境变量管理**：切换配置时自动导出目标 provider 的环境变量
2. **环境变量隔离**：`ctm-use` 会在切换成功后先 `unset` 所有 CTM 管理过的环境变量，再应用目标配置，避免旧配置残留
3. **独立清理函数**：`ctm-unset` 可以单独清理所有 CTM 管理过的环境变量
4. **自动 Shell 初始化**：`ctm init` 自动写入 `~/.bashrc` 或 `~/.zshrc`，并生成 `ctm-use` / `ctm-unset`
5. **清晰的状态显示**：当前配置用绿色和 `*` 标记
6. **简洁的代码结构**：使用 Cobra 框架，代码组织清晰

## 注意事项

由于操作系统限制，子进程无法直接修改父进程的环境变量，因此需要使用 `eval "$(ctm use <index>)"` / `eval "$(ctm unset)"`，或使用 `ctm init` 生成的 `ctm-use <index>` / `ctm-unset` 来应用到当前 shell。

推荐先执行 `ctm init`，重新打开 shell 或 source 对应 rc 文件后，使用 `ctm-use <index>` 完成切换，使用 `ctm-unset` 清理 CTM 管理的环境变量。
