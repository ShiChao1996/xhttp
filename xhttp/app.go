/*
 * Revision History:
 *     Initial: 2018/05/25        ShiChao
 */

package xhttp

import (
	"net/http"
	"sync"
	"net"
	"github.com/urfave/negroni"
)

type App struct {
	listener    *net.Listener
	server      *http.Server
	router      *Router
	middlewares []negroni.Handler
	config      *Config
	tlsConfig   *TLSConfig
	pool        sync.Pool
	stop        chan bool
}

func New() (app *App) {
	app = &App{
		server: new(http.Server),
		router: new(Router),
		stop:   make(chan bool),
	}

	app.pool.New = func() interface{} {
		return newContext()
	}

	return
}

func (app *App) GetContext() Context {
	return app.pool.Get().(Context)
}

func (app *App) ReleaseContext(ctx context) {
	app.pool.Put(ctx)
}

func (app *App) Use(fn negroni.Handler) {
	app.middlewares = append(app.middlewares, fn)
}
