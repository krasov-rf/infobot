package infobotdb

type OptionsInfoBot struct {
	Id     int    `db:"id"`
	UserId int64  `db:"user_id"`
	SiteId int    `db:"site_id"`
	Domain string `db:"domain"`

	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

func WithId(id int) InfoBotOptionFunc {
	return func(cs *OptionsInfoBot) {
		cs.Id = id
	}
}
func WithUserId(user_id int64) InfoBotOptionFunc {
	return func(cs *OptionsInfoBot) {
		cs.UserId = user_id
	}
}
func WithSiteId(site_id int) InfoBotOptionFunc {
	return func(cs *OptionsInfoBot) {
		cs.SiteId = site_id
	}
}
func WithDomain(domain string) InfoBotOptionFunc {
	return func(cs *OptionsInfoBot) {
		cs.Domain = domain
	}
}
func WithOffset(offset int) InfoBotOptionFunc {
	return func(cs *OptionsInfoBot) {
		cs.Offset = offset
	}
}

type InfoBotOptionFunc func(*OptionsInfoBot)

func NewInfoBotOptions(options ...InfoBotOptionFunc) *OptionsInfoBot {

	service := &OptionsInfoBot{
		Limit: QUERY_LIMIT,
	}
	for _, option := range options {
		option(service)
	}
	return service
}
