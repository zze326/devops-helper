package resource

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	v_resource "github.com/zze326/devops-helper/app/models/view/resource"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

type Host struct {
	gormc.Controller
	writeBuffer   bytes.Buffer
	readBuffer    bytes.Buffer
	session       *ssh.Session
	sessionBuffer bytes.Buffer
	host          *o_resource.Host
	operatorName  string
	operatorID    int
	isSaveSession bool
}

func (c Host) Add(req v_resource.AddHostReq) revel.Result {
	if len(req.GroupIDs) == 0 {
		return results.JsonErrorMsg("至少选择一个分组")
	}

	exists, err := utils.DBExistsByRaw(c.DB, `select 1 from host t1 
	left join host_group_host t2 on t1.id = t2.host_id
	left join host_group t3 on t2.host_group_id = t3.id
	where t1.deleted_at is null and t1.name = ? and t2.host_group_id in (?)`, req.Name, req.GroupIDs)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("所选分组下已存在该主机名称")
	}

	hostModel := new(o_resource.Host)
	if err := copier.Copy(hostModel, &req); err != nil {
		return results.JsonError(err)
	}

	var hostGroupModels []*o_resource.HostGroup
	if err := c.DB.Where("id in (?)", req.GroupIDs).Find(&hostGroupModels).Error; err != nil {
		return results.JsonError(err)
	}

	hostModel.HostGroups = hostGroupModels
	if err := c.DB.Create(hostModel).Error; err != nil {
		return results.JsonError(err)
	}

	return results.JsonOk()
}

func (c Host) Delete(id int) revel.Result {
	return results.JsonError(c.DB.Delete(&o_resource.Host{}, id).Error)
}

func (c Host) Get(id int) revel.Result {
	hostModel := new(o_resource.Host)
	if err := c.DB.Preload("HostGroups").Where("id = ?", id).First(hostModel).Error; err != nil {
		return results.JsonError(err)
	}
	hostModel.Password = strings.Repeat("*", utf8.RuneCountInString(hostModel.Password))
	hostModel.PrivateKey = ""
	return results.JsonOkData(hostModel)
}

func (c Host) GetPasswordAndPrivateKey(id int) revel.Result {
	hostModel := new(o_resource.Host)
	if err := c.DB.Where("id = ?", id).First(hostModel).Error; err != nil {
		return results.JsonError(err)
	}
	return results.JsonOkData(struct {
		Password   string `json:"password"`
		PrivateKey string `json:"private_key"`
	}{
		Password:   hostModel.Password,
		PrivateKey: hostModel.PrivateKey,
	})
}

func (c Host) Edit(req v_resource.EditHostReq) revel.Result {
	exists, err := utils.DBExistsByRaw(c.DB, `select 1 from host t1 
	left join host_group_host t2 on t1.id = t2.host_id
	left join host_group t3 on t2.host_group_id = t3.id
	where t1.deleted_at is null and t1.id != ? and t1.name = ? and t2.host_group_id in (?)`, req.ID, req.Name, req.GroupIDs)
	if err != nil {
		return results.JsonError(err)
	}
	if exists {
		return results.JsonErrorMsg("所选分组下已存在该主机名称")
	}

	hostModel := new(o_resource.Host)
	if err := c.DB.Where("id = ?", req.ID).First(hostModel).Error; err != nil {
		return results.JsonError(err)
	}

	var hostGroupModels []*o_resource.HostGroup
	if err := c.DB.Where("id in (?)", req.GroupIDs).Find(&hostGroupModels).Error; err != nil {
		return results.JsonError(err)
	}

	hostModel.Name = req.Name
	hostModel.Host = req.Host
	hostModel.Port = req.Port
	hostModel.Username = req.Username
	hostModel.UseKey = req.UseKey
	hostModel.PrivateKey = req.PrivateKey
	hostModel.Password = req.Password

	hostModel.Desc = req.Desc

	return results.JsonError(c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(hostModel).Association("HostGroups").Replace(hostGroupModels); err != nil {
			return err
		}
		txScope := tx.Model(hostModel).Select("*")
		if !req.UpdatePasswordAndPrivateKey {
			txScope = txScope.Omit("private_key", "password")
		}
		if err := txScope.Updates(hostModel).Error; err != nil {
			return err
		}
		return nil
	}))
}

// ListPage 分页查询
func (c Host) ListPage(pager *utils.Pager, hostGroupID int) revel.Result {
	var hostModels []*o_resource.Host
	pager.Order = "id desc"
	if hostGroupID > 0 {
		pager.Wheres = append(pager.Wheres, utils.WhereClause{
			Logic: "and",
			Columns: []utils.ColumnClause{
				{Column: "id",
					Op:    "->",
					Value: c.DB.Table("host_group_host").Select("host_id").Where("host_group_id = ?", hostGroupID),
				},
			},
		})
	}

	total, err := utils.Paginate[o_resource.Host](c.DB, pager, &hostModels)
	if err != nil {
		return results.JsonError(err)
	}

	for _, model := range hostModels {
		model.Password = ""
		model.PrivateKey = ""
	}

	return results.JsonOkData(results.PageData{
		Total: total,
		Items: hostModels,
	})
}

// TestSSH 测试 SSH 连接
func (c Host) TestSSH(id int) revel.Result {
	hostModel := new(o_resource.Host)
	if err := c.DB.Where("id = ?", id).First(hostModel).Error; err != nil {
		return results.JsonError(err)
	}

	client, err := hostModel.SSHClient()
	if err != nil {
		g.Logger.Errorf("%v", err)
		return results.JsonErrorMsgf("SSH 连接失败，请检查主机配置[%s]", err.Error())
	}
	defer func() {
		if err := client.Close(); err != nil {
			g.Logger.Errorf("%v", err)
		}
	}()
	return results.JsonOk()
}

// DownloadFile 下载文件
func (c Host) DownloadFile(id int, path string) revel.Result {
	hostModel := new(o_resource.Host)
	if err := c.DB.Where("id = ?", id).First(hostModel).Error; err != nil {
		return results.JsonError(err)
	}
	client, err := hostModel.SFTPClient()
	if err != nil {
		return results.JsonError(err)
	}

	file, err := client.OpenFile(path, os.O_RDONLY)
	if err != nil {
		return results.JsonError(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			g.Logger.Errorf("%v", err)
		}
	}()

	var length int64 = -1
	if size, err := file.Seek(0, 2); err == nil {
		if _, err = file.Seek(0, 0); err == nil {
			length = size
		}
	}
	c.Response.Out.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(path)))
	if err := c.Response.GetStreamWriter().WriteStream(filepath.Base(path), length, time.Now(), file); err != nil {
		return results.JsonError(err)
	}

	return nil
}
