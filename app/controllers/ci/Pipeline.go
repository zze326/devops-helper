package ci

import (
	"github.com/jinzhu/copier"
	"github.com/revel/revel"
	o_ci "github.com/zze326/devops-helper/app/models/orm/ci"
	v_ci "github.com/zze326/devops-helper/app/models/view/ci"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"gorm.io/gorm"
)

type Pipeline struct {
	gormc.Controller
}

// ListPage 分页查询
func (c Pipeline) ListPage(pager *utils.Pager) revel.Result {
	var pipelineModels []*o_ci.Pipeline
	pager.Order = "id desc"
	total, err := utils.Paginate[o_ci.Pipeline](c.DB, pager, &pipelineModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: pipelineModels,
	})
}

// Add 创建流水线
func (c Pipeline) Add(req v_ci.AddPipelineReq) revel.Result {
	pipelineModel := new(o_ci.Pipeline)
	if err := copier.Copy(pipelineModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.DB.First(pipelineModel, "name = ?", req.Name).Error
	if err == nil {
		return results.JsonErrorMsg("流水线名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if err = c.DB.Create(pipelineModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 根据 ID 获取
func (c Pipeline) Get(id int) revel.Result {
	pipelineModel := new(o_ci.Pipeline)
	if err := c.DB.First(pipelineModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(pipelineModel)
}

// Edit 编辑
func (c Pipeline) Edit(req v_ci.EditPipelineReq) revel.Result {
	exists, err := utils.DBExists[o_ci.Pipeline](c.DB, "id != ? and name = ?", req.ID, req.Name)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("流水线名称已存在")
	}
	pipelineModel := new(o_ci.Pipeline)
	if err := c.DB.First(pipelineModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	pipelineModel.Name = req.Name
	pipelineModel.Desc = req.Desc

	if err := c.DB.Updates(pipelineModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除流水线
func (c Pipeline) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_ci.Pipeline{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Run 运行流水线
func (c Pipeline) Run(id int) revel.Result {
	if err := c.DB.Delete(&o_ci.Pipeline{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}
