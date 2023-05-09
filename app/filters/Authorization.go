package filters

import (
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	"github.com/zze326/devops-helper/app/results"
	"net/http"
	"strings"
)

var AuthorizationFilter = func(c *revel.Controller, fc []revel.Filter) {
	requestPath := strings.ToLower(c.Request.GetPath())
	requestUser := c.Params.Get("_requestUsername")
	result, err := g.Enforcer.Enforce(requestUser, strings.ToLower(requestPath), "ANY")
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		c.Result = results.JsonError(err)
		return
	}

	if !result {
		c.Response.SetStatus(http.StatusForbidden)
		c.Result = results.JsonErrorMsg("没有权限")
		return
	}

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
