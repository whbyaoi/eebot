package model

type Request struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

type GroupMessageParams struct {
	GroupID    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

type PrivateMessageParams struct {
	GroupMessageParams

	UserID int64 `json:"user_id"`
}
