package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/revel/revel"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"golang.org/x/net/context"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"strings"
	"unicode/utf8"
)

type Kubernetes struct {
	gormc.Controller
	resizeEvent         chan remotecommand.TerminalSize
	cancelWebsocketFunc context.CancelFunc
	writeBuffer         bytes.Buffer
	readBuffer          bytes.Buffer
}

// Terminal Web 终端
func (c Kubernetes) Terminal(namespace, podName, containerName string, configID int) revel.Result {
	c.Log.Info("建立 Websocket 连接")
	var (
		kubeConfig string
	)

	if err := c.DB.Raw("select config from k8s_configs where id = ?", configID).Scan(&kubeConfig).Error; err != nil {
		return results.JsonError(err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		c.Log.Errorf("Failed to create Kubernetes config: %+v", err)
		return nil
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.Log.Errorf("Failed to create Kubernetes clientset: %+v", err)
		return nil
	}
	req := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"sh", "-c", "TERM=xterm-256color /bin/bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		c.Log.Errorf("Failed to create Kubernetes exec: %+v", err)
		return nil
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	c.cancelWebsocketFunc = cancelFunc

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  &c,
		Stdout: &c,
		Stderr: &c,
		Tty:    true,
	})
	if err != nil && err != context.Canceled {
		c.Log.Errorf("Failed to execute command in Kubernetes Pod: %+v", err)
		return nil
	}

	c.Log.Info("断开 Websocket 连接")
	return nil
}

// Read 接收 Web 终端的命令
func (c *Kubernetes) Read(p []byte) (n int, err error) {
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
			c.Log.Infof("输入命令：%s", cmdStr)
			c.readBuffer.Reset()
		}
		c.readBuffer.Write([]byte(xtermMsg.Input))
		copy(p, fmt.Sprintf("%s", xtermMsg.Input))
		return len(xtermMsg.Input), nil
	} else if xtermMsg.MsgType == "resize" {
		c.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else {
		// 关闭连接
		c.cancelWebsocketFunc()
	}

	return 0, nil
}

// Write 响应到 Web 终端
func (c *Kubernetes) Write(p []byte) (n int, err error) {
	msgBytes := p
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

func (c *Kubernetes) Next() (size *remotecommand.TerminalSize) {
	ret := <-c.resizeEvent
	size = &ret
	return
}

func (c Kubernetes) TailLog(namespace, podName, containerName string, configID int) revel.Result {
	c.Log.Info("建立 Websocket 连接")
	var (
		kubeConfig string
	)

	if err := c.DB.Raw("select config from k8s_configs where id = ?", configID).Scan(&kubeConfig).Error; err != nil {
		return results.JsonError(err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		c.Log.Errorf("Failed to create Kubernetes config: %+v", err)
		return nil
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.Log.Errorf("Failed to create Kubernetes clientset: %+v", err)
		return nil
	}

	line := int64(300)
	req := clientSet.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{
		Container: containerName,
		Follow:    true,
		TailLines: &line,
	})
	ctx, cancelFunc := context.WithCancel(context.Background())
	c.cancelWebsocketFunc = cancelFunc
	stream, err := req.Stream(ctx)
	if err != nil {
		c.Log.Errorf("Failed to create log stream: %+v", err)
		return nil
	}

	scanner := bufio.NewScanner(stream)
	defer func(stream io.ReadCloser) {
		err := stream.Close()
		if err != nil {
			c.Log.Errorf("Failed to close log stream: %+v", err)
		}
	}(stream)

	go func() {
		var msg string
		_ = c.Request.WebSocket.MessageReceive(&msg)
		cancelFunc()
	}()

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			c.Log.Infof("cancelled")
			cancelFunc()
		default:
			err := c.Request.WebSocket.MessageSend(scanner.Text())
			if err != nil {
				c.Log.Errorf("Failed to send log message err: %+v", err)
				cancelFunc()
			}
		}
	}

	c.Log.Info("断开 Websocket 连接")
	return nil
}
