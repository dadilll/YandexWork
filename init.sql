CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- Создаем индекс для быстрого доступа к пользователю по логину
CREATE INDEX idx_users_login ON users(login);

CREATE TABLE user_tasks (
    user_id INTEGER REFERENCES users(id),
    task_id VARCHAR(36) REFERENCES tasks(id) ON DELETE CASCADE
);

-- Создаем индекс для быстрого доступа к задачам пользователя
CREATE INDEX idx_user_tasks_user_id ON user_tasks(user_id);

CREATE TABLE locks (
    id TEXT PRIMARY KEY,
    status TEXT
);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    expression TEXT,
    status TEXT,
    result REAL
);