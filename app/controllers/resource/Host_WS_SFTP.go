package resource

import (
	"bytes"
	"encoding/json"
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/utils"
	"math"
	"os"
	"path"
	"strings"
)

func (c Host) SFTPFileManage(id int) revel.Result {
	g.Logger.Infof("建立 SFTP 文件管理器 Websocket 连接")
	hostModel := new(o_resource.Host)
	if err := c.DB.Where("id = ?", id).First(hostModel).Error; err != nil {
		return results.JsonError(err)
	}

	sftpClient, err := hostModel.SFTPClient()

	if err != nil {
		g.Logger.Errorf("%v", err)
		return results.JsonErrorMsgf("SSH 连接失败，请检查主机配置[%s]", err.Error())
	}
	defer sftpClient.Close()

	type fileOperate struct {
		Type            string `json:"type"`
		Path            string `json:"path"`
		Filename        string `json:"filename"`
		Data            []byte `json:"data"`
		ChunkStart      int64  `json:"chunk_start"`
		ChunkEnd        int64  `json:"chunk_end"`
		TotalSize       int64  `json:"total_size"`
		ShowHiddenFiles bool   `json:"show_hidden_files"`
	}

	type wsResp struct {
		Type      string `json:"type"`
		Success   bool   `json:"success"`
		Data      any    `json:"data"`
		Path      string `json:"path"`
		TotalSize int64  `json:"total_size"`
		ChunkEnd  int64  `json:"chunk_end"`
		Msg       string `json:"msg"`
	}

	type fileinfo struct {
		Name    string `json:"name"`
		Mode    string `json:"mode"`
		Size    int64  `json:"size"`
		ModTime string `json:"mtime"`
		AbsPath string `json:"abs_path"`
		IsDir   bool   `json:"is_dir"`
	}
	var msgBytesBuffer bytes.Buffer
	for {
		// 接受字节
		var msgBytes []byte
		// 接收消息
		if err := c.Request.WebSocket.MessageReceive(&msgBytes); err != nil {
			g.Logger.Errorf("websocket message receive error: %s", err.Error())
			break
		}

		msg := new(fileOperate)
		if msgBytes[0] != '{' {
			msgBytesBuffer.Write(msgBytes)
			if msgBytes[len(msgBytes)-1] != '}' {
				continue
			}
			if json.Unmarshal(msgBytesBuffer.Bytes(), &msg) != nil {
				continue
			}
			msgBytesBuffer.Reset()
		} else {
			if msgBytes[len(msgBytes)-1] != '}' {
				msgBytesBuffer.Write(msgBytes)
				continue
			}
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				msgBytesBuffer.Write(msgBytes)
				if json.Unmarshal(msgBytesBuffer.Bytes(), &msg) != nil {
					continue
				}
				msgBytesBuffer.Reset()
			}
		}

		if msg.Type == "exit" {
			break
		}

		switch msg.Type {
		// 列出目录下的文件
		case "list":
			files, err := sftpClient.ReadDir(msg.Path)
			if err != nil {
				g.Logger.Errorf("read dir error: %s", err.Error())
				break
			}
			var fileinfos []*fileinfo
			for _, file := range files {
				if !msg.ShowHiddenFiles && strings.HasPrefix(file.Name(), ".") {
					continue
				}
				fileinfos = append(fileinfos, &fileinfo{
					Name:    file.Name(),
					Mode:    utils.FileMode(file.Mode()),
					Size:    file.Size(),
					ModTime: file.ModTime().Format("2006-01-02 15:04:05"),
					AbsPath: path.Join(msg.Path, file.Name()),
					IsDir:   file.IsDir(),
				})
			}

			if err := c.Request.WebSocket.MessageSendJSON(&wsResp{
				Type: "listData",
				Data: fileinfos,
				Path: msg.Path,
			}); err != nil {
				g.Logger.Errorf("websocket message send error: %s", err.Error())
				break
			} // 列出目录下的文件
		// 上传文件
		case "uploadFileChunk":
			file, err := sftpClient.OpenFile(path.Join(msg.Path, msg.Filename), os.O_CREATE|os.O_WRONLY)
			if err != nil {
				if err2 := c.Request.WebSocket.MessageSendJSON(&wsResp{
					Type:    "uploadFileChunk",
					Success: false,
					Msg:     err.Error(),
				}); err2 != nil {
					g.Logger.Errorf("websocket message send error: %s", err2.Error())
					break
				}
				continue
			}
			defer file.Close()

			_, err = file.WriteAt(msg.Data, msg.ChunkStart)
			if err != nil {
				g.Logger.Errorf("写入文件时出错：", err)
				break
			}

			if msg.ChunkEnd >= msg.TotalSize {
				if err := c.Request.WebSocket.MessageSendJSON(&wsResp{
					Type:    "uploadFileChunk",
					Success: true,
					Path:    msg.Path,
				}); err != nil {
					g.Logger.Errorf("websocket message send error: %s", err.Error())
					break
				}
			} else {
				g.Logger.Infof("正在上传文件: %s 到 %s, 进度: %.0f%%, 偏移量: %d", msg.Filename, msg.Path, math.Round(float64(msg.ChunkEnd)/float64(msg.TotalSize)*100), msg.ChunkEnd)
				if err := c.Request.WebSocket.MessageSendJSON(&wsResp{
					Type:      "uploadingFileChunk",
					Success:   true,
					ChunkEnd:  msg.ChunkEnd,
					Path:      msg.Path,
					TotalSize: msg.TotalSize,
				}); err != nil {
					g.Logger.Errorf("websocket message send error: %s", err.Error())
					break
				}
			}
		// 删除文件
		case "delete":
			if err := sftpClient.Remove(msg.Path); err != nil {
				if err2 := c.Request.WebSocket.MessageSendJSON(&wsResp{
					Type:    "delete",
					Success: false,
					Msg:     err.Error(),
				}); err2 != nil {
					g.Logger.Errorf("websocket message send error: %s", err2.Error())
					break
				}
				g.Logger.Errorf("删除文件失败：", err)
				continue
			}
			if err2 := c.Request.WebSocket.MessageSendJSON(&wsResp{
				Type:    "delete",
				Success: true,
			}); err2 != nil {
				g.Logger.Errorf("websocket message send error: %s", err2.Error())
				break
			}
		}

	}
	g.Logger.Infof("断开 SFTP 文件管理器 Websocket 连接")
	return results.JsonOk()
}
