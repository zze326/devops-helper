package orm

import (
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
		&o_resource.Host{},
		&o_resource.HostGroup{},
		&o_resource.HostGroupPermission{},
		&o_resource.HostTerminalSession{},
	}
}
