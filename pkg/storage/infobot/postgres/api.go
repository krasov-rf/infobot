package infobotdb_pg

import (
	"context"
	"database/sql"
	"errors"

	er "github.com/krasov-rf/infobot/pkg/errors"
	"github.com/krasov-rf/infobot/pkg/serializers"
	infobotdb "github.com/krasov-rf/infobot/pkg/storage/infobot"
)

// Зарегестрировать нового телеграм пользователя
func (d *InfoBotDb) TelegramUserRegister(ctx context.Context, user *serializers.UserSerializer) error {
	sqlt := `
		INSERT INTO tg_users (user_id, user_name, first_name, last_name)
		VALUES (:user_id, :user_name, :first_name, :last_name)
	`
	_, err := d.NamedExecContext(ctx, sqlt, user)
	if err != nil {
		return err
	}
	return nil
}

// Получить пользователя
func (d *InfoBotDb) TelegramUserGet(ctx context.Context, user_id int64) (*serializers.UserSerializer, error) {
	var res serializers.UserSerializer
	sqlt := `
		SELECT user_id, user_name, first_name, last_name
		FROM tg_users
		WHERE user_id = $1
	`
	err := d.GetContext(ctx, &res, sqlt, user_id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// вывести сайты
func (d *InfoBotDb) MonitoringSites(ctx context.Context, opt *infobotdb.OptionsInfoBot) ([]*serializers.SiteSerializer, int, error) {
	sqlt, err := infobotdb.Template("sites", `
		SELECT 
			s.id, s.url, s.working, s.status_code, s.secret_key, s.last_checked_at, 
			t.monitoring, t.duration_minutes,
			COUNT(*) OVER (PARTITION BY tg_user_id) AS count_user_sites
		FROM tg_user_sites t
		LEFT JOIN sites s ON t.site_id = s.id
		WHERE t.tg_user_id = :user_id
			{{ if .Id }} AND s.id = :id {{ end }}

		{{ if .Limit }}OFFSET :offset LIMIT :limit {{ end }}
	`, opt)
	if err != nil {
		return nil, 0, err
	}

	row, err := d.NamedQueryContext(ctx, sqlt, opt)
	if err != nil {
		return nil, 0, err
	}

	limit := 5
	if opt.Limit != 0 {
		limit = opt.Limit
	}

	var cnt int
	res := make([]*serializers.SiteSerializer, 0, limit)
	for row.Next() {
		var s serializers.SiteSerializer
		err = row.Scan(
			&s.Id,
			&s.Url, &s.Working, &s.StatusCode,
			&s.SecretKey, &s.LastCheckedAt,
			&s.Monitoring, &s.DurationMinutes,
			&cnt,
		)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, &s)
	}

	return res, cnt, nil
}

// Добавить новый сайт для мониторинга
func (d *InfoBotDb) MonitoringSiteAdd(
	ctx context.Context,
	user_id int64,
	site_url string,
	working bool,
	status_code int,
) (*serializers.SiteSerializer, error) {
	site := serializers.SiteSerializer{
		Url:             site_url,
		Working:         working,
		StatusCode:      status_code,
		Monitoring:      false,
		DurationMinutes: 15,
	}

	tx, err := d.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.Get(&site.Id, `SELECT id FROM sites WHERE url = $1`, site.Url)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		site.SecretKey, err = infobotdb.GenerateSecretKey(64)
		if err != nil {
			return nil, err
		}

		row := tx.QueryRowContext(ctx, `
				INSERT INTO sites (url, working, status_code, secret_key)
				VALUES ($1, $2, $3, $4)
				RETURNING id
			`,
			site.Url, site.Working, site.StatusCode,
			site.SecretKey,
		)
		if err = row.Err(); err != nil {
			return nil, err
		}
		err = row.Scan(&site.Id)
		if err != nil {
			return nil, err
		}
	}

	var ex bool
	err = tx.Get(&ex, `
		SELECT EXISTS(
			SELECT 1 FROM tg_user_sites 
			WHERE site_id = $1 and 
			      tg_user_id = $2
		)
	`, site.Id, user_id)
	if err != nil {
		return nil, err
	}
	if ex {
		return nil, er.ErrorExist
	}

	_, err = tx.Exec(`
		INSERT INTO tg_user_sites 
		(site_id, tg_user_id, monitoring, duration_minutes)
		VALUES ($1, $2, $3, $4)
	`, site.Id, user_id, site.Monitoring, site.DurationMinutes)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &site, nil
}

// Обновить данные сайта
func (d *InfoBotDb) MonitoringSiteUpdate(ctx context.Context, user_id int64, site *serializers.SiteSerializer) (*serializers.SiteSerializer, error) {
	sqlt := `
		UPDATE tg_user_sites
		SET monitoring = :monitoring,
			duration_minutes = :duration_minutes
		WHERE site_id = :site_id and tg_user_id = :user_id
	`
	_, err := d.NamedExecContext(ctx, sqlt, map[string]any{
		"site_id":          site.Id,
		"user_id":          user_id,
		"monitoring":       site.Monitoring,
		"duration_minutes": site.DurationMinutes,
	})
	if err != nil {
		return nil, err
	}
	return site, nil
}

// Удаление сайта из мониторинговой таблицы
func (d *InfoBotDb) MonitoringSiteDelete(ctx context.Context, user_id int64, site_id int) error {
	sqlt := `
		DELETE FROM tg_user_sites
		WHERE site_id = :site_id and tg_user_id = :user_id
	`
	_, err := d.NamedExecContext(ctx, sqlt, map[string]any{
		"site_id": site_id,
		"user_id": user_id,
	})

	if err != nil {
		return err
	}
	return err
}

// Вывести пользовательские обращения
func (d *InfoBotDb) Feedbacks(ctx context.Context, opt *infobotdb.OptionsInfoBot) ([]*serializers.FeedbackSerializer, int, error) {
	sqlt, err := infobotdb.Template("feedbacks", `
		SELECT 
			f.id, f.name, f.contact, f.message, f.created_at,
			COUNT(*) OVER (PARTITION BY f.site_id) AS count_feedbacks
		FROM feedbacks f
		LEFT JOIN sites s ON f.site_id = s.id
		WHERE 1=1 
			{{ if .Id }} AND f.id = :id {{ end }}
			{{ if .SiteId }} AND f.site_id = :site_id {{ end }}

		{{ if .Limit }}OFFSET :offset LIMIT :limit {{ end }}
	`, opt)
	if err != nil {
		return nil, 0, err
	}

	row, err := d.NamedQueryContext(ctx, sqlt, opt)
	if err != nil {
		return nil, 0, err
	}

	var cnt int
	res := make([]*serializers.FeedbackSerializer, 0, opt.Limit)
	for row.Next() {
		var s serializers.FeedbackSerializer
		err = row.Scan(
			&s.Id, &s.Name,
			&s.Contact, &s.Message,
			&s.CreatedAt, &cnt,
		)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, &s)
	}

	return res, cnt, nil
}

// Добавить пользовательское обращение
func (d *InfoBotDb) FeedbackInsert(ctx context.Context, user *serializers.FeedbackSerializer) error {
	sqlt := `
		INSERT INTO feedbacks (site_id, name, contact, message, created_at)
		VALUES (:site_id, :name, :contact, :message, :created_at)
	`
	_, err := d.NamedExecContext(ctx, sqlt, user)
	if err != nil {
		return err
	}
	return nil
}
