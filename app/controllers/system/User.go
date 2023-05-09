package system

import (
	"fmt"
	_ "github.com/asaskevich/govalidator"
	"github.com/jinzhu/copier"
	"github.com/revel/revel"
	o_system "github.com/zze326/devops-helper/app/models/orm/system"
	v_system "github.com/zze326/devops-helper/app/models/view/system"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"gorm.io/gorm"
)

type User struct {
	gormc.Controller
}

// ListPage 分页查询
func (c User) ListPage(pager *utils.Pager) revel.Result {
	var userModels []*o_system.User
	pager.Order = "id desc"
	total, err := utils.Paginate[o_system.User](c.DB, pager, &userModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: userModels,
	})
}

// Add 创建用户
func (c User) Add(req v_system.AddUserReq) revel.Result {
	userModel := new(o_system.User)
	if err := copier.Copy(userModel, &req); err != nil {
		return results.JsonError(err)
	}
	err := c.DB.First(userModel, "username = ?", req.Username).Error
	if err == nil {
		return results.JsonErrorMsg("用户名已存在")
	} else {
		if err != gorm.ErrRecordNotFound {
			return results.JsonError(err)
		}
	}

	defaultPassword := "123456"
	userModel.Password = utils.EncodeMD5(defaultPassword)

	if len(req.RoleIDs) > 0 {
		var roleModels []*o_system.Role
		if err := c.DB.Find(&roleModels, req.RoleIDs).Error; err != nil {
			return results.JsonError(err)
		}
		userModel.Roles = roleModels
	}

	if err = c.DB.Create(userModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkMsg("创建成功，初始密码为：" + defaultPassword)
}

// Get 根据 ID 获取
func (c User) Get(id int) revel.Result {
	userModel := new(o_system.User)
	if err := c.DB.Preload("Roles").First(userModel, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(userModel)
}

// Edit 编辑用户
func (c User) Edit(req v_system.EditUserReq) revel.Result {
	userModel := new(o_system.User)
	if err := c.DB.First(userModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}
	userModel.Username = req.Username
	userModel.RealName = req.RealName
	userModel.Phone = req.Phone
	userModel.Email = req.Email

	var roles []*o_system.Role
	if len(req.RoleIDs) > 0 {
		if err := c.DB.Find(&roles, req.RoleIDs).Error; err != nil {
			return results.JsonError(err)
		}
	}
	if err := c.DB.Model(userModel).Association("Roles").Replace(roles); err != nil {
		return results.JsonError(err)
	}
	if err := c.DB.Updates(userModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除用户
func (c User) Delete(id int) revel.Result {
	if err := c.DB.Delete(&o_system.User{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// ResetPwd 重置密码
func (c User) ResetPwd(req v_system.ResetUserPasswordReq) revel.Result {
	userModel := new(o_system.User)
	if err := c.DB.First(userModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	userModel.Password = utils.EncodeMD5(req.NewPassword)
	if err := c.DB.Save(userModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkMsg(fmt.Sprintf("重置密码成功，新密码：%s", req.NewPassword))
}

// ListAll 获取所有用户
func (c User) ListAll() revel.Result {
	var userModels []*o_system.User
	if err := c.DB.Omit("Password").Find(&userModels).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(userModels)
}
