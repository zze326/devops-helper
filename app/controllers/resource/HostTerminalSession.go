package resource

import (
	"bytes"
	"github.com/revel/revel"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
)

// HostTerminalSession 主机终端会话
type HostTerminalSession struct {
	gormc.Controller
	writeBuffer bytes.Buffer
}

// ListPage 分页查询
func (c HostTerminalSession) ListPage(pager *utils.Pager) revel.Result {
	var sessionModels []*o_resource.HostTerminalSession
	pager.Order = "id desc"
	pager.OmitColumns = []string{"data"}
	total, err := utils.Paginate[o_resource.HostTerminalSession](c.DB, pager, &sessionModels)
	if err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: sessionModels,
	})
}

// CheckSessionFile 会话文件是否存在
func (c HostTerminalSession) CheckSessionFile(id int) revel.Result {
	sessionModel := new(o_resource.HostTerminalSession)
	if err := c.DB.Where("id = ?", id).First(sessionModel).Error; err != nil {
		return results.JsonError(err)
	}

	return results.JsonOkData(struct {
		Exists bool `json:"exists"`
	}{utils.FileExists(sessionModel.Filepath)})
}
