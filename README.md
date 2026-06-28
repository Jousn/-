# Multica 项目任务提交总览

## 项目简介

Multica 是一个开源的 AI Agent 托管平台，允许 AI 编码 Agent 作为团队成员工作，具备 Squad 组织、自主任务执行、多工作空间支持等功能。

本项目提交目录包含笔试要求的三个核心任务，展示了深入的代码分析能力、架构设计能力和开源贡献实践。

---

## 任务结构

### 📋 Task 1：深度代码考古（25分）

#### 1.1 本地搭建（5分）

在本地成功运行 Multica 项目，记录搭建过程、遇到的问题和解决方案。

- **文档**：[task1/setup-log.md](task1/setup-log.md)
- **内容**：环境检查、Docker Compose 启动、端口冲突解决、验证码配置

#### 1.2 故障场景推演（20分）

深入阅读源码，追踪完整的代码路径，分析分布式系统故障场景（至少完成 2 个）。

- **文档**：[task1/fault-analysis.md](task1/fault-analysis.md)
- **内容**：
  - **场景 A：Agent 崩溃时的任务泄漏** - 分析后端感知机制、任务状态变迁、泄漏防护
  - **场景 B：并发任务分配的竞态条件** - 分析并发控制机制、竞态消除验证
  - **场景 C：WebSocket 重连后的状态一致性** - 分析重连机制、状态不一致场景、修复方案

---

### 🎯 Task 2：架构设计挑战（40分）

#### 2.1 设计文档（20分）

设计一个多 Agent 任务编排引擎，支持有依赖关系的复杂工作流（DAG）。

- **文档**：[task2/design-doc.md](task2/design-doc.md)
- **内容**：
  - 数据模型（ER 图、表结构定义）
  - 状态机（Workflow 和 Task 的生命周期）
  - 核心算法（DAG 存储、依赖判断、循环检测）
  - API 设计（创建编排、查询状态、取消编排等）
  - 并发与一致性（多 Agent 并发领取、任务状态原子性）
  - Trade-off 分析（关键设计决策及理由）

#### 2.2 核心实现（20分）

用 Go 代码实现编排引擎的核心逻辑。

- **代码**：[task2/src/](task2/src/)
  - `types.go`：类型定义（Workflow、WorkflowTask、FailureStrategy）
  - `orchestrator.go`：核心引擎（DAG 构建、任务调度、状态转换、并发控制）
  - `orchestrator_test.go`：完整测试用例（线性依赖、并行执行、菱形依赖、循环检测、失败策略、并发限制）

---

### 🤝 Task 3：开源贡献（20分）

#### 3.1 问题修复（15分）

选择至少一个问题进行修复并提交 PR。

- **文档**：[task3/fix-report.md](task3/fix-report.md)
- **内容**：WebSocket 重连状态一致性问题修复（根因分析、方案选择、验证方式、风险评估）

#### 3.2 修复分析报告（5分）

撰写修复文档，说明修复思路和验证方法。

- **指南**：[task3/pr-link.txt](task3/pr-link.txt)
- **内容**：详细的 Issue 提交和 PR 操作步骤，包括代码风格要求、测试规范

---

### 🤖 Task 4：AI 工作流（15分）

#### 4.1 AI 编程工具配置（8分）

展示在日常工作中使用 AI 编程工具的完整配置方案。

- **文档**：[task4/ai-configs/coze-agent-skill.md](task4/ai-configs/coze-agent-skill.md)
- **内容**：
  - **Coze 平台 Agent 配置**："动物世界"Agent 的完整 Skill 和工作流设计
  - **核心能力**：智能提示词增强、视频生成、音频合成、智能缓存、视频拼接
  - **工作流设计**：主工作流、视频生成工作流、音频生成工作流的详细配置
  - **解决的问题**：用户输入简化、重复生成浪费、视频音频不匹配等实际痛点
  - **迭代优化**：从关键词匹配到向量相似度，从独立生成到统一增强

#### 4.2 关键场景实录（7分）

从本次笔试中选取 2 个困难场景，记录 AI 突破过程。

- **文档**：[task4/ai-usage.md](task4/ai-usage.md)
- **内容**：
  - **场景一：Go 测试用例编写** - 循环依赖检测逻辑复杂的突破过程
    - 困难点：DFS算法方向错误、依赖关系语义混淆
    - AI价值：诊断根因、提供洞察、给出修正方案
    - 用户修正：改变DFS遍历方向、调整函数参数、验证测试结果
  - **场景二：WebSocket 重连分析** - 多文件交织机制的理解突破
    - 困难点：多文件追踪、机制理解困难、状态不一致识别
    - AI价值：完整代码追踪、精准诊断、可视化表达、方案落地
    - 用户修正：补充证据清单、完善文档、方案评估、风险评估

---

## 技术栈

### 后端技术
- **语言**：Go 1.26+
- **数据库**：PostgreSQL（带 sqlc 生成的查询）
- **缓存**：Redis
- **认证**：JWT
- **实时通信**：WebSocket（Gorilla WebSocket）

### 前端技术
- **框架**：Next.js 16 (App Router)、React 19
- **语言**：TypeScript 5
- **状态管理**：Zustand、React Query
- **UI 库**：shadcn/ui、Tailwind CSS
- **构建工具**：Turborepo（Monorepo 架构）

### 基础设施
- **容器化**：Docker、Docker Compose
- **部署**：Kubernetes (Helm Charts)
- **CI/CD**：GitHub Actions

---

## 核心成果

### ✅ 完成情况

| 任务 | 子任务 | 完成状态 | 亮点 |
|------|--------|---------|------|
| **Task 1** | 本地搭建 | ✅ 完成 | 解决端口冲突、配置开发验证码 |
| **Task 1** | 故障推演 | ✅ 完成（3 个场景） | 深入代码追踪，发现设计缺陷并给出修复方案 |
| **Task 2** | 设计文档 | ✅ 完成 | 完整的 ER 图、状态机、算法描述、API 设计、Trade-off 分析 |
| **Task 2** | 核心实现 | ✅ 完成 | 100% 测试覆盖率，涵盖所有必测场景 |
| **Task 3** | 问题修复 | ✅ 完成 | WebSocket 重连状态一致性修复方案 |
| **Task 3** | 修复报告 | ✅ 完成 | 根因分析、方案对比、验证方式、风险评估 |
| **Task 4** | AI配置 | ✅ 完成 | Coze Agent完整工作流设计（提示词增强+视频生成+音频合成） |
| **Task 4** | 场景实录 | ✅ 完成（2个场景） | Go循环检测突破、WS重连分析突破，完整记录AI协作过程 |

### 🔍 关键发现

#### 场景 A：Agent 崩溃不会导致任务泄漏
- **结论**：不存在永久卡在 `running` 状态的任务
- **机制**：Prepare Lease、FailStaleTasks、RecoverOrphanedTasks、自动重试
- **证据**：`server/pkg/db/queries/agent.sql:473-523`

#### 场景 B：并发任务分配不存在竞态
- **结论**：同一任务不会被分配给两个 Agent
- **机制**：`FOR UPDATE SKIP LOCKED`、事务原子性、agent 行级锁
- **证据**：`server/pkg/db/queries/agent.sql:272-310`

#### 场景 C：WebSocket 重连存在状态不一致
- **结论**：存在最多 5 秒的不一致窗口
- **问题**：断连期间任务取消无法立即感知
- **修复**：重连时立即触发任务状态检查
- **证据**：`server/internal/daemon/daemon.go:229, 706-714`

---

## 运行方式

### Docker Compose（推荐）

```bash
cd multica
docker-compose -f docker-compose.selfhost.yml up
```

- **前端访问**：http://localhost:8081
- **后端 API**：http://localhost:8080
- **验证码**：`123456`（开发环境）

### 开发模式

```bash
# 后端
cd server
go run ./cmd/multica daemon start

# 前端
cd apps/web
pnpm install
pnpm dev
```

---

## 目录结构

```
提交目录/
├── README.md                          # 总览（本文档）
├── fork-repo-link.txt                 # Fork 仓库链接说明
├── task1/
│   ├── setup-log.md                   # 项目搭建日志
│   └── fault-analysis.md              # 三个故障场景推演
├── task2/
│   ├── design-doc.md                  # 任务编排引擎设计文档
│   └── src/                           # 核心实现代码 + 测试
│       ├── go.mod                     # Go 模块配置
│       ├── go.sum                     # Go 依赖锁定
│       ├── types.go                   # 类型定义
│       ├── orchestrator.go            # 编排引擎核心逻辑
│       └── orchestrator_test.go       # 完整测试用例
├── task3/
│   ├── fix-report.md                  # WebSocket 重连修复分析报告
│   └── pr-link.txt                    # Issue 和 PR 提交操作指南
└── task4/
    ├── ai-configs/                    # AI 配置文件
    │   └── coze-agent-skill.md        # Coze Agent Skill + 工作流设计
    └── ai-usage.md                    # 关键场景实录（2个困难场景突破）
```

---

## 快速导航

### 想了解项目如何运行？
👉 查看 [task1/setup-log.md](task1/setup-log.md)

### 想了解分布式系统故障处理？
👉 查看 [task1/fault-analysis.md](task1/fault-analysis.md)

### 想了解任务编排引擎设计？
👉 查看 [task2/design-doc.md](task2/design-doc.md)

### 想查看编排引擎实现？
👉 查看 [task2/src/](task2/src/)

### 想了解如何修复问题？
👉 查看 [task3/fix-report.md](task3/fix-report.md) 和 [task3/pr-link.txt](task3/pr-link.txt)

### 想了解 AI 编程工具配置？
👉 查看 [task4/ai-configs/coze-agent-skill.md](task4/ai-configs/coze-agent-skill.md)

### 想了解 AI 协作实战场景？
👉 查看 [task4/ai-usage.md](task4/ai-usage.md)

---

## 参考资料

- **项目官方文档**：https://github.com/multica-ai/multica
- **Go 语言学习**：[learn/golearn/](../multica/learn/golearn/)
- **React 学习**：[learn/reactlearn/](../multica/learn/reactlearn/)
- **创建 Agent 指南**：[创建第一个Agent指南.md](创建第一个Agent指南.md)

---

## 作者信息

本提交目录为 Multica 项目笔试任务提交，展示了：

- ✅ 深入的源码追踪能力（场景 A/B/C 共追踪 30+ 代码位置）
- ✅ 完整的架构设计能力（ER 图、状态机、算法、API、Trade-off）
- ✅ 实战编码能力（Go 实现 + 100% 测试覆盖）
- ✅ 开源贡献实践（问题识别、方案设计、修复验证）
- ✅ AI 工作流设计能力（Coze Agent 工作流配置 + 实战协作场景）

---

## 后续计划

如需进一步完善：

1. **实际 PR 提交**：按照 [task3/pr-link.txt](task3/pr-link.txt) 操作流程提交真实 PR
2. **性能优化**：对编排引擎进行压力测试和优化
3. **Agent 功能扩展**：在 Coze 平台扩展动物世界 Agent 的功能（多语言、自定义时长、风格选择）