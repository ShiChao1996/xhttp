/*
 * Revision History:
 *     Initial: 2018/05/27        ShiChao
 */

package xhttp

import (
	"net/http"
	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

func NewRouter() (r *Router) {
	r = &Router{
		router: &mux.Router{},
	}
	r.router.NotFoundHandler = http.NotFoundHandler()
	r.router.MethodNotAllowedHandler = MethodNotAllowedHandler()

	return
}

func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	})
}

func (router *Router) Get(path string, f http.HandlerFunc) {
	router.router.HandleFunc(path, f).Methods(GET)
}

func (router *Router) Post(path string, f http.HandlerFunc) {
	router.router.HandleFunc(path, f).Methods(POST)
}
