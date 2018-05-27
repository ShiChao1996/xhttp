/*
 * Revision History:
 *     From: https://github.com/labstack/echo/blob/master/middleware/middleware.go
 *     Modify:  2017/12/19        Jia Chenhui
 */

package middleware

import (
	"github.com/fengyfei/gu/libs/http/server"
)

// Skipper defines a function to skip middleware. Returning true skips processing
// the middleware.
type Skipper func(*server.Context) bool

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(*server.Context) bool {
	return false
}
