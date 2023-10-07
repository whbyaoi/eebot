package model

// PrivateMessage 私聊消息
type PrivateMessage struct {
	MessageBase

	TargetID   int64 `json:"target_id"`   // 接收者 QQ 号
	TempSource int   `json:"temp_source"` // 临时会话来源
}

// GroupMessage 群聊消息
type GroupMessage struct {
	MessageBase

	GroupID   int64     `json:"group_id"`  //群号
	Anonymous anonymous `json:"anonymous"` // 匿名信息, 如果不是匿名消息则为 null
}

type anonymous struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}
