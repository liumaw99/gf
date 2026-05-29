-- ============================================
-- 初始迁移：Lagom 数据库 Schema
-- ============================================

-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- 用户表
-- ============================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    apple_id VARCHAR(255) UNIQUE,
    google_id VARCHAR(255) UNIQUE,
    email VARCHAR(255),
    name VARCHAR(100),
    avatar_url TEXT,
    subscription_tier VARCHAR(20) DEFAULT 'free',
    subscription_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- 名人角色表
-- ============================================
CREATE TABLE celebrity_characters (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    category VARCHAR(50) NOT NULL,
    era VARCHAR(50),
    nationality VARCHAR(50),
    famous_for TEXT[] DEFAULT '{}',
    avatar_url TEXT,
    avatar_lottie_url TEXT,
    style_tags TEXT[] DEFAULT '{}',
    distillate JSONB NOT NULL DEFAULT '{}',
    quality_score DECIMAL(3,2) DEFAULT 0.0,
    status VARCHAR(20) DEFAULT 'active',
    is_premium BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- 伙伴配置表
-- ============================================
CREATE TABLE companion_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL DEFAULT 'Lagom',
    avatar_type VARCHAR(20) DEFAULT 'fox',
    personality_type VARCHAR(20) DEFAULT 'gentle',
    encouragement_style VARCHAR(20) DEFAULT 'soft',
    voice_type VARCHAR(20),
    growth_level INT DEFAULT 1,
    relationship_score DECIMAL(3,2) DEFAULT 0.0,
    character_id VARCHAR(50) REFERENCES celebrity_characters(id),
    custom_name VARCHAR(100),
    custom_avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id)
);

-- ============================================
-- 用户角色配置表
-- ============================================
CREATE TABLE user_character_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    character_id VARCHAR(50) NOT NULL REFERENCES celebrity_characters(id),
    custom_name VARCHAR(100),
    custom_avatar_url TEXT,
    intimacy_level INT DEFAULT 1,
    intimacy_score INT DEFAULT 0,
    unlocked_stories TEXT[] DEFAULT '{}',
    unlocked_expressions TEXT[] DEFAULT '{}',
    unlocked_deep_night_mode BOOLEAN DEFAULT false,
    total_messages INT DEFAULT 0,
    total_goals_completed INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_favorite BOOLEAN DEFAULT false,
    switch_count_this_month INT DEFAULT 0,
    last_switch_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, character_id)
);

-- ============================================
-- 习惯表
-- ============================================
CREATE TABLE habits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    priority INT DEFAULT 3,
    frequency_type VARCHAR(20) DEFAULT 'daily',
    target_days INT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    current_phase INT DEFAULT 1,
    hsi DECIMAL(5,2) DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- 微行动表
-- ============================================
CREATE TABLE micro_actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    phase INT NOT NULL,
    day_in_phase INT NOT NULL,
    action_text VARCHAR(300) NOT NULL,
    estimated_duration INT,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- 每日记录表
-- ============================================
CREATE TABLE daily_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    mood_score INT CHECK (mood_score BETWEEN 1 AND 5),
    reflection_note TEXT,
    summary TEXT,
    UNIQUE(user_id, date)
);

-- ============================================
-- 每日完成记录表
-- ============================================
CREATE TABLE daily_completions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    completed_at TIMESTAMP DEFAULT NOW(),
    quality_score INT CHECK (quality_score BETWEEN 1 AND 5),
    note TEXT,
    UNIQUE(user_id, habit_id, date)
);

-- ============================================
-- 专注会话表
-- ============================================
CREATE TABLE focus_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    duration INT NOT NULL,
    focus_quality_score DECIMAL(3,2),
    distractions TEXT[] DEFAULT '{}',
    context VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- 对话历史表
-- ============================================
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    character_id VARCHAR(50) REFERENCES celebrity_characters(id),
    role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    is_complete BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- 索引
-- ============================================
CREATE INDEX idx_users_apple ON users(apple_id);
CREATE INDEX idx_users_google ON users(google_id);
CREATE INDEX idx_companion_user ON companion_configs(user_id);
CREATE INDEX idx_user_char_user ON user_character_configs(user_id);
CREATE INDEX idx_user_char_active ON user_character_configs(user_id, is_active);
CREATE INDEX idx_habits_user ON habits(user_id);
CREATE INDEX idx_habits_user_active ON habits(user_id, is_active);
CREATE INDEX idx_micro_actions_habit ON micro_actions(habit_id);
CREATE INDEX idx_daily_logs_user_date ON daily_logs(user_id, date DESC);
CREATE INDEX idx_daily_completions_user_date ON daily_completions(user_id, date DESC);
CREATE INDEX idx_focus_sessions_user ON focus_sessions(user_id, start_time DESC);
CREATE INDEX idx_chat_messages_user_created ON chat_messages(user_id, created_at DESC);
CREATE INDEX idx_chat_messages_character ON chat_messages(user_id, character_id, created_at DESC);
CREATE INDEX idx_celebrity_category ON celebrity_characters(category);
CREATE INDEX idx_celebrity_status ON celebrity_characters(status);
