package o_resource

import "github.com/zze326/devops-helper/app/modules/gormc"

type HostTerminalSession struct {
	gormc.Model
	HostID       int    `gorm:"comment:主机 ID" json:"host_id"`
	HostAddr     string `gorm:"comment:主机名或IP" json:"host"`
	HostName     string `gorm:"comment:主机名" json:"host_name"`
	OperatorID   int    `gorm:"comment:操作人 ID" json:"operator_id"`
	OperatorName string `gorm:"comment:操作人" json:"operator"`
	Data         []byte `gorm:"type:LONGBLOB;comment:会话数据" json:"data"`
}

func (HostTerminalSession) TableName() string {
	return "host_terminal_session"
}
