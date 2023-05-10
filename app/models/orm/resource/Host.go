package o_resource

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/zze326/devops-helper/app/modules/gormc"
	"golang.org/x/crypto/ssh"
	"time"
)

type Host struct {
	gormc.Model
	Name        string       `gorm:"comment:名称" json:"name"`
	Host        string       `gorm:"comment:主机名或IP" json:"host"`
	Port        int          `gorm:"comment:端口" json:"port"`
	Username    string       `gorm:"comment:用户名" json:"username"`
	Password    string       `gorm:"comment:密码" json:"password"`
	PrivateKey  string       `gorm:"size:4096;comment:私钥;" json:"private_key"`
	UseKey      bool         `gorm:"comment:是否使用公钥连接" json:"use_key"`
	Desc        string       `gorm:"comment:描述" json:"desc"`
	HostGroups  []*HostGroup `gorm:"many2many:host_group_host;" json:"host_groups"`
	SaveSession bool         `gorm:"comment:是否保存会话;default:0;" json:"save_session"`
}

func (Host) TableName() string {
	return "host"
}

func (m *Host) SFTPClient() (*sftp.Client, error) {
	client, err := m.SSHClient()
	if err != nil {
		return nil, err
	}
	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}
	return sftpClient, nil
}
func (m *Host) SSHClient() (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod
	if m.UseKey {
		signer, err := ssh.ParsePrivateKey([]byte(m.PrivateKey))
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else {
		authMethods = append(authMethods, ssh.Password(m.Password))
	}

	config := &ssh.ClientConfig{
		User:            m.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", m.Host, m.Port), config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
