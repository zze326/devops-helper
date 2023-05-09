package filters

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"net/http"
	"strings"
)

var AuthenticationFilter = func(c *revel.Controller, fc []revel.Filter) {
	var jwtToken string
	if c.Request.Method == "WS" {
		jwtToken = c.Params.Get("token")
	} else {
		authHeader := c.Request.Header.Get("Authorization")
		if utils.IsEmpty(authHeader) {
			jwtToken = c.Params.Get("token")
			if utils.IsEmpty(jwtToken) {
				c.Response.SetStatus(http.StatusUnauthorized)
				c.Result = results.JsonError(fmt.Errorf("未登录"))
				return
			} else {
				authHeader = fmt.Sprintf("Bearer %s", jwtToken)
			}
		}

		jwtParts := strings.Split(authHeader, " ")
		if len(jwtParts) != 2 || jwtParts[0] != "Bearer" {
			c.Response.SetStatus(http.StatusUnauthorized)
			c.Result = results.JsonError(fmt.Errorf("token 格式错误"))
			return
		}
		jwtToken = jwtParts[1]
	}

	secret := revel.Config.StringDefault("jwt.secret", "")

	token, err := jwt.Parse(jwtToken, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Result = results.JsonError(fmt.Errorf("token 解析失败, err: %s", err.Error()))
		return
	}

	if !token.Valid {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Result = results.JsonError(fmt.Errorf("token 无效"))
		return
	}

	username := *utils.GetFromClaims[string]("username", token.Claims)
	userRealName := *utils.GetFromClaims[string]("user_real_name", token.Claims)
	userID := uint(*utils.GetFromClaims[float64]("user_id", token.Claims))

	c.Params.Add("_requestUsername", fmt.Sprintf("%s", username))
	c.Params.Add("_requestUserID", fmt.Sprintf("%d", userID))
	c.Params.Add("_requestUserRealName", fmt.Sprintf("%s", userRealName))
	fc[0](c, fc[1:]) // Execute the next filter stage.
}
