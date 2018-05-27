/*
 * Revision History:
 *     Initial: 2017/12/19        Feng Yifei
 */

package middleware

import (
	"github.com/urfave/negroni"
)

// NegroniLoggerHandler returns a logging handler.
func NegroniLoggerHandler() negroni.Handler {
	return negroni.NewLogger()
}
