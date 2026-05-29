# Lagom 技术实现方案指南

> 从代码层面指导开发实施，包含关键模块的伪代码、数据结构、算法实现  
> 面向：实际编写代码的开发者

---

## 一、项目初始化命令

### 1.1 后端初始化

```bash
# 创建项目目录
mkdir -p /Users/fidoo/Desktop/liumaw99/gf/server
cd /Users/fidoo/Desktop/liumaw99/gf/server

# 初始化 Go 模块
go mod init github.com/lagom/lagom-server

# 创建目录结构
mkdir -p cmd/server config internal/{app,handler,service,repository,middleware,pkg/{jwt,encrypt,ai,push},cron} migrations

# 安装核心依赖
go get github.com/gin-gonic/gin@v1.10.0
go get github.com/jackc/pgx/v5@v5.7.1
go get github.com/golang-jwt/jwt/v5@v5.2.1
go get github.com/spf13/viper@v1.19.0
go get github.com/redis/go-redis/v9@v9.6.1
go get github.com/robfig/cron/v3@v3.0.1
go get go.uber.org/zap@v1.27.0

# 安装迁移工具
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 创建 Makefile
cat > Makefile << 'EOF'
.PHONY: build run test migrate-up migrate-down docker-build

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

migrate-up:
	migrate -path migrations -database "postgresql://lagom:lagom@localhost:5432/lagom?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://lagom:lagom@localhost:5432/lagom?sslmode=disable" down

docker-build:
	docker build -t lagom-server:latest .

.DEFAULT_GOAL := run
EOF
```

### 1.2 前端初始化

```bash
# 创建 Flutter 项目
flutter create --org com.lagom --project-name lagom_app /Users/fidoo/Desktop/liumaw99/gf/app/lagom_app

cd /Users/fidoo/Desktop/liumaw99/gf/app/lagom_app

# 添加依赖
cat >> pubspec.yaml << 'EOF'
dependencies:
  flutter_riverpod: ^2.5.1
  dio: ^5.7.0
  drift: ^2.21.0
  sqlite3_flutter_libs: ^0.5.24
  shared_preferences: ^2.3.2
  flutter_secure_storage: ^9.2.2
  flutter_local_notifications: ^17.2.3
  workmanager: ^0.5.2
  fl_chart: ^0.68.0
  lottie: ^3.1.2
  freezed_annotation: ^2.4.4
  json_annotation: ^4.9.0
  go_router: ^14.2.0
  uuid: ^4.5.1
  intl: ^0.19.0
  logger: ^2.4.0

dependency_overrides:

dev_dependencies:
  build_runner: ^2.4.13
  freezed: ^2.5.7
  json_serializable: ^6.8.0
  riverpod_generator: ^2.4.3
  drift_dev: ^2.21.0
EOF

flutter pub get

# 生成代码
flutter pub run build_runner build --delete-conflicting-outputs
```

---

## 二、后端核心实现

### 2.1 主入口

```go
// cmd/server/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/lagom/lagom-server/internal/app"
    "github.com/lagom/lagom-server/config"
)

func main() {
    // 加载配置
    cfg, err := config.Load("config/config.yaml")
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // 创建应用
    application, err := app.New(cfg)
    if err != nil {
        log.Fatalf("failed to create app: %v", err)
    }

    // 启动服务
    go func() {
        if err := application.Run(); err != nil {
            log.Printf("server error: %v", err)
        }
    }()

    // 优雅关闭
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := application.Shutdown(ctx); err != nil {
        log.Printf("shutdown error: %v", err)
    }
}
```

### 2.2 应用初始化

```go
// internal/app/app.go
package app

import (
    "context"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/lagom/lagom-server/config"
    "github.com/lagom/lagom-server/internal/handler"
    "github.com/lagom/lagom-server/internal/middleware"
    "github.com/lagom/lagom-server/internal/service"
    "github.com/lagom/lagom-server/internal/repository"
    "github.com/redis/go-redis/v9"
    "github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
    router *gin.Engine
    server *http.Server
    db     *pgxpool.Pool
    redis  *redis.Client
}

func New(cfg *config.Config) (*App, error) {
    // 数据库连接
    db, err := pgxpool.New(context.Background(), cfg.Database.DSN)
    if err != nil {
        return nil, fmt.Errorf("database connection failed: %w", err)
    }

    // Redis 连接
    rdb := redis.NewClient(&redis.Options{
        Addr: cfg.Redis.Addr,
    })

    // 仓库层
    userRepo := repository.NewUserRepo(db)
    habitRepo := repository.NewHabitRepo(db)
    chatRepo := repository.NewChatRepo(db)
    characterRepo := repository.NewCharacterRepo(db, rdb)
    userCharRepo := repository.NewUserCharacterRepo(db)

    // 服务层
    authService := service.NewAuthService(cfg, userRepo)
    aiService := service.NewAIService(cfg, chatRepo, characterRepo)
    habitService := service.NewHabitService(habitRepo)
    characterService := service.NewCharacterService(characterRepo, userCharRepo, userRepo)

    // 处理器
    authHandler := handler.NewAuthHandler(authService)
    chatHandler := handler.NewChatHandler(aiService)
    habitHandler := handler.NewHabitHandler(habitService)
    characterHandler := handler.NewCharacterHandler(characterService)

    // 路由
    router := gin.New()
    router.Use(gin.Recovery())
    router.Use(middleware.Logger())
    router.Use(middleware.CORS())

    // 健康检查
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // API 路由组
    api := router.Group("/api/v1")
    {
        // 公开路由
        auth := api.Group("/auth")
        {
            auth.POST("/apple", authHandler.AppleSignIn)
            auth.POST("/google", authHandler.GoogleSignIn)
            auth.POST("/refresh", authHandler.RefreshToken)
        }

        // 需要认证的路由
        authorized := api.Group("")
        authorized.Use(middleware.JWTAuth(cfg.JWT.Secret))
        {
            // 聊天
            authorized.POST("/chat", chatHandler.Chat)
            authorized.GET("/chat/history", chatHandler.GetHistory)

            // 习惯
            authorized.GET("/habits", habitHandler.List)
            authorized.POST("/habits", habitHandler.Create)
            authorized.POST("/habits/:id/complete", habitHandler.Complete)

            // 名人角色
            authorized.GET("/characters", characterHandler.List)
            authorized.GET("/characters/:id", characterHandler.GetDetail)
            authorized.POST("/characters/recommend", characterHandler.Recommend)
            authorized.GET("/user/characters", characterHandler.GetUserCharacters)
            authorized.POST("/user/characters", characterHandler.SelectCharacter)
            authorized.POST("/user/characters/:id/switch", characterHandler.SwitchCharacter)
        }
    }

    return &App{
        router: router,
        server: &http.Server{
            Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
            Handler: router,
        },
        db:    db,
        redis: rdb,
    }, nil
}

func (a *App) Run() error {
    return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
    a.db.Close()
    a.redis.Close()
    return a.server.Shutdown(ctx)
}
```

### 2.3 AI 服务核心实现

```go
// internal/service/ai_service.go
package service

import (
    "bufio"
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// AIClient 接口
type AIClient interface {
    StreamChat(ctx context.Context, messages []Message) (<-chan string, error)
}

// DeepSeekClient 实现
type DeepSeekClient struct {
    apiKey string
    baseURL string
    client  *http.Client
}

func NewDeepSeekClient(apiKey string) *DeepSeekClient {
    return &DeepSeekClient{
        apiKey:  apiKey,
        baseURL: "https://api.deepseek.com/v1",
        client:  &http.Client{Timeout: 60 * time.Second},
    }
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    Stream      bool      `json:"stream"`
    Temperature float64   `json:"temperature"`
    MaxTokens   int       `json:"max_tokens"`
}

type ChatResponse struct {
    Choices []struct {
        Delta struct {
            Content string `json:"content"`
        } `json:"delta"`
    } `json:"choices"`
}

func (c *DeepSeekClient) StreamChat(ctx context.Context, messages []Message) (<-chan string, error) {
    reqBody := ChatRequest{
        Model:       "deepseek-chat",
        Messages:    messages,
        Stream:      true,
        Temperature: 0.7,
        MaxTokens:   300,
    }

    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", 
        c.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        return nil, fmt.Errorf("API error: %s", string(body))
    }

    ch := make(chan string, 10)

    go func() {
        defer close(ch)
        defer resp.Body.Close()

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            line := scanner.Text()
            if !strings.HasPrefix(line, "data: ") {
                continue
            }

            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" {
                return
            }

            var chatResp ChatResponse
            if err := json.Unmarshal([]byte(data), &chatResp); err != nil {
                continue
            }

            if len(chatResp.Choices) > 0 {
                content := chatResp.Choices[0].Delta.Content
                if content != "" {
                    select {
                    case ch <- content:
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }
    }()

    return ch, nil
}

// AIService 业务层
type AIService struct {
    primary    AIClient
    fallback   AIClient
    chatRepo   ChatRepo
    charRepo   CharacterRepo
    memoryRepo MemoryRepo
}

func NewAIService(cfg *config.Config, chatRepo ChatRepo, charRepo CharacterRepo) *AIService {
    return &AIService{
        primary:    NewDeepSeekClient(cfg.AI.DeepSeekKey),
        fallback:   NewDeepSeekClient(cfg.AI.FallbackKey), // 或 GPT-4o-mini
        chatRepo:   chatRepo,
        charRepo:   charRepo,
        memoryRepo: NewMemoryRepo(), // 初始化
    }
}

func (s *AIService) Chat(ctx context.Context, userID string, characterID *string, userMessage string) (<-chan string, error) {
    // 构建 system prompt
    systemPrompt := s.buildSystemPrompt(ctx, userID, characterID)

    // 获取历史
    history := s.chatRepo.GetRecent(ctx, userID, characterID, 10)

    // 组装 messages
    messages := []Message{{Role: "system", Content: systemPrompt}}
    for _, h := range history {
        messages = append(messages, Message{Role: h.Role, Content: h.Content})
    }
    messages = append(messages, Message{Role: "user", Content: userMessage})

    // 调用 AI
    stream, err := s.primary.StreamChat(ctx, messages)
    if err != nil {
        stream, err = s.fallback.StreamChat(ctx, messages)
        if err != nil {
            return nil, err
        }
    }

    // 异步保存
    go s.saveResponse(ctx, userID, characterID, userMessage, stream)

    return stream, nil
}

func (s *AIService) buildSystemPrompt(ctx context.Context, userID string, characterID *string) string {
    var prompt strings.Builder

    // 基础 prompt
    prompt.WriteString("你是 Lagom 的 AI 伙伴，帮助用户养成自律习惯。\n")
    prompt.WriteString("回复要温暖、简短（不超过100字），善用emoji增加温度。\n")
    prompt.WriteString("失败时给予共情和鼓励，绝不批评。\n\n")

    // 名人角色注入
    if characterID != nil && *characterID != "" {
        character := s.charRepo.GetByID(ctx, *characterID)
        if character != nil {
            prompt.WriteString(fmt.Sprintf("你当前扮演的是基于%s风格蒸馏的 AI 伙伴。\n", character.Name))
            prompt.WriteString("风格指令：\n")
            for _, directive := range character.Distillate.Style.Directives {
                prompt.WriteString("- " + directive + "\n")
            }
            prompt.WriteString("价值观指令：\n")
            for _, directive := range character.Distillate.Values.Directives {
                prompt.WriteString("- " + directive + "\n")
            }
        }
    }

    return prompt.String()
}

func (s *AIService) saveResponse(ctx context.Context, userID string, characterID *string, userMsg string, stream <-chan string) {
    var fullResponse strings.Builder
    for delta := range stream {
        fullResponse.WriteString(delta)
    }

    // 保存用户消息
    s.chatRepo.Save(ctx, userID, characterID, "user", userMsg)
    // 保存 AI 回复
    s.chatRepo.Save(ctx, userID, characterID, "assistant", fullResponse.String())
    // 更新记忆
    s.memoryRepo.Update(ctx, userID, userMsg, fullResponse.String())
}
```

### 2.4 聊天 Handler（SSE）

```go
// internal/handler/chat_handler.go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/lagom/lagom-server/internal/service"
)

type ChatHandler struct {
    aiService *service.AIService
}

func NewChatHandler(aiService *service.AIService) *ChatHandler {
    return &ChatHandler{aiService: aiService}
}

type ChatRequest struct {
    Message     string  `json:"message" binding:"required"`
    CharacterID *string `json:"character_id"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
    userID := c.GetString("user_id")
    
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // SSE 设置
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("X-Accel-Buffering", "no")

    stream, err := h.aiService.Chat(c.Request.Context(), userID, req.CharacterID, req.Message)
    if err != nil {
        c.SSEvent("error", err.Error())
        return
    }

    // 流式输出
    c.Stream(func(w io.Writer) bool {
        select {
        case delta, ok := <-stream:
            if !ok {
                c.SSEvent("done", "")
                return false
            }
            c.SSEvent("delta", delta)
            return true
        case <-c.Request.Context().Done():
            return false
        }
    })
}
```

---

## 三、前端核心实现

### 3.1 Riverpod 架构

```dart
// lib/providers/auth_provider.dart

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lagom_app/core/api/api_client.dart';

// API Client Provider
final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient(baseUrl: 'https://api.lagom.app');
});

// 认证状态
@Riverpod(keepAlive: true)
class AuthNotifier extends _$AuthNotifier {
  @override
  AuthState build() {
    // 启动时检查本地 token
    _checkAuth();
    return const AuthState.loading();
  }

  Future<void> _checkAuth() async {
    final token = await SecureStorage.read('access_token');
    if (token != null) {
      state = const AuthState.authenticated();
    } else {
      state = const AuthState.unauthenticated();
    }
  }

  Future<void> signInWithApple() async {
    // Apple Sign In 逻辑
    final credential = await SignInWithApple.getAppleIDCredential(...);
    final result = await ref.read(apiClientProvider).post('/auth/apple', 
      data: {'id_token': credential.identityToken});
    
    await SecureStorage.write('access_token', result['access_token']);
    await SecureStorage.write('refresh_token', result['refresh_token']);
    
    state = const AuthState.authenticated();
  }

  Future<void> signOut() async {
    await SecureStorage.delete('access_token');
    await SecureStorage.delete('refresh_token');
    state = const AuthState.unauthenticated();
  }
}
```

### 3.2 聊天核心逻辑

```dart
// lib/features/companion/companion_notifier.dart

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lagom_app/core/api/api_client.dart';

// 消息模型
@freezed
class ChatMessage with _$ChatMessage {
  const factory ChatMessage({
    required String id,
    required String role, // 'user' | 'assistant'
    required String content,
    DateTime? createdAt,
    bool? isComplete,
  }) = _ChatMessage;
}

// 伙伴状态
@freezed
class CompanionState with _$CompanionState {
  const factory CompanionState({
    @Default([]) List<ChatMessage> messages,
    @Default(false) bool isTyping,
    String? currentCharacterId,
    String? companionName,
    String? avatarUrl,
    @Default(1) int intimacyLevel,
  }) = _CompanionState;
}

@Riverpod(keepAlive: true)
class CompanionNotifier extends _$CompanionNotifier {
  late ApiClient _api;

  @override
  CompanionState build() {
    _api = ref.read(apiClientProvider);
    _loadHistory();
    return const CompanionState();
  }

  Future<void> _loadHistory() async {
    try {
      final response = await _api.get('/chat/history');
      final messages = (response.data as List)
          .map((m) => ChatMessage(
                id: m['id'],
                role: m['role'],
                content: m['content'],
                createdAt: DateTime.parse(m['created_at']),
                isComplete: true,
              ))
          .toList();
      state = state.copyWith(messages: messages);
    } catch (e) {
      // 首次使用，无历史记录
    }
  }

  Future<void> sendMessage(String content) async {
    // 1. 添加用户消息到列表
    final userMsg = ChatMessage(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      role: 'user',
      content: content,
      createdAt: DateTime.now(),
      isComplete: true,
    );
    state = state.copyWith(
      messages: [...state.messages, userMsg],
      isTyping: true,
    );

    // 2. 创建 AI 消息占位
    final aiMsgId = '${DateTime.now().millisecondsSinceEpoch}_ai';
    final aiMsg = ChatMessage(
      id: aiMsgId,
      role: 'assistant',
      content: '',
      isComplete: false,
    );
    state = state.copyWith(messages: [...state.messages, aiMsg]);

    // 3. 发起 SSE 请求
    try {
      final stream = await _api.postStream('/chat', data: {
        'message': content,
        'character_id': state.currentCharacterId,
      });

      final buffer = StringBuffer();
      await for (final event in stream) {
        if (event['type'] == 'delta') {
          buffer.write(event['content']);
          _updateAiMessage(aiMsgId, buffer.toString(), false);
        } else if (event['type'] == 'done') {
          _updateAiMessage(aiMsgId, buffer.toString(), true);
        }
      }
    } catch (e) {
      _updateAiMessage(aiMsgId, '抱歉，我遇到了一点问题，稍后再试好吗？', true);
    } finally {
      state = state.copyWith(isTyping: false);
    }
  }

  void _updateAiMessage(String id, String content, bool isComplete) {
    final updatedMessages = state.messages.map((m) {
      if (m.id == id) {
        return m.copyWith(content: content, isComplete: isComplete);
      }
      return m;
    }).toList();
    state = state.copyWith(messages: updatedMessages);
  }

  Future<void> switchCharacter(String? characterId) async {
    state = state.copyWith(currentCharacterId: characterId);
    // 清空对话历史，重新加载
    await _loadHistory();
  }
}
```

### 3.3 API Client（SSE 支持）

```dart
// lib/core/api/api_client.dart

import 'package:dio/dio.dart';

class ApiClient {
  final Dio _dio;

  ApiClient({required String baseUrl}) : _dio = Dio(BaseOptions(
    baseUrl: baseUrl,
    connectTimeout: const Duration(seconds: 30),
    receiveTimeout: const Duration(seconds: 60),
  )) {
    // 添加拦截器
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        // 添加 JWT Token
        final token = await SecureStorage.read('access_token');
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
      onError: (error, handler) async {
        // Token 过期，尝试刷新
        if (error.response?.statusCode == 401) {
          final refreshed = await _refreshToken();
          if (refreshed) {
            error.requestOptions.headers['Authorization'] = 
              'Bearer ${await SecureStorage.read('access_token')}';
            final response = await _dio.fetch(error.requestOptions);
            handler.resolve(response);
            return;
          }
        }
        handler.next(error);
      },
    ));
  }

  Future<Response> get(String path, {Map<String, dynamic>? query}) async {
    return _dio.get(path, queryParameters: query);
  }

  Future<Response> post(String path, {dynamic data}) async {
    return _dio.post(path, data: data);
  }

  // SSE 流式请求
  Stream<Map<String, dynamic>> postStream(String path, {dynamic data}) async* {
    final response = await _dio.post(
      path,
      data: data,
      options: Options(
        responseType: ResponseType.stream,
        headers: {'Accept': 'text/event-stream'},
      ),
    );

    final stream = response.data as Stream<List<int>>;
    final transformer = Utf8Decoder().bind(stream);
    
    await for (final line in transformer.transform(const LineSplitter())) {
      if (line.startsWith('data: ')) {
        final jsonStr = line.substring(6);
        if (jsonStr == '[DONE]') break;
        yield jsonDecode(jsonStr) as Map<String, dynamic>;
      }
    }
  }

  Future<bool> _refreshToken() async {
    try {
      final refreshToken = await SecureStorage.read('refresh_token');
      final response = await _dio.post('/auth/refresh', data: {
        'refresh_token': refreshToken,
      });
      await SecureStorage.write('access_token', response.data['access_token']);
      return true;
    } catch (e) {
      return false;
    }
  }
}
```

---

## 四、本地数据库（Drift）

```dart
// lib/core/database/app_database.dart

import 'package:drift/drift.dart';
import 'package:drift_flutter/drift_flutter.dart';

part 'app_database.g.dart';

// 用户表
class Users extends Table {
  TextColumn get id => text()();
  TextColumn get name => text().nullable()();
  TextColumn get email => text().nullable()();
  TextColumn get subscriptionTier => text().withDefault(const Constant('free'))();
  DateTimeColumn get createdAt => dateTime().withDefault(currentDateAndTime)();
  
  @override
  Set<Column> get primaryKey => {id};
}

// 习惯表
class Habits extends Table {
  TextColumn get id => text()();
  TextColumn get userId => text()();
  TextColumn get title => text()();
  TextColumn get category => text()();
  TextColumn get description => text().nullable()();
  IntColumn get priority => integer().withDefault(const Constant(3))();
  BoolColumn get isActive => boolean().withDefault(const Constant(true))();
  RealColumn get hsi => real().withDefault(const Constant(0.0))();
  DateTimeColumn get createdAt => dateTime().withDefault(currentDateAndTime)();
  
  @override
  Set<Column> get primaryKey => {id};
}

// 每日记录表
class DailyLogs extends Table {
  TextColumn get id => text()();
  TextColumn get userId => text()();
  DateTimeColumn get date => dateTime()();
  IntColumn get moodScore => integer().nullable()();
  TextColumn get reflectionNote => text().nullable()();
  
  @override
  Set<Column> get primaryKey => {id};
}

// 对话消息表
class ChatMessages extends Table {
  TextColumn get id => text()();
  TextColumn get userId => text()();
  TextColumn get characterId => text().nullable()();
  TextColumn get role => text()();
  TextColumn get content => text()();
  BoolColumn get isComplete => boolean().withDefault(const Constant(true))();
  DateTimeColumn get createdAt => dateTime().withDefault(currentDateAndTime)();
  
  @override
  Set<Column> get primaryKey => {id};
}

// 本地角色配置表
class LocalCharacterConfigs extends Table {
  TextColumn get id => text()();
  TextColumn get userId => text()();
  TextColumn get characterId => text()();
  TextColumn get customName => text().nullable()();
  IntColumn get intimacyLevel => integer().withDefault(const Constant(1))();
  IntColumn get intimacyScore => integer().withDefault(const Constant(0))();
  BoolColumn get isActive => boolean().withDefault(const Constant(false))();
  DateTimeColumn get lastSyncAt => dateTime().nullable()();
  
  @override
  Set<Column> get primaryKey => {id};
}

@DriftDatabase(tables: [Users, Habits, DailyLogs, ChatMessages, LocalCharacterConfigs])
class AppDatabase extends _$AppDatabase {
  AppDatabase() : super(_openConnection());

  @override
  int get schemaVersion => 1;

  static QueryExecutor _openConnection() {
    return driftDatabase(name: 'lagom_database');
  }

  // 习惯相关查询
  Future<List<Habit>> getAllHabits() => select(habits).get();
  Future<void> insertHabit(HabitsCompanion habit) => into(habits).insert(habit);
  Future<void> updateHsi(String habitId, double newHsi) {
    return update(habits).replace(
      HabitsCompanion(id: Value(habitId), hsi: Value(newHsi)),
    );
  }

  // 对话相关查询
  Future<List<ChatMessage>> getRecentMessages(String userId, {int limit = 50}) {
    return (select(chatMessages)
      ..where((m) => m.userId.equals(userId))
      ..orderBy([(m) => OrderingTerm.desc(m.createdAt)])
      ..limit(limit))
      .get();
  }

  Future<void> insertMessage(ChatMessagesCompanion msg) => into(chatMessages).insert(msg);

  // 角色配置相关
  Future<List<LocalCharacterConfig>> getCharacterConfigs(String userId) {
    return (select(localCharacterConfigs)
      ..where((c) => c.userId.equals(userId)))
      .get();
  }

  Future<void> upsertCharacterConfig(LocalCharacterConfigsCompanion config) {
    return into(localCharacterConfigs).insertOnConflictUpdate(config);
  }
}
```

---

## 五、名人蒸馏引擎（Python PoC）

```python
# tools/distillation/distill_character.py
"""
名人蒸馏引擎原型
输入：原始语料文件（txt/md/json）
输出：CharacterDistillate JSON
"""

import json
import re
from dataclasses import dataclass, asdict
from typing import List, Dict
from collections import Counter

@dataclass
class StyleProfile:
    sentence_patterns: Dict[str, List[str]]
    vocabulary_fingerprint: Dict
    tonality_profile: Dict

@dataclass
class ValuesProfile:
    discipline_view: Dict
    failure_attitude: Dict
    time_philosophy: Dict
    motivation_style: Dict

@dataclass
class CharacterDistillate:
    version: str
    basic_info: Dict
    style_profile: StyleProfile
    values_profile: ValuesProfile

class DistillationEngine:
    def __init__(self):
        self.corpus = ""
        self.segments = []
    
    def load_corpus(self, filepath: str):
        """加载原始语料"""
        with open(filepath, 'r', encoding='utf-8') as f:
            self.corpus = f.read()
        self._segment()
    
    def _segment(self):
        """将语料分段（按段落/章节）"""
        # 简单按空行分段
        self.segments = [s.strip() for s in self.corpus.split('\n\n') if s.strip()]
    
    def extract_vocabulary_fingerprint(self) -> Dict:
        """提取词汇指纹"""
        # 分词（简单按字/词分，实际应用 jieba）
        words = re.findall(r'[\u4e00-\u9fa5]{2,4}', self.corpus)
        
        # 统计词频
        word_freq = Counter(words)
        
        # 提取特色词汇（排除通用词）
        common_words = {'我们', '他们', '可以', '这个', '一个', '没有', '什么'}
        signature_words = [
            {"word": w, "frequency": f}
            for w, f in word_freq.most_common(100)
            if w not in common_words
        ][:20]
        
        # 提取常用短语
        phrases = re.findall(r'[\u4e00-\u9fa5]{2,6}[，。！？]', self.corpus)
        phrase_freq = Counter(phrases)
        signature_phrases = [p for p, _ in phrase_freq.most_common(10)]
        
        return {
            "signature_words": signature_words,
            "signature_phrases": signature_phrases,
        }
    
    def extract_sentence_patterns(self) -> Dict[str, List[str]]:
        """提取句式模式"""
        patterns = {
            "openings": [],
            "transitions": [],
            "emphasis": [],
            "closings": [],
        }
        
        # 提取开头模式
        opening_patterns = re.findall(r'(^[\u4e00-\u9fa5]{2,8}[，。])', self.corpus, re.MULTILINE)
        patterns["openings"] = list(set(opening_patterns))[:10]
        
        # 提取过渡词
        transitions = re.findall(r'(话说回来|不过|也就是说|换句话说|总而言之)[，。]', self.corpus)
        patterns["transitions"] = list(set(transitions))
        
        return patterns
    
    def extract_values(self) -> Dict:
        """提取价值观（基于关键词匹配 + 上下文）"""
        values = {
            "discipline_view": {},
            "failure_attitude": {},
            "time_philosophy": {},
            "motivation_style": {},
        }
        
        # 查找与"自律/习惯/坚持"相关的段落
        discipline_keywords = ['自律', '习惯', '坚持', '每天', '重复', '节奏']
        discipline_segments = self._find_relevant_segments(discipline_keywords)
        
        if discipline_segments:
            values["discipline_view"] = {
                "key_quotes": discipline_segments[:3],
                "approach": self._classify_approach(discipline_segments)
            }
        
        # 查找与"失败/困难"相关的段落
        failure_keywords = ['失败', '困难', '挫折', '放弃', '不顺']
        failure_segments = self._find_relevant_segments(failure_keywords)
        
        if failure_segments:
            values["failure_attitude"] = {
                "key_quotes": failure_segments[:3],
                "approach": self._classify_failure_approach(failure_segments)
            }
        
        return values
    
    def _find_relevant_segments(self, keywords: List[str]) -> List[str]:
        """查找包含关键词的段落"""
        results = []
        for seg in self.segments:
            if any(kw in seg for kw in keywords):
                results.append(seg)
        return results
    
    def _classify_approach(self, segments: List[str]) -> str:
        """分类自律方式"""
        text = " ".join(segments)
        if "节奏" in text or "自然" in text:
            return "gentle_persistence"
        elif "强迫" in text or "必须" in text:
            return "aggressive_discipline"
        else:
            return "balanced"
    
    def _classify_failure_approach(self, segments: List[str]) -> str:
        """分类失败应对方式"""
        text = " ".join(segments)
        if "接受" in text or "没关系" in text:
            return "acceptance"
        elif "重新" in text or "再来" in text:
            return "reframing"
        else:
            return "neutral"
    
    def distill(self, basic_info: Dict) -> CharacterDistillate:
        """执行完整蒸馏流程"""
        vocab = self.extract_vocabulary_fingerprint()
        patterns = self.extract_sentence_patterns()
        values = self.extract_values()
        
        style = StyleProfile(
            sentence_patterns=patterns,
            vocabulary_fingerprint=vocab,
            tonality_profile={"overall_tone": "calm_reflective"}  # 简化版
        )
        
        values_profile = ValuesProfile(
            discipline_view=values.get("discipline_view", {}),
            failure_attitude=values.get("failure_attitude", {}),
            time_philosophy={},
            motivation_style={},
        )
        
        return CharacterDistillate(
            version="1.0",
            basic_info=basic_info,
            style_profile=style,
            values_profile=values_profile
        )

# 使用示例
if __name__ == "__main__":
    engine = DistillationEngine()
    engine.load_corpus("corpus/haruki_murakami_running.txt")
    
    distillate = engine.distill({
        "name": "村上春树",
        "name_en": "Haruki Murakami",
        "category": "writer",
        "famous_for": ["novelist", "runner", "jazz_lover"]
    })
    
    with open("distillates/haruki_murakami_v1.json", "w", encoding="utf-8") as f:
        json.dump(asdict(distillate), f, ensure_ascii=False, indent=2)
    
    print("Distillation complete!")
```

---

## 六、本地通知实现

```dart
// lib/core/notifications/notification_service.dart

import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:timezone/timezone.dart' as tz;

class NotificationService {
  static final _notifications = FlutterLocalNotificationsPlugin();

  static Future<void> init() async {
    const androidSettings = AndroidInitializationSettings('@mipmap/ic_launcher');
    const iosSettings = DarwinInitializationSettings(
      requestAlertPermission: true,
      requestBadgePermission: true,
      requestSoundPermission: true,
    );
    
    await _notifications.initialize(
      const InitializationSettings(android: androidSettings, iOS: iosSettings),
    );
  }

  // 显示即时通知
  static Future<void> show({
    required int id,
    required String title,
    required String body,
    String? payload,
  }) async {
    const androidDetails = AndroidNotificationDetails(
      'lagom_channel',
      'Lagom 通知',
      channelDescription: 'Lagom 自律伙伴的日常提醒',
      importance: Importance.high,
      priority: Priority.high,
    );
    const iosDetails = DarwinNotificationDetails();
    
    await _notifications.show(
      id,
      title,
      body,
      const NotificationDetails(android: androidDetails, iOS: iosDetails),
      payload: payload,
    );
  }

  // 定时通知（每日一件小事）
  static Future<void> scheduleDailyReminder({
    required int id,
    required String title,
    required String body,
    required Time time,
  }) async {
    await _notifications.zonedSchedule(
      id,
      title,
      body,
      _nextInstanceOfTime(time),
      const NotificationDetails(
        android: AndroidNotificationDetails(
          'daily_reminder',
          '每日提醒',
          importance: Importance.high,
        ),
        iOS: DarwinNotificationDetails(),
      ),
      androidScheduleMode: AndroidScheduleMode.exactAllowWhileIdle,
      matchDateTimeComponents: DateTimeComponents.time,
    );
  }

  static tz.TZDateTime _nextInstanceOfTime(Time time) {
    final now = tz.TZDateTime.now(tz.local);
    var scheduled = tz.TZDateTime(tz.local, now.year, now.month, now.day, 
      time.hour, time.minute);
    if (scheduled.isBefore(now)) {
      scheduled = scheduled.add(const Duration(days: 1));
    }
    return scheduled;
  }

  // 取消所有通知
  static Future<void> cancelAll() => _notifications.cancelAll();
}
```

---

## 七、关键算法实现

### 7.1 HSI（习惯强度指数）计算

```dart
// lib/core/algorithms/hsi_calculator.dart

class HSICalculator {
  /// 计算习惯强度指数 (Habit Strength Index)
  /// 
  /// 参数：
  /// - completions: 最近30天的完成记录 [bool]
  /// - qualityScores: 对应的完成质量 [1-5]
  /// - moodScores: 完成时的心情 [1-5]
  /// 
  /// 返回：0-100 的 HSI 值
  static double calculate({
    required List<bool> completions,
    required List<int> qualityScores,
    required List<int> moodScores,
  }) {
    assert(completions.length == 30);
    
    // 1. 完成频率 (40%)
    final completionRate = completions.where((c) => c).length / completions.length;
    final frequencyScore = completionRate * 100;
    
    // 2. 完成质量 (30%)
    final avgQuality = qualityScores.isEmpty 
      ? 3.0 
      : qualityScores.reduce((a, b) => a + b) / qualityScores.length;
    final qualityScore = (avgQuality / 5.0) * 100;
    
    // 3. 情绪一致性 (30%)
    // 完成时心情愉悦的比例
    final positiveMoods = moodScores.where((m) => m >= 4).length;
    final moodConsistency = moodScores.isEmpty
      ? 0.5
      : positiveMoods / moodScores.length;
    final moodScore = moodConsistency * 100;
    
    // 加权计算
    final hsi = frequencyScore * 0.4 + qualityScore * 0.3 + moodScore * 0.3;
    return hsi.clamp(0.0, 100.0);
  }

  /// 获取阶段名称
  static String getStage(double hsi) {
    if (hsi < 30) return "萌芽期";
    if (hsi < 60) return "成长期";
    if (hsi < 85) return "稳定期";
    return "内化期";
  }

  /// 获取阶段 emoji
  static String getStageEmoji(double hsi) {
    if (hsi < 30) return "🌱";
    if (hsi < 60) return "🌿";
    if (hsi < 85) return "🌳";
    return "✨";
  }
}
```

### 7.2 每日推荐算法

```dart
// lib/core/algorithms/daily_recommender.dart

class DailyRecommender {
  /// 推荐"今天的一件小事"
  static HabitMicroAction recommend({
    required List<Habit> activeHabits,
    required UserMood? currentMood,
    required DateTime now,
  }) {
    final scoredHabits = activeHabits.map((habit) {
      double score = 0;

      // 紧急度：连续未完成天数
      score += habit.missedDays * 10;

      // 优先级：用户设定
      score += habit.priority * 5;

      // 时间匹配度
      score += _timeMatchScore(habit.bestTime, now) * 3;

      // 情绪适配度
      score += _moodMatchScore(habit.category, currentMood) * 2;

      // HSI 加成：HSI 低的 habit 需要更多关注
      if (habit.hsi < 30) score += 15;

      return ScoredHabit(habit, score);
    }).toList();

    // 排序取最高
    scoredHabits.sort((a, b) => b.score.compareTo(a.score));
    final topHabit = scoredHabits.first.habit;

    return topHabit.currentMicroAction;
  }

  static double _timeMatchScore(TimeOfDay? bestTime, DateTime now) {
    if (bestTime == null) return 0.5;
    
    final currentHour = now.hour;
    final bestHour = bestTime.hour;
    final diff = (currentHour - bestHour).abs();
    
    if (diff <= 1) return 1.0;
    if (diff <= 3) return 0.7;
    return 0.3;
  }

  static double _moodMatchScore(String category, UserMood? mood) {
    if (mood == null) return 0.5;
    
    // 情绪低落时，推荐简单/治愈类 habit
    if (mood == UserMood.tired || mood == UserMood.stressed) {
      if (category == 'mindfulness' || category == 'health') return 1.0;
      if (category == 'fitness') return 0.3;
    }
    
    // 情绪高涨时，推荐挑战类 habit
    if (mood == UserMood.energetic) {
      if (category == 'fitness' || category == 'learning') return 1.0;
    }
    
    return 0.5;
  }
}
```

---

## 八、环境配置模板

### 8.1 后端配置

```yaml
# config/config.yaml

server:
  port: 8080
  mode: "release" # debug / release

database:
  dsn: "postgresql://lagom:password@localhost:5432/lagom?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-here" # 生产环境使用环境变量
  access_ttl: 15m
  refresh_ttl: 7d

ai:
  deepseek_key: "${DEEPSEEK_API_KEY}" # 从环境变量读取
  fallback_key: "${OPENAI_API_KEY}"
  max_tokens: 300
  temperature: 0.7

oauth:
  apple:
    client_id: "com.lagom.app"
    team_id: "YOUR_TEAM_ID"
    key_id: "YOUR_KEY_ID"
    private_key: "${APPLE_PRIVATE_KEY}"
  google:
    client_id: "YOUR_GOOGLE_CLIENT_ID"
    client_secret: "${GOOGLE_CLIENT_SECRET}"

push:
  fcm:
    credentials: "${FCM_CREDENTIALS}"
```

### 8.2 环境变量模板

```bash
# .env
DB_PASSWORD=your_db_password
DEEPSEEK_API_KEY=sk-xxxxxxxx
OPENAI_API_KEY=sk-xxxxxxxx
APPLE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----
..."
GOOGLE_CLIENT_SECRET=xxxxxxxx
FCM_CREDENTIALS={"type":"service_account",...}
```

### 8.3 Docker Compose

```yaml
# docker-compose.yml

version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    env_file: .env
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    volumes:
      - ./config:/app/config

  postgres:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_DB: lagom
      POSTGRES_USER: lagom
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
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

---

## 九、开发命令速查

```bash
# ===== 后端 =====
# 运行
make run

# 测试
make test

# 数据库迁移
make migrate-up
make migrate-down

# 生成新迁移
migrate create -ext sql -dir migrations -seq add_character_table

# Docker 构建
make docker-build

# ===== 前端 =====
# 运行
cd app/lagom_app && flutter run

# 构建
flutter build ios
flutter build apk

# 生成代码（Riverpod/Freezed/Drift）
flutter pub run build_runner build --delete-conflicting-outputs

#  watch 模式（开发时自动生成）
flutter pub run build_runner watch --delete-conflicting-outputs

# 测试
flutter test

# ===== 蒸馏工具 =====
# 运行蒸馏
cd tools/distillation && python distill_character.py --input corpus/murakami.txt --output distillates/

# 验证蒸馏结果
python validate_distillate.py --distillate distillates/haruki_murakami_v1.json
```
