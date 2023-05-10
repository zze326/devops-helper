package o_resource

import (
	"github.com/zze326/devops-helper/app/modules/gormc"
	"time"
)

type HostTerminalSession struct {
	gormc.Model
	HostID           int       `gorm:"comment:主机 ID" json:"host_id"`
	HostAddr         string    `gorm:"comment:主机名或IP" json:"host_addr"`
	HostName         string    `gorm:"comment:主机名" json:"host_name"`
	OperatorID       int       `gorm:"comment:操作人 ID" json:"operator_id"`
	OperatorName     string    `gorm:"comment:操作人用户名" json:"operator_name"`
	OperatorRealName string    `gorm:"comment:操作人真实姓名" json:"operator_real_name"`
	StartTime        time.Time `gorm:"comment:会话开始时间" json:"start_time"`
	Filepath         string    `gorm:"comment:会话文件路径" json:"filepath"`
}

func (HostTerminalSession) TableName() string {
	return "host_terminal_session"
}
