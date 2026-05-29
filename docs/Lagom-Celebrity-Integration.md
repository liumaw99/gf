# Lagom 名人系统与现有架构集成方案

> 将"AI 蒸馏名人"功能无缝集成到 Lagom 现有架构中  
> 目标：不改变现有核心流程，通过扩展点和配置化实现名人角色切换

---

## 一、集成架构总览

```
┌─────────────────────────────────────────────────────────────────┐
│                    Lagom App（Flutter）                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐      │
│  │  原有系统    │     │  扩展层      │     │  名人系统    │      │
│  │  Base Layer │ ←──→│ Extension   │ ←──→│ Celebrity   │      │
│  │             │     │ Layer       │     │ Layer       │      │
│  └─────────────┘     └─────────────┘     └─────────────┘      │
│         │                   │                   │               │
│    • 聊天界面          • 角色选择器           • 蒸馏档案加载    │
│    • 习惯引擎          • 角色配置管理          • 风格注入       │
│    • 专注计时          • 切换逻辑              • 知识库检索     │
│    • 成长档案          • 亲密度系统            • 一致性校验     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    Lagom Server（Go）                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐      │
│  │  Base API   │     │  Extension  │     │  Celebrity  │      │
│  │             │     │ Middleware  │     │ Service     │      │
│  └─────────────┘     └─────────────┘     └─────────────┘      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**设计原则**：
1. **向后兼容**：不启用名人功能时，系统完全按原有逻辑运行
2. **配置驱动**：名人角色通过配置文件/数据库动态加载，无需改代码
3. **无缝切换**：用户可在"基础伙伴"和"名人伙伴"之间自由切换
4. **数据隔离**：每个角色的亲密度、对话历史独立存储

---

## 二、数据模型扩展

### 2.1 扩展后的数据库 Schema

```sql
-- ============================================
-- 名人角色表（新增）
-- ============================================
CREATE TABLE celebrity_characters (
    id VARCHAR(50) PRIMARY KEY,              -- haruki_murakami_v1
    name VARCHAR(100) NOT NULL,              -- 村上春树
    name_en VARCHAR(100),                    -- Haruki Murakami
    category VARCHAR(50) NOT NULL,           -- writer / entrepreneur / athlete / philosopher / fictional / scientist
    era VARCHAR(50),                         -- contemporary / historical
    nationality VARCHAR(50),
    famous_for TEXT[],                       -- ["novelist", "runner", "jazz_lover"]
    avatar_url TEXT,                         -- 角色形象图
    avatar_lottie_url TEXT,                  -- Lottie 动画（可选）
    
    -- 风格标签（用于推荐）
    style_tags TEXT[],                       -- ["calm", "reflective", "literary"]
    
    -- 档案存储（JSONB，完整的蒸馏档案）
    distillate JSONB NOT NULL,
    
    -- 质量评分
    quality_score DECIMAL(3,2) DEFAULT 0.0,
    
    -- 状态
    status VARCHAR(20) DEFAULT 'active',     -- active / inactive / beta
    
    -- 权限
    is_premium BOOLEAN DEFAULT false,        -- 是否付费专属
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_celebrity_category ON celebrity_characters(category);
CREATE INDEX idx_celebrity_status ON celebrity_characters(status);
CREATE INDEX idx_celebrity_style_tags ON celebrity_characters USING GIN(style_tags);

-- ============================================
-- 用户角色配置表（新增）
-- ============================================
CREATE TABLE user_character_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- 角色引用
    character_id VARCHAR(50) REFERENCES celebrity_characters(id),
    
    -- 用户自定义（覆盖默认）
    custom_name VARCHAR(100),                -- 用户给角色的昵称
    custom_avatar_url TEXT,                  -- 用户上传的形象
    
    -- 亲密度系统
    intimacy_level INT DEFAULT 1,            -- 1-50
    intimacy_score INT DEFAULT 0,            -- 累计分数
    
    -- 成长档案
    unlocked_stories TEXT[],                 -- 已解锁的个人故事
    unlocked_expressions TEXT[],             -- 已解锁的风格表达
    unlocked_deep_night_mode BOOLEAN DEFAULT false,
    
    -- 对话统计
    total_messages INT DEFAULT 0,
    total_goals_completed_together INT DEFAULT 0,
    
    -- 状态
    is_active BOOLEAN DEFAULT true,          -- 当前是否激活
    is_favorite BOOLEAN DEFAULT false,
    
    -- 切换限制（免费用户）
    switch_count_this_month INT DEFAULT 0,
    last_switch_date DATE,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(user_id, character_id)
);

CREATE INDEX idx_user_character_user ON user_character_configs(user_id);
CREATE INDEX idx_user_character_active ON user_character_configs(user_id, is_active);

-- ============================================
-- 对话历史表扩展（兼容原有表）
-- ============================================
-- 原有 chat_messages 表增加 character_id 字段
ALTER TABLE chat_messages ADD COLUMN character_id VARCHAR(50) 
    REFERENCES celebrity_characters(id);

-- 创建按角色查询的索引
CREATE INDEX idx_chat_messages_character ON chat_messages(user_id, character_id, created_at DESC);

-- ============================================
-- 角色推荐测试记录（新增）
-- ============================================
CREATE TABLE character_recommendation_tests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    answers JSONB NOT NULL,                  -- 5题答案
    recommended_character_id VARCHAR(50) REFERENCES celebrity_characters(id),
    user_accepted BOOLEAN,                   -- 用户是否接受推荐
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 2.2 扩展后的数据模型（Dart）

```dart
// models/celebrity_character.dart

@freezed
class CelebrityCharacter with _$CelebrityCharacter {
  const factory CelebrityCharacter({
    required String id,
    required String name,
    String? nameEn,
    required CharacterCategory category,
    required String era,
    required String nationality,
    required List<String> famousFor,
    required String avatarUrl,
    String? avatarLottieUrl,
    required List<String> styleTags,
    required CharacterDistillate distillate,
    required double qualityScore,
    required CharacterStatus status,
    required bool isPremium,
  }) = _CelebrityCharacter;

  factory CelebrityCharacter.fromJson(Map<String, dynamic> json) =>
      _$CelebrityCharacterFromJson(json);
}

enum CharacterCategory {
  writer('作家', Icons.menu_book),
  entrepreneur('创业者', Icons.rocket_launch),
  athlete('运动员', Icons.sports),
  philosopher('哲学家', Icons.psychology),
  fictional('虚构人物', Icons.auto_stories),
  scientist('科学家', Icons.science),
  polymath('通才', Icons.lightbulb);

  final String label;
  final IconData icon;
  const CharacterCategory(this.label, this.icon);
}

// models/user_character_config.dart

@freezed
class UserCharacterConfig with _$UserCharacterConfig {
  const factory UserCharacterConfig({
    required String id,
    required String userId,
    required String characterId,
    String? customName,
    String? customAvatarUrl,
    required int intimacyLevel,
    required int intimacyScore,
    required List<String> unlockedStories,
    required List<String> unlockedExpressions,
    required bool unlockedDeepNightMode,
    required int totalMessages,
    required int totalGoalsCompletedTogether,
    required bool isActive,
    required bool isFavorite,
    required int switchCountThisMonth,
    DateTime? lastSwitchDate,
  }) = _UserCharacterConfig;

  factory UserCharacterConfig.fromJson(Map<String, dynamic> json) =>
      _$UserCharacterConfigFromJson(json);
}

// 扩展 CompanionConfig，增加 character 相关字段
@freezed
class CompanionConfig with _$CompanionConfig {
  const factory CompanionConfig({
    // ... 原有字段 ...
    
    // 新增：名人角色相关
    String? characterId,           // 如果非空，表示使用名人角色
    bool get isCelebrity => characterId != null;
    
    // 用户自定义覆盖
    String? customName,
    String? customAvatarUrl,
  }) = _CompanionConfig;
}
```

---

## 三、API 扩展

### 3.1 新增 API 端点

```go
// ============================================
// 角色管理 API
// ============================================

// GET /api/v1/characters
// 获取所有可用角色列表
// Query: category, search, page, limit
// Response: [{ id, name, category, avatarUrl, styleTags, qualityScore, isPremium }]

// GET /api/v1/characters/:id
// 获取角色详情（不含完整蒸馏档案，仅展示信息）
// Response: { id, name, category, famousFor, avatarUrl, distillate_preview }

// GET /api/v1/characters/:id/distillate
// 获取完整蒸馏档案（需认证，用于客户端 prompt 组装）
// Response: 完整 CharacterDistillate JSON

// POST /api/v1/characters/recommend
// 角色推荐测试
// Body: { answers: ["A", "B", "C", "D", "A"] }
// Response: { recommendedCharacterId, matchScore, alternativeIds }

// ============================================
// 用户角色配置 API
// ============================================

// GET /api/v1/user/characters
// 获取用户已解锁的角色列表
// Response: [{ configId, characterId, name, intimacyLevel, isActive, isFavorite }]

// POST /api/v1/user/characters
// 选择/解锁新角色
// Body: { characterId, customName? }
// Response: { configId, character, remainingSwitches? }

// PUT /api/v1/user/characters/:configId
// 更新角色配置（昵称、头像等）
// Body: { customName?, customAvatarUrl?, isFavorite? }

// POST /api/v1/user/characters/:configId/switch
// 切换当前激活角色
// Response: { success, newActiveCharacter, remainingSwitches? }

// DELETE /api/v1/user/characters/:configId
// 删除角色配置（保留历史对话，仅删除配置）

// ============================================
// 扩展聊天 API（兼容原有）
// ============================================

// POST /api/v1/chat
// Body 增加可选字段: { message, characterId? }
// 如果不传 characterId，使用用户当前激活的角色
// 如果用户当前是基础伙伴，按原有逻辑处理

// ============================================
// 亲密度系统 API
// ============================================

// GET /api/v1/user/characters/:configId/intimacy
// 获取亲密度详情
// Response: { level, score, nextLevelScore, unlockedFeatures, progress }

// POST /api/v1/user/characters/:configId/intimacy/claim
// 领取亲密度奖励（解锁故事、表达等）
```

### 3.2 扩展后的 AI 聊天服务

```go
// internal/service/ai_service.go — 扩展版

type AIService struct {
    primary      AIClient
    fallback     AIClient
    memoryRepo   MemoryRepo
    chatRepo     ChatRepo
    characterRepo CharacterRepo  // 新增
}

func (s *AIService) Chat(ctx context.Context, userID string, characterID *string, message string) (<-chan string, error) {
    
    // 1. 获取用户画像 + 记忆摘要
    profile := s.memoryRepo.GetUserProfile(ctx, userID)
    memorySummary := s.memoryRepo.GetRecentSummary(ctx, userID)
    
    // 2. 构建 System Prompt（关键扩展点）
    var systemPrompt string
    
    if characterID != nil && *characterID != "" {
        // ===== 名人角色模式 =====
        character := s.characterRepo.GetByID(ctx, *characterID)
        distillate := character.Distillate
        
        // 组装名人专用 System Prompt
        systemPrompt = buildCelebritySystemPrompt(distillate, profile, memorySummary)
        
    } else {
        // ===== 基础伙伴模式（原有逻辑） =====
        companionConfig := s.getCompanionConfig(ctx, userID)
        systemPrompt = buildBaseSystemPrompt(companionConfig, profile, memorySummary)
    }
    
    // 3. 获取最近对话历史
    history := s.chatRepo.GetRecentHistory(ctx, userID, characterID, 10)
    
    // 4. 调用 AI（流式）— 原有逻辑不变
    stream, err := s.primary.StreamChat(ctx, systemPrompt, history, message)
    if err != nil {
        stream, err = s.fallback.StreamChat(ctx, systemPrompt, history, message)
    }
    
    // 5. 异步保存 + 更新记忆 — 原有逻辑不变
    go func() {
        fullResponse := collectStream(stream)
        s.chatRepo.SaveMessage(ctx, userID, characterID, "assistant", fullResponse)
        s.memoryRepo.UpdateSummary(ctx, userID, message, fullResponse)
        
        // 新增：更新名人角色亲密度
        if characterID != nil {
            s.characterRepo.UpdateIntimacy(ctx, userID, *characterID, 1)
        }
    }()
    
    return stream, nil
}

// 构建名人 System Prompt
func buildCelebritySystemPrompt(distillate CharacterDistillate, profile UserProfile, memory string) string {
    return fmt.Sprintf(`你是基于%s的思维方式和语言风格蒸馏而成的 AI 伙伴。

【风格指令】
%s

【价值观指令】
%s

【认知模式指令】
%s

【知识库注入】
%s

【互动约束】
- 你不代表%s本人，而是基于公开资料风格化生成的 AI 伙伴
- 避免涉及政治、品牌代言、未经证实的私人生活
- 使用%s的风格鼓励用户养成习惯、保持自律
- 回复长度控制在 100 字以内
- 善用个人故事和比喻

用户画像：
%s

长期记忆摘要：
%s
`, distillate.BasicInfo.Name,
   strings.Join(distillate.StyleProfile.StyleDirectives, "\n"),
   strings.Join(distillate.ValuesProfile.ValuesDirectives, "\n"),
   strings.Join(distillate.CognitiveProfile.CognitiveDirectives, "\n"),
   formatKnowledgeInjection(distillate.KnowledgeProfile),
   distillate.BasicInfo.Name,
   distillate.BasicInfo.Name,
   profile.Summary,
   memory)
}
```

---

## 四、前端集成方案

### 4.1 导航扩展

```dart
// 底部 Tab Bar 增加角色切换入口

┌────────┬────────┬────────┬────────┬────────┐
│  💬    │  🌱    │  🎯    │  📖    │  🎭    │
│  伙伴   │  习惯   │  专注   │  成长   │  角色   │  ← 新增"角色"Tab
└────────┴────────┴────────┴────────┴────────┘

// 或：在"伙伴"Tab 内集成角色选择
```

### 4.2 角色选择页面

```dart
// features/characters/character_selection_screen.dart

class CharacterSelectionScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final characters = ref.watch(availableCharactersProvider);
    final userConfigs = ref.watch(userCharacterConfigsProvider);
    
    return Scaffold(
      appBar: AppBar(title: Text('选择你的伙伴')),
      body: Column(
        children: [
          // 顶部：推荐测试入口
          CharacterRecommendationCard(),
          
          // 分类筛选
          CategoryFilterChips(),
          
          // 角色网格
          Expanded(
            child: characters.when(
              data: (list) => CharacterGrid(
                characters: list,
                userConfigs: userConfigs,
                onSelect: (character) => _onSelectCharacter(context, ref, character),
              ),
              loading: () => ShimmerCharacterGrid(),
              error: (e, _) => ErrorWidget(e),
            ),
          ),
        ],
      ),
    );
  }
  
  void _onSelectCharacter(BuildContext context, WidgetRef ref, CelebrityCharacter character) {
    // 检查是否已解锁
    final config = ref.read(userCharacterConfigsProvider)
        .valueOrNull
        ?.firstWhereOrNull((c) => c.characterId == character.id);
    
    if (config != null) {
      // 已解锁，直接切换
      ref.read(companionNotifierProvider.notifier).switchCharacter(character.id);
      context.pop();
    } else {
      // 未解锁，显示详情页
      context.push('/characters/${character.id}');
    }
  }
}
```

### 4.3 聊天界面适配

```dart
// features/companion/companion_screen.dart — 扩展版

class CompanionScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final companionState = ref.watch(companionNotifierProvider);
    final currentCharacter = companionState.currentCharacter;
    
    return Scaffold(
      appBar: AppBar(
        title: Row(
          children: [
            // 角色头像（可点击切换）
            GestureDetector(
              onTap: () => context.push('/characters'),
              child: CircleAvatar(
                backgroundImage: NetworkImage(
                  currentCharacter?.avatarUrl ?? companionState.defaultAvatarUrl,
                ),
              ),
            ),
            SizedBox(width: 8),
            // 角色名称
            Text(currentCharacter?.name ?? 'Lagom'),
            // 亲密度标识
            if (companionState.intimacyLevel > 1)
              IntimacyBadge(level: companionState.intimacyLevel),
          ],
        ),
      ),
      body: ChatInterface(
        messages: companionState.messages,
        isTyping: companionState.isTyping,
        // 根据角色定制气泡颜色
        bubbleColor: currentCharacter?.themeColor ?? LagomColors.orange,
        onSend: (msg) => ref.read(companionNotifierProvider.notifier).sendMessage(msg),
      ),
    );
  }
}
```

### 4.4 亲密度系统 UI

```dart
// features/characters/intimacy_progress_card.dart

class IntimacyProgressCard extends StatelessWidget {
  final UserCharacterConfig config;
  final CelebrityCharacter character;
  
  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text('与 ${character.name} 的亲密度'),
                Spacer(),
                Text('Lv.${config.intimacyLevel}'),
              ],
            ),
            SizedBox(height: 8),
            // 进度条
            LinearProgressIndicator(
              value: config.intimacyScore / _nextLevelScore(config.intimacyLevel),
              backgroundColor: Colors.grey[200],
              valueColor: AlwaysStoppedAnimation(LagomColors.orange),
            ),
            SizedBox(height: 8),
            // 已解锁内容
            Wrap(
              spacing: 8,
              children: [
                if (config.intimacyLevel >= 5)
                  Chip(label: Text('个人故事')),
                if (config.intimacyLevel >= 15)
                  Chip(label: Text('专属昵称')),
                if (config.unlockedDeepNightMode)
                  Chip(label: Text('深夜模式')),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
```

---

## 五、后端集成方案

### 5.1 服务层扩展

```go
// internal/service/character_service.go

type CharacterService struct {
    characterRepo CharacterRepo
    userCharRepo  UserCharacterRepo
    distillRepo   DistillateRepo
}

// 获取可用角色列表
func (s *CharacterService) ListCharacters(ctx context.Context, userID string, filter CharacterFilter) ([]CharacterDTO, error) {
    characters := s.characterRepo.List(ctx, filter)
    
    // 标记用户已解锁的角色
    userConfigs := s.userCharRepo.GetByUser(ctx, userID)
    configMap := make(map[string]*UserCharacterConfig)
    for _, c := range userConfigs {
        configMap[c.CharacterID] = c
    }
    
    var result []CharacterDTO
    for _, c := range characters {
        dto := CharacterDTO{
            Character:     c,
            IsUnlocked:    configMap[c.ID] != nil,
            IsActive:      configMap[c.ID] != nil && configMap[c.ID].IsActive,
            IntimacyLevel: configMap[c.ID]?.IntimacyLevel ?? 0,
        }
        result = append(result, dto)
    }
    
    return result, nil
}

// 切换角色
func (s *CharacterService) SwitchCharacter(ctx context.Context, userID, configID string) error {
    // 1. 获取用户配置
    config := s.userCharRepo.GetByID(ctx, configID)
    if config.UserID != userID {
        return errors.New("unauthorized")
    }
    
    // 2. 检查切换限制（免费用户）
    user := s.userRepo.GetByID(ctx, userID)
    if user.SubscriptionTier == "free" {
        if config.SwitchCountThisMonth >= 2 {
            return errors.New("switch limit reached for free user")
        }
    }
    
    // 3. 将当前激活角色设为非激活
    s.userCharRepo.DeactivateAll(ctx, userID)
    
    // 4. 激活新角色
    s.userCharRepo.Activate(ctx, configID)
    
    // 5. 更新切换统计
    s.userCharRepo.IncrementSwitchCount(ctx, configID)
    
    return nil
}

// 角色推荐测试
func (s *CharacterService) RecommendCharacter(ctx context.Context, answers []string) (*RecommendationResult, error) {
    // 简单的规则匹配算法
    scores := make(map[string]int)
    
    // Q1: 回应方式偏好
    switch answers[0] {
    case "A": scores["haruki_murakami_v1"] += 3; scores["lao_tzu_v1"] += 2
    case "B": scores["elon_musk_v1"] += 3; scores["kobe_bryant_v1"] += 2
    case "C": scores["charlie_munger_v1"] += 3; scores["dumbledore_v1"] += 2
    case "D": scores["kobe_bryant_v1"] += 3; scores["eliud_kipchoge_v1"] += 2
    }
    
    // Q2: 学习方式
    switch answers[1] {
    case "A": scores["haruki_murakami_v1"] += 3
    case "B": scores["elon_musk_v1"] += 3
    case "C": scores["charlie_munger_v1"] += 3
    case "D": scores["steve_jobs_v1"] += 3
    }
    
    // Q3: 打动的话
    switch answers[2] {
    case "A": scores["haruki_murakami_v1"] += 3; scores["su_shi_v1"] += 2
    case "B": scores["elon_musk_v1"] += 3
    case "C": scores["charlie_munger_v1"] += 3
    case "D": scores["kobe_bryant_v1"] += 3
    }
    
    // Q4: 理想早晨
    switch answers[3] {
    case "A": scores["haruki_murakami_v1"] += 3; scores["lao_tzu_v1"] += 2
    case "B": scores["elon_musk_v1"] += 3; scores["charlie_munger_v1"] += 2
    case "C": scores["charlie_munger_v1"] += 3
    case "D": scores["kobe_bryant_v1"] += 3; scores["eliud_kipchoge_v1"] += 2
    }
    
    // Q5: 激励方式
    switch answers[4] {
    case "A": scores["haruki_murakami_v1"] += 3; scores["dumbledore_v1"] += 2
    case "B": scores["elon_musk_v1"] += 3
    case "C": scores["charlie_munger_v1"] += 3
    case "D": scores["kobe_bryant_v1"] += 3
    }
    
    // 找出最高分
    var bestID string
    var bestScore int
    for id, score := range scores {
        if score > bestScore {
            bestScore = score
            bestID = id
        }
    }
    
    // 找备选
    var alternatives []string
    for id, score := range scores {
        if id != bestID && score >= bestScore-2 {
            alternatives = append(alternatives, id)
        }
    }
    
    return &RecommendationResult{
        RecommendedCharacterID: bestID,
        MatchScore:             float64(bestScore) / 15.0, // max 15
        AlternativeIDs:         alternatives,
    }, nil
}
```

### 5.2 蒸馏档案管理

```go
// internal/repository/character_repo.go

type CharacterRepo interface {
    // CRUD
    GetByID(ctx context.Context, id string) (*CelebrityCharacter, error)
    List(ctx context.Context, filter CharacterFilter) ([]*CelebrityCharacter, error)
    Create(ctx context.Context, character *CelebrityCharacter) error
    Update(ctx context.Context, character *CelebrityCharacter) error
    
    // 批量导入
    BatchImport(ctx context.Context, characters []*CelebrityCharacter) error
}

// 蒸馏档案存储（JSONB 列）
type CelebrityCharacter struct {
    ID          string
    Name        string
    Distillate  CharacterDistillate  // 存储为 JSONB
    // ... 其他字段
}
```

---

## 六、Prompt 注入系统

### 6.1 动态 Prompt 组装

```go
// internal/pkg/ai/prompt_builder.go

type PromptBuilder struct {
    baseTemplate    string
    characterRepo   CharacterRepo
}

func (b *PromptBuilder) Build(ctx context.Context, userID string, characterID *string, message string) (string, error) {
    
    // 基础 prompt（所有角色共用）
    prompt := b.baseTemplate
    
    if characterID != nil {
        // 加载角色蒸馏档案
        character := b.characterRepo.GetByID(ctx, *characterID)
        distillate := character.Distillate
        
        // 注入风格层
        prompt += b.buildStyleLayer(distillate.StyleProfile)
        
        // 注入价值观层
        prompt += b.buildValuesLayer(distillate.ValuesProfile)
        
        // 注入认知层
        prompt += b.buildCognitiveLayer(distillate.CognitiveProfile)
        
        // 注入知识层（RAG）
        knowledge := b.retrieveRelevantKnowledge(ctx, *characterID, message)
        prompt += b.buildKnowledgeLayer(knowledge)
        
        // 注入互动约束
        prompt += b.buildConstraintLayer(distillate.InteractionConstraints)
    }
    
    return prompt, nil
}

// RAG：检索相关知识
func (b *PromptBuilder) retrieveRelevantKnowledge(ctx context.Context, characterID, query string) KnowledgeSnippet {
    // 1. 生成 query embedding
    embedding := b.aiClient.Embed(ctx, query)
    
    // 2. 在角色知识库中向量检索
    results := b.vectorStore.SimilaritySearch(ctx, characterID, embedding, 3)
    
    // 3. 组装为提示片段
    return KnowledgeSnippet{
        Stories:    results.Stories,
        Analogies:  results.Analogies,
        Expressions: results.Expressions,
    }
}
```

### 6.2 Prompt 版本管理

```yaml
# config/prompt_versions.yaml

prompt_templates:
  base_v1:
    template: |
      你是 Lagom 的 AI 伙伴，帮助用户养成自律习惯...
    
  celebrity_v1:
    template: |
      {{BASE_PROMPT}}
      
      你当前扮演的角色是：{{CHARACTER_NAME}}
      
      【风格层】
      {{STYLE_LAYER}}
      
      【价值观层】
      {{VALUES_LAYER}}
      
      【认知层】
      {{COGNITIVE_LAYER}}
      
      【相关知识】
      {{KNOWLEDGE_LAYER}}
      
      【约束】
      {{CONSTRAINT_LAYER}}
      
      记住：你不是{{CHARACTER_NAME}}本人，而是基于公开资料风格化生成的 AI 伙伴。
      你的核心目标是：用{{CHARACTER_NAME}}的风格帮助用户养成自律习惯。
```

---

## 七、缓存与性能优化

### 7.1 角色档案缓存

```go
// Redis 缓存策略

// 1. 蒸馏档案缓存（几乎不变，长期缓存）
// Key: character:distillate:{character_id}
// TTL: 7 days
// Value: 完整蒸馏档案 JSON

// 2. 角色列表缓存（按分类）
// Key: characters:list:{category}:{page}:{limit}
// TTL: 1 hour
// Value: 角色列表 DTO

// 3. 用户角色配置缓存
// Key: user:{user_id}:characters
// TTL: 15 minutes
// Value: 用户所有角色配置

// 4. 亲密度排行榜（可选）
// Key: characters:intimacy:leaderboard
// TTL: 5 minutes
// Value: 全局亲密度排行（匿名）
```

### 7.2 AI 调用优化

```go
// 减少 API 成本的策略

type AIOptimizer struct {
    cache AIResponseCache
}

// 意图缓存：常见意图的回复缓存
func (o *AIOptimizer) CheckIntentCache(ctx context.Context, characterID, message string) (*CachedResponse, bool) {
    // 对"早安""晚安""完成打卡"等高频率意图，使用模板回复
    intent := o.classifyIntent(message)
    if intent.IsCommon && intent.Confidence > 0.9 {
        return o.cache.Get(ctx, characterID, intent.Type), true
    }
    return nil, false
}

// 响应缓存：相同上下文+相同问题的回复缓存（短时）
func (o *AIOptimizer) CheckResponseCache(ctx context.Context, characterID, message string) (*CachedResponse, bool) {
    // Key: hash(character_id + system_prompt_hash + last_3_messages + message)
    // TTL: 5 minutes
    return o.cache.Get(ctx, o.buildCacheKey(characterID, message)), true
}
```

---

## 八、数据迁移策略

### 8.1 已有用户的数据兼容

```go
// 迁移脚本

func MigrateExistingUsers(db *sqlx.DB) error {
    // 1. 为所有现有用户创建"默认伙伴"配置
    //    （保持原有体验不变）
    
    // 2. 现有对话历史标记为"base_companion"
    _, err := db.Exec(`
        UPDATE chat_messages 
        SET character_id = 'base_companion' 
        WHERE character_id IS NULL
    `)
    
    // 3. 在 companion_configs 表中，为现有用户插入 base companion 记录
    _, err = db.Exec(`
        INSERT INTO user_character_configs 
        (user_id, character_id, is_active, created_at)
        SELECT id, 'base_companion', true, NOW()
        FROM users
        WHERE id NOT IN (
            SELECT user_id FROM user_character_configs
        )
    `)
    
    return err
}
```

### 8.2 灰度发布策略

```
Phase 1: 内测（1周）
- 仅向 5% 用户展示名人功能入口
- 收集反馈、修复 bug

Phase 2: 软发布（1周）
- 向 30% 用户开放
- 监控 AI 调用成本、用户留存变化

Phase 3: 全量发布
- 向所有用户开放
- 配合营销活动（"选择你的自律导师"）
```

---

## 九、测试策略

### 9.1 单元测试

```go
// Test: 角色推荐算法
func TestCharacterRecommendation(t *testing.T) {
    service := NewCharacterService(/* mock repos */)
    
    // 测试：选择A（温柔理解）为主的用户
    answers := []string{"A", "A", "A", "A", "A"}
    result, err := service.RecommendCharacter(context.Background(), answers)
    
    assert.NoError(t, err)
    assert.Equal(t, "haruki_murakami_v1", result.RecommendedCharacterID)
    assert.True(t, result.MatchScore > 0.8)
}

// Test: Prompt 构建
func TestPromptBuilder(t *testing.T) {
    builder := NewPromptBuilder(/* mock */)
    
    prompt, err := builder.Build(ctx, "user1", stringPtr("haruki_murakami_v1"), "我今天不想跑步")
    
    assert.NoError(t, err)
    assert.Contains(t, prompt, "村上春树")
    assert.Contains(t, prompt, "找到自己的节奏")
}
```

### 9.2 集成测试

```go
// Test: 完整对话流程
func TestCelebrityChatFlow(t *testing.T) {
    // 1. 用户选择村上春树
    // 2. 发送消息
    // 3. 验证回复中包含村上风格特征
    // 4. 验证亲密度 +1
}
```

### 9.3 盲测流程

```
1. 准备 10 个测试问题
2. 对每个问题，生成：村上春树版、马斯克版、通用版回复
3. 邀请 20 名用户（知道是村上春树，但不知道哪个回复是哪版）
4. 用户选择"最像村上春树"的回复
5. 目标：村上版的选中率 > 70%
```
