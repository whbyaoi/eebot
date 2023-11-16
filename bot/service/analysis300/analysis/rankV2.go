package analysis

import (
	"eebot/bot/service/analysis300/db"
	"slices"
	"sort"
	"time"
)

var attrs = []string{
	"ActualTotal", "Total", "Win", "WinRate", "AvgHit", "AvgKill", "AvgDeath",
	"AvgAssist", "AvgTower", "AvgPutEye", "AvgDestryEye", "AvgMoney", "AvgMoney",
	"AvgMakeDamage", "AvgTakeDamage", "AvgHitPerMinite", "AvgKillPerMinite", "AvgDeathPerMinite",
	"AvgAssistPerMinite", "AvgTowerPerMinite", "AvgPutEyePerMinite", "AvgDestryEyePerMinite", "AvgMoneyPerMinite",
	"AvgMakeDamagePerMinite", "AvgTakeDamagePerMinite", "AvgMoneyConversionRate", "AvgUsedTime", "AvgJJL", "Score",
}

var weightTran = map[string]int{
	"AvgHitPerMinite":        4,  // 补刀
	"AvgKillPerMinite":       6,  // k
	"AvgDeathPerMinite":      8,  // d
	"AvgAssistPerMinite":     10, // a
	"AvgTowerPerMinite":      12, // 推塔
	"AvgPutEyePerMinite":     14, // 插眼
	"AvgDestryEyePerMinite":  16, // 排眼
	"AvgMoneyPerMinite":      18, // 经济
	"AvgMakeDamagePerMinite": 21, // 输出
	"AvgTakeDamagePerMinite": 24, // 承伤
	"AvgMoneyConversionRate": 26, // 转换率
}

// 数据
type HeroData struct {
	PlayerID               uint64
	HeroID                 int
	ActualTotal            float64 // 实际场次
	Total                  float64 // 参与计算场次
	Win                    float64 // 参与计算场次的胜场
	WinRate                float64 // 参与计算场次的胜率 0.xxxx
	AvgHit                 float64
	AvgKill                float64
	AvgDeath               float64
	AvgAssist              float64
	AvgTower               float64
	AvgPutEye              float64
	AvgDestryEye           float64
	AvgMoney               float64
	AvgMakeDamage          float64
	AvgTakeDamage          float64
	AvgHitPerMinite        float64
	AvgKillPerMinite       float64
	AvgDeathPerMinite      float64
	AvgAssistPerMinite     float64
	AvgTowerPerMinite      float64
	AvgPutEyePerMinite     float64
	AvgDestryEyePerMinite  float64
	AvgMoneyPerMinite      float64
	AvgMakeDamagePerMinite float64
	AvgTakeDamagePerMinite float64
	AvgMoneyConversionRate float64
	AvgUsedTime            float64
	AvgJJL                 float64
	Score                  float64

	Rank Rank
}

func (hd *HeroData) get(attr string) float64 {
	switch attr {
	case "ActualTotal":
		return hd.ActualTotal
	case "Total":
		return hd.Total
	case "Win":
		return hd.Win
	case "WinRate":
		return hd.WinRate
	case "AvgHit":
		return hd.AvgHit
	case "AvgKill":
		return hd.AvgKill
	case "AvgDeath":
		return hd.AvgDeath
	case "AvgAssist":
		return hd.AvgAssist
	case "AvgTower":
		return hd.AvgTower
	case "AvgPutEye":
		return hd.AvgPutEye
	case "AvgDestryEye":
		return hd.AvgDestryEye
	case "AvgMoney":
		return hd.AvgMoney
	case "AvgMakeDamage":
		return hd.AvgMakeDamage
	case "AvgTakeDamage":
		return hd.AvgTakeDamage
	case "AvgHitPerMinite":
		return hd.AvgHitPerMinite
	case "AvgKillPerMinite":
		return hd.AvgKillPerMinite
	case "AvgDeathPerMinite":
		return hd.AvgDeathPerMinite
	case "AvgAssistPerMinite":
		return hd.AvgAssistPerMinite
	case "AvgTowerPerMinite":
		return hd.AvgTowerPerMinite
	case "AvgPutEyePerMinite":
		return hd.AvgPutEyePerMinite
	case "AvgDestryEyePerMinite":
		return hd.AvgDestryEyePerMinite
	case "AvgMoneyPerMinite":
		return hd.AvgMoneyPerMinite
	case "AvgMakeDamagePerMinite":
		return hd.AvgMakeDamagePerMinite
	case "AvgTakeDamagePerMinite":
		return hd.AvgTakeDamagePerMinite
	case "AvgMoneyConversionRate":
		return hd.AvgMoneyConversionRate
	case "AvgUsedTime":
		return hd.AvgUsedTime
	case "AvgJJL":
		return hd.AvgJJL
	case "Score":
		return hd.Score
	default:
		return 0
	}
}

func (hd *HeroData) set(attr string, value float64) {
	switch attr {
	case "ActualTotal":
		hd.ActualTotal = value
	case "Total":
		hd.Total = value
	case "Win":
		hd.Win = value
	case "WinRate":
		hd.WinRate = value
	case "AvgHit":
		hd.AvgHit = value
	case "AvgKill":
		hd.AvgKill = value
	case "AvgDeath":
		hd.AvgDeath = value
	case "AvgAssist":
		hd.AvgAssist = value
	case "AvgTower":
		hd.AvgTower = value
	case "AvgPutEye":
		hd.AvgPutEye = value
	case "AvgDestryEye":
		hd.AvgDestryEye = value
	case "AvgMoney":
		hd.AvgMoney = value
	case "AvgMakeDamage":
		hd.AvgMakeDamage = value
	case "AvgTakeDamage":
		hd.AvgTakeDamage = value
	case "AvgHitPerMinite":
		hd.AvgHitPerMinite = value
	case "AvgKillPerMinite":
		hd.AvgKillPerMinite = value
	case "AvgDeathPerMinite":
		hd.AvgDeathPerMinite = value
	case "AvgAssistPerMinite":
		hd.AvgAssistPerMinite = value
	case "AvgTowerPerMinite":
		hd.AvgTowerPerMinite = value
	case "AvgPutEyePerMinite":
		hd.AvgPutEyePerMinite = value
	case "AvgDestryEyePerMinite":
		hd.AvgDestryEyePerMinite = value
	case "AvgMoneyPerMinite":
		hd.AvgMoneyPerMinite = value
	case "AvgMakeDamagePerMinite":
		hd.AvgMakeDamagePerMinite = value
	case "AvgTakeDamagePerMinite":
		hd.AvgTakeDamagePerMinite = value
	case "AvgMoneyConversionRate":
		hd.AvgMoneyConversionRate = value
	case "AvgUsedTime":
		hd.AvgUsedTime = value
	case "AvgJJL":
		hd.AvgJJL = value
	case "Score":
		hd.Score = value
	}
}

type Rank struct {
	ActualTotal            float64
	Total                  float64
	Win                    float64
	WinRate                float64
	AvgHit                 float64
	AvgKill                float64
	AvgDeath               float64
	AvgAssist              float64
	AvgTower               float64
	AvgPutEye              float64
	AvgDestryEye           float64
	AvgMoney               float64
	AvgMakeDamage          float64
	AvgTakeDamage          float64
	AvgHitPerMinite        float64
	AvgKillPerMinite       float64
	AvgDeathPerMinite      float64
	AvgAssistPerMinite     float64
	AvgTowerPerMinite      float64
	AvgPutEyePerMinite     float64
	AvgDestryEyePerMinite  float64
	AvgMoneyPerMinite      float64
	AvgMakeDamagePerMinite float64
	AvgTakeDamagePerMinite float64
	AvgMoneyConversionRate float64
	AvgUsedTime            float64
	AvgJJL                 float64
	Score                  float64

	PlayerCount float64
}

func (hd *Rank) set(attr string, value float64) {
	switch attr {
	case "ActualTotal":
		hd.ActualTotal = value
	case "Total":
		hd.Total = value
	case "Win":
		hd.Win = value
	case "WinRate":
		hd.WinRate = value
	case "AvgHit":
		hd.AvgHit = value
	case "AvgKill":
		hd.AvgKill = value
	case "AvgDeath":
		hd.AvgDeath = value
	case "AvgAssist":
		hd.AvgAssist = value
	case "AvgTower":
		hd.AvgTower = value
	case "AvgPutEye":
		hd.AvgPutEye = value
	case "AvgDestryEye":
		hd.AvgDestryEye = value
	case "AvgMoney":
		hd.AvgMoney = value
	case "AvgMakeDamage":
		hd.AvgMakeDamage = value
	case "AvgTakeDamage":
		hd.AvgTakeDamage = value
	case "AvgHitPerMinite":
		hd.AvgHitPerMinite = value
	case "AvgKillPerMinite":
		hd.AvgKillPerMinite = value
	case "AvgDeathPerMinite":
		hd.AvgDeathPerMinite = value
	case "AvgAssistPerMinite":
		hd.AvgAssistPerMinite = value
	case "AvgTowerPerMinite":
		hd.AvgTowerPerMinite = value
	case "AvgPutEyePerMinite":
		hd.AvgPutEyePerMinite = value
	case "AvgDestryEyePerMinite":
		hd.AvgDestryEyePerMinite = value
	case "AvgMoneyPerMinite":
		hd.AvgMoneyPerMinite = value
	case "AvgMakeDamagePerMinite":
		hd.AvgMakeDamagePerMinite = value
	case "AvgTakeDamagePerMinite":
		hd.AvgTakeDamagePerMinite = value
	case "AvgMoneyConversionRate":
		hd.AvgMoneyConversionRate = value
	case "AvgUsedTime":
		hd.AvgUsedTime = value
	case "AvgJJL":
		hd.AvgJJL = value
	case "Score":
		hd.Score = value
	}
}

// GetRankFromTop 返回降序评分前top
func GetRankFromTop(HeroID int, fv int, top int) ([]*HeroData, int) {
	slice, sorted := getRank(HeroID, fv)
	if len(sorted["Score"]) < top {
		slices.Reverse[[]*HeroData](sorted["Score"])
		return sorted["Score"], len(slice)
	} else {
		slices.Reverse[[]*HeroData](sorted["Score"][len(slice)-top:])
		return sorted["Score"][len(slice)-top:], len(slice)
	}
}

func GetRankFromPlayers(HeroID int, fv int, PlayerID []uint64) (players map[uint64]*HeroData, n int) {
	slice, _ := getRank(HeroID, fv)
	players = map[uint64]*HeroData{}
	for i := range slice {
		if slices.Contains(PlayerID, slice[i].PlayerID) {
			players[slice[i].PlayerID] = slice[i]
		}
	}
	return players, len(slice)
}

func CalculateData(idToData map[uint64][]*db.Player, fv int, HeroID int) (heroDataSlice []*HeroData) {
	for playerID, plays := range idToData {
		heroData := &HeroData{
			PlayerID:    playerID,
			HeroID:      HeroID,
			ActualTotal: float64(len(plays)),
		}
		for _, play := range plays {
			if play.FV < fv {
				break
			}
			if heroData.Total >= 50 {
				break
			}
			heroData.Total++
			if play.Result == 1 || play.Result == 3 {
				heroData.Win++
			}
			minute := float64(play.UsedTime) / 60
			heroData.AvgHit += float64(play.KillUnit)
			heroData.AvgKill += float64(play.KillPlayer)
			heroData.AvgDeath += float64(play.Death)
			heroData.AvgAssist += float64(play.Assist)
			heroData.AvgTower += float64(play.DestoryTower)
			heroData.AvgPutEye += float64(play.PutEyes)
			heroData.AvgDestryEye += float64(play.DestoryEyes)
			heroData.AvgMoney += float64(play.TotalMoney)
			heroData.AvgMakeDamage += play.MakeDamagePercent * float64(play.MakeDamageSide)
			heroData.AvgTakeDamage += float64(play.TakeDamageSide) * play.TakeDamagePercent
			heroData.AvgHitPerMinite += float64(play.KillUnit) / minute
			heroData.AvgKillPerMinite += float64(play.KillPlayer) / minute
			heroData.AvgDeathPerMinite += float64(play.Death) / minute
			heroData.AvgAssistPerMinite += float64(play.Assist) / minute
			heroData.AvgTowerPerMinite += float64(play.DestoryTower) / minute
			heroData.AvgPutEyePerMinite += float64(play.PutEyes) / minute
			heroData.AvgDestryEyePerMinite += float64(play.DestoryEyes) / minute
			heroData.AvgMoneyPerMinite += float64(play.TotalMoney) / minute
			heroData.AvgMakeDamagePerMinite += play.MakeDamagePercent * float64(play.MakeDamageSide) / minute
			heroData.AvgTakeDamagePerMinite += float64(play.TakeDamageSide) * play.TakeDamagePercent / minute
			heroData.AvgMoneyConversionRate += float64(play.MakeDamageSide) * play.MakeDamagePercent / float64(play.TotalMoney) * 100
			heroData.AvgJJL += float64(play.FV)
			heroData.AvgUsedTime += float64(play.UsedTime)
		}
		if heroData.Total >= 5 {
			heroDataSlice = append(heroDataSlice, heroData)
		}
	}
	for _, heroData := range heroDataSlice {
		heroData.WinRate = heroData.Win / heroData.Total
		heroData.AvgHit /= heroData.Total
		heroData.AvgKill /= heroData.Total
		heroData.AvgDeath /= heroData.Total
		heroData.AvgAssist /= heroData.Total
		heroData.AvgTower /= heroData.Total
		heroData.AvgPutEye /= heroData.Total
		heroData.AvgDestryEye /= heroData.Total
		heroData.AvgMoney /= heroData.Total
		heroData.AvgMakeDamage /= heroData.Total
		heroData.AvgTakeDamage /= heroData.Total
		heroData.AvgHitPerMinite /= heroData.Total
		heroData.AvgKillPerMinite /= heroData.Total
		heroData.AvgDeathPerMinite /= heroData.Total
		heroData.AvgAssistPerMinite /= heroData.Total
		heroData.AvgTowerPerMinite /= heroData.Total
		heroData.AvgPutEyePerMinite /= heroData.Total
		heroData.AvgDestryEyePerMinite /= heroData.Total
		heroData.AvgMoneyPerMinite /= heroData.Total
		heroData.AvgMakeDamagePerMinite /= heroData.Total
		heroData.AvgTakeDamagePerMinite /= heroData.Total
		heroData.AvgMoneyConversionRate /= heroData.Total
		heroData.AvgUsedTime /= heroData.Total
		heroData.AvgJJL /= heroData.Total
	}
	return
}

func getRank(HeroID int, fv int) ([]*HeroData, map[string][]*HeroData) {
	var players []db.Player
	start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Unix() - ExpiryDate
	db.SqlDB.Model(db.Player{}).Where("create_time > ? and hero_id = ?", start, HeroID).Order("create_time desc").Find(&players)
	idToRecord := map[uint64][]*db.Player{}
	for i := range players {
		idToRecord[players[i].PlayerID] = append(idToRecord[players[i].PlayerID], &players[i])
	}

	heroDataSlice := CalculateData(idToRecord, fv, HeroID)
	clear(players)
	clear(idToRecord)
	// 排序
	sortedData := map[string][]*HeroData{}
	for _, attr := range attrs[:len(attrs)-1] {
		sort.Slice(heroDataSlice, func(i, j int) bool {
			return heroDataSlice[i].get(attr) < heroDataSlice[j].get(attr)
		})
		tmp := make([]*HeroData, len(heroDataSlice))
		copy(tmp, heroDataSlice)
		sortedData[attr] = tmp
	}
	// 计算权重
	heroWeight := MergeImportance(HeroFactor[db.HeroIDToName[HeroID]])
	for attr, s := range sortedData {
		for rank, data := range s {
			data.Score += float64(rank) / float64(len(s)-1) * heroWeight[weightTran[attr]] * 100
			data.Rank.set(attr, float64(rank)/float64(len(s)-1)*100)
		}
	}
	for _, data := range heroDataSlice {
		data.Score = data.Score * (0.7 + min(data.Total, MaxPlayTimes)/MaxPlayTimes*0.3) * (0.7 + data.Rank.WinRate/100*0.3)
	}
	attr := "Score"
	sort.Slice(heroDataSlice, func(i, j int) bool {
		return heroDataSlice[i].get(attr) < heroDataSlice[j].get(attr)
	})
	tmp := make([]*HeroData, len(heroDataSlice))
	copy(tmp, heroDataSlice)
	sortedData[attr] = tmp
	for rank, heroData := range sortedData[attr] {
		heroData.Rank.set(attr, float64(rank)/float64(len(sortedData[attr])-1)*100)
	}
	return heroDataSlice, sortedData
}
