package resource

import (
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	o_system "github.com/zze326/devops-helper/app/models/orm/system"
	v_resource "github.com/zze326/devops-helper/app/models/view/resource"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"gorm.io/gorm"
)

type HostGroup struct {
	gormc.Controller
}

// Add 创建主机组
func (c HostGroup) Add(req v_resource.AddHostGroupReq) revel.Result {
	exists, err := utils.DBExists[o_resource.HostGroup](c.DB, "name = ? and parent_id = ?", req.Name, req.ParentID)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("服务器组名称已存在")
	}

	if req.ParentID != 0 {
		exists, err = utils.DBExists[o_resource.HostGroup](c.DB, req.ParentID)
		if err != nil {
			return results.JsonError(err)
		}
		if !exists {
			return results.JsonErrorMsg("父级服务器组不存在")
		}
	}

	serverGroupModel := new(o_resource.HostGroup)
	serverGroupModel.ParentID = req.ParentID
	serverGroupModel.Name = req.Name
	if err := c.DB.Create(serverGroupModel).Error; err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(serverGroupModel)
}

// ListTreeWithHosts 获取用户可访问的服务器组树
func (c HostGroup) ListTreeWithHosts(_requestUserID int) revel.Result {
	// 获取用户角色
	userModel := new(o_system.User)
	if err := g.DB.Preload("Roles").First(userModel, _requestUserID).Error; err != nil {
		return results.JsonError(err)
	}

	var roleIDs []int
	for _, role := range userModel.Roles {
		roleIDs = append(roleIDs, role.ID)
	}

	var serverGroupModels []*o_resource.HostGroup
	// 从数据库中获取所有权限记录
	if err := c.DB.Preload("HostGroupPermissions", func(db *gorm.DB) *gorm.DB {
		if userModel.IsSuper() {
			return db
		}
		return db.Where("(type = 1 and ref_id = ?) or (type = 2 and ref_id in (?))", _requestUserID, roleIDs).Select("id")
	}).Preload("Hosts", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "host", "desc")
	}).Find(&serverGroupModels).Error; err != nil {
		return results.JsonError(err)
	}

	// 构建 ID 到权限映射的 map
	serverGroupMap := make(map[int]*o_resource.HostGroup, len(serverGroupModels))
	for _, serverGroup := range serverGroupModels {
		if !userModel.IsSuper() && len(serverGroup.HostGroupPermissions) == 0 {
			serverGroup.Hosts = nil
		}
		serverGroupMap[serverGroup.ID] = serverGroup
	}

	// 找出根节点并返回
	var roots []*o_resource.HostGroup
	// 对每个权限，将其作为子节点添加到对应父节点的 Children 列表中
	for _, serverGroup := range serverGroupModels {
		parent, ok := serverGroupMap[serverGroup.ParentID]
		if !ok {
			// 如果该权限没有父节点，则认为它是根节点
			roots = append(roots, serverGroup)
			continue
		}
		parent.Children = append(parent.Children, serverGroup)
	}

	roots = filterTree(roots)

	return results.JsonOkData(roots)
}

// 遍历 roots，删除没有子节点并且没有主机的节点
func filterTree(roots []*o_resource.HostGroup) []*o_resource.HostGroup {
	var result []*o_resource.HostGroup
	for _, root := range roots {
		if hasHosts(root) {
			root.Children = filterTree(root.Children)
			result = append(result, root)
		}
	}
	return result
}

// 遍历所有子节点，判断是否有主机
func hasHosts(serverGroup *o_resource.HostGroup) bool {
	if len(serverGroup.Hosts) > 0 {
		return true
	}
	for _, child := range serverGroup.Children {
		if hasHosts(child) {
			return true
		}
	}
	return false
}

// ListTree 获取服务器组树
func (c HostGroup) ListTree() revel.Result {
	var serverGroupModels []*o_resource.HostGroup
	// 从数据库中获取所有权限记录
	if err := c.DB.Preload("HostGroupPermissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("id")
	}).Preload("Hosts", func(db *gorm.DB) *gorm.DB {
		return db.Select("id")
	}).Find(&serverGroupModels).Error; err != nil {
		return results.JsonError(err)
	}

	// 构建 ID 到权限映射的 map
	serverGroupMap := make(map[int]*o_resource.HostGroup, len(serverGroupModels))
	for _, serverGroup := range serverGroupModels {
		serverGroupMap[serverGroup.ID] = serverGroup
	}

	// 找出根节点并返回
	var roots []*o_resource.HostGroup
	// 对每个权限，将其作为子节点添加到对应父节点的 Children 列表中
	for _, serverGroup := range serverGroupModels {
		parent, ok := serverGroupMap[serverGroup.ParentID]
		if !ok {
			// 如果该权限没有父节点，则认为它是根节点
			roots = append(roots, serverGroup)
			continue
		}
		parent.Children = append(parent.Children, serverGroup)
	}
	return results.JsonOkData(roots)
}

// Rename 重命名服务器组
func (c HostGroup) Rename(req v_resource.RenameHostGroupReq) revel.Result {
	hostGroupModel := new(o_resource.HostGroup)
	if err := c.DB.Model(&o_resource.HostGroup{}).Where("id = ?", req.ID).First(hostGroupModel).Error; err != nil {
		return results.JsonError(err)
	}

	exists, err := utils.DBExists[o_resource.HostGroup](c.DB, "id != ? and parent_id = ? and name = ?", req.ID, hostGroupModel.ParentID, req.Name)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("服务器组名称已存在")
	}
	if err := c.DB.Model(&o_resource.HostGroup{}).Where("id = ?", req.ID).Update("name", req.Name).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

// Delete 删除服务器组
func (c HostGroup) Delete(id int) revel.Result {
	exists, err := utils.DBExists[o_resource.HostGroup](c.DB, "parent_id = ?", id)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("请移除子分组后再尝试删除")
	}

	if err := c.DB.Delete(&o_resource.HostGroup{}, id).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}
