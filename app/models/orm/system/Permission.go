package o_system

import (
	"github.com/zze326/devops-helper/app/modules/gormc"
	"gorm.io/gorm"
)

type Permission struct {
	gormc.Model
	Code           string           `gorm:"type:varchar(48);comment:权限编码;" json:"code"`
	Name           string           `gorm:"comment:权限名称" json:"name"`
	ParentID       int              `gorm:"comment:父级权限ID" json:"parent_id"`
	Parent         *Permission      `gorm:"foreignKey:ParentID" json:"parent"`
	Children       []*Permission    `gorm:"foreignKey:ParentID" json:"children"`
	Menus          []*Menu          `gorm:"many2many:permission_menu;" json:"menu"`
	FrontendRoutes []*FrontendRoute `gorm:"many2many:permission_frontend_route;" json:"frontend_routes"`
	BackendRoutes  []*BackendRoute  `gorm:"many2many:permission_backend_route;" json:"backend_routes"`
}

func (Permission) TableName() string {
	return "permission"
}

func (Permission) ListTree(db *gorm.DB, preloadMenus bool, preloadFrontendRoutes bool, preloadBackendRoutes bool) ([]*Permission, map[int]*Permission, error) {
	var permissions []*Permission

	scopes := db
	if preloadMenus {
		// 预加载菜单
		scopes = db.Preload("Menus", "enabled = ?", true)
	}

	if preloadFrontendRoutes {
		// 预加载前端路由
		scopes = db.Preload("FrontendRoutes")
	}

	if preloadBackendRoutes {
		// 预加载后端路由
		scopes = db.Preload("BackendRoutes")
	}

	// 从数据库中获取所有权限记录
	if err := scopes.Find(&permissions).Error; err != nil {
		return nil, nil, err
	}

	// 构建 ID 到权限映射的 map
	permissionMap := make(map[int]*Permission, len(permissions))
	for _, permission := range permissions {
		permissionMap[permission.ID] = permission
	}

	// 找出根节点并返回
	var roots []*Permission
	// 对每个权限，将其作为子节点添加到对应父节点的 Children 列表中
	for _, permission := range permissions {
		parent, ok := permissionMap[permission.ParentID]
		if !ok {
			// 如果该权限没有父节点，则认为它是根节点
			roots = append(roots, permission)
			continue
		}
		parent.Children = append(parent.Children, permission)
	}

	return roots, permissionMap, nil
}
