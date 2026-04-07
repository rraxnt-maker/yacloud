-- migrations/001_create_profile_tables.sql

-- Таблица профилей пользователей
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    full_name VARCHAR(255),
    bio TEXT,
    avatar_url TEXT,
    phone VARCHAR(50),
    location VARCHAR(255),
    website VARCHAR(255),
    company VARCHAR(255),
    position VARCHAR(255),
    birth_date DATE,
    gender VARCHAR(20),
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Таблица настроек пользователя
CREATE TABLE user_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Notification settings
    email_notifications BOOLEAN DEFAULT true,
    push_notifications BOOLEAN DEFAULT true,
    slack_notifications BOOLEAN DEFAULT false,
    
    -- Privacy settings
    profile_visibility VARCHAR(20) DEFAULT 'public', -- public, private, contacts_only
    show_online_status BOOLEAN DEFAULT true,
    show_last_seen BOOLEAN DEFAULT true,
    show_email BOOLEAN DEFAULT false,
    show_phone BOOLEAN DEFAULT false,
    
    -- Theme & UI
    theme VARCHAR(20) DEFAULT 'light', -- light, dark, system
    compact_view BOOLEAN DEFAULT false,
    
    -- Security
    two_factor_enabled BOOLEAN DEFAULT false,
    session_timeout_minutes INTEGER DEFAULT 60,
    
    -- Preferences
    default_dashboard VARCHAR(50) DEFAULT 'home',
    items_per_page INTEGER DEFAULT 20,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_user_profiles_name ON user_profiles(full_name);
CREATE INDEX idx_user_profiles_location ON user_profiles(location);
CREATE INDEX idx_user_settings_visibility ON user_settings(profile_visibility);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_profiles_updated_at 
    BEFORE UPDATE ON user_profiles 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at 
    BEFORE UPDATE ON user_settings 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
