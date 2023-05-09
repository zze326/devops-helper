package startups

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/zze326/devops-helper/app/g"
	o_system "github.com/zze326/devops-helper/app/models/orm/system"
	"strings"
	"time"
)

func initCasbin() {
	go func() {
		var firstLoop = true
		for {
			if err := RefreshCasbin(); err != nil {
				if firstLoop {
					g.Logger.Fatal(err.Error())
				} else {
					g.Logger.Error(err.Error())
				}
			}
			firstLoop = false
			time.Sleep(2 * time.Minute)
		}
	}()
}

func RefreshCasbin() error {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act) || g(r.sub, "admin")`)

	enforcer, err := newCasbinEnforcer(m)
	if err != nil {
		return err
	}
	g.Enforcer = enforcer
	g.Logger.Infof("Casbin enforcer updated")
	return nil
}

func newCasbinEnforcer(m model.Model) (*casbin.Enforcer, error) {
	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, err
	}

	_, permissionMap, err := o_system.Permission{}.ListTree(g.DB, false, false, true)
	if err != nil {
		return nil, err
	}

	defaultPermission := new(o_system.Permission)
	if err := g.DB.First(defaultPermission, "code = ?", "system-default").Error; err != nil {
		return nil, err
	}

	var roles []*o_system.Role
	if err := g.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}

	for _, role := range roles {
		role.Permissions = append(role.Permissions, defaultPermission)
		if err := loopAddPolicy(enforcer, role.Code, role.Permissions, permissionMap); err != nil {
			return nil, err
		}
	}

	var users []*o_system.User
	if err := g.DB.Preload("Roles").Find(&users).Error; err != nil {
		return nil, err
	}

	for _, user := range users {
		var roles []string
		for _, role := range user.Roles {
			roles = append(roles, role.Code)
		}
		if _, err := enforcer.AddRolesForUser(user.Username, roles); err != nil {
			return nil, err
		}
		g.Logger.Debugf("Casbin roles added: %s, %v", user.Username, roles)
	}
	return enforcer, nil
}

func loopAddPolicy(enforcer *casbin.Enforcer, roleCode string, permissions []*o_system.Permission, permissionMap map[int]*o_system.Permission) error {
	for _, permission := range permissions {
		tmpPermission := permissionMap[permission.ID]
		for _, backendRoute := range tmpPermission.BackendRoutes {
			if enforcer.HasPolicy(roleCode, backendRoute.Path, "ANY") {
				continue
			}
			if _, err := enforcer.AddPolicy(roleCode, strings.ToLower(backendRoute.Path), "ANY"); err != nil {
				return err
			}
			g.Logger.Debugf("Casbin policy added: %s, %s, %s", roleCode, backendRoute.Path, "ANY")
		}

		if len(tmpPermission.Children) > 0 {
			if err := loopAddPolicy(enforcer, roleCode, tmpPermission.Children, permissionMap); err != nil {
				return err
			}
		}
	}
	return nil
}
