package resource

import (
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	"github.com/zze326/devops-helper/app/results"
	"time"
)

func (c Host) TerminalSessionLog(id int) revel.Result {
	g.Logger.Infof("建立查看 Web 终端记录 Websocket 连接")
	sessionLogModel := new(o_resource.HostTerminalSession)
	if err := c.DB.Where("id = ?", id).First(sessionLogModel).Error; err != nil {
		return results.JsonError(err)
	}
	// 控制速度发送 sessionLogModel.Data
	msgBytes := sessionLogModel.Data
	for len(msgBytes) > 0 {
		if len(msgBytes) > 8 {
			c.Request.WebSocket.MessageSend(string(msgBytes[:8]))
			msgBytes = msgBytes[8:]
		} else {
			c.Request.WebSocket.MessageSend(string(msgBytes))
			msgBytes = nil
		}
		time.Sleep(60 * time.Millisecond)
	}

	//err = c.Request.WebSocket.MessageSend(string(msgBytes))

	g.Logger.Infof("关闭查看 Web 终端记录 Websocket 连接")
	return results.JsonOk()
}
