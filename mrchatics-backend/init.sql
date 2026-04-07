-- Создаем таблицу
CREATE TABLE IF NOT EXISTS user_statuses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'offline',
    custom_status TEXT,
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индексы
CREATE INDEX IF NOT EXISTS idx_user_statuses_user_id ON user_statuses(user_id);
CREATE INDEX IF NOT EXISTS idx_user_statuses_status ON user_statuses(status);
CREATE INDEX IF NOT EXISTS idx_user_statuses_last_seen ON user_statuses(last_seen);

-- Создаем таблицу истории
CREATE TABLE IF NOT EXISTS status_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    old_status VARCHAR(50),
    new_status VARCHAR(50),
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Очищаем таблицы
TRUNCATE TABLE user_statuses CASCADE;

-- Добавляем тестовых пользователей
INSERT INTO user_statuses (user_id, status, custom_status) VALUES
(1, 'online', 'Working on project'),
(2, 'away', 'In a meeting'),
(3, 'busy', 'Deep focus mode'),
(4, 'offline', 'Taking a break'),
(5, 'online', 'Available for chat'),
(6, 'busy', 'Code review'),
(7, 'away', 'Lunch break'),
(8, 'online', 'Helping team'),
(9, 'offline', 'End of day'),
(10, 'online', 'Ready to help'),
(11, 'busy', 'Important task'),
(12, 'away', 'Coffee break'),
(13, 'online', 'Taking requests'),
(14, 'offline', 'Out of office'),
(15, 'online', 'Working remotely');

-- Выводим результат
SELECT COUNT(*) as total_users FROM user_statuses;
SELECT user_id, status, custom_status FROM user_statuses ORDER BY user_id;
