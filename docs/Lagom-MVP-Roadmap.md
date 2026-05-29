# Lagom MVP 开发路线图

> 技术选型确认：Flutter + Riverpod / Go + Gin + PostgreSQL / 云端 AI 优先  
> 目标：从 0 到 MVP 上架，预计 2 个月核心功能 + 1 个月打磨测试

---

## 一、项目里程碑

```
Month 1        Month 2        Month 3
Week 1-2       Week 3-4       Week 5-6       Week 7-8       Week 9-12
├──────────────┼──────────────┼──────────────┼──────────────┼──────────────┤
│ 基础设施     │ 核心功能     │ 核心功能     │ 集成测试     │ 上架准备     │
│ 搭建         │ 开发(上)     │ 开发(下)     │ & 优化       │ & 发布       │
└──────────────┴──────────────┴──────────────┴──────────────┴──────────────┘
```

---

## 二、Phase 1: 基础设施搭建（Week 1-2）

### 2.1 后端基础设施

| 任务 | 详细内容 | 预计时间 |
|------|----------|----------|
| 项目脚手架 | Go 项目初始化、目录结构、Makefile | 1 天 |
| 配置管理 | Viper 配置加载、环境变量、配置文件模板 | 0.5 天 |
| 数据库 | PostgreSQL 部署、迁移系统（golang-migrate）、初始 Schema | 1 天 |
| Redis | Redis 部署、连接封装 | 0.5 天 |
| Gin 框架 | 路由注册、中间件链（CORS、日志、恢复） | 0.5 天 |
| JWT 认证 | JWT 生成/验证、Refresh Token 机制 | 1 天 |
| OAuth2 | Apple Sign In + Google Sign In 集成 | 1.5 天 |
| API 基础 | 统一响应格式、参数校验、错误处理 | 0.5 天 |
| 健康检查 | /health 端点、数据库连通性检查 | 0.5 天 |
| Docker | Dockerfile、docker-compose 开发环境 | 0.5 天 |

**Phase 1 交付物**：
- [ ] 可运行的 Go 服务，支持注册/登录
- [ ] 本地开发环境一键启动（docker-compose up）
- [ ] API 文档（Swagger/OpenAPI）

---

### 2.2 前端基础设施

| 任务 | 详细内容 | 预计时间 |
|------|----------|----------|
| Flutter 项目 | 初始化、目录结构、包管理 | 0.5 天 |
| 主题系统 | 色彩、字体、圆角、阴影定义 | 0.5 天 |
| Riverpod 架构 | Provider 组织、状态管理规范 | 0.5 天 |
| 本地数据库 | Drift/SQLite 配置、表定义、DAO | 1 天 |
| 安全存储 | flutter_secure_storage 封装 | 0.5 天 |
| HTTP 客户端 | Dio 配置、拦截器、错误处理 | 0.5 天 |
| 路由系统 | GoRouter 配置、页面导航 | 0.5 天 |
| 通知系统 | flutter_local_notifications 初始化 | 0.5 天 |
| 多语言 | intl 配置、中文/英文字符串 | 0.5 天 |

**Phase 1 交付物**：
- [ ] 可运行的 Flutter App，支持 iOS/Android 模拟器
- [ ] 统一的 UI 组件库（按钮、卡片、输入框）

---

## 三、Phase 2: 核心功能开发 — 上（Week 3-4）

### 3.1 后端任务

| 任务 | 详细内容 | 预计时间 |
|------|----------|----------|
| 用户模块 | 用户 CRUD、资料管理、设置 | 1 天 |
| 伙伴配置 | CompanionConfig 模型、API | 0.5 天 |
| AI 代理服务 | DeepSeek/Qwen API 封装、SSE 流式响应 | 1.5 天 |
| System Prompt | 伙伴性格 Prompt 模板、用户画像注入 | 1 天 |
| 对话历史 | ChatMessage 存储、分页查询 | 0.5 天 |
| 习惯模型 | Habit / MicroAction CRUD | 1 天 |
| 计划生成 | AI 渐进计划生成 Prompt、解析 | 1.5 天 |

**关键实现：SSE 流式对话**

```go
// handler/chat_handler.go
func (h *ChatHandler) Chat(c *gin.Context) {
    userID := c.GetString("user_id")
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse(err))
        return
    }
    
    // SSE 头部
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    
    stream, err := h.aiService.Chat(c.Request.Context(), userID, req.Message)
    if err != nil {
        c.SSEvent("error", err.Error())
        return
    }
    
    for delta := range stream {
        c.SSEvent("delta", delta)
        c.Writer.Flush()
    }
    
    c.SSEvent("done", "")
}
```

---

### 3.2 前端任务

| 任务 | 详细内容 | 预计时间 |
|------|----------|----------|
| Onboarding 流程 | 欢迎页、形象选择、5 步对话引导、伙伴诞生 | 2 天 |
| 聊天界面 | 消息列表、气泡样式、输入框、发送逻辑 | 2 天 |
| SSE 消息流 | 接收 SSE、打字机效果、流式渲染 | 1 天 |
| 伙伴头像 | 6 个基础形象 Lottie/图片、表情状态 | 1 天 |
| 首页布局 | Tab Bar、聊天主界面、导航结构 | 1 天 |

**关键实现：流式消息渲染**

```dart
// companion_notifier.dart
class CompanionNotifier extends StateNotifier<CompanionState> {
  final ApiClient _api;
  
  Future<void> sendMessage(String content) async {
    // 1. 立即显示用户消息
    state = state.copyWith(
      messages: [...state.messages, Message.user(content)],
    );
    
    // 2. 添加"正在输入"状态
    state = state.copyWith(isTyping: true);
    
    // 3. 发起 SSE 请求
    final buffer = StringBuffer();
    await for (final delta in _api.streamChat(content)) {
      buffer.write(delta);
      // 实时更新 AI 消息
      final messages = [...state.messages];
      // 替换或追加最后一条 AI 消息
      state = state.copyWith(
        messages: _updateLastAiMessage(messages, buffer.toString()),
      );
    }
    
    state = state.copyWith(isTyping: false);
  }
}
```

---

## 四、Phase 3: 核心功能开发 — 下（Week 5-6）

### 4.1 后端任务

| 任务 | 详细内容 | 预计时间 |
|------|----------|------|
| 每日推荐 | DailyOneThing 推荐算法 | 1 天 |
| 完成记录 | DailyCompletion 记录、HSI 计算 | 1 天 |
| 家园系统 | 家园元素与 HSI 的映射关系 | 0.5 天 |
| 专注计时 | FocusSession 模型、计时 API | 0.5 天 |
| 推送服务 | FCM + APNs 封装、定时提醒 | 1 天 |
| 数据同步 | 增量同步 API、冲突检测 | 1.5 天 |
| 数据导出 | 用户数据导出（JSON/CSV） | 0.5 天 |

### 4.2 前端任务

| 任务 | 详细内容 | 预计时间 |
|------|----------|----------|
| 习惯页面 | 习惯列表、HSI 进度条、详情页 | 1.5 天 |
| 每日任务卡片 | 嵌入聊天流的任务卡片、完成/跳过逻辑 | 1 天 |
| 家园页面 | 场景渲染、元素生长状态、伙伴互动 | 1.5 天 |
| 专注页面 | 计时器、场景选择、专注中界面 | 1.5 天 |
| 设置页面 | 账户、通知、隐私、数据导出 | 1 天 |
| 本地通知 | 早安、任务提醒、晚间反思提醒 | 1 天 |

---

## 五、Phase 4: 集成测试与优化（Week 7-8）

### 5.1 测试计划

| 测试类型 | 内容 | 时间 |
|----------|------|------|
| 单元测试 | 核心算法（HSI 计算、推荐算法） | 2 天 |
| 集成测试 | API 端到端测试、前后端联调 | 2 天 |
| AI 质量测试 | Prompt 调优、回复质量评估 | 2 天 |
| 用户体验测试 | 邀请 5-10 位目标用户内测 | 3 天 |
| 性能测试 | AI 响应延迟、App 启动时间 | 1 天 |

### 5.2 AI Prompt 调优清单

```
□ 确保回复不超过 100 字
□ 失败时绝不用批评语气
□ 正确识别用户情绪并调整语气
□ 不在用户忙碌时推荐任务
□ 用户说"不想做"时提供降低门槛选项
□ 用户说"我做到了"时给予具体表扬
□ 避免重复的回复模板
□ 正确使用 emoji，不过度
```

### 5.3 性能优化

| 优化项 | 目标 |
|--------|------|
| App 冷启动 | < 2 秒 |
| AI 首字响应 | < 3 秒 |
| 页面切换 | < 300ms |
| 本地数据库查询 | < 50ms |
| 内存占用 | < 150MB |

---

## 六、Phase 5: 上架准备与发布（Week 9-12）

### 6.1 上架准备

| 任务 | 详细内容 | 时间 |
|------|----------|------|
| 应用商店资料 | 截图、描述、关键词、隐私政策 | 3 天 |
| 图标与启动图 | App Icon（各尺寸）、启动页 | 2 天 |
| 隐私政策 | GDPR/个保法合规、隐私政策页面 | 2 天 |
| 用户协议 | 服务条款、付费协议 | 1 天 |
| 测试账号 | 提供给审核员的测试账号 | 0.5 天 |
| 应用审核 | iOS App Store 审核（通常 1-3 天） | 3 天 |
| | Android Google Play 审核（通常 1-3 天） | 3 天 |
| | 国内应用商店（小米/华为/OPPO/vivo） | 各 3-5 天 |

### 6.2 应用商店截图规范

```
必需截图（iOS: 6.5" + 5.5", Android: 多种尺寸）:

1. 聊天首页（展示 AI 伙伴对话 + 今日任务卡片）
2. 习惯页面（展示 HSI 进度 + 家园预览）
3. 专注页面（展示计时器 + 场景选择）
4. 伙伴创建（展示 onboarding 的温馨感）
5. 成长报告（展示温暖的周报样式）
6. 设置/隐私（展示数据本地化、导出功能）
```

### 6.3 发布策略

```
Week 9-10: TestFlight / Google Play 内测
  └── 邀请 20-50 位种子用户
  └── 收集反馈、修复 bug

Week 11: 软发布（Soft Launch）
  └── 仅限部分国家/地区
  └── 监控崩溃率、留存率

Week 12: 正式发布
  └── 全量上架
  └── 配合小红书 / X 产品介绍
```

---

## 七、技术债务与后续迭代

### 7.1 MVP 已知技术债务

| 债务 | 原因 | 偿还计划 |
|------|------|----------|
| 纯云端 AI | 端侧模型不成熟 | V2 引入端侧意图识别 |
| 简单同步 | MVP 优先功能 | V2 引入 Operational Transform |
| 无语音 | 开发时间限制 | V2 接入 TTS/STT |
| 基础通知 | 未做智能触发 | V2 引入地理围栏、行为触发 |
| 单服务器 | 成本考虑 | 用户 > 1万时考虑负载均衡 |

### 7.2 版本规划

| 版本 | 时间 | 核心功能 |
|------|------|----------|
| v1.0 (MVP) | Month 3 | AI 伙伴对话、微习惯、专注计时、家园 |
| v1.1 | Month 4 | 语音交互、情绪感知、智能提醒 |
| v1.2 | Month 5 | 成长报告、好友陪伴模式 |
| v2.0 | Month 6 | 端侧 AI、深度数据分析、Web 端 |
| v2.1 | Month 7+ | 社区功能、合作伙伴内容、API 开放平台 |

---

## 八、开发检查清单

### 每日开发 Checklist

```
□ 今天提交代码前运行了测试？
□ 新增 API 有对应文档？
□ UI 改动在不同尺寸屏幕上测试过？
□ AI Prompt 改动有记录？
□ 敏感信息未提交到代码库？
□ 本地化字符串已更新？
```

### MVP 完成标准

```
□ iOS 和 Android 均可正常安装运行
□ 新用户可完成 Onboarding 创建伙伴
□ 可与 AI 伙伴进行自然对话
□ 可创建习惯并获得每日推荐
□ 可标记完成/跳过
□ 可进行专注计时
□ 可查看家园可视化
□ 支持 Apple/Google 登录
□ 支持本地数据存储
□ 付费用户支持云端同步
□ App Store / Google Play 审核通过
□ 崩溃率 < 1%
□ AI 响应延迟 < 5 秒（90分位）
```

---

## 九、风险应对

| 风险 | 概率 | 应对策略 |
|------|------|----------|
| AI API 成本超预期 | 中 | 设置每日调用上限、缓存常见回复 |
| App Store 审核被拒 | 中 | 提前阅读审核指南、预留修改时间 |
| 开发进度延期 | 高 | MVP 范围可缩减（优先聊天+习惯） |
| AI 回复质量不稳定 | 中 | 丰富的 fallback 回复库、人工审核机制 |
| 用户反馈冷淡 | 中 | 早期种子用户深度访谈、快速迭代 |
| 服务器宕机 | 低 | 基础监控告警、数据自动备份 |
