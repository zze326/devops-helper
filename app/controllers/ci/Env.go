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

type CIEnv struct {
	gormc.Controller
}

// ListPage 分页查询
func (c CIEnv) ListPage(pager *utils.Pager) revel.Result {
	var envModels []*o_ci.Env
	pager.Order = "id desc"
	total, err := utils.Paginate[o_ci.Env](c.DB, pager, &envModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: envModels,
	})
}

// Add 创建
func (c CIEnv) Add(req v_ci.AddEnvReq) revel.Result {
	envModel := new(o_ci.Env)
	if err := copier.Copy(envModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.DB.First(envModel, "name = ?", req.Name).Error
	if err == nil {
		return results.JsonErrorMsg("环境名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if err = c.DB.Create(envModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 根据 ID 获取
func (c CIEnv) Get(id int) revel.Result {
	envModel := new(o_ci.Env)
	if err := c.DB.First(envModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(envModel)
}

// Edit 编辑
func (c CIEnv) Edit(req v_ci.EditEnvReq) revel.Result {
	exists, err := utils.DBExists[o_ci.Env](c.DB, "id != ? and name = ?", req.ID, req.Name)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("环境名称已存在")
	}
	envModel := new(o_ci.Env)
	if err := c.DB.First(envModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	envModel.Name = req.Name
	envModel.Image = req.Image
	envModel.K8sSecretName = req.K8sSecretName

	if err := c.DB.Updates(envModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除
func (c CIEnv) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_ci.Env{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// ListAll 查询所有
func (c CIEnv) ListAll() revel.Result {
	var envModels []*o_ci.Env
	if err := c.DB.Find(&envModels).Order("id desc").Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(envModels)
}
