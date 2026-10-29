# CTM - Claude Token Manager

一个用 Go 编写的命令行工具，用于管理多个 Claude 配置，方便终端使用。

## 项目结构

```
ctm/
├── main.go              # 程序入口
├── cmd/                 # Cobra 命令定义
│   ├── root.go          # 根命令
│   ├── list.go          # list 命令
│   ├── use.go           # use 命令
│   ├── unset.go         # unset 命令
│   ├── show.go          # show 命令
│   ├── config.go        # config 命令
│   ├── init.go          # init 命令
│   ├── completion.go    # 补全候选（provider 索引）
│   └── shell.go         # shell 引号处理
├── logic/               # 业务逻辑
│   ├── types.go         # 数据结构定义
│   ├── config.go        # CTM 配置文件操作
│   └── settings.go      # Claude Code settings.json 读写
├── config.json.example  # 配置文件示例
└── go.mod
```

## 安装

```bash
make build
```

二进制会生成在 `bin/ctm`。也可以直接使用 `go build -o bin/ctm .`。

## 配置文件

CTM 的配置文件位于 `~/.ctm/config.json`，格式如下：

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

切换时，CTM 只操作 `~/.claude/settings.json` 的 `env` 字段，其余字段（如 `permissions`、`model`、`enabledPlugins`）原样保留。写入后顶层字段会按字母序重排，内容不变。

## 使用方法

### 1. 列出所有配置

```bash
ctm list
```

当前使用的配置会标记为绿色并带 `*` 前缀。

### 2. 初始化 Shell

```bash
ctm init
```

`ctm init` 会根据当前 shell 自动写入：

- zsh：`~/.zshrc`
- bash：`~/.bashrc`

它会做两件事：

1. 把 ctm 所在目录加入 `PATH`，让 `ctm` 可直接调用；
2. 生成补全脚本到 `~/.ctm/completion.<shell>`，并在 rc 文件中 source 它。

写入内容示例：

```bash
# >>> ctm init >>>
export PATH='/path/to/ctm-dir':"$PATH"

if [ -f '/home/user/.ctm/completion.bash' ]; then
    . '/home/user/.ctm/completion.bash'
fi
# <<< ctm init <<<
```

重复执行 `ctm init` 时，会更新 `# >>> ctm init >>>` 和 `# <<< ctm init <<<` 之间的配置块，不会重复追加。

执行后重新打开 shell，或手动 source 对应配置文件：

```bash
source ~/.zshrc  # zsh
source ~/.bashrc # bash
```

### 3. 切换配置

```bash
ctm use 1
```

CTM 会以「隔离」方式切换：

1. 把 `settings.json` 中所有 CTM 管理过的环境变量名（即全部 provider 出现过的 key）从 `env` 中删除；
2. 写入目标 provider 的环境变量；
3. 更新 `~/.ctm/config.json` 的 `currentIndex`。

这样旧 provider 独有的变量不会残留。不属于任何 provider 的变量（例如你自己手写在 `env` 里的）不会被触碰。

### 4. 清理 CTM 管理的环境变量

```bash
ctm unset
```

把 `settings.json` 的 `env` 中所有 CTM 管理过的变量删除（`env` 删空后整个字段会被移除）。

### 5. 查看配置详情

查看当前配置：

```bash
ctm show
```

查看指定索引的配置：

```bash
ctm show 0
```

### 6. 生成补全脚本

`ctm init` 已自动安装补全。如需手动生成：

```bash
ctm completion bash    # 输出到 stdout，可自行 source
ctm completion zsh
```

### 7. 查看配置文件路径

```bash
ctm config
```

### 8. 帮助信息

```bash
ctm help
```

## 补全

安装补全后：

- `ctm <TAB>` 补全子命令；
- `ctm use <TAB>` / `ctm show <TAB>` 补全 provider 索引及其名称。

补全脚本支持 bash、zsh、fish、powershell（通过 `ctm completion <shell>` 生成）；`ctm init` 自动安装仅支持 bash 和 zsh。

## 特性

1. **自动环境变量管理**：切换配置时把目标 provider 的环境变量写入 Claude Code 的 `settings.json`
2. **环境变量隔离**：`ctm use` 会先删除所有由 CTM 管理过的变量，再写入目标配置，避免旧配置残留
3. **独立清理**：`ctm unset` 可以单独清理所有 CTM 管理过的变量
4. **自动 Shell 初始化**：`ctm init` 写入 PATH 并安装补全
5. **命令补全**：子命令与 provider 索引均可补全
6. **安全写入**：`settings.json` 采用临时文件 + rename 原子写入，保留原文件权限
7. **清晰的状态显示**：当前配置用绿色和 `*` 标记
8. **简洁的代码结构**：使用 Cobra 框架，代码组织清晰

## 注意事项

- CTM 只修改 `~/.claude/settings.json`，不会改动当前 shell 的环境变量。切换后新启动的 Claude Code 及其子进程会拿到新值。
- **删除变量需要重启才彻底生效**：Claude Code 对 settings.json 中新增/修改的变量会较快生效，但删除的变量要到下次启动 Claude Code 才真正从环境中消失。
- `settings.json` 中 `env` 的值必须是字符串，CTM 会按此写入。
- 若 `settings.json` 不是合法 JSON，CTM 会报错退出，不会覆盖原文件。
