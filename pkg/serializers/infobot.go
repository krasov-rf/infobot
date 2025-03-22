package serializers

import (
	"sync"
	"time"
)

type UserSerializer struct {
	UserId    int64  `db:"user_id"`
	UserName  string `db:"user_name"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type SiteSerializer struct {
	Id              int    `db:"id"`
	Url             string `db:"url"`
	StatusCode      int    `db:"status_code"`
	Monitoring      bool   `db:"monitoring"`
	DurationMinutes int    `db:"duration_minutes"`
}

type ConversionSerializer struct {
	Id        int       `db:"id"`
	SiteId    int       `db:"site_id"`
	Name      string    `db:"name"`
	Contact   string    `db:"contact"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type User struct {
	sync.Mutex
	chatId     int64
	action     ACTION_TYPE
	actionSite *SiteSerializer

	offset int
}

// Установить userId (он же chatId), пользователя
func (u *User) SetChatId(chatId int64) {
	u.Lock()
	defer u.Unlock()
	u.chatId = chatId
}

// Получить userId (он же chatId), пользователя
func (u *User) GetChatId() int64 {
	u.Lock()
	defer u.Unlock()
	return u.chatId
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

type Users struct {
	sync.Mutex
	data map[int64]*User
}

func NewUsers() *Users {
	return &Users{
		data: make(map[int64]*User),
	}
}

func (u *Users) Get(userId int64) *User {
	u.Lock()
	defer u.Unlock()

	user, ok := u.data[userId]
	if !ok {
		user = &User{}
		user.SetChatId(userId)
		u.data[userId] = user
	}
	return user
}
