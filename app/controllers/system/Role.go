package system

import (
	"github.com/jinzhu/copier"
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/models/orm/system"
	"github.com/zze326/devops-helper/app/models/view/system"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"gorm.io/gorm"
)

type Role struct {
	gormc.TxnController
}

// Add 创建角色
func (c Role) Add(req v_system.AddRoleReq) revel.Result {
	roleModel := new(o_system.Role)
	if err := copier.Copy(roleModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.Txn.First(roleModel, "name = ?", req.Name).Error
	if err == nil {
		return results.JsonErrorMsg("角色名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	err = c.Txn.First(roleModel, "code = ?", req.Code).Error
	if err == nil {
		return results.JsonErrorMsg("角色代码已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	var permissionModels []*o_system.Permission
	if err := c.Txn.Find(&permissionModels, req.PermissionIDs).Error; err != nil {
		return results.JsonError(err)
	}

	roleModel.Permissions = permissionModels

	if err = c.Txn.Create(roleModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Edit 编辑角色
func (c Role) Edit(req v_system.EditRoleReq) revel.Result {
	roleModel := new(o_system.Role)
	if err := c.Txn.First(roleModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	roleModel.Name = req.Name
	roleModel.Code = req.Code

	var permissionModels []*o_system.Permission
	if err := c.Txn.Find(&permissionModels, req.PermissionIDs).Error; err != nil {
		return results.JsonError(err)
	}

	if err := c.Txn.Model(roleModel).Association("Permissions").Replace(permissionModels); err != nil {
		return results.JsonError(err)
	}

	if err := c.Txn.Updates(roleModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除角色
func (c Role) Delete(id int) revel.Result {
	if err := c.Txn.Delete(&o_system.Role{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 根据 ID 获取
func (c Role) Get(id int) revel.Result {
	roleModel := new(o_system.Role)
	if err := c.Txn.Preload("Permissions").First(roleModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(roleModel)
}

// ListPage 分页查询
func (c Role) ListPage(pager *utils.Pager) revel.Result {
	var roleModels []*o_system.Role
	pager.Order = "id desc"
	total, err := utils.Paginate[o_system.Role](c.Txn, pager, &roleModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: roleModels,
	})
}

// ListAll 查询所有
func (c Role) ListAll() revel.Result {
	var roleModels []*o_system.Role
	if err := c.Txn.Find(&roleModels).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(roleModels)
}
