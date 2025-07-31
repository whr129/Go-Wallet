package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	util "github.com/whr129/go-wallet/pkg/util"
)

const (
	X_USER_ID = "X-User-ID"
	X_EMAIL   = "X-Email"
	X_ROLE    = "X-Role"
)

func NewReverseProxy(target string) gin.HandlerFunc {
	url, err := url.Parse(target)

	if err != nil {
		panic(fmt.Sprintf("failed to parse target URL %s: %v", target, err))
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(ctx *gin.Context) {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				ctx.Error(err)
				ctx.Abort()
			}
		}()

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = ctx.Param("path")
			req.Header = ctx.Request.Header

			if userID, err := util.GetXUserID(ctx); userID != 0 && err {
				req.Header.Set(X_USER_ID, fmt.Sprintf("%d", userID))
			}
			if email, err := util.GetXEmail(ctx); email != "" && !err {
				req.Header.Set(X_EMAIL, email)
			}
			if role, err := util.GetXRole(ctx); role != "" && !err {
				req.Header.Set(X_ROLE, role)
			}

		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
