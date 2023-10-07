package controller

import (
	"eebot/bot/model"
	"eebot/bot/service"
	"fmt"
	"strings"
)

func Repeat(gm model.GroupMessage) (err error) {
	var sentMessage string
	slices := strings.Split(gm.RawMessage, " ")
	prefix := fmt.Sprintf("[CQ:at,qq=%d] ", gm.UserID)
	if len(slices) > 2 {
		sentMessage = strings.Join(slices[2:], " ")
	}

	return service.Reply(sentMessage, prefix, gm.GroupID)
}
