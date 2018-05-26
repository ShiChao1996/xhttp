/*
 * Revision History:
 *     Initial: 2018/05/25        ShiChao
 */

package xhttp

import (
	"github.com/gorilla/mux"
)

// Router register routes to be matched and dispatched to a handler.
type Router struct {
	router     *mux.Router
	errHandler func(*Context)
}
