package system

import (
	"github.com/jinzhu/copier"
	"github.com/revel/revel"
	o_system "github.com/zze326/devops-helper/app/models/orm/system"
	v_system "github.com/zze326/devops-helper/app/models/view/system"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
)

type DataDict struct {
	gormc.Controller
}

// ListPage 分页查询
func (c DataDict) ListPage(pager *utils.Pager) revel.Result {
	var dataDictModels []*o_system.DataDict
	pager.Order = "id desc"
	total, err := utils.Paginate[o_system.DataDict](c.DB, pager, &dataDictModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: dataDictModels,
	})
}

// Add 创建
func (c DataDict) Add(req v_system.AddDataDictReq) revel.Result {
	exists, err := utils.DBExists[o_system.DataDict](c.DB, "type_code = ?", req.TypeCode)
	if err != nil {
		return results.JsonError(err)
	}

	if exists {
		return results.JsonErrorMsg("数据字典类型代码已存在")
	}

	exists, err = utils.DBExists[o_system.DataDict](c.DB, "name = ?", req.Name)
	if err != nil {
		return results.JsonError(err)
	}

	if exists {
		return results.JsonErrorMsg("数据字典名称已存在")
	}

	dataDictModel := new(o_system.DataDict)
	if err := copier.Copy(dataDictModel, &req); err != nil {
		return results.JsonError(err)
	}

	if err = c.DB.Create(dataDictModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 根据 ID 获取
func (c DataDict) Get(id int) revel.Result {
	dataDictModel := new(o_system.DataDict)
	if err := c.DB.First(dataDictModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(dataDictModel)
}

// Edit 编辑
func (c DataDict) Edit(req v_system.EditDataDictReq) revel.Result {
	exists, err := utils.DBExists[o_system.DataDict](c.DB, "id != ? and name = ?", req.ID, req.Name)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("数据字典名称已存在")
	}

	dataDictModel := new(o_system.DataDict)
	if err := c.DB.First(dataDictModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	dataDictModel.Name = req.Name
	dataDictModel.TypeCode = req.TypeCode

	if err := c.DB.Updates(dataDictModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除
func (c DataDict) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_system.DataDict{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}
