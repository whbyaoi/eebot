package model

type Heart struct {
	MetaEventBase

	Status   interface{} `json:"status"`   // 应用程序状态(这里不需要所以为空接口)
	Interval int64       `json:"interval"` // 距离上一次心跳包的时间(单位是毫秒)
}

type status struct {
	AppInitialized bool
	AppEnabled     bool
	PluginsGood    bool
	AppGood        bool
	Online         bool
	Stat           interface{}
}

type LifeCycle struct {
	MetaEventBase

	SubType string `json:"sub_type"` // 子类型 enable, disable, connect
}
