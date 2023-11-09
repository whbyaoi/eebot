package message

import (
	"eebot/bot/controller"
	"eebot/bot/model"
	"eebot/bot/service"
	"fmt"
	"strings"
)

func GroupMessageHub(gm model.GroupMessage) (err error) {
	var action string
	slices := strings.Split(CutSpace(gm.RawMessage), " ")
	if len(slices) > 1 {
		action = slices[1]
	}
	switch action {
	case "复读机", "repeat":
		err = controller.Repeat(gm)
	case "300":
		err = controller.AnalysisHub(slices[1:], true, gm.UserID, gm.GroupID)
	default:
		prefix := fmt.Sprintf("[CQ:at,qq=%d] \n", gm.UserID)
		msg := "未知服务名：" + slices[1]
		msg += "\n 目前支持服务名：300"
		service.Reply(msg, prefix, gm.GroupID)
	}
	return
}

// CutSpace 替换多余的空格
func CutSpace(raw string) string {
	for strings.Contains(raw, "  ") {
		raw = strings.ReplaceAll(raw, "  ", " ")
	}
	return raw
}
