package service

import (
	"eebot/bot/model"
	"eebot/ws"
)

// Reply 回复消息
//
//	msg: 消息主体
//	prefix: 群聊at动作前缀，私聊为空
//	id: 群id或者qq名id
func Reply(msg string, prefix string, id int64) error {
	var req model.Request
	if prefix != "" {
		req = model.Request{
			Action: "send_group_msg",
			Params: model.GroupMessageParams{
				GroupID:    id,
				Message:    msg,
				AutoEscape: false,
			},
		}
	} else {
		req = model.Request{
			Action: "send_private_msg",
			Params: model.PrivateMessageParams{
				GroupMessageParams: model.GroupMessageParams{
					Message:    msg,
					AutoEscape: false,
				},
				UserID: id,
			},
		}
	}

	return ws.Send(req)
}
