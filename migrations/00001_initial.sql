-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_last_checked_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_checked_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE "tg_users" (
  "user_id" int4 PRIMARY key,
  "user_name" varchar(200) NOT NULL,
  "first_name" varchar(200) NOT NULL,
  "last_name" varchar(200) NOT NULL
);
COMMENT ON TABLE tg_users IS 'Телграмовские пользователи';
COMMENT ON COLUMN tg_users.user_id IS 'Телеграмовский идентификатор пользователя';
COMMENT ON COLUMN tg_users.user_name IS 'Его пользовательское имя';
COMMENT ON COLUMN tg_users.first_name IS 'Имя';
COMMENT ON COLUMN tg_users.last_name IS 'Фамилия';

CREATE TABLE "sites" (
  "id" SERIAL4 PRIMARY key,
  "url" varchar(200) NOT NULL UNIQUE,
  "working" bool NOT NULL,
  "status_code" int NOT NULL,
  "secret_key" varchar(200) NOT NULL,
  "last_checked_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE sites IS 'Добавленные пользователями сайты для мониторинга';
COMMENT ON COLUMN sites.id IS 'Внутренний идентификатор сайта';
COMMENT ON COLUMN sites.url IS 'Ссылка на сайт';
COMMENT ON COLUMN sites.working IS 'Работает ли сайт';
COMMENT ON COLUMN sites.status_code IS 'Код ответа сайта';
COMMENT ON COLUMN sites.secret_key IS 'Секретный ключ для gRPC';

CREATE TRIGGER set_last_checked_at
BEFORE UPDATE ON sites
FOR EACH ROW
EXECUTE FUNCTION update_last_checked_at_column();

CREATE TABLE "tg_user_sites" (
  "id" SERIAL4 PRIMARY key,
  "tg_user_id" integer NOT NULL,
  "site_id" integer NOT NULL,
  "monitoring" bool NOT NULL,
  "duration_minutes" int NOT NULL,
  "is_staff" bool NOT NULL DEFAULT false,
  "is_super" bool NOT NULL DEFAULT false
);
COMMENT ON TABLE tg_user_sites IS 'Связь между пользователем и сайтом';
COMMENT ON COLUMN tg_user_sites.id IS 'Идентификатор связи';
COMMENT ON COLUMN tg_user_sites.tg_user_id is 'Пользовательский идентификатор телеграм';
COMMENT ON COLUMN tg_user_sites.site_id is 'Идентификатор сайта';
COMMENT ON COLUMN tg_user_sites.monitoring IS 'Мониторить сайт или нет';
COMMENT ON COLUMN tg_user_sites.duration_minutes is 'Период, в течении которого производится проверка сайта на доступность';
COMMENT ON COLUMN tg_user_sites.is_staff is 'Является ли пользователем сотрудником этого сайта';
COMMENT ON COLUMN tg_user_sites.is_super IS 'Является ли пользователь админом этого сайта';


CREATE UNIQUE INDEX "tg_user_site" ON "tg_user_sites" ("tg_user_id", "site_id");
COMMENT ON INDEX tg_user_site IS 'Ограничение связи m2m между пользователем и сайтом на уникальность';

ALTER TABLE "tg_user_sites" ADD FOREIGN KEY ("site_id") REFERENCES "sites" ("id");
ALTER TABLE "tg_user_sites" ADD FOREIGN KEY ("tg_user_id") REFERENCES "tg_users" ("user_id");


CREATE TABLE "feedbacks" (
  "id" SERIAL4 PRIMARY KEY,
  "site_id" integer NOT NULL,
  "name" varchar(300) NOT NULL,
  "contact" varchar(300) NOT NULL,
  "message" varchar(500) NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE feedbacks IS 'Обращения от пользователей с обслуживаемых сайтов';
COMMENT ON COLUMN feedbacks.id IS 'Идентификатор обращения';
COMMENT ON COLUMN feedbacks.site_id is 'Идентификатор сайта';
COMMENT ON COLUMN feedbacks.name is 'Имя обратившегося пользователя';
COMMENT ON COLUMN feedbacks.contact is 'Контакты дл связи';
COMMENT ON COLUMN feedbacks.message is 'Дополнительное сообщение';
COMMENT ON COLUMN feedbacks.created_at is 'Дата обращения';

ALTER TABLE "feedbacks" ADD FOREIGN KEY ("site_id") REFERENCES "sites" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feedbacks;
DROP TABLE IF EXISTS tg_user_sites;
DROP TABLE IF EXISTS sites;
DROP TABLE IF EXISTS telegram_users;
-- +goose StatementEnd
