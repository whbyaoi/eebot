package model

type base struct {
	Time     int64  `json:"time"`      // 事件发生的unix时间戳
	SelfID   int64  `json:"self_id"`   // 收到事件的机器人的 QQ 号
	PostType string `json:"post_type"` // 表示该上报的类型, 消息, 消息发送, 请求, 通知, 或元事件
}

type MessageBase struct {
	base

	MessageType string      `json:"message_type"` // 消息类型
	SubType     string      `json:"sub_type"`     // 表示消息的子类型
	MessageID   int32       `json:"message_id"`   // 消息 ID
	UserID      int64       `json:"user_id"`      // 发送者 QQ 号
	Message     interface{} `json:"-"`            // 一个消息链(不知道是啥意思暂且屏蔽)
	RawMessage  string      `json:"raw_message"`  // CQ 码格式的消息
	Font        int         `json:"font"`         // 字体
	Sender      Sender      `json:"sender"`       // 发送者信息
}

type RequestBase struct {
	base
	RequestType string `json:"request_type"` // 请求类型
}

type NoticeBase struct {
	base
	NoticeType string `json:"notice_type"` // 通知类型
}

type MetaEventBase struct {
	MetaEventType string `json:"meta_event_type"` // 元数据类型
}

type Sender struct {
	UserID   int64  `json:"user_id"`  // 发送者 QQ 号
	NickName string `json:"nickname"` // 昵称
	Sex      string `json:"sex"`      // 性别, male 或 female 或 unknown
	Age      int32  `json:"age"`      // 年龄

	// 临时qq群特有
	GroupID int64 `json:"group_id"` // 临时群消息来源群号

	// 群聊特有
	Card  string `json:"card"`  // 群名片／备注
	Area  string `json:"area"`  // 地区
	Level string `json:"level"` // 成员等级
	Role  string `json:"role"`  // 角色, owner 或 admin 或 member
	Title string `json:"title"` // 专属头衔
}
