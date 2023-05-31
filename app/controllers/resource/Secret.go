package resource

import (
	"github.com/jinzhu/copier"
	"github.com/revel/revel"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	v_resource "github.com/zze326/devops-helper/app/models/view/resource"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"gorm.io/gorm"
)

type Secret struct {
	gormc.Controller
}

// ListPage 分页查询
func (c Secret) ListPage(pager *utils.Pager) revel.Result {
	var secretModels []*o_resource.Secret
	pager.Order = "id desc"
	total, err := utils.Paginate[o_resource.Secret](c.DB, pager, &secretModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: secretModels,
	})
}

// Add 创建
func (c Secret) Add(req v_resource.AddSecretReq) revel.Result {
	secretModel := new(o_resource.Secret)
	if err := copier.Copy(secretModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.DB.First(secretModel, "name = ? and type = ?", req.Name, req.Type).Error
	if err == nil {
		return results.JsonErrorMsg("秘钥名称已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	if err = c.DB.Create(secretModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Get 根据 ID 获取
func (c Secret) Get(id int) revel.Result {
	secretModel := new(o_resource.Secret)
	if err := c.DB.First(secretModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(secretModel)
}

// Edit 编辑
func (c Secret) Edit(req v_resource.EditSecretReq) revel.Result {
	exists, err := utils.DBExists[o_resource.Secret](c.DB, "id != ? and name = ? and type = ?", req.ID, req.Name, req.Type)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("秘钥名称已存在")
	}

	secretModel := new(o_resource.Secret)
	if err := c.DB.First(secretModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	secretModel.Name = req.Name
	secretModel.Type = req.Type

	if secretModel.UseContent() {
		secretModel.Content = req.Content
	} else {
		secretModel.Username = req.Username
		secretModel.Password = req.Password
	}

	if err := c.DB.Updates(secretModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除
func (c Secret) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_resource.Secret{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}
