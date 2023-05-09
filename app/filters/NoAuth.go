package filters

import (
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/utils"
	"strings"
)

var NoAuthFilter = func(c *revel.Controller, fc []revel.Filter) {
	if utils.InSlice[string](strings.ToLower(c.Request.URL.Path), []string{
		"/app/login",
		"/app/ping",
	}) {
		fc[2](c, fc[3:]) // 跳过认证和授权 2 个 filter
		return
	}
	fc[0](c, fc[1:])
}
