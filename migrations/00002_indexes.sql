-- +goose Up
-- +goose StatementBegin

-- Индексы для tg_users
CREATE INDEX idx_tg_users_user_name ON tg_users (user_name);

-- Индексы для sites
CREATE INDEX idx_sites_working ON sites (working);
CREATE INDEX idx_sites_status_code ON sites (status_code);

-- Индексы для tg_user_sites
CREATE INDEX idx_tg_user_sites_tg_user_id ON tg_user_sites (tg_user_id);
CREATE INDEX idx_tg_user_sites_site_id ON tg_user_sites (site_id);

-- Индексы для feedbacks
CREATE INDEX idx_feedbacks_site_id ON feedbacks (site_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Индексы для tg_users
DROP INDEX IF EXISTS idx_tg_users_user_name;

-- Индексы для sites
DROP INDEX IF EXISTS idx_sites_working;
DROP INDEX IF EXISTS idx_sites_status_code;

-- Индексы для tg_user_sites
DROP INDEX IF EXISTS idx_tg_user_sites_tg_user_id;
DROP INDEX IF EXISTS idx_tg_user_sites_site_id;

-- Индексы для feedbacks
DROP INDEX IF EXISTS idx_feedbacks_site_id;

-- +goose StatementEnd
