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

			if userID := util.GetXUserID(ctx); userID != 0 {
				req.Header.Set(X_USER_ID, fmt.Sprintf("%d", userID))
			}
			if email := util.GetXEmail(ctx); email != "" {
				req.Header.Set(X_EMAIL, email)
			}

		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
