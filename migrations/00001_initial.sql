-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tg_users" (
  "user_id" int4 PRIMARY key,
  "user_name" varchar(200) NOT NULL,
  "first_name" varchar(200) NOT NULL,
  "last_name" varchar(200) NOT NULL
);

CREATE TABLE "sites" (
  "id" SERIAL4 PRIMARY key,
  "url" varchar(200) NOT NULL UNIQUE
);

CREATE TABLE "tg_user_sites" (
  "id" SERIAL4 PRIMARY key,
  "tg_user_id" integer NOT NULL,
  "site_id" integer NOT NULL,
  "working" bool NOT NULL,
  "status_code" int NOT NULL,
  "monitoring" bool NOT NULL,
  "duration_minutes" int NOT NULL,
  "is_staff" bool NOT NULL
);

CREATE UNIQUE INDEX "tg_user_site" ON "tg_user_sites" ("tg_user_id", "site_id");

ALTER TABLE "tg_user_sites" ADD FOREIGN KEY ("tg_user_id") REFERENCES "tg_users" ("user_id");

ALTER TABLE "tg_user_sites" ADD FOREIGN KEY ("site_id") REFERENCES "sites" ("id");

CREATE TABLE "feedbacks" (
  "id" SERIAL4 PRIMARY KEY,
  "site_id" integer NOT NULL,
  "name" varchar(300) NOT NULL,
  "contact" varchar(300) NOT NULL,
  "message" varchar(500) NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "feedbacks" ADD FOREIGN KEY ("site_id") REFERENCES "sites" ("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feedbacks;
DROP TABLE IF EXISTS tg_user_sites;
DROP TABLE IF EXISTS sites;
DROP TABLE IF EXISTS telegram_users;
-- +goose StatementEnd
