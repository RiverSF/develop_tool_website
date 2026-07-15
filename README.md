# 开发者工具网站

## OpenSpec Skills 使用说明

本项目已接入 [OpenSpec](https://github.com/Fission-AI/OpenSpec)，用于在写代码前先对齐需求与方案。

**前置条件：** 已安装 Node.js ≥ 20.19.0，并全局安装 CLI：

```bash
npm install -g @fission-ai/openspec@latest
```

初始化（本仓库已完成，新环境可重跑）：

```bash
openspec init --tools cursor
```

重启 Cursor 后，可在对话中使用下方斜杠命令（或由 Agent 自动匹配对应 Skill）。

### 推荐工作流

```text
/opsx:explore  →  /opsx:propose  →  /opsx:apply  →  /opsx:sync  →  /opsx:archive
   （可选）         生成方案           实现任务      （可选）        归档合并
```

### Skills / 斜杠命令

| 命令 | Skill | 用途 |
|------|-------|------|
| `/opsx:explore` | `openspec-explore` | 探索想法、排查问题、澄清需求；只读思考，不写业务代码 |
| `/opsx:propose` | `openspec-propose` | 一次性创建变更，生成 proposal / design / specs / tasks |
| `/opsx:update` | `openspec-update-change` | 修订已有变更的规划产物，保持文档一致；不改代码 |
| `/opsx:apply` | `openspec-apply-change` | 按 `tasks.md` 逐项实现 |
| `/opsx:sync` | `openspec-sync-specs` | 将变更中的 delta specs 合并进主 specs，不归档 |
| `/opsx:archive` | `openspec-archive-change` | 实现完成后归档变更，并更新主 specs |

Skill 定义见 `.cursor/skills/`，斜杠命令见 `.cursor/commands/`。

### 使用示例

```text
/opsx:explore 暗色模式怎么接现有样式最干净？
/opsx:propose add-dark-mode
/opsx:apply
/opsx:archive
```

也可带变更名：

```text
/opsx:apply add-dark-mode
/opsx:update add-dark-mode
/opsx:sync add-dark-mode
/opsx:archive add-dark-mode
```

### 目录结构

```text
openspec/
├── specs/          # 主规格（系统当前行为的真相源）
├── changes/        # 进行中的变更（每个变更一个目录）
│   └── <name>/
│       ├── proposal.md
│       ├── design.md
│       ├── tasks.md
│       └── specs/  # delta specs
└── config.yaml
```

### 常用 CLI

```bash
openspec list                          # 列出进行中的变更
openspec show <change>                 # 查看变更详情
openspec validate <change>             # 校验规格格式
openspec update                        # 刷新本项目的 Agent 指令 / 命令
```

更多说明见 [OpenSpec 文档](https://github.com/Fission-AI/OpenSpec)。

## GitNexus Skills 使用说明

本项目已接入 [GitNexus](https://github.com/abhigyanpatwari/GitNexus)：把代码库索引成知识图谱，通过 MCP 工具给 Agent 完整的调用链 / 依赖 / 影响面上下文。

**前置条件：** Node.js ≥ 20，并全局安装 CLI（推荐，MCP 启动更快）：

```bash
npm install -g gitnexus@latest
```

本仓库索引与 Cursor 集成（已完成，新环境可重跑）：

```bash
# 在仓库根目录建立/刷新索引
gitnexus analyze

# 一次性配置 Cursor MCP + 全局 skills
gitnexus setup -c cursor
```

重启 Cursor 后即可使用 MCP 工具；Agent 会按任务自动匹配对应 Skill。

### Skills

| Skill | 何时使用 |
|-------|----------|
| `gitnexus-exploring` | 理解架构、追踪执行流、探索陌生代码（如「X 怎么工作？」） |
| `gitnexus-debugging` | 排查 bug、追踪错误来源（如「为什么失败？」「错误从哪来？」） |
| `gitnexus-impact-analysis` | 改代码前分析影响面（如「改 X 会波及什么？」） |
| `gitnexus-refactoring` | 安全重命名 / 抽取 / 拆分 / 移动代码 |
| `gitnexus-pr-review` | Review PR、评估合并风险与测试缺口 |
| `gitnexus-cli` | 需要跑 analyze / status / clean / wiki 等 CLI |
| `gitnexus-guide` | 查询 GitNexus 工具、资源、图谱 schema |
| `gitnexus-pdg-query` | 查询控制/数据依赖（需先 `gitnexus analyze --pdg`） |
| `gitnexus-taint-analysis` | 污点分析 / source→sink 数据流（需 `--pdg`） |

- 仓库内副本：`.claude/skills/gitnexus/`
- Cursor 全局 skills：`~/.cursor/skills/gitnexus-*`

### 常用自然语言示例

```text
这个仓库的认证流程怎么走？
改 UserService.validate 会影响哪些调用方？
帮我安全地把这个函数重命名
Review 一下当前改动的影响面
重新索引这个仓库
```

### 常用 CLI

```bash
gitnexus analyze                 # 建立/增量更新索引
gitnexus analyze --skills        # 额外生成按模块划分的 repo-specific skills
gitnexus analyze --pdg           # 启用 PDG / 污点分析层
gitnexus status                  # 查看索引状态
gitnexus setup -c cursor         # 配置 Cursor MCP
gitnexus doctor                  # 诊断运行时能力（如 FTS 扩展）
```

索引数据在 `.gitnexus/`（已忽略，勿提交）。改完大量代码后若 Agent 提示 index stale，再跑一次 `gitnexus analyze`。

更多说明见 [GitNexus 文档](https://github.com/abhigyanpatwari/GitNexus)。

## Agency Agents（中文专家角色库）使用说明

本项目已接入 [agency-agents-zh](https://github.com/jnMetaCode/agency-agents-zh)：把专家角色安装为 Cursor Project Rules（`.mdc`），对话时按描述自动匹配，或用 `@` 手动引用。

**安装位置：** `.cursor/rules/*.mdc`（项目级）

> 官方建议只保留约 10–20 个规则，避免全量安装稀释自动匹配。本仓库已精选 **18 个开发相关** 角色（见下表）；营销 / 销售 / 游戏等已剔除。

### 当前保留的角色

| 规则文件 | 角色 |
|----------|------|
| `engineering-backend-architect.mdc` | 后端架构师 |
| `engineering-frontend-developer.mdc` | 前端开发者 |
| `engineering-software-architect.mdc` | 软件架构师 |
| `engineering-code-reviewer.mdc` | 代码审查员 |
| `engineering-security-engineer.mdc` | 安全工程师 |
| `engineering-devops-automator.mdc` | DevOps 自动化 |
| `engineering-sre.mdc` | 站点可靠性工程师 |
| `engineering-database-optimizer.mdc` | 数据库优化 |
| `engineering-git-workflow-master.mdc` | Git 工作流 |
| `engineering-technical-writer.mdc` | 技术文档 |
| `engineering-codebase-onboarding-engineer.mdc` | 代码库上手 |
| `engineering-minimal-change-engineer.mdc` | 最小改动工程 |
| `engineering-ai-engineer.mdc` | AI 工程师 |
| `engineering-incident-response-commander.mdc` | 故障响应 |
| `testing-api-tester.mdc` | API 测试 |
| `testing-performance-benchmarker.mdc` | 性能基准 |
| `testing-reality-checker.mdc` | 交付验收 |
| `specialized-mcp-builder.mdc` | MCP 构建 |

### 如何使用

1. `.cursor/rules/` 下的 `.mdc` 会被 Cursor 自动识别（`alwaysApply: false`，按 `description` **智能匹配**）。
2. 在 Chat / Agent 中正常提问即可，例如：
   ```text
   帮我审查这个组件的性能问题   → 倾向匹配前端开发者
   这段代码有安全漏洞吗         → 倾向匹配安全工程师
   设计一下这个 API 的后端架构   → 倾向匹配后端架构师
   ```
3. 也可在对话里用 `@规则名` 手动指定某个智能体。
4. 在 **Cursor Settings**（`Ctrl+,`）→ **Rules** → **Project Rules** 中查看 / 开关已安装规则。

### 安装 / 增补角色（Windows）

全量转换后，只拷贝需要的 `.mdc` 到本项目（推荐），避免再次全量灌入：

```powershell
git clone --depth 1 https://github.com/jnMetaCode/agency-agents-zh.git
cd agency-agents-zh
powershell -ExecutionPolicy Bypass -File .\scripts\convert.ps1 -Tool cursor

# 按需拷贝到本项目
Copy-Item .\integrations\cursor\rules\engineering-xxx.mdc D:\Project\develop_tool_website\.cursor\rules\
```

若仍用官方一键安装（会装入全部规则，装完请再精选删除）：

```powershell
cd D:\Project\develop_tool_website
powershell -ExecutionPolicy Bypass -File D:\Project\agency-agents-zh\scripts\install.ps1 -Tool cursor
```

macOS / Linux：`./scripts/convert.sh --tool cursor`，再按需 `cp` 到 `.cursor/rules/`。

更多角色说明见 [agency-agents-zh README](https://github.com/jnMetaCode/agency-agents-zh)。
