/*
 * Revision History:
 *     Initial: 2017/12/19        Feng Yifei
 */

package middleware

import (
	"net/http"

	"github.com/fengyfei/gu/libs/logger"
	"github.com/urfave/negroni"
)

// NegroniRecoverHandler returns a handler for recover from a http request.
func NegroniRecoverHandler() negroni.Handler {
	fn := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer recoverFunc(w)
		next.ServeHTTP(w, r)
	}
	return negroni.HandlerFunc(fn)
}

func recoverFunc(w http.ResponseWriter) {
	if err := recover(); err != nil {
		logger.Error("Recovered from panic in http handler:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
