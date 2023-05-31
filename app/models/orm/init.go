package orm

import (
	o_ci "github.com/zze326/devops-helper/app/models/orm/ci"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	o_system "github.com/zze326/devops-helper/app/models/orm/system"
)

func ModelTypesToMigrate() []interface{} {
	return []interface{}{
		&o_system.User{},
		&o_system.FrontendRoute{},
		&o_system.BackendRoute{},
		&o_system.Menu{},
		&o_system.Permission{},
		&o_system.Role{},
		&o_system.DataDict{},
		&o_system.DataDictItem{},
		&o_resource.Host{},
		&o_resource.HostGroup{},
		&o_resource.HostGroupPermission{},
		&o_resource.HostTerminalSession{},
		&o_resource.Secret{},
		&o_ci.Pipeline{},
		&o_ci.EnvRef{},
		&o_ci.Env{},
		&o_ci.Stage{},
		&o_ci.Task{},
	}
}
