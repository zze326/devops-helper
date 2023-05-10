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

type BackendRoute struct {
	gormc.Controller
}

// Add 创建路由
func (c BackendRoute) Add(req v_system.AddBackendRouteReq) revel.Result {
	backendRouteModel := new(o_system.BackendRoute)
	if err := copier.Copy(backendRouteModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.DB.First(backendRouteModel, "name = ?", req.Name).Error
	if err == nil {
		return results.JsonErrorMsg("路由名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	err = c.DB.First(backendRouteModel, "path = ?", req.Path).Error
	if err == nil {
		return results.JsonErrorMsg("路由路径已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if err = c.DB.Create(backendRouteModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Edit 编辑路由
func (c BackendRoute) Edit(req v_system.EditBackendRouteReq) revel.Result {
	routeModel := new(o_system.BackendRoute)
	if err := c.DB.First(routeModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	routeModel.Name = req.Name
	routeModel.Path = req.Path
	if err := c.DB.Save(routeModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除路由
func (c BackendRoute) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_system.BackendRoute{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 根据 ID 获取
func (c BackendRoute) Get(id int) revel.Result {
	routeModel := new(o_system.BackendRoute)
	if err := c.DB.First(routeModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(routeModel)
}

// ListPage 分页查询
func (c BackendRoute) ListPage(pager *utils.Pager) revel.Result {
	var backendRouteModels []*o_system.BackendRoute
	pager.Order = "path asc,id desc"
	total, err := utils.Paginate[o_system.BackendRoute](c.DB, pager, &backendRouteModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: backendRouteModels,
	})
}

// ListAll 查询所有
func (c BackendRoute) ListAll() revel.Result {
	var backendRouteModels []*o_system.BackendRoute
	if err := c.DB.Order("path asc").Find(&backendRouteModels).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(backendRouteModels)
}
