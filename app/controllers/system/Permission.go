package system

import (
	"github.com/jinzhu/copier"
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/models/orm/system"
	"github.com/zze326/devops-helper/app/models/view/system"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"gorm.io/gorm"
)

type Permission struct {
	gormc.TxnController
}

// Add 添加权限
func (c Permission) Add(req v_system.AddPermissionReq) revel.Result {
	permissionModel := new(o_system.Permission)
	if err := copier.Copy(permissionModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.Txn.First(permissionModel, "name = ?", req.Name).Error
	if err == nil {
		return results.JsonErrorMsg("权限名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	err = c.Txn.First(permissionModel, "code = ?", req.Code).Error
	if err == nil {
		return results.JsonErrorMsg("权限代码已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if len(req.ChildrenIDs) > 0 {
		var children []*o_system.Permission
		if err := c.Txn.Find(&children, req.ChildrenIDs).Error; err != nil {
			return results.JsonError(err)
		}
		permissionModel.Children = children
	}

	if len(req.MenuIDs) > 0 {
		var menus []*o_system.Menu
		if err := c.Txn.Find(&menus, req.MenuIDs).Error; err != nil {
			return results.JsonError(err)
		}
		permissionModel.Menus = menus
	}

	if len(req.BackendRouteIDs) > 0 {
		var backendRoutes []*o_system.BackendRoute
		if err := c.Txn.Find(&backendRoutes, req.BackendRouteIDs).Error; err != nil {
			return results.JsonError(err)
		}
		permissionModel.BackendRoutes = backendRoutes
	}

	if len(req.FrontendRouteIDs) > 0 {
		var frontendRoutes []*o_system.FrontendRoute
		if err := c.Txn.Find(&frontendRoutes, req.FrontendRouteIDs).Error; err != nil {
			return results.JsonError(err)
		}
		permissionModel.FrontendRoutes = frontendRoutes
	}

	if err = c.Txn.Create(permissionModel).Error; err != nil {
		return results.JsonError(err)
	}

	return results.JsonOk()
}

// Edit 编辑权限
func (c Permission) Edit(req v_system.EditPermissionReq) revel.Result {
	permissionModel := new(o_system.Permission)
	err := c.Txn.First(permissionModel, "id != ? and name = ?", req.ID, req.Name).Error
	if err == nil {
		return results.JsonErrorMsg("权限名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	err = c.Txn.First(permissionModel, "id != ? and code = ?", req.ID, req.Code).Error
	if err == nil {
		return results.JsonErrorMsg("权限代码已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if err := c.Txn.First(permissionModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	permissionModel.Name = req.Name
	permissionModel.Code = req.Code
	permissionModel.ParentID = req.ParentID

	var children []*o_system.Permission
	if len(req.ChildrenIDs) > 0 {
		if err := c.Txn.Find(&children, req.ChildrenIDs).Error; err != nil {
			return results.JsonError(err)
		}
	}
	if err := c.Txn.Model(permissionModel).Association("Children").Replace(children); err != nil {
		return results.JsonError(err)
	}

	var menus []*o_system.Menu
	if len(req.MenuIDs) > 0 {
		if err := c.Txn.Find(&menus, req.MenuIDs).Error; err != nil {
			return results.JsonError(err)
		}
	}
	if err := c.Txn.Model(permissionModel).Association("Menus").Replace(menus); err != nil {
		return results.JsonError(err)
	}

	var backendRoutes []*o_system.BackendRoute
	if len(req.BackendRouteIDs) > 0 {
		if err := c.Txn.Find(&backendRoutes, req.BackendRouteIDs).Error; err != nil {
			return results.JsonError(err)
		}
	}
	if err := c.Txn.Model(permissionModel).Association("BackendRoutes").Replace(backendRoutes); err != nil {
		return results.JsonError(err)
	}

	var frontendRoutes []*o_system.FrontendRoute
	if len(req.FrontendRouteIDs) > 0 {
		if err := c.Txn.Find(&frontendRoutes, req.FrontendRouteIDs).Error; err != nil {
			return results.JsonError(err)
		}
	}
	if err := c.Txn.Model(permissionModel).Association("FrontendRoutes").Replace(frontendRoutes); err != nil {
		return results.JsonError(err)
	}

	if err = c.Txn.Updates(permissionModel).Error; err != nil {
		return results.JsonError(err)
	}

	return results.JsonOk()
}

// ListAllTop 获取所有顶级权限
func (c Permission) ListAllTop() revel.Result {
	var permissions []*o_system.Permission
	if err := c.Txn.Find(&permissions, "parent_id = 0").Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(permissions)
}

// ListTree 获取权限树
func (c Permission) ListTree() revel.Result {
	roots, _, err := o_system.Permission{}.ListTree(c.Txn, false, false, false)
	if err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(roots)
}

// Get 获取权限
func (c Permission) Get(id int) revel.Result {
	var permission o_system.Permission
	if err := c.Txn.Preload("Children").Preload("Menus").Preload("FrontendRoutes").Preload("BackendRoutes").First(&permission, id).Error; err != nil {
		return results.JsonError(err)
	}

	var (
		childrenIDs      []int
		menuIDs          []int
		frontendRouteIDs []int
		backendRouteIDs  []int
	)

	for _, child := range permission.Children {
		childrenIDs = append(childrenIDs, child.ID)
	}

	for _, menu := range permission.Menus {
		menuIDs = append(menuIDs, menu.ID)
	}

	for _, frontendRoute := range permission.FrontendRoutes {
		frontendRouteIDs = append(frontendRouteIDs, frontendRoute.ID)
	}

	for _, backendRoute := range permission.BackendRoutes {
		backendRouteIDs = append(backendRouteIDs, backendRoute.ID)
	}
	resp := new(v_system.GetPermissionResp)
	if err := copier.Copy(&resp, &permission); err != nil {
		return results.JsonError(err)
	}
	resp.Children = permission.Children
	resp.ChildrenIDs = childrenIDs
	resp.MenuIDs = menuIDs
	resp.FrontendRouteIDs = frontendRouteIDs
	resp.BackendRouteIDs = backendRouteIDs
	return results.JsonOkData(resp)
}

// Delete 删除权限
func (c Permission) Delete(id int) revel.Result {
	permissionModel := new(o_system.Permission)

	if err := c.Txn.First(permissionModel, id).Error; err != nil {
		return results.JsonError(err)
	}

	if err := c.Txn.Model(permissionModel).Association("Children").Clear(); err != nil {
		return results.JsonError(err)
	}
	if err := c.Txn.Model(permissionModel).Association("Menus").Clear(); err != nil {
		return results.JsonError(err)
	}
	if err := c.Txn.Model(permissionModel).Association("FrontendRoutes").Clear(); err != nil {
		return results.JsonError(err)
	}
	if err := c.Txn.Model(permissionModel).Association("BackendRoutes").Clear(); err != nil {
		return results.JsonError(err)
	}

	if err := c.Txn.Delete(&o_system.Permission{}, id).Error; err != nil {
		return results.JsonError(err)
	}

	return results.JsonOk()
}
