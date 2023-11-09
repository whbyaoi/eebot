package router

import (
	"eebot/bot/model"
	"eebot/bot/router/message"
	"eebot/g"
	"eebot/ws"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func WsMessageHandler(b []byte) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		g.Logger.Errorf("Unmarshal to map[string]interface{} 错误")
		return
	}

	// 应该是确认返回消息
	if _, ok := m["post_type"]; !ok {
		return
	}
	postType, ok := m["post_type"].(string)
	if !ok {
		g.Logger.Errorf("ws消息格式错误：post_type 字段值无法转换为string类型")
		return
	}

	switch postType {
	// 暂时只处理message消息
	case "message":
		_ = messageHandler(b)
	default:
	}
}

func messageHandler(b []byte) (err error) {
	var messageBase model.MessageBase
	err = json.Unmarshal(b, &messageBase)
	if err != nil {
		return err
	}

	t0 := time.Now()
	var source int64
	if messageBase.MessageType == "group" {
		if !isAtMe(messageBase.RawMessage, messageBase.SelfID) || isDev() {
			return
		}
		var groupMessage model.GroupMessage
		err = json.Unmarshal(b, &groupMessage)
		if err != nil {
			return err
		}
		g.Logger.Printf("收到 %d 群聊消息：%s", groupMessage.GroupID, groupMessage.RawMessage)
		source = groupMessage.GroupID
		err = message.GroupMessageHub(groupMessage)
	} else {
		var privateMessage model.PrivateMessage
		err = json.Unmarshal(b, &privateMessage)
		if err != nil {
			return err
		}
		g.Logger.Printf("收到 %d 私聊消息：%s", privateMessage.UserID, privateMessage.RawMessage)
		source = privateMessage.UserID
		err = message.PrivateMessageHub(privateMessage)
	}

	// 推送
	if g.Config.GetBool("report") {
		go func() {
			ws.Send(model.Request{
				Action: "send_private_msg",
				Params: model.PrivateMessageParams{
					GroupMessageParams: model.GroupMessageParams{
						Message:    fmt.Sprintf("%s: 处理 %d 群聊消息 %s 完毕", time.Now().Format(time.TimeOnly), source, messageBase.RawMessage),
						AutoEscape: false,
					},
					UserID: g.Config.GetInt64("report-id"),
				},
			})
		}()
	}
	if err != nil {
		g.Logger.Errorf("处理来自 %d 的消息 %s 时错误：%s，耗时 %v", source, messageBase.RawMessage, err.Error(), time.Since(t0))
	} else {
		g.Logger.Infof("处理来自 %d 的消息完毕：%s，耗时 %v", source, messageBase.RawMessage, time.Since(t0))
	}

	return
}

func isAtMe(rawMessage string, qq int64) bool {
	return strings.HasPrefix(rawMessage, fmt.Sprintf("[CQ:at,qq=%d]", qq))
}

func isDev() bool {
	return g.Config.GetBool("dev")
}

// FormatJson 格式化Json以便更容器查看, 如果m格式错误则返回空字符串
func FormatJson(m interface{}, indent bool) string {
	var b []byte
	var err error
	if !indent {
		b, err = json.Marshal(m)
	} else {
		b, err = json.MarshalIndent(m, "", "  ")
	}
	if err != nil {
		return ""
	}
	return string(b)
}
