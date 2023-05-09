package filters

import (
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	o_system "github.com/zze326/devops-helper/app/models/orm/system"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"net/http"
	"strconv"
	"strings"
)

var TerminalAuthFilter = func(c *revel.Controller, fc []revel.Filter) {
	if !utils.InSlice[string](strings.ToLower(c.Request.URL.Path), []string{
		"/host/testssh",
		"/host/terminal",
	}) {
		fc[0](c, fc[1:])
		return
	}
	requetUsername := c.Params.Get("_requestUsername")
	requestUserID, err := strconv.Atoi(c.Params.Get("_requestUserID"))
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		c.Result = results.JsonError(err)
		return
	}

	if requetUsername == "admin" {
		fc[0](c, fc[1:])
		return
	}

	hostID, err := strconv.Atoi(c.Params.Get("id"))
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		c.Result = results.JsonError(err)
		return
		//c.Result = results.JsonErrorMsg("您没有权限访问该主机")
	}

	hasPermission, err := checkUserTerminalPermission(hostID, requestUserID)
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		c.Result = results.JsonError(err)
		return
	}

	if !hasPermission {
		c.Response.SetStatus(http.StatusForbidden)
		c.Result = results.JsonError(err)
		return
	}

	fc[0](c, fc[1:])
}

func checkUserTerminalPermission(hostID, userID int) (bool, error) {
	hostModel := new(o_resource.Host)
	if err := g.DB.Preload("HostGroups.HostGroupPermissions").First(hostModel, hostID).Error; err != nil {
		return false, err
	}

	userModel := new(o_system.User)
	if err := g.DB.Preload("Roles").First(userModel, userID).Error; err != nil {
		return false, err
	}

	var roleIDs []int
	for _, role := range userModel.Roles {
		roleIDs = append(roleIDs, role.ID)
	}

	for _, hostGroup := range hostModel.HostGroups {
		for _, permission := range hostGroup.HostGroupPermissions {
			if permission.Type == 1 && permission.RefID == userID {
				return true, nil
			} else if permission.Type == 2 && utils.InSlice(permission.RefID, roleIDs) {
				return true, nil
			}
		}
	}

	return false, nil
}
