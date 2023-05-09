package resource

import (
	"database/sql"
	"fmt"
	"github.com/revel/revel"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	v_resource "github.com/zze326/devops-helper/app/models/view/resource"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"gorm.io/gorm"
)

type HostGroupPermission struct {
	gormc.Controller
}

// Edit 编辑主机组权限
func (c HostGroupPermission) Edit(req v_resource.EditHostGroupPermissionReq) revel.Result {
	exists, err := utils.DBExists[o_resource.HostGroupPermission](c.DB, "id != ? and type = ? and ref_id = ?", req.ID, req.Type, req.RefID)
	if err != nil {
		return results.JsonError(err)
	}

	if exists {
		var permissionType string
		if req.Type == 1 {
			permissionType = "用户"
		} else if req.Type == 2 {
			permissionType = "角色"
		} else {
			return results.JsonOkMsg("不存在的权限类型")
		}
		return results.JsonErrorMsgf("已存在该%s权限", permissionType)
	}

	hostGroupPermissionModel := new(o_resource.HostGroupPermission)
	if err := c.DB.First(hostGroupPermissionModel, req.ID).Error; err != nil {
		return results.JsonError(err)
	}

	hostGroupPermissionModel.Type = req.Type
	hostGroupPermissionModel.RefID = req.RefID

	var hostGroupModels []*o_resource.HostGroup
	if len(req.HostGroupIDs) > 0 {
		if err := c.DB.Find(&hostGroupModels, req.HostGroupIDs).Error; err != nil {
			return results.JsonError(err)
		}
	}

	if err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(hostGroupPermissionModel).Association("HostGroups").Replace(hostGroupModels); err != nil {
			return err
		}

		if err := tx.Select("*").Updates(hostGroupPermissionModel).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return results.JsonError(err)
	}

	return results.JsonOk()
}

// Get 获取主机组权限
func (c HostGroupPermission) Get(id int) revel.Result {
	hostGroupPermissionModel := new(o_resource.HostGroupPermission)
	if err := c.DB.Preload("HostGroups").First(hostGroupPermissionModel, id).Error; err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(hostGroupPermissionModel)
}

// Add 创建主机组权限
func (c HostGroupPermission) Add(req v_resource.AddHostGroupPermissionReq) revel.Result {
	exists, err := utils.DBExists[o_resource.HostGroupPermission](c.DB, "type = ? and ref_id = ?", req.Type, req.RefID)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		var permissionType string
		if req.Type == 1 {
			permissionType = "用户"
		} else if req.Type == 2 {
			permissionType = "角色"
		} else {
			return results.JsonOkMsg("不存在的权限类型")
		}
		return results.JsonErrorMsgf("已存在该%s权限", permissionType)
	}

	var hostGroupModels []*o_resource.HostGroup
	if err := c.DB.Find(&hostGroupModels, "id in ?", req.HostGroupIDs).Error; err != nil {
		return results.JsonError(err)
	}

	hostGroupPermission := new(o_resource.HostGroupPermission)
	hostGroupPermission.Type = req.Type
	hostGroupPermission.RefID = req.RefID
	hostGroupPermission.HostGroups = hostGroupModels // 这里会自动创建关联表

	if err := c.DB.Create(hostGroupPermission).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// ListPage 分页查询主机组权限
func (c HostGroupPermission) ListPage(pager *utils.Pager, hostGroupID int) revel.Result {
	var (
		respItems []*v_resource.ListPageHostGroupPermissionItem
		total     int
	)

	selectFieldsStr := `t1.id,
		t1.type,
		t1.ref_id,
		case when t1.type = 1 then t2.real_name else t3.name end as ref_name,
		GROUP_CONCAT(DISTINCT(CONCAT(t5.id,'_',t5.name)) SEPARATOR ',') as host_group_names_str`

	whereStr := "where t1.deleted_at is null"

	mainSql := `select %s
		from host_group_permission t1 
	left join user t2 on t1.type = 1 and t1.ref_id = t2.id
	left join role t3 on t1.type = 2 and t1.ref_id = t3.id
	left join host_group_permission_host_group t4 on t1.id = t4.host_group_permission_id
	left join host_group t5 on t4.host_group_id = t5.id
	%s
	group by t1.id, t1.type
	order by t1.id desc limit @offset, @limit`

	var nameArgs []any
	if hostGroupID > 0 {
		whereStr += " and t5.id = @host_group_id"
		nameArgs = append(nameArgs, sql.Named("host_group_id", hostGroupID))
	}

	if !utils.IsEmpty(pager.Search) {
		whereStr += " and (t2.username like @search or t2.real_name like @search or t3.name like @search or t3.code like @search)"
		nameArgs = append(nameArgs, sql.Named("search", "%"+pager.Search+"%"))
	}

	querySql := fmt.Sprintf(mainSql, selectFieldsStr, whereStr)
	countSql := fmt.Sprintf(mainSql, "count(*)", whereStr)

	nameArgs = append(nameArgs, sql.Named("offset", pager.Offset()))
	nameArgs = append(nameArgs, sql.Named("limit", pager.Limit()))

	if err := c.DB.Raw(countSql, nameArgs...).Scan(&total).Error; err != nil {
		return results.JsonError(err)
	}

	if err := c.DB.Raw(querySql, nameArgs...).Scan(&respItems).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(results.PageData{
		Total: total,
		Items: respItems,
	})
}

// Delete 删除主机组权限
func (c HostGroupPermission) Delete(id int) revel.Result {
	return results.JsonError(c.DB.Delete(&o_resource.HostGroupPermission{}, id).Error)
}
