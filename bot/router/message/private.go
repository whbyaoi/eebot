package message

import (
	"eebot/bot/controller"
	"eebot/bot/model"
	"strings"
)

func PrivateMessageHub(pm model.PrivateMessage) (err error) {
	var action string
	slices := strings.Split(CutSpace(pm.RawMessage), " ")
	if len(slices) > 0 {
		action = slices[0]
	}

	switch action {
	case "300":
		err = controller.AnalysisHub(slices, false, pm.UserID, pm.UserID)
	default:

	}

	return
}

