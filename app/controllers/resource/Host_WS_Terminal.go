package resource

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"path"
	"strings"
	"time"
	"unicode/utf8"
)

// Terminal Web 终端
func (c Host) Terminal(_requestUserID int, _requestUsername string, _requestUserRealName string, id int) revel.Result {
	g.Logger.Infof("建立 Web 终端 Websocket 连接")
	c.isSaveSession = revel.Config.BoolDefault("host.terminal.savesession", false)
	hostModel := new(o_resource.Host)
	if err := c.DB.Where("id = ?", id).First(hostModel).Error; err != nil {
		return results.JsonError(err)
	}
	c.host = hostModel
	c.operatorID = _requestUserID
	c.operatorName = _requestUsername
	c.operatorRealName = _requestUserRealName

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

	session, err := client.NewSession()
	if err != nil {
		g.Logger.Errorf("%v", err)
		return results.JsonError(err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			g.Logger.Errorf("%v", err)
		}
	}()
	c.session = session
	c.startTime = time.Now()
	defer c.saveSession()

	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err)
		os.Exit(1)
	}

	fd := int(os.Stdin.Fd())
	session.Stdout = &c
	session.Stderr = &c
	session.Stdin = &c
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	termWidth, termHeight, err := terminal.GetSize(fd)
	err = session.RequestPty("xterm", termHeight, termWidth, modes)
	if err != nil {
		g.Logger.Errorf("%v", err)
		return results.JsonError(err)
	}

	err = session.Shell()
	if err != nil {
		g.Logger.Errorf("%v", err)
		return results.JsonError(err)
	}

	c.sessionFilepath = path.Join(revel.Config.StringDefault("host.terminal.sessionfiledir", "host-sessions"), fmt.Sprintf("%d", c.host.ID), fmt.Sprintf("%d.sessionb", c.startTime.UnixMicro()))
	if c.isSaveSession {
		c.sessionFile, err = utils.OpenOrCreateFile(c.sessionFilepath)
		if err != nil {
			g.Logger.Errorf("创建会话文件失败: %v", err)
		} else {
			defer c.sessionFile.Close()
		}
	}

	err = session.Wait()
	if err != nil {
		g.Logger.Errorf("%v", err)
		return results.JsonError(err)
	}
	g.Logger.Infof("关闭 Web 终端 Websocket 连接")
	return results.JsonOk()
}

// Read 接收 Web 终端的命令
func (c *Host) Read(p []byte) (n int, err error) {
	type xtermMessage struct {
		MsgType string `json:"type"`
		Input   string `json:"input"`
		Rows    uint16 `json:"rows"`
		Cols    uint16 `json:"cols"`
	}

	var xtermMsg xtermMessage
	err = c.Request.WebSocket.MessageReceiveJSON(&xtermMsg)
	if err != nil {
		return 0, err
	}
	if xtermMsg.MsgType == "input" {
		if cmdStr := strings.TrimSpace(c.readBuffer.String()); xtermMsg.Input == "\r" && len(cmdStr) > 0 {
			c.hasInput = true
			c.Log.Infof("输入命令：%s", cmdStr)
			c.readBuffer.Reset()
		}
		if !utils.InSlice[string](xtermMsg.Input, []string{"\t", "\r", ""}) {
			c.readBuffer.Write([]byte(xtermMsg.Input))
		}
		copy(p, fmt.Sprintf("%s", xtermMsg.Input))
		return len(xtermMsg.Input), nil
	} else if xtermMsg.MsgType == "resize" {
		g.Logger.Infof("resize: cols=%d, rows=%d", xtermMsg.Cols, xtermMsg.Rows)
		// 改变终端大小
		if err = c.session.WindowChange(int(xtermMsg.Rows), int(xtermMsg.Cols)); err != nil {
			g.Logger.Errorf("改变终端大小失败: %s", err.Error())
		}
	} else if xtermMsg.MsgType == "close" {
		g.Logger.Infof("关闭 Web 终端")
		if err := c.session.Close(); err != nil {
			g.Logger.Errorf("关闭 Web 终端失败: %s", err.Error())
		}
		return 0, io.EOF
	}
	return 0, nil
}

// Write 响应到 Web 终端
func (c *Host) Write(p []byte) (n int, err error) {
	msgBytes := p

	if _, err := c.sessionFile.Write(p); err != nil {
		g.Logger.Errorf("写入会话记录到文件失败, err: %v", err)
	}
	if !utf8.Valid(msgBytes) {
		c.writeBuffer.Write(msgBytes)
		return len(p), nil
	} else {
		if c.writeBuffer.Len() > 0 {
			c.writeBuffer.Write(msgBytes)
			msgBytes = c.writeBuffer.Bytes()
			c.writeBuffer.Reset()
		}
	}

	err = c.Request.WebSocket.MessageSend(string(msgBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	return len(p), nil
}

// saveSession 保存会话数据到数据库
func (c *Host) saveSession() {
	if c.isSaveSession && !c.hasInput {
		if err := utils.EnsureFileNotExists(c.sessionFilepath); err != nil {
			g.Logger.Errorf("删除空会话文件失败: %v", err)
		}
	}

	if !c.isSaveSession || !c.hasInput {
		return
	}

	session := new(o_resource.HostTerminalSession)
	session.HostID = c.host.ID
	session.HostAddr = c.host.Host
	session.HostName = c.host.Name
	session.OperatorName = c.operatorName
	session.StartTime = c.startTime
	session.OperatorID = c.operatorID
	session.OperatorRealName = c.operatorRealName
	session.Filepath = c.sessionFilepath

	if err := c.DB.Save(session).Error; err != nil {
		g.Logger.Errorf("保存会话数据到数据库失败: %s", err.Error())
		return
	}
}
