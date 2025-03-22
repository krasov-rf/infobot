package infobot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UpdateType int

const (
	UPDATE_MESSAGE UpdateType = iota
	UPDATE_CALLBACK
)

type Router struct {
	routes     map[UpdateType][]route
	middleware []func(handlerFunc) handlerFunc
}

type route struct {
	path    string
	handler handlerFunc
	params  []string
}

type handlerFunc func(BotContext, tgbotapi.Update)

func NewRouter() *Router {
	return &Router{
		routes: make(map[UpdateType][]route),
	}
}

func (r *Router) RouteMessage(path string, handler handlerFunc) {
	r.addUpdate(UPDATE_MESSAGE, path, handler)
}

func (r *Router) RouteCallback(path string, handler handlerFunc) {
	r.addUpdate(UPDATE_CALLBACK, path, handler)
}

func (r *Router) addUpdate(method UpdateType, path string, handler handlerFunc) {
	parts := strings.Split(path, "/")
	var params []string
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			params = append(params, part[1:])
		}
	}

	r.routes[method] = append(r.routes[method], route{
		path:    path,
		handler: handler,
		params:  params,
	})
}

func (r *Router) Use(middleware ...func(handlerFunc) handlerFunc) {
	r.middleware = append(r.middleware, middleware...)
}

func (r *Router) Run(updates tgbotapi.UpdatesChannel) error {
	return nil
}

func (r *Router) handleUpdate(update tgbotapi.Update) {
	var (
		route  string
		ctx    context.Context
		botCtx BotContext
	)

	ctx = context.Background()

	// обработчики нажатий на кнопку

	if update.CallbackQuery != nil {
		data := strings.Split(update.CallbackQuery.Data, "|")
		route = data[0]
		if len(data) == 2 {
			d := data[1]
			ctx = context.WithValue(ctx, CTX_KEY_DATA, d)
		}
		botCtx.Context = ctx

		for _, ex_route := range r.routes[UPDATE_CALLBACK] {
			if ex_route.path == route {
				r.applyMiddleware(ex_route.handler)(botCtx, update)
				return
			}
		}
		return
	}

	// обработчики сообщений

	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		route = update.Message.Command()
	} else {
		route = update.Message.Text
	}

	for _, ex_route := range r.routes[UPDATE_MESSAGE] {
		if ex_route.path == route {
			f := r.applyMiddleware(ex_route.handler)
			f(botCtx, update)
			return
		}
	}
}

func (r *Router) applyMiddleware(handler handlerFunc) handlerFunc {
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	return handler
}
