/***
只负责ws的建立，收发消息
***/

package ws

import (
	"eebot/g"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/gorilla/websocket"
)

var WsClient *websocket.Conn

func InitWebsocket() (err error) {
	host := g.Config.GetString("cqhttp")
	u := url.URL{
		Scheme: "ws",
		Host:   host,
	}

	header := http.Header{}
	header.Set("Authorization", g.Config.GetString("access-token"))

	g.Logger.Infof("开始建立ws连接：ws://%s", host)
	WsClient, _, err = websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		g.Logger.Errorf("建立ws失败：%s", err.Error())
		return
	}
	g.Logger.Infof("建立ws成功，开始处理消息")
	return
}

func Read(msgHandle func([]byte)) (err error) {
	defer func() {
		if r := recover(); r != nil {
			g.Logger.Errorf("panic: %s", string(debug.Stack()))
			err = fmt.Errorf("panic: %s", string(debug.Stack()))
		}
	}()
	for {
		_, message, err := WsClient.ReadMessage()
		if err != nil {
			g.Logger.Errorf("ws读取消息错误：%s", err.Error())
			return err
		}
		go msgHandle(message)
	}
}

func Send(req interface{}) error {
	err := WsClient.WriteJSON(req)
	if err != nil {
		g.Logger.Errorf("ws发送消息错误：%s", err.Error())
	}
	return err
}
