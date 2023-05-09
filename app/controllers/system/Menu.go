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

type Menu struct {
	gormc.Controller
}

// Add 创建菜单
func (c Menu) Add(req v_system.AddMenuReq) revel.Result {
	menuModel := new(o_system.Menu)
	if err := copier.Copy(menuModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.DB.First(menuModel, "name = ?", req.Name).Error
	if err == nil {
		return results.JsonErrorMsg("菜单名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if err = c.DB.Create(menuModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 获取菜单详情
func (c Menu) Get(id int) revel.Result {
	menuModel := new(o_system.Menu)
	if err := c.DB.First(menuModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(menuModel)
}

// ListTree 菜单树
func (c Menu) ListTree() revel.Result {
	var menuModels []*o_system.Menu
	if err := c.DB.Preload("Children", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
				return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
					return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
						return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
							return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
								return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
									return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
										return db.Order("sort asc").Preload("Children", func(db *gorm.DB) *gorm.DB {
											return db.Order("sort asc")
										})
									})
								})
							})
						})
					})
				})
			})
		})
	}).Where("(parent_id = 0 or parent_id is null)").Order("sort").Find(&menuModels).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(menuModels)
}

// Edit 编辑菜单
func (c Menu) Edit(req v_system.UpdateMenuReq) revel.Result {
	menuModel := new(o_system.Menu)
	if err := c.DB.First(menuModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	menuModel.Name = req.Name
	menuModel.Icon = req.Icon
	menuModel.RouteID = req.RouteID
	menuModel.ParentID = req.ParentID
	menuModel.Sort = req.Sort
	menuModel.Enabled = req.Enabled
	if err := c.DB.Save(menuModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// UpdateStatus 更新菜单状态 启用/禁用
func (c Menu) UpdateStatus(id int, enabled bool) revel.Result {
	menuModel := new(o_system.Menu)
	if err := c.DB.First(menuModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	menuModel.Enabled = enabled
	if err := c.DB.Save(menuModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除菜单
func (c Menu) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_system.Menu{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}
