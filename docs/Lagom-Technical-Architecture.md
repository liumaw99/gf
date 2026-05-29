# Lagom 技术架构详细设计

> 技术选型（已确认）：Flutter + Riverpod / Go + Gin + PostgreSQL / 云端 AI 优先  
> 架构原则：单体服务、简单可维护、独立开发者友好

---

## 一、技术栈总览

```
┌─────────────────────────────────────────────────────────────┐
│                        客户端层                              │
│  Flutter 3.x + Dart 3                                       │
│  ├── iOS / Android（首发）                                   │
│  └── macOS / Windows / Web（后续）                           │
├─────────────────────────────────────────────────────────────┤
│                        接入层                                │
│  CDN + HTTPS + WSS                                          │
│  └── Cloudflare / 阿里云 CDN（静态资源+API 加速）             │
├─────────────────────────────────────────────────────────────┤
│                        服务端层                              │
│  Go 1.22+ 单体服务（单进程单容器）                            │
│  ├── Gin Web 框架                                            │
│  ├── 内置模块：用户/习惯/对话/AI代理/推送/定时任务              │
│  └── 部署：Docker + 单云服务器                               │
├─────────────────────────────────────────────────────────────┤
│                        数据层                                │
│  PostgreSQL 15+（主存储）                                    │
│  ├── pgvector（向量存储，AI 记忆）                            │
│  └── Redis 7+（缓存 + 会话 + 限流）                          │
├─────────────────────────────────────────────────────────────┤
│                        AI 层                                 │
│  云端 API 优先（MVP）                                        │
│  ├── DeepSeek-V3 / Qwen-Max（主模型）                        │
│  ├── GPT-4o-mini（Fallback）                                 │
│  └── 讯飞/阿里云（语音 TTS/STT，V2）                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 二、前端架构（Flutter）

### 2.1 项目结构

```
lagom_app/
├── android/                      # Android 平台配置
├── ios/                          # iOS 平台配置
├── lib/
│   ├── main.dart                 # 应用入口
│   ├── app.dart                  # MaterialApp + 主题配置
│   ├── config/                   # 配置管理
│   │   ├── constants.dart        # 常量（API地址、超时等）
│   │   ├── themes.dart           # 主题定义
│   │   └── routes.dart           # 路由表
│   ├── core/                     # 核心基础设施
│   │   ├── api/                  # HTTP 客户端封装
│   │   │   ├── dio_client.dart   # Dio 配置（拦截器、超时）
│   │   │   ├── api_exception.dart
│   │   │   └── interceptors/     # JWT、日志、重试拦截器
│   │   ├── database/             # 本地数据库
│   │   │   ├── app_database.dart     # Drift/SQLite 配置
│   │   │   ├── daos/                 # Data Access Objects
│   │   │   └── tables/               # 表定义
│   │   ├── storage/              # 本地存储
│   │   │   ├── secure_storage.dart   # 敏感数据（flutter_secure_storage）
│   │   │   └── prefs_storage.dart    # 轻量配置（shared_preferences）
│   │   ├── notifications/        # 本地通知
│   │   │   └── notification_service.dart
│   │   └── ai/                   # AI 相关封装
│   │       ├── chat_stream.dart      # SSE 流式对话
│   │       └── voice_service.dart    # 语音服务（V2）
│   ├── features/                 # 功能模块（按特性组织）
│   │   ├── companion/            # AI 伙伴模块
│   │   │   ├── companion_screen.dart
│   │   │   ├── widgets/
│   │   │   └── companion_notifier.dart
│   │   ├── habits/               # 习惯模块
│   │   │   ├── habits_screen.dart
│   │   │   ├── daily_one_thing.dart
│   │   │   └── habits_notifier.dart
│   │   ├── focus/                # 专注模块
│   │   │   ├── focus_screen.dart
│   │   │   ├── focus_timer.dart
│   │   │   └── focus_notifier.dart
│   │   ├── journal/              # 成长档案模块
│   │   │   ├── journal_screen.dart
│   │   │   └── reflection_flow.dart
│   │   ├── onboarding/           # 首次引导
│   │   │   ├── welcome_screen.dart
│   │   │   └── companion_creation_flow.dart
│   │   └── settings/             # 设置模块
│   │       ├── settings_screen.dart
│   │       └── privacy_settings.dart
│   ├── models/                   # 数据模型（跨模块共享）
│   │   ├── user.dart
│   │   ├── companion.dart
│   │   ├── habit.dart
│   │   ├── daily_log.dart
│   │   └── focus_session.dart
│   └── providers/                # 全局 Riverpod Providers
│       ├── auth_provider.dart
│       ├── companion_provider.dart
│       └── habit_provider.dart
├── test/                         # 单元测试
└── pubspec.yaml
```

### 2.2 核心依赖

```yaml
name: lagom
description: AI-powered self-discipline companion
publish_to: 'none'
version: 1.0.0+1

environment:
  sdk: '>=3.4.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter

  # 状态管理
  flutter_riverpod: ^2.5.1
  riverpod_annotation: ^2.3.5

  # 网络
  dio: ^5.7.0
  retrofit: ^4.4.0

  # 本地数据库
  drift: ^2.21.0
  sqlite3_flutter_libs: ^0.5.24

  # 本地存储
  shared_preferences: ^2.3.2
  flutter_secure_storage: ^9.2.2

  # 通知
  flutter_local_notifications: ^17.2.3
  timezone: ^0.9.4

  # 后台任务
  workmanager: ^0.5.2

  # UI 组件
  fl_chart: ^0.68.0          # 图表
  shimmer: ^3.0.0            # 骨架屏
  flutter_animate: ^4.5.0    # 动画
  lottie: ^3.1.2             # Lottie 动画

  # 语音（V2）
  flutter_tts: ^4.0.2
  speech_to_text: ^6.6.0

  # 工具
  freezed_annotation: ^2.4.4
  json_annotation: ^4.9.0
  uuid: ^4.5.1
  intl: ^0.19.0
  logger: ^2.4.0

  # 依赖注入（配合 Riverpod）
  get_it: ^7.7.0

dependency_overrides:
  # 如需覆盖版本

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^4.0.0
  build_runner: ^2.4.13
  freezed: ^2.5.7
  json_serializable: ^6.8.0
  riverpod_generator: ^2.4.3
  retrofit_generator: ^9.1.0
  drift_dev: ^2.21.0
```

### 2.3 状态管理设计

```dart
// 核心状态架构

// 1. 认证状态
@Riverpod(keepAlive: true)
class AuthNotifier extends _$AuthNotifier {
  @override
  AuthState build() => const AuthState.unauthenticated();

  Future<void> signInWithApple() async { ... }
  Future<void> signInWithGoogle() async { ... }
  Future<void> signOut() async { ... }
}

// 2. 伙伴状态
@Riverpod(keepAlive: true)
class CompanionNotifier extends _$CompanionNotifier {
  @override
  Future<Companion> build() async {
    // 优先从本地加载
    final local = await _localDb.getCompanion();
    if (local != null) return local;
    // 新用户返回默认
    return Companion.defaultCompanion();
  }

  Future<void> sendMessage(String content) async { ... }
  Future<void> updateMood(UserMood mood) async { ... }
}

// 3. 习惯状态
@Riverpod(keepAlive: true)
class HabitNotifier extends _$HabitNotifier {
  @override
  Future<List<Habit>> build() => _localDb.getAllHabits();

  Future<void> completeAction(String habitId) async { ... }
  Future<void> createHabit(String goalDescription) async { ... }
}

// 4. 每日一件小事（计算属性）
@riverpod
Future<MicroAction?> dailyOneThing(Ref ref) async {
  final habits = await ref.watch(habitNotifierProvider.future);
  final mood = ref.watch(companionNotifierProvider).valueOrNull?.detectedMood;
  return DailyOneThingRecommender.recommend(habits, mood);
}
```

### 2.4 本地数据流设计

```
【离线优先架构】

用户操作 → 本地数据库（Drift/SQLite）→ UI 即时更新
                ↓
         后台同步队列
                ↓
         网络可用时 → 服务端同步
                ↓
         冲突解决（最后写入优先 + 时间戳）

核心原则：
- 所有读写操作先操作本地数据库
- 网络请求异步进行，不阻塞 UI
- 同步失败时自动重试，用户无感知
- 付费用户的云端同步在后台静默完成
```

---

## 三、后端架构（Go 单体服务）

### 3.1 项目结构

```
lagom-server/
├── cmd/
│   └── server/
│       └── main.go                 # 应用入口
├── config/
│   ├── config.go                   # 配置结构体
│   └── config.yaml                 # 配置文件模板
├── internal/
│   ├── app/                        # 应用初始化
│   │   ├── app.go                  # 依赖注入、生命周期管理
│   │   └── router.go               # 路由注册
│   ├── handler/                    # HTTP 处理器
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── habit_handler.go
│   │   ├── chat_handler.go         # SSE 流式对话
│   │   ├── sync_handler.go
│   │   └── health_handler.go
│   ├── service/                    # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── habit_service.go
│   │   ├── ai_service.go           # AI 代理服务
│   │   ├── sync_service.go
│   │   └── notification_service.go
│   ├── repository/                 # 数据访问层
│   │   ├── user_repo.go
│   │   ├── habit_repo.go
│   │   ├── chat_repo.go
│   │   └── memory_repo.go          # AI 记忆存储
│   ├── model/                      # 数据模型
│   │   ├── user.go
│   │   ├── habit.go
│   │   ├── companion.go
│   │   └── chat.go
│   ├── middleware/                 # 中间件
│   │   ├── auth.go                 # JWT 认证
│   │   ├── cors.go
│   │   ├── rate_limit.go           # 限流
│   │   └── logger.go               # 请求日志
│   ├── pkg/                        # 内部工具包
│   │   ├── jwt/                    # JWT 工具
│   │   ├── encrypt/                # 加密工具
│   │   ├── ai/                     # AI 客户端封装
│   │   │   ├── deepseek.go
│   │   │   ├── qwen.go
│   │   │   └── fallback.go
│   │   └── push/                   # 推送封装
│   │       ├── fcm.go
│   │       └── apns.go
│   └── cron/                       # 定时任务
│       └── scheduler.go            # 每日提醒、报告生成
├── migrations/                     # 数据库迁移
│   └── 001_init.sql
├── Dockerfile
├── go.mod
└── Makefile
```

### 3.2 核心依赖

```go
// go.mod
module github.com/lagom/lagom-server

go 1.22

require (
    // Web 框架
    github.com/gin-gonic/gin v1.10.0
    
    // 配置
    github.com/spf13/viper v1.19.0
    
    // 数据库
    github.com/jackc/pgx/v5 v5.7.1        // PostgreSQL driver
    github.com/jackc/pgxvector v0.0.3    // pgvector 支持
    
    // ORM（轻量，可选 raw SQL）
    github.com/jmoiron/sqlx v1.4.0
    
    // Redis
    github.com/redis/go-redis/v9 v9.6.1
    
    // JWT
    github.com/golang-jwt/jwt/v5 v5.2.1
    
    // OAuth2
    golang.org/x/oauth2 v0.23.0
    
    // 验证
    github.com/go-playground/validator/v10 v10.22.1
    
    // 日志
    go.uber.org/zap v1.27.0
    
    // 定时任务
    github.com/robfig/cron/v3 v3.0.1
    
    // 测试
    github.com/stretchr/testify v1.9.0
)
```

### 3.3 API 设计

#### 认证相关

```http
POST   /api/v1/auth/apple          # Apple 登录
POST   /api/v1/auth/google         # Google 登录
POST   /api/v1/auth/refresh        # 刷新 Token
DELETE /api/v1/auth/logout         # 退出登录
```

#### 用户相关

```http
GET    /api/v1/user/profile        # 获取用户信息
PUT    /api/v1/user/profile        # 更新用户信息
GET    /api/v1/user/settings       # 获取设置
PUT    /api/v1/user/settings       # 更新设置
```

#### 习惯相关

```http
GET    /api/v1/habits              # 获取习惯列表
POST   /api/v1/habits              # 创建习惯（AI生成计划）
GET    /api/v1/habits/:id          # 获取习惯详情
PUT    /api/v1/habits/:id          # 更新习惯
DELETE /api/v1/habits/:id          # 删除习惯
POST   /api/v1/habits/:id/complete # 标记完成
POST   /api/v1/habits/:id/skip     # 跳过（缓冲日）
```

#### AI 对话（核心）

```http
POST   /api/v1/chat                # 发送消息（SSE 流式返回）
GET    /api/v1/chat/history        # 获取对话历史
DELETE /api/v1/chat/history        # 清空对话（保留记忆）
```

SSE 响应格式：
```
data: {"type":"delta","content":"嗨"}
data: {"type":"delta","content="！"}
data: {"type":"done"}
```

#### 同步相关

```http
POST   /api/v1/sync                # 上传本地变更
GET    /api/v1/sync                # 获取服务器变更
POST   /api/v1/sync/resolve        # 冲突解决
```

#### 数据导出

```http
GET    /api/v1/export              # 导出用户所有数据
```

### 3.4 AI 代理服务设计

```go
// internal/service/ai_service.go

type AIService struct {
    primary    AIClient      // DeepSeek / Qwen
    fallback   AIClient      // GPT-4o-mini
    memoryRepo MemoryRepo    // 长期记忆存储
    chatRepo   ChatRepo      // 对话历史
}

func (s *AIService) Chat(ctx context.Context, userID string, message string) (<-chan string, error) {
    // 1. 获取用户画像 + 记忆摘要
    profile := s.memoryRepo.GetUserProfile(ctx, userID)
    memorySummary := s.memoryRepo.GetRecentSummary(ctx, userID)
    
    // 2. 构建 System Prompt
    systemPrompt := buildSystemPrompt(profile, memorySummary)
    
    // 3. 获取最近对话历史
    history := s.chatRepo.GetRecentHistory(ctx, userID, 10)
    
    // 4. 调用 AI（流式）
    stream, err := s.primary.StreamChat(ctx, systemPrompt, history, message)
    if err != nil {
        stream, err = s.fallback.StreamChat(ctx, systemPrompt, history, message)
    }
    
    // 5. 异步保存完整回复 + 更新记忆
    go func() {
        fullResponse := collectStream(stream)
        s.chatRepo.SaveMessage(ctx, userID, "assistant", fullResponse)
        s.memoryRepo.UpdateSummary(ctx, userID, message, fullResponse)
    }()
    
    return stream, nil
}

// System Prompt 构建
func buildSystemPrompt(profile UserProfile, memory string) string {
    return fmt.Sprintf(`你是 Lagom 的 AI 伙伴，名字是 %s。

用户画像：
- 性格：%s
- 鼓励风格：%s
- 当前目标：%s

长期记忆摘要：
%s

沟通原则：
1. 每次回复不超过 100 字，保持轻松自然
2. 根据用户的情绪状态调整语气
3. 失败时给予共情，不批评
4. 适当使用emoji增加温度
5. 在合适的时候推荐"今天的一件小事"
`, profile.CompanionName, profile.Personality, 
   profile.EncouragementStyle, profile.CurrentGoals, memory)
}
```

### 3.5 数据模型（数据库 Schema）

```sql
-- migrations/001_init.sql

-- 用户表
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    apple_id VARCHAR(255) UNIQUE,
    google_id VARCHAR(255) UNIQUE,
    email VARCHAR(255),
    name VARCHAR(100),
    avatar_url TEXT,
    subscription_tier VARCHAR(20) DEFAULT 'free', -- free, premium
    subscription_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 伙伴配置表
CREATE TABLE companion_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    avatar_type VARCHAR(20) NOT NULL, -- fox, cat, sloth, lion, rabbit, bear
    personality_type VARCHAR(20) NOT NULL, -- gentle, energetic, humorous, calm
    encouragement_style VARCHAR(20) NOT NULL, -- soft, direct, humor, silent
    voice_type VARCHAR(20),
    growth_level INT DEFAULT 1,
    relationship_score DECIMAL(3,2) DEFAULT 0.0, -- 0.0 to 1.0
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id)
);

-- 习惯表
CREATE TABLE habits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    category VARCHAR(50) NOT NULL, -- health, learning, fitness, mindfulness, productivity
    description TEXT,
    priority INT DEFAULT 3, -- 1-5
    frequency_type VARCHAR(20) DEFAULT 'daily', -- daily, weekly, custom
    target_days INT[] DEFAULT '{}', -- 对于 weekly/custom，记录目标星期几
    is_active BOOLEAN DEFAULT true,
    current_phase INT DEFAULT 1, -- 当前阶段（渐进计划）
    hsi DECIMAL(5,2) DEFAULT 0.0, -- 习惯强度指数
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 微行动表
CREATE TABLE micro_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    phase INT NOT NULL,
    day_in_phase INT NOT NULL,
    action_text VARCHAR(300) NOT NULL,
    estimated_duration INT, -- 预计完成时间（分钟）
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 每日记录表
CREATE TABLE daily_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    mood_score INT CHECK (mood_score BETWEEN 1 AND 5),
    reflection_note TEXT,
    summary TEXT, -- AI 生成的今日摘要
    UNIQUE(user_id, date)
);

-- 每日完成记录表
CREATE TABLE daily_completions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    completed_at TIMESTAMP DEFAULT NOW(),
    quality_score INT CHECK (quality_score BETWEEN 1 AND 5), -- 完成质量自评
    note TEXT,
    UNIQUE(user_id, habit_id, date)
);

-- 对话历史表
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引优化查询
CREATE INDEX idx_chat_messages_user_created ON chat_messages(user_id, created_at DESC);
CREATE INDEX idx_daily_logs_user_date ON daily_logs(user_id, date DESC);
CREATE INDEX idx_daily_completions_user_date ON daily_completions(user_id, date DESC);
CREATE INDEX idx_habits_user_active ON habits(user_id, is_active);

-- pgvector 扩展（用于 AI 记忆向量存储）
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE memory_vectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    embedding vector(1536), -- 根据使用的 embedding 模型调整维度
    memory_type VARCHAR(50) NOT NULL, -- fact, preference, goal, reflection
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_memory_vectors_user ON memory_vectors(user_id);
```

---

## 四、AI 架构设计

### 4.1 AI 服务流程

```
用户发送消息
    ↓
[意图识别] 轻量分类器 / 规则匹配
    ├── 闲聊 → 直接走 LLM 对话
    ├── 打卡完成 → 记录 + 伙伴回应
    ├── 创建习惯 → 调用计划生成
    ├── 情绪表达 → 情绪分析 + 共情回应
    └── 未知 → 走 LLM 兜底
    ↓
[上下文组装]
    ├── System Prompt（伙伴性格 + 用户画像）
    ├── 长期记忆摘要（pgvector 检索）
    └── 最近 10 轮对话
    ↓
[流式调用 LLM API]
    ├── 主模型：DeepSeek-V3
    ├── 超时 / 失败 → 降级 GPT-4o-mini
    └── 全部失败 → 预设 fallback 回复
    ↓
[返回 SSE 流]
    ↓
[异步后处理]
    ├── 保存完整回复到数据库
    ├── 提取关键事实 → 更新记忆向量
    └── 情绪分析 → 更新用户 mood
```

### 4.2 记忆管理

```go
// 长期记忆系统

type MemoryManager struct {
    vectorStore *pgvector.Store
}

// 保存新记忆
func (m *MemoryManager) SaveMemory(ctx context.Context, userID, content, memoryType string) error {
    // 1. 生成 embedding（调用 embedding API）
    embedding := m.generateEmbedding(content)
    
    // 2. 存入向量数据库
    return m.vectorStore.Insert(ctx, userID, content, embedding, memoryType)
}

// 检索相关记忆
func (m *MemoryManager) RetrieveRelevant(ctx context.Context, userID, query string, limit int) ([]Memory, error) {
    // 1. 生成 query 的 embedding
    queryEmbedding := m.generateEmbedding(query)
    
    // 2. 向量相似度搜索
    return m.vectorStore.SimilaritySearch(ctx, userID, queryEmbedding, limit)
}

// 定期生成记忆摘要（减少上下文长度）
func (m *MemoryManager) GenerateSummary(ctx context.Context, userID string) (string, error) {
    // 每周运行一次
    // 将最近 7 天的记忆聚合成一段摘要
    memories := m.vectorStore.GetRecent(ctx, userID, 7*24*time.Hour)
    summary := m.callLLMToSummarize(memories)
    return summary, nil
}
```

### 4.3 成本估算

| 功能 | 单次调用成本 | 月均调用量（1000用户） | 月成本 |
|------|-------------|----------------------|--------|
| 对话（DeepSeek-V3） | ~¥0.003 / 1K tokens | 300K 次 | ~¥900 |
| 对话（GPT-4o-mini fallback） | ~$0.0006 / 1K tokens | 30K 次 | ~$18 |
| Embedding | ~¥0.001 / 1K tokens | 50K 次 | ~¥50 |
| 语音 TTS（V2） | ~¥0.02 / 次 | 100K 次 | ~¥2000 |
| **合计** | | | **~¥3000-4000 / 月** |

> 注：1000 付费用户按 $3.99/月 = $3990 收入，可覆盖 AI 成本。前期免费用户阶段需控制成本。

---

## 五、部署架构

### 5.1 单体部署方案

```
【独立开发者最小成本部署】

阿里云 / 腾讯云 / DigitalOcean
├── ECS / Droplet: 2C4G
│   └── Docker 运行 lagom-server
│   └── 内置 PostgreSQL（或使用云数据库）
│
├── 云数据库 PostgreSQL: 最小规格
│   └── 支持 pgvector 扩展
│
├── Redis: 云托管最小规格（或本地部署）
│
└── CDN: 静态资源加速

预估月成本：
- 服务器：¥100-200
- 数据库：¥50-100
- CDN：¥20-50
- AI API：按量计费
- 总计：¥200-400（基础）+ AI 调用费用
```

### 5.2 Docker Compose（开发/小流量）

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - REDIS_HOST=redis
      - DEEPSEEK_API_KEY=${DEEPSEEK_API_KEY}
    depends_on:
      - postgres
      - redis

  postgres:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_DB: lagom
      POSTGRES_USER: lagom
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"

volumes:
  postgres_data:
  redis_data:
```

### 5.3 CI/CD 流程

```
GitHub 推送
    ↓
GitHub Actions
    ├── 运行测试
    ├── 构建 Docker 镜像
    ├── 推送到镜像仓库
    └── SSH 到服务器执行 docker-compose pull && up -d
```

---

## 六、安全设计

### 6.1 认证安全

| 措施 | 实现 |
|------|------|
| JWT | RS256 签名，Access Token 15 分钟过期 |
| Refresh Token | 7 天过期，存储于 Redis 黑名单 |
| OAuth2 | Apple Sign In + Google Sign In，不存储密码 |
| 设备绑定 | Refresh Token 绑定设备指纹 |

### 6.2 数据安全

| 措施 | 实现 |
|------|------|
| 传输加密 | 全站 HTTPS，TLS 1.3 |
| 存储加密 | 云端数据 AES-256-GCM |
| 本地数据 | iOS Keychain / Android Keystore |
| 数据脱敏 | 日志中脱敏用户 ID 和敏感内容 |
| 数据删除 | 用户注销时 30 天内物理删除 |

### 6.3 API 安全

| 措施 | 实现 |
|------|------|
| 限流 | 单 IP 100 req/min，单用户 60 req/min |
| AI 调用限流 | 单用户 50 次/天（免费），无限（付费） |
| 输入验证 | 严格参数校验，防 SQL 注入/XSS |
| CORS | 仅允许 App 域名 |

---

## 七、监控与日志

### 7.1 监控指标

```
应用层面：
- HTTP 请求 QPS / 延迟 / 错误率
- AI API 调用成功率 / 延迟 / 成本
- 活跃用户数 / 留存率（需埋点）

业务层面：
- 每日 habit 完成率
- AI 对话轮次 / 满意度
- 付费转化率 / 续订率

基础设施：
- CPU / 内存 / 磁盘
- 数据库连接数 / 查询耗时
- Redis 内存使用
```

### 7.2 日志规范

```go
// 结构化日志（Zap）
logger.Info("chat completed",
    zap.String("user_id", userID),
    zap.Int("tokens_used", tokens),
    zap.Duration("latency", latency),
    zap.String("model", model),
)

// 错误日志
logger.Error("ai request failed",
    zap.String("user_id", userID),
    zap.Error(err),
    zap.String("fallback_model", fallbackModel),
)
```

---

## 八、扩展性考虑

### 8.1 未来可能的拆分点

当前单体架构在以下情况时考虑拆分：

| 信号 | 拆分方案 |
|------|----------|
| AI 调用量 > 10K/分钟 | AI 服务独立部署，支持水平扩展 |
| 用户量 > 100万 | 用户服务拆分，支持分库分表 |
| 实时推送需求增加 | 引入 WebSocket 服务或 MQTT |
| 多端同步需求复杂 | 独立 Sync 服务，支持 OT 算法 |

### 8.2 端侧 AI 演进路线

```
Phase 1（MVP）: 纯云端 AI
Phase 2（V2）: 简单意图识别端侧化（减少 API 调用）
Phase 3（V3）: 端侧小模型（Gemma 2B）处理基础对话
Phase 4（V4）: 端侧处理 80% 对话，云端仅处理复杂推理
```
