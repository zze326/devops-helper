package system

import (
	"github.com/jinzhu/copier"
	"github.com/revel/revel"
	o_system "github.com/zze326/devops-helper/app/models/orm/system"
	v_system "github.com/zze326/devops-helper/app/models/view/system"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"gorm.io/gorm"
)

type DataDictItem struct {
	gormc.Controller
}

// ListPage 分页查询
func (c DataDictItem) ListPage(pager *utils.Pager) revel.Result {
	var dataDictItemModels []*o_system.DataDictItem
	pager.Order = "sort"
	total, err := utils.Paginate[o_system.DataDictItem](c.DB, pager, &dataDictItemModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: dataDictItemModels,
	})
}

// Add 创建
func (c DataDictItem) Add(req v_system.AddDataDictItemReq) revel.Result {
	dataDictItemModel := new(o_system.DataDictItem)
	if err := copier.Copy(dataDictItemModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.DB.First(dataDictItemModel, "data_dict_id = ? and value = ?", req.DataDictID, req.Value).Error
	if err == nil {
		return results.JsonErrorMsg("数据字典项代码已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if err = c.DB.Create(dataDictItemModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 根据 ID 获取
func (c DataDictItem) Get(id int) revel.Result {
	dataDictItemModel := new(o_system.DataDictItem)
	if err := c.DB.First(dataDictItemModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(dataDictItemModel)
}

// Edit 编辑
func (c DataDictItem) Edit(req v_system.EditDataDictItemReq) revel.Result {
	exists, err := utils.DBExists[o_system.DataDictItem](c.DB, "id != ? and value = ?", req.ID, req.Value)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("数据字典项代码已存在")
	}

	dataDictItemModel := new(o_system.DataDictItem)
	if err := c.DB.First(dataDictItemModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	dataDictItemModel.Label = req.Label
	dataDictItemModel.Value = req.Value
	dataDictItemModel.Sort = req.Sort

	if err := c.DB.Updates(dataDictItemModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除
func (c DataDictItem) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_system.DataDictItem{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// ListByDataDictTypeCode 根据数据字典 Code 获取数据字典项
func (c DataDictItem) ListByDataDictTypeCode(typeCode string) revel.Result {
	dataDictModel := new(o_system.DataDict)
	if err := c.DB.First(dataDictModel, "type_code = ?", typeCode).Error; err != nil {
		return results.JsonError(err)
	}

	var dataDictItemModels []*o_system.DataDictItem
	if err := c.DB.Where("data_dict_id = ?", dataDictModel.ID).Order("sort").Find(&dataDictItemModels).Error; err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(dataDictItemModels)
}
