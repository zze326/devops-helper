package resource

import (
	"github.com/revel/revel"
	"github.com/zze326/devops-helper/app/g"
	o_resource "github.com/zze326/devops-helper/app/models/orm/resource"
	"github.com/zze326/devops-helper/app/results"
	"io"
	"os"
	"time"
	"unicode/utf8"
)

func (c Host) TerminalSessionLog(id int) revel.Result {
	g.Logger.Infof("建立回放 Web 终端记录 Websocket 连接")

	defer func() {
		if err := recover(); err != nil {
			g.Logger.Errorf("err: %v", err)
		}
	}()

	sessionModel := new(o_resource.HostTerminalSession)
	if err := c.DB.Where("id = ?", id).First(sessionModel).Error; err != nil {
		return results.JsonError(err)
	}

	file, err := os.Open(sessionModel.Filepath)
	if err != nil {
		g.Logger.Errorf("打开会话文件失败：%v", err)
		return results.JsonError(err)
	}

	type msgObj struct {
		Type     string  `json:"type"`
		Speed    float32 `json:"speed"`
		Progress float32 `json:"progress"`
	}

	type respObj struct {
		Total int64  `json:"total"`
		Sent  int64  `json:"sent"`
		Data  string `json:"data"`
		Clear bool   `json:"clear"`
	}

	closeChan := make(chan bool)
	stopChan := make(chan bool)
	resumeChan := make(chan bool)
	speedChan := make(chan float32)
	progressChan := make(chan float32)

	defer func() {
		defer close(closeChan)
		defer close(stopChan)
		defer close(resumeChan)
		defer close(speedChan)
		defer close(progressChan)
	}()

	go func() {
		for {
			msg := new(msgObj)
			if err := c.Request.WebSocket.MessageReceiveJSON(&msg); err != nil {
				g.Logger.Errorf("接收消息失败: %v", err)
				break
			}
			if msg.Type == "close" {
				closeChan <- true // 当退出循环时，向通道发送信号
				break
			}
			if msg.Type == "pause" {
				stopChan <- true
			}
			if msg.Type == "continue" {
				resumeChan <- true
			}
			if msg.Type == "speed" {
				if msg.Speed > 0 {
					g.Logger.Infof("speed: %v", msg.Speed)
					speedChan <- msg.Speed
				}
			}
			if msg.Type == "progress" {
				g.Logger.Infof("progress: %v", msg.Progress)
				progressChan <- msg.Progress
			}
		}
	}()

	stat, err := file.Stat()
	if err != nil {
		return results.JsonError(err)
	}

	g.Logger.Infof("%d", stat.Size())

	var (
		total            = stat.Size()
		sent     int64   = 0
		speed    float32 = 1
		oneSend  int64   = 24
		msgBytes         = make([]byte, oneSend)
	)
loop:
	for {
		select {
		case <-closeChan:
			break loop // 关闭循环
		case <-stopChan:
			<-resumeChan
		case speed = <-speedChan:
		case progress := <-progressChan:
			sendPoint := int64(float32(total) * progress)
			tmpMsgBytes := make([]byte, sendPoint)
			if _, err := file.Seek(0, 0); err != nil {
				g.Logger.Errorf("%v", err)
				break loop
			}

			if _, err := io.ReadFull(file, tmpMsgBytes); err != nil {
				g.Logger.Errorf("%v", err)
				break loop
			}

			if err := c.Request.WebSocket.MessageSendJSON(respObj{
				Total: total,
				Sent:  sent,
				Data:  string(tmpMsgBytes),
				Clear: true,
			}); err != nil {
				g.Logger.Errorf("%v", err)
				break loop
			}
			//allMsgBytes = fileBytes[sendPoint:]
			sent = sendPoint
		default:
			sendLen, err := io.ReadFull(file, msgBytes)
			if err == io.EOF {
				break loop
			} else if err != nil && err != io.ErrUnexpectedEOF {
				return results.JsonError(err)
			}

			sent += int64(sendLen)

			if !utf8.Valid(msgBytes) {
				c.writeBuffer.Write(msgBytes)
				continue loop
			} else {
				if c.writeBuffer.Len() > 0 {
					c.writeBuffer.Write(msgBytes)
					if !utf8.Valid(c.writeBuffer.Bytes()) {
						continue loop
					}
					msgBytes = c.writeBuffer.Bytes()
					c.writeBuffer.Reset()
				}
			}

			if err := c.Request.WebSocket.MessageSendJSON(respObj{
				Total: total,
				Sent:  sent,
				Data:  string(msgBytes[:sendLen]),
			}); err != nil {
				g.Logger.Errorf("%v", err)
				break loop
			}

			time.Sleep(time.Duration((100 / speed) * float32(time.Millisecond)))
		}
	}
	g.Logger.Infof("关闭回放 Web 终端记录 Websocket 连接")
	return results.JsonOk()
}
