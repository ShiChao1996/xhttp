/*
 * Revision History:
 *     Initial: 2018/05/25        ShiChao
 */

package xhttp

import (
	"net"
	"net/http"
	"fmt"
	"context"
	"sync"
	"time"

	"github.com/urfave/negroni"
)

type HandleFunc func(ctx Context) error
type FilterFunc func(ctx Context) bool

type Server struct {
	listener    net.Listener
	server      *http.Server
	router      *Router
	middlewares []negroni.Handler
	pool        sync.Pool
	errHandler  HandleFunc
	stop        chan bool
}

func New() (server *Server) {
	server = &Server{
		server: new(http.Server),
		router: NewRouter(),
		stop:   make(chan bool),
	}

	server.pool.New = func() interface{} {
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

func (server *Server) Use(fn negroni.Handler) {
	server.middlewares = append(server.middlewares, fn)
}

func (server *Server) Get(path string, f HandleFunc, filters ...FilterFunc) {
	server.router.Get(path, server.wrapHandlerFunc(f, ))
}

func (server *Server) Post(path string, f HandleFunc, filters ...FilterFunc) {
	server.router.Post(path, server.wrapHandlerFunc(f))
}

func (server *Server) wrapHandlerFunc(f HandleFunc, filters ...FilterFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := server.pool.Get().(Context)
		ctx.Reset(r, w)
		defer server.pool.Put(ctx)

		for _, filter := range filters {
			if pass := filter(ctx); !pass {
				return
			}
		}

		if err := f(ctx); err != nil {
			if h := server.errHandler; h != nil {
				h(ctx)
			} else {
				server.errHandler = defaultErrHandler
				defaultErrHandler(ctx)
			}
		}
	}
}

func (server *Server) init() {
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	for _, m := range server.middlewares {
		n.Use(m)
	}

	n.UseHandler(server.router.router)

	server.server.Handler = n
}

func (server *Server) Run(addr string) {
	l, err := net.Listen(TCP, addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.listener = l

	server.init()

	err = server.server.Serve(l)
	fmt.Println(err)
}

func (server *Server) RunTLS(addr string, certFile, keyFile string) {
	l, err := net.Listen(TCP, addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.listener = l

	server.init()

	err = server.server.ServeTLS(l, certFile, keyFile)
	fmt.Println(err)
}

func (server *Server) GracefulClose() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	if err := server.server.Shutdown(ctx); err != nil {
		server.server.Close()
	}
	cancel()
}
