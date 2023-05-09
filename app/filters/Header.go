package filters

import (
	"github.com/revel/revel"
	"net/http"
)

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// 允许跨域
	c.Response.Out.Header().Add("Access-Control-Allow-Origin", "*")
	c.Response.Out.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Response.Out.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if c.Request.Method == "OPTIONS" {
		c.Response.SetStatus(http.StatusOK)
		return
	}
	fc[0](c, fc[1:]) // Execute the next filter stage.
}
