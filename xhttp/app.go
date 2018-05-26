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
	"github.com/gorilla/mux"
	"fmt"
)

type HandleFunc func(ctx Context) error

type App struct {
	listener    net.Listener
	server      *http.Server
	router      *mux.Router
	middlewares []negroni.Handler
	tlsConfig   *TLSConfig
	pool        sync.Pool
	errHandler  HandleFunc
	stop        chan bool
}

func New() (app *App) {
	app = &App{
		server: new(http.Server),
		router: &mux.Router{},
		stop:   make(chan bool),
	}

	app.pool.New = func() interface{} {
		return newContext()
	}

	return
}

func defaultErrHandler(ctx Context) error {
	// todo: add log
	res := ctx.Response()
	http.Error(res, "405 method not allowed", http.StatusMethodNotAllowed)
	return nil
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

func (app *App) wrapHandlerFunc(f HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.pool.Get().(Context)
		ctx.Reset(r, w)
		defer app.pool.Put(ctx)

		if err := f(ctx); err != nil {
			if h := app.errHandler; h != nil {
				h(ctx)
			} else {
				app.errHandler = defaultErrHandler
				defaultErrHandler(ctx)
			}
		}
	}
}

func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	})
}

func (app *App) Get(path string, f HandleFunc) {
	app.router.HandleFunc(path, app.wrapHandlerFunc(f)).Methods(GET)
}

func (app *App) Post(path string, f HandleFunc) {
	app.router.HandleFunc(path, app.wrapHandlerFunc(f)).Methods(POST)
}

func (app *App) initServer() {
	app.router.NotFoundHandler = http.NotFoundHandler()
	app.router.MethodNotAllowedHandler = MethodNotAllowedHandler()

	n := negroni.New()
	for _, m := range app.middlewares {
		n.Use(m)
	}

	n.UseHandler(app.router)
	app.server.Handler = n
}

func (app *App) ListenAndServe(addr string) {
	l, err := net.Listen(TCP, addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	app.listener = l

	app.initServer()

	err = app.server.Serve(l)
	fmt.Println(err)
}
