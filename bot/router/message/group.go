package message

import (
	"eebot/bot/controller"
	"eebot/bot/model"
	"strings"
)

func GroupMessageHub(gm model.GroupMessage) (err error) {
	var action string
	slices := strings.Split(gm.RawMessage, " ")
	if len(slices) > 1 {
		action = slices[1]
	}
	switch action {
	case "复读机", "repeat":
		err = controller.Repeat(gm)
	case "300":
		err = controller.AnalysisHub(slices[1:], true, gm.GroupID)
	default:
		// 无视错误消息
	}

	return
}
