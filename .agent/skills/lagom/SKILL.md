# Lagom 项目开发规范

> 项目级 Skill，用于规范 Lagom（AI 自律伙伴 App）的全部开发活动  
> 技术栈：Flutter + Riverpod / Go + Gin + PostgreSQL / 云端 AI（DeepSeek）  
> 适用范围：`/Users/fidoo/Desktop/liumaw99/gf` 目录及其子目录

---

## 一、项目概述

### 1.1 产品定位

Lagom 是一款基于 AI 情感陪伴的自律养成应用。核心差异化：**不是打卡工具，而是"懂你节奏的 AI 伙伴"**。用户可选择基于真实名人/虚构人物蒸馏而成的自律伙伴，以该人物的思维方式、语言风格陪伴成长。

### 1.2 核心功能

- **AI 伙伴对话**：自然语言交互，情绪感知，主动关怀
- **微习惯引擎**：每日只推荐"一件小事"，降低执行门槛
- **名人蒸馏系统**：12+ 角色（村上春树、马斯克、芒格、科比等），风格化对话
- **专注空间**：沉浸式计时，伙伴陪伴，白噪音场景
- **成长档案**：AI 驱动的反思与报告，非数据焦虑型可视化

### 1.3 技术选型（已确认，不可随意更改）

| 层级 | 技术 | 版本 | 说明 |
|------|------|------|------|
| 前端 | Flutter | 3.x+ | 跨平台 iOS/Android |
| 前端状态管理 | Riverpod | 2.x+ | 配合 StateNotifier |
| 前端本地存储 | Drift (SQLite) | 2.21+ | 主数据存储 |
| 前端轻量存储 | shared_preferences | 2.3+ | 配置项 |
| 前端安全存储 | flutter_secure_storage | 9.2+ | Token/密钥 |
| 后端 | Go | 1.22+ | 单体服务 |
| 后端框架 | Gin | 1.10+ | Web 框架 |
| 数据库 | PostgreSQL | 15+ | 主存储 |
| 向量扩展 | pgvector | - | AI 记忆存储 |
| 缓存 | Redis | 7+ | 会话/限流/缓存 |
| AI 主模型 | DeepSeek-V3 | - | 云端 API |
| AI Fallback | GPT-4o-mini | - | 降级方案 |

---

## 二、架构原则

### 2.1 后端原则（Go）

```
1. 单体架构：所有模块在一个进程内，函数调用而非 RPC
2. 分层清晰：handler → service → repository，单向依赖
3. 依赖注入：通过构造函数注入，不用全局变量
4. 错误处理：统一包装，不裸抛原始错误给客户端
5. 上下文传递：所有 IO 操作接受 context.Context
6. 配置驱动：名人角色通过 JSONB/配置文件动态加载，不改代码
```

### 2.2 前端原则（Flutter）

```
1. 功能模块化：按 feature 组织代码，非按类型
2. 状态集中：全局状态用 Riverpod，局部状态用 StatefulWidget
3. 离线优先：所有操作先写本地 DB，后台同步
4. UI 与逻辑分离：Screen → Notifier → Service → Repository
5. 响应式：使用 AsyncValue 处理异步状态（loading/error/data）
```

### 2.3 AI 服务原则

```
1. 流式响应：SSE (Server-Sent Events) 实现打字机效果
2. Prompt 分层：基础 prompt + 角色 prompt（名人）+ 上下文
3. 记忆管理：最近 10 轮对话 + pgvector 长期记忆摘要
4. 降级策略：主模型失败 → fallback 模型 → 预设模板回复
5. 成本控制：意图缓存 + 响应缓存 + 每日调用上限
```

---

## 三、目录结构规范

### 3.1 后端目录（`server/`）

```
server/
├── cmd/server/
│   └── main.go                 # 唯一入口
├── config/
│   ├── config.go               # 配置结构体
│   └── config.yaml             # 配置文件模板
├── internal/
│   ├── app/
│   │   ├── app.go              # 应用初始化、依赖注入
│   │   └── router.go           # 路由注册
│   ├── handler/                # HTTP 处理器（仅解析请求/组装响应）
│   │   ├── auth_handler.go
│   │   ├── chat_handler.go
│   │   ├── habit_handler.go
│   │   └── character_handler.go
│   ├── service/                # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── ai_service.go       # AI 对话核心
│   │   ├── habit_service.go
│   │   └── character_service.go
│   ├── repository/             # 数据访问层
│   │   ├── user_repo.go
│   │   ├── habit_repo.go
│   │   ├── chat_repo.go
│   │   └── character_repo.go
│   ├── model/                  # 数据模型（仅结构体，无逻辑）
│   │   ├── user.go
│   │   ├── habit.go
│   │   └── character.go
│   ├── middleware/             # 中间件
│   │   ├── auth.go             # JWT 认证
│   │   ├── cors.go
│   │   ├── rate_limit.go
│   │   └── logger.go
│   ├── pkg/                    # 内部工具包（可复用，无业务逻辑）
│   │   ├── jwt/
│   │   ├── encrypt/
│   │   ├── ai/                 # AI 客户端封装
│   │   └── push/               # FCM/APNs 封装
│   └── cron/                   # 定时任务
│       └── scheduler.go
├── migrations/                 # 数据库迁移文件
│   └── 001_init.sql
├── Dockerfile
├── Makefile
└── go.mod
```

**禁止**：
- ❌ handler 直接调用 repository（必须通过 service）
- ❌ service 之间互相调用（需要则提取到 pkg 或上层协调）
- ❌ model 包引入其他 internal 包

### 3.2 前端目录（`app/lagom_app/`）

```
app/lagom_app/
├── android/                    # Android 平台配置
├── ios/                        # iOS 平台配置
├── lib/
│   ├── main.dart               # 入口
│   ├── app.dart                # MaterialApp + 主题
│   ├── config/                 # 配置
│   │   ├── constants.dart      # API 地址、超时等
│   │   ├── themes.dart         # LagomTheme
│   │   └── routes.dart         # GoRouter 路由表
│   ├── core/                   # 核心基础设施
│   │   ├── api/
│   │   │   ├── dio_client.dart
│   │   │   └── interceptors/
│   │   ├── database/
│   │   │   ├── app_database.dart
│   │   │   └── daos/
│   │   ├── storage/
│   │   ├── notifications/
│   │   └── ai/
│   ├── features/               # 按功能模块组织
│   │   ├── companion/          # AI 伙伴（聊天）
│   │   │   ├── companion_screen.dart
│   │   │   ├── widgets/
│   │   │   └── companion_notifier.dart
│   │   ├── habits/             # 习惯/微习惯
│   │   ├── focus/              # 专注空间
│   │   ├── journal/            # 成长档案
│   │   ├── characters/         # 名人角色系统
│   │   ├── onboarding/         # 首次引导
│   │   └── settings/           # 设置
│   ├── models/                 # 跨模块共享模型
│   │   ├── user.dart
│   │   ├── companion.dart
│   │   ├── habit.dart
│   │   └── character.dart
│   └── providers/              # 全局 Riverpod Providers
│       ├── auth_provider.dart
│       ├── companion_provider.dart
│       └── character_provider.dart
├── test/                       # 测试
└── pubspec.yaml
```

**禁止**：
- ❌ feature 之间直接 import（通过 models/ 或 providers/ 共享）
- ❌ UI 层直接调用 API Client（必须通过 Notifier）
- ❌ 在 Widget build 方法中发起异步请求

---

## 四、命名规范

### 4.1 Go 命名

```go
// 文件命名：snake_case.go
auth_handler.go
habit_service.go
user_repo.go

// 接口命名：动词 + er / 名词
// ❌ bad: IRepo
// ✅ good:
type UserRepo interface {}
type TokenGenerator interface {}

// 结构体命名：PascalCase，表意明确
type User struct {}
type CharacterDistillate struct {}

// 方法命名：动词开头，简洁
type UserRepo interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

// 变量命名
// ❌ bad: u, h, c
// ✅ good:
user := &User{}
habit := &Habit{}
character := &CelebrityCharacter{}

// 错误变量：以 Err 开头
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidToken = errors.New("invalid token")

// 常量：ALL_CAPS
const MaxHabitsFree = 3
const DefaultFocusDuration = 25
```

### 4.2 Dart/Flutter 命名

```dart
// 文件命名：snake_case.dart
companion_screen.dart
habit_notifier.dart
daily_recommender.dart

// 类命名：PascalCase
class CompanionScreen extends ConsumerWidget {}
class HabitNotifier extends StateNotifier<HabitState> {}

// 方法/函数：camelCase
Future<void> sendMessage(String content) async {}
void updateHsi(String habitId, double newHsi) {}

// 变量：camelCase
final companionState = ref.watch(companionNotifierProvider);
final dailyAction = recommender.recommend(habits);

// 常量：camelCase（Dart 风格）
const maxHabitsFree = 3;
const defaultFocusDuration = Duration(minutes: 25);

// Provider 命名：小写 + Provider 后缀
final authNotifierProvider = StateNotifierProvider<AuthNotifier, AuthState>(...);
final companionNotifierProvider = StateNotifierProvider<CompanionNotifier, CompanionState>(...);

// Freezed 模型：文件名与类名一致
// habit.dart → class Habit with _$Habit
```

---

## 五、代码规范

### 5.1 Go 代码规范

```go
// 1. 每个函数不超过 50 行，超过则拆分
// 2. 每个文件不超过 300 行，超过则拆分
// 3. 错误处理：立即返回，不嵌套

// ✅ Good
func (s *HabitService) CreateHabit(ctx context.Context, userID string, goal string) (*Habit, error) {
    habit, err := s.ai.GeneratePlan(ctx, goal)
    if err != nil {
        return nil, fmt.Errorf("generate plan: %w", err)
    }
    
    if err := s.repo.Create(ctx, habit); err != nil {
        return nil, fmt.Errorf("save habit: %w", err)
    }
    
    return habit, nil
}

// ❌ Bad：嵌套过深
func (s *HabitService) CreateHabit(ctx context.Context, userID string, goal string) (*Habit, error) {
    habit, err := s.ai.GeneratePlan(ctx, goal)
    if err == nil {
        if err := s.repo.Create(ctx, habit); err == nil {
            return habit, nil
        } else {
            return nil, err
        }
    } else {
        return nil, err
    }
}

// 4. context 作为第一个参数
func DoSomething(ctx context.Context, arg string) error {}

// 5. 结构体初始化用字段名
// ✅ Good
user := &model.User{
    Name:  "Alice",
    Email: "alice@example.com",
}

// 6. 日志使用结构化日志
logger.Info("habit created",
    zap.String("user_id", userID),
    zap.String("habit_id", habit.ID),
    zap.Int("phases", len(habit.Phases)),
)
```

### 5.2 Dart/Flutter 代码规范

```dart
// 1. Widget 构建方法不超过 60 行，拆分小 Widget
// ✅ Good
class CompanionScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: _buildAppBar(context, ref),
      body: _buildBody(context, ref),
      bottomNavigationBar: _buildInputBar(context, ref),
    );
  }

  Widget _buildAppBar(BuildContext context, WidgetRef ref) { ... }
  Widget _buildBody(BuildContext context, WidgetRef ref) { ... }
  Widget _buildInputBar(BuildContext context, WidgetRef ref) { ... }
}

// 2. Notifier 中不直接操作 UI
// ✅ Good：在 Notifier 中只处理状态
class CompanionNotifier extends StateNotifier<CompanionState> {
  Future<void> sendMessage(String content) async {
    state = state.copyWith(isTyping: true);
    // ... API 调用
    state = state.copyWith(isTyping: false);
  }
}

// 3. 使用 const 构造函数
// ✅ Good
const SizedBox(height: 16)
const LagomButton(text: '完成')

// 4. 字符串用单引号（除非包含单引号）
// ✅ Good
const title = '今天的一件小事';

// 5. 异步状态用 AsyncValue
// ✅ Good
final habitsAsync = ref.watch(habitsNotifierProvider);

return habitsAsync.when(
  data: (habits) => HabitList(habits: habits),
  loading: () => const ShimmerHabitList(),
  error: (err, _) => ErrorWidget(message: err.toString()),
);

// 6. 颜色/样式从 theme 获取
// ✅ Good
final theme = Theme.of(context);
final primaryColor = theme.colorScheme.primary;

// ❌ Bad
const primaryColor = Color(0xFFF4A261);
```

---

## 六、API 设计规范

### 6.1 URL 规范

```
POST   /api/v1/auth/apple          # 认证
POST   /api/v1/auth/google
POST   /api/v1/auth/refresh

GET    /api/v1/user/profile        # 用户
PUT    /api/v1/user/profile
GET    /api/v1/user/settings

GET    /api/v1/habits              # 习惯
POST   /api/v1/habits
GET    /api/v1/habits/:id
PUT    /api/v1/habits/:id
DELETE /api/v1/habits/:id
POST   /api/v1/habits/:id/complete
POST   /api/v1/habits/:id/skip

POST   /api/v1/chat                # 对话（SSE 流式）
GET    /api/v1/chat/history

GET    /api/v1/characters          # 名人角色
GET    /api/v1/characters/:id
POST   /api/v1/characters/recommend
GET    /api/v1/user/characters
POST   /api/v1/user/characters
POST   /api/v1/user/characters/:id/switch
```

### 6.2 响应格式

```json
// 成功
{
  "code": 0,
  "message": "success",
  "data": { ... }
}

// 列表
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [ ... ],
    "total": 100,
    "page": 1,
    "limit": 20
  }
}

// 错误
{
  "code": 1001,
  "message": "user not found"
}
```

### 6.3 SSE 格式

```
data: {"type":"delta","content":"嗨"}
data: {"type":"delta","content":"！"}
data: {"type":"done"}
```

### 6.4 Postman Collection 维护规范

**文件位置**：`docs/Lagom-API.postman_collection.json`

**触发条件**：以下任一情况发生时，**必须同步更新** Postman Collection：
- 新增/删除/修改路由（URL、Method）
- 修改请求参数（字段增减、类型变化、校验规则）
- 修改响应结构（data 字段变化、新增/删除字段）
- 新增/修改错误码

**更新规则**：
```
1. 保持现有 JSON 结构（Collection v2.1 格式）
2. 新增接口：在对应 folder 内新增 item，复制同类请求的 header/body 格式
3. 修改接口：更新 url、body、description，同步更新 Tests 脚本中的断言
4. 删除接口：从 JSON 中移除对应 item，同步更新 Tests 中引用的变量
5. 变量变更：如新增请求级变量，添加到 collectionVariables
6. 中文化：所有 name、description 使用中文
```

**Tests 脚本规范**（每个请求必须包含）：
```javascript
// 状态码断言
pm.test('Status code is 200', function () {
  pm.response.to.have.status(200);
});

// 业务码断言
var jsonData = pm.response.json();
pm.test('Business code is 0', function () {
  pm.expect(jsonData.code).to.eql(0);
});

// Token 自动提取（Register/Login/Refresh 专用）
if (jsonData.data && jsonData.data.access_token) {
  pm.collectionVariables.set('accessToken', jsonData.data.access_token);
  pm.collectionVariables.set('refreshToken', jsonData.data.refresh_token);
}
```

**自动化指令**（在会话中触发）：
> "更新 Postman 集合" / "同步 Postman" / "接口变了更新 Postman"

当收到上述指令时，执行以下步骤：
1. 扫描 `internal/app/router.go` 获取最新路由列表
2. 扫描 `internal/handler/*_handler.go` 获取请求/响应 DTO
3. 对比现有 Postman JSON，找出差异
4. 更新 JSON 文件，保持格式一致
5. 提交时一并提交 JSON 变更

---

## 七、数据库规范

### 7.1 命名

```sql
-- 表名：snake_case，复数
users, habits, chat_messages, character_configs

-- 字段名：snake_case
id, user_id, created_at, updated_at, is_active

-- 主键：UUID
id UUID PRIMARY KEY DEFAULT gen_random_uuid()

-- 外键：{table}_id
user_id UUID REFERENCES users(id)

-- 时间戳：created_at, updated_at
created_at TIMESTAMP DEFAULT NOW()
updated_at TIMESTAMP DEFAULT NOW()

-- 布尔值：is_{adjective}
is_active, is_premium, is_complete
```

### 7.2 索引规范

```sql
-- 外键自动创建索引
-- 查询频繁的字段加索引
CREATE INDEX idx_chat_messages_user_created ON chat_messages(user_id, created_at DESC);

-- 复合索引：查询条件顺序
CREATE INDEX idx_habits_user_active ON habits(user_id, is_active);

-- JSONB 查询用 GIN 索引
CREATE INDEX idx_character_distillate ON celebrity_characters USING GIN(distillate);
```

### 7.3 迁移规范

```bash
# 创建新迁移
migrate create -ext sql -dir migrations -seq add_focus_sessions

# 生成 002_add_focus_sessions.up.sql / 002_add_focus_sessions.down.sql

# 执行迁移
make migrate-up

# 回滚
make migrate-down
```

---

## 八、AI/Prompt 开发规范

### 8.1 System Prompt 结构

```
【基础层】Lagom 通用指令（所有角色共用）
├── 你是 Lagom 的 AI 伙伴
├── 回复不超过 100 字
├── 温暖、共情、不批评
└── 适当使用 emoji

【角色层】名人蒸馏指令（按角色注入）
├── 风格指令（句式、词汇、语气）
├── 价值观指令（自律观、失败观、时间观）
├── 认知模式指令（分析框架、决策风格）
└── 知识库注入（个人故事、比喻库）

【上下文层】动态组装
├── 用户画像
├── 长期记忆摘要
└── 最近 10 轮对话
```

### 8.2 Prompt 调优流程

```
1. 准备 10 个测试场景（完成/失败/情绪/求助/闲聊）
2. 为每个场景写"理想回复"
3. 调整 System Prompt，使 AI 输出接近理想回复
4. 盲测：20 人对比"通用版"vs"名人版"
5. 胜率 > 70% 才算合格
6. 记录每次变更和效果，版本化 Prompt
```

### 8.3 成本优化

```go
// 意图缓存：高频意图用模板回复
var intentTemplates = map[string][]string{
    "morning_greeting": {
        "早啊！今天想从哪件小事开始？",
        "早上好～准备好今天的挑战了吗？",
    },
    "complete_encouragement": {
        "太棒了！又完成了一件小事～",
        "做得很好！给自己一点掌声👏",
    },
}

// 响应缓存：相同上下文相同问题的回复缓存 5 分钟
// 每日调用上限：免费用户 50 次，付费用户无限制
```

---

## 九、安全规范

### 9.1 认证

```
- JWT: RS256 签名，Access Token 15 分钟过期
- Refresh Token: 7 天过期，Redis 黑名单
- 设备绑定: Refresh Token 绑定设备指纹
- OAuth2: Apple Sign In + Google Sign In，不存储密码
```

### 9.2 数据安全

```
- 传输: 全站 HTTPS，TLS 1.3
- 存储: 云端 AES-256-GCM，本地 iOS Keychain / Android Keystore
- 日志: 脱敏用户 ID 和敏感内容
- 删除: 用户注销 30 天内物理删除
```

### 9.3 API 安全

```
- 限流: 单 IP 100 req/min，单用户 60 req/min
- AI 限流: 免费用户 50 次/天
- 输入验证: 严格参数校验，防 SQL 注入/XSS
- CORS: 仅允许 App 域名
```

---

## 十、测试规范

### 10.1 后端测试

```go
// 单元测试：核心业务逻辑
func TestHSICalculator(t *testing.T) { ... }
func TestDailyRecommender(t *testing.T) { ... }
func TestCharacterRecommendation(t *testing.T) { ... }

// 集成测试：API 端到端
func TestChatFlow(t *testing.T) { ... }
func TestHabitCRUD(t *testing.T) { ... }

// 目标覆盖率：> 70%
```

### 10.2 前端测试

```dart
// Widget 测试：UI 组件
// 集成测试：页面流程
// 目标：核心流程有测试覆盖
```

### 10.3 AI 质量测试

```
盲测流程：
1. 准备 10 个场景 × 12 个角色 = 120 组测试用例
2. 每组：名人版 vs 通用版，用户不知道哪个是哪个
3. 用户选择"更像该名人"的回复
4. 目标胜率：> 70%
5. 未达标角色：继续调优 Prompt
```

---

## 十一、Git 提交规范

```
<type>(<scope>): <subject>

<body>

<footer>
```

**type**：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具

**scope**：`frontend`, `backend`, `ai`, `character`, `habit`, `focus`

**示例**：
```
feat(character): 添加村上春树蒸馏档案

- 提取词汇指纹和句式模式
- 添加 3 个典型场景回复示例
- 盲测胜率 78%

Closes #12
```

---

## 十二、开发检查清单

### 每次提交前

```
□ 代码格式化（gofmt / dart format）
□ 测试通过（go test / flutter test）
□ 无编译错误
□ 敏感信息未提交（密钥、密码）
□ 新增 API 有文档注释
□ 新增字段有数据库迁移
```

### 每次 PR 前

```
□ 功能自测通过
□ 代码 Review（Self-Review）
□ 无 TODO/FIXME 遗留（或已记录 Issue）
□ 性能无明显退化
□ 移动端 UI 测试（不同尺寸）
```

### 发布前

```
□ 所有 P0 功能完成
□ 崩溃率 < 1%
□ AI 延迟 < 5s（90分位）
□ 盲测胜率 > 70%
□ 隐私政策更新
□ 应用商店资料准备
```

---

## 十三、关键业务规则速查

### 13.1 微习惯引擎

```
- 每天只推荐"一件小事"
- 完成即成功，不追求"全部完成"
- 缓冲日：每周 2 天，不完成不扣分
- HSI = 频率×0.4 + 质量×0.3 + 情绪×0.3
- 失败不展示"连续失败 X 天"
```

### 13.2 名人系统

```
- 免费用户：初始锁定 1 个角色 7 天，每月限切换 2 次
- 付费用户：同时拥有 3 个常驻角色，无限切换
- 亲密度：对话 +1，完成目标 +3，反思 +2
- 解锁：Lv5 个人故事，Lv15 专属昵称，Lv31 深夜模式
- 版权：风格而非内容，不直接引用原文，标注"AI 风格化生成"
```

### 13.3 AI 对话

```
- 回复不超过 100 字
- 失败时共情不批评
- 根据情绪状态调整语气
- 不展示"连续失败"负面数据
- 高频意图用模板缓存
```

---

## 十四、参考文档

项目全部设计文档位于 `/docs/`：

| 文档 | 用途 |
|------|------|
| `Lagom-Product-Design.md` | 功能规格（单一事实来源） |
| `Lagom-Technical-Architecture.md` | 技术架构详细设计 |
| `Lagom-UI-UX-Design.md` | UI/UX 设计规范 |
| `Lagom-Celebrity-Characters.md` | 12 个首发角色档案 |
| `Lagom-Celebrity-Integration.md` | 名人系统集成方案 |
| `Lagom-Implementation-Guide.md` | 代码级实现指南 |
| `Lagom-Development-Plan.md` | 开发路线图 |
