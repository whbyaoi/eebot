package message

import (
	"eebot/bot/controller"
	"eebot/bot/model"
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
		// 无视错误消息
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
