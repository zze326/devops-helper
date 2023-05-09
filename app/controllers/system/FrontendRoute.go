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

type FrontendRoute struct {
	gormc.Controller
}

// Add 创建路由
func (c FrontendRoute) Add(req v_system.AddFrontendRouteReq) revel.Result {
	frontendRouteModel := new(o_system.FrontendRoute)
	if err := copier.Copy(frontendRouteModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.DB.First(frontendRouteModel, "name = ?", req.Name).Error
	if err == nil {
		return results.JsonErrorMsg("路由名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	err = c.DB.First(frontendRouteModel, "route_path = ?", req.RoutePath).Error
	if err == nil {
		return results.JsonErrorMsg("路由路径已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if err = c.DB.Create(frontendRouteModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Edit 编辑路由
func (c FrontendRoute) Edit(req v_system.EditFrontendRouteReq) revel.Result {
	routeModel := new(o_system.FrontendRoute)
	if err := c.DB.First(routeModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	routeModel.RoutePath = req.RoutePath
	routeModel.Name = req.Name
	routeModel.ComponentPath = req.ComponentPath
	routeModel.Title = req.Title
	if err := c.DB.Save(routeModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除路由
func (c FrontendRoute) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_system.FrontendRoute{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 根据 ID 获取
func (c FrontendRoute) Get(id int) revel.Result {
	routeModel := new(o_system.FrontendRoute)
	if err := c.DB.First(routeModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(routeModel)
}

// ListPage 分页查询
func (c FrontendRoute) ListPage(pager *utils.Pager) revel.Result {
	var frontendRouteModels []*o_system.FrontendRoute
	pager.Order = "id desc"
	total, err := utils.Paginate[o_system.FrontendRoute](c.DB, pager, &frontendRouteModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: frontendRouteModels,
	})
}

// ListAll 查询所有
func (c FrontendRoute) ListAll() revel.Result {
	var frontendRouteModels []*o_system.FrontendRoute
	if err := c.DB.Find(&frontendRouteModels).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(frontendRouteModels)
}
