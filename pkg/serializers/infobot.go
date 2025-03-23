package serializers

import (
	"sync"
	"time"

	"github.com/lib/pq"
)

// Телграмовские пользователь
type UserSerializer struct {
	// пользовательский тг ид
	UserId int64 `db:"user_id"`
	// username
	UserName string `db:"user_name"`
	// имя пользоователя
	FirstName string `db:"first_name"`
	// фамилия пользователя
	LastName string `db:"last_name"`
}

// Cайты для мониторинга
type SiteSerializer struct {
	// идентификатор сайта
	Id int `db:"id"`
	// ссылка на сайт
	Url string `db:"url"`
	// работает ли сайт
	Working bool `db:"working"`
	// код ответа в ходе проверок на доступность
	StatusCode int `db:"status_code"`
	// проверяется ли на достпность шедулером
	Monitoring bool `db:"monitoring"`
	// период в течении которого производится проверка
	DurationMinutes int `db:"duration_minutes"`
	// секретный ключ для оставления фидбеков от пользователей по gRPC
	SecretKey string `db:"secret_key"`
	// время последней проверки сайта на доступность
	LastCheckedAt time.Time `db:"last_checked_at"`
}

// Обращения от пользователей с обслуживаемых сайтов
type FeedbackSerializer struct {
	// идентифатор обращения
	Id int `db:"id"`
	// идентификатор сайта с которого оно пришло
	SiteId int `db:"site_id"`
	// Имя обратившегося пользователя
	Name string `db:"name"`
	// Контакты для обратной связи с ним
	Contact string `db:"contact"`
	// Дополнительное сообщение которое он написал
	Message string `db:"message"`
	// Страница с которой плител фидбек
	FeedbackUrl string `db:"feedback_url"`
	// Дата создания обращения
	CreatedAt time.Time `db:"created_at"`
}

// закэшированый пользователь
type User struct {
	sync.Mutex
	// его идентификатор
	userId int64
	// действие которое в данный момент выполняет
	action ACTION_TYPE
	// сайт над которым в данный момент производит работу
	actionSite *SiteSerializer

	// смещение для пагинационных сообщений
	offset int
}

// сайты у которых подошло время для проверки на доступность
type SiteForChecked struct {
	// идентификатор сайта
	Id int `db:"id"`
	// ссылка на сайт
	Url string `db:"url"`
	// работает ли сайт
	Working bool `db:"working"`
	// старый код ответа сайта
	StatusCode int `db:"status_code"`
	// пользователи слушающие сайт
	TgUsers pq.Int64Array `db:"tg_users"`
}

// Установить userId, пользователя
func (u *User) SetUserId(chatId int64) {
	u.Lock()
	defer u.Unlock()
	u.userId = chatId
}

// Получить userId, пользователя
func (u *User) GetUserId() int64 {
	u.Lock()
	defer u.Unlock()
	return u.userId
}

// Установить действие, которое в данный момент выполняет пользователь
func (u *User) SetAction(action ACTION_TYPE) {
	u.Lock()
	defer u.Unlock()
	u.action = action
}

// Получить действие, которое в данный момент выполняет пользователь
func (u *User) GetAction() ACTION_TYPE {
	u.Lock()
	defer u.Unlock()
	return u.action
}

// Установить сайт, над которым в данный момент работает пользователь
func (u *User) SetActionSite(actionSite *SiteSerializer) {
	u.Lock()
	defer u.Unlock()
	u.actionSite = actionSite
}

// Получить сайт, над которым в данный момент работает пользователь
func (u *User) GetActionSite() *SiteSerializer {
	u.Lock()
	defer u.Unlock()
	return u.actionSite
}

// Установить новые значения offset, limit
func (u *User) SetOffset(offset int) {
	u.Lock()
	defer u.Unlock()
	u.offset = offset
}

// Получить offset и limit, для отрисовки пагинации
func (u *User) GetOffset() int {
	u.Lock()
	defer u.Unlock()
	return u.offset
}

// закэшированные пользователи
type Users struct {
	sync.Mutex
	data map[int64]*User
}

// создать пользовательский кэш
func NewUsers() *Users {
	return &Users{
		data: make(map[int64]*User),
	}
}

// получить или создать и получить пользователя, если он не был закэширован
func (u *Users) Get(userId int64) *User {
	u.Lock()
	defer u.Unlock()

	user, ok := u.data[userId]
	if !ok {
		user = &User{}
		user.SetUserId(userId)
		u.data[userId] = user
	}
	return user
}
