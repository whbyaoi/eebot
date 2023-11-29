package db

type PlayerPartition struct {
	ID      uint64 `gorm:"primaryKey"`
	MatchID string

	UsedTime             uint64  // 战绩所用时间
	CreateTime           uint64  // 游戏结束时的时间戳
	PlayerID             uint64  `json:"PlayerID"`    // 玩家id
	Name                 string  `json:"RN"`          // 名称
	AM                   []int   `json:"AM" gorm:"-"` // AwardMoney 8个一组，一组代表一分钟金钱的变化。其中 0 到 7 分别代表： 0-自然金钱增长，1-补刀，8-总
	Ep                   []int   `json:"Ep" gorm:"-"`
	KM                   []int   `json:"KM" gorm:"-"`
	MD                   []int   `json:"MD" gorm:"-"`
	SummonerLevel        int     `json:"SummonerLevel" gorm:"-"`
	HeroID               int     `json:"HeroID" gorm:"primaryKey"`
	HeroLv               int     `json:"HeroLv"`
	Side                 int     `json:"Side"`
	Result               int     `json:"Result"` // 1-win, 2-lose
	FirstWin             int     `json:"FirstWin"`
	SummonerSkill1       int     `json:"SummonerSkill1"`
	SummonerSkill2       int     `json:"SummonerSkill2"`
	TotalMoney           int     `json:"TotalMoney"`
	KillUnit             int     `json:"KillUnit"`
	KillPlayer           int     `json:"KillPlayer"`
	Death                int     `json:"Death"`
	Assist               int     `json:"Assist"`
	ConKillMax           int     `json:"ConKillMax"`
	MulKillMax           int     `json:"MulKillMax"`
	DestoryTower         int     `json:"DestoryTower"`
	Treat                int     `json:"Treat"`
	PutEyes              int     `json:"PutEyes"`
	DestoryEyes          int     `json:"DestoryEyes"`
	AwardExp             int     `json:"AwardExp" gorm:"-"`
	JumpMax              int     `json:"JumpMax" gorm:"-"`
	WinCount             int     `json:"WinCount" gorm:"-"`
	TotalCount           int     `json:"TotalCount" gorm:"-"`
	KGGoal               int     `json:"KGGoal" gorm:"-"`
	IsSurrender          int     `json:"IsSurrender" gorm:"-"`
	Elo                  int     `json:"Elo"`
	LUpTime              []int   `json:"LUpTime" gorm:"-"`
	LUpSO                []int   `json:"LUpSO" gorm:"-"`
	SD                   []int   `json:"SD" gorm:"-"`
	ED                   []int   `json:"ED" gorm:"-"`
	TD                   []int   `json:"TD" gorm:"-"`
	PD                   []int   `json:"PD" gorm:"-"`
	PE                   []int   `json:"PE" gorm:"-"`
	MG                   []int   `json:"MG" gorm:"-"`
	MS                   []int   `json:"MS" gorm:"-"`
	HS                   int     `json:"HS" gorm:"-"`
	Guid                 int     `json:"Guid" gorm:"-"`
	FV                   int     `json:"FV"` // fight value 竞技力
	TT                   []int   `json:"TT" gorm:"-"`
	M25071               int     `json:"M25071" gorm:"-"`
	M25072               int     `json:"M25072" gorm:"-"`
	M25078               int     `json:"M25078" gorm:"-"`
	TotalMoneySide       int     `json:"TotalMoneySide"`
	TotalMoneyPercent    float64 `json:"TotalMoneyPercent"`
	MakeDamageSide       int     `json:"MakeDamageSide"`
	MakeDamagePercent    float64 `json:"MakeDamagePercent"`
	TakeDamageSide       int     `json:"TakeDamageSide"`
	TakeDamagePercent    float64 `json:"TakeDamagePercent"`
	DamageConversionRate float64 `json:"DamageConversionRate"`
}

func ToPartition(p Player) PlayerPartition {
	return PlayerPartition{
		ID:                   p.ID,
		MatchID:              p.MatchID,
		UsedTime:             p.UsedTime,
		CreateTime:           p.CreateTime,
		PlayerID:             p.PlayerID,
		Name:                 p.Name,
		HeroID:               p.HeroID,
		HeroLv:               p.HeroLv,
		Side:                 p.Side,
		Result:               p.Result,
		FirstWin:             p.FirstWin,
		SummonerSkill1:       p.SummonerSkill1,
		SummonerSkill2:       p.SummonerSkill2,
		TotalMoney:           p.TotalMoney,
		KillUnit:             p.KillUnit,
		KillPlayer:           p.KillPlayer,
		Death:                p.Death,
		Assist:               p.Assist,
		ConKillMax:           p.ConKillMax,
		MulKillMax:           p.MulKillMax,
		DestoryTower:         p.DestoryTower,
		Treat:                p.Treat,
		PutEyes:              p.PutEyes,
		DestoryEyes:          p.DestoryEyes,
		Elo:                  p.Elo,
		FV:                   p.FV,
		TotalMoneySide:       p.TotalMoneySide,
		TotalMoneyPercent:    p.TotalMoneyPercent,
		MakeDamageSide:       p.MakeDamageSide,
		MakeDamagePercent:    p.MakeDamagePercent,
		TakeDamageSide:       p.TakeDamageSide,
		TakeDamagePercent:    p.TakeDamagePercent,
		DamageConversionRate: p.DamageConversionRate,
	}
}

func (PlayerPartition) TableName() string {
	return "players_partition"
}

type MatchPartition struct {
	ID         uint64            `gorm:"primaryKey"`
	MatchID    string            `json:"MTID"`                                                   // 比赛id
	MatchIDInt uint64            `json:"MTIDInt"`                                                // int格式比赛id，不知道为什么和上面的不一样
	MID        uint64            `json:"MID"`                                                    // 比赛类型
	UsedTime   uint64            `json:"UsedTime"`                                               // 所用时间
	Players    []PlayerPartition `json:"Players" gorm:"foreignKey:match_id;references:match_id"` // 玩家
	CreateTime uint64            `json:"CreateTime"`                                             // 游戏结束时的时间戳
	KillTrees  string            `json:"KillTrees" gorm:"-"`                                     // 击杀树
}

func (MatchPartition) TableName() string {
	return "matches"
}
