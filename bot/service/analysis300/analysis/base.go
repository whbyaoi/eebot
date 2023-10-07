package analysis

import (
	"eebot/bot/service/analysis300/db"
	"fmt"
	"sort"
)

var PlayerListKey = "300analysis:player_list"

// WinOrLoseAnalysis 胜负分析
//
//	return: result [][3]int
//			- result[*][0] -- 己方均分
//			- result[*][1] -- 对方均分
//			- result[*][2] -- 输赢，1=赢，2=输
//
//			diff    相对均分偏移
//			fvRange 竞技力范围
func WinOrLoseAnalysis(PlayerID uint64) (result [][3]int, diff int, fvRange [2]int, fvNow int) {
	matchIds, sides := getMatchIdsAndSides(PlayerID)

	fvRange[0] = 2500
	fvRange[1] = 0
	var maxTimestamps uint64
	for i := range matchIds {
		var selfFV int
		var localPlayers []db.Player
		db.SqlDB.Model(&db.Player{}).Where("match_id = ?", matchIds[i]).Find(&localPlayers)

		// TODO 效率有点低
		var match db.Match
		db.SqlDB.Model(&db.Match{}).Where("match_id = ?", matchIds[i]).Find(&match)

		var tmp [3]int
		fvSum1 := 0 // 己方团分
		fvSum2 := 0 // 对面团分
		for j := range localPlayers {
			if localPlayers[j].PlayerID == PlayerID {
				tmp[2] = localPlayers[j].Result
				selfFV = localPlayers[j].FV
				fvRange[0] = min(localPlayers[j].FV, fvRange[0])
				fvRange[1] = max(localPlayers[j].FV, fvRange[1])
				if match.CreateTime >= maxTimestamps {
					maxTimestamps = match.CreateTime
					fvNow = selfFV
				}
			}
			if localPlayers[j].Side == sides[i] {
				fvSum1 += localPlayers[j].FV
			} else {
				fvSum2 += localPlayers[j].FV
			}
		}
		tmp[0] = fvSum1 / 7
		tmp[1] = fvSum2 / 7

		diff += selfFV - (tmp[0]+tmp[1])/2
		result = append(result, tmp)
	}
	if len(matchIds) != 0 {
		diff /= len(matchIds)
	}
	return
}

// HeroAnalysis 英雄分析(常用英雄及其数据)
//
//	return:
//		result[*][0] -- 英雄id
//		result[*][1] -- 场次
//		result[*][2] -- 胜场
//		result[*][3] -- 场均补刀
//		result[*][4] -- 场均每分均刀
//		result[*][5] -- 场均击杀
//		result[*][6] -- 场均每分均击杀
//		result[*][7] -- 场均死亡
//		result[*][8] -- 场均每分均死亡
//		result[*][9] -- 场均助攻
//		result[*][10] -- 场均每分均助攻
//		result[*][11] -- 场均推塔
//		result[*][12] -- 场均每分均推塔
//		result[*][13] -- 场均插眼
//		result[*][14] -- 场均每分均插眼
//		result[*][15] -- 场均排眼
//		result[*][16] -- 场均每分均排眼
//		result[*][17] -- 场均经济
//		result[*][18] -- 场均每分均经济
//		result[*][19] -- 场均经济占比
//		result[*][20] -- 场均输出
//		result[*][21] -- 场均每分均输出
//		result[*][22] -- 场均输出占比
//		result[*][23] -- 场均承伤
//		result[*][24] -- 场均每分均承伤
//		result[*][25] -- 场均承伤占比
//		result[*][26] -- 场均转换率
//		result[*][27] -- 场均耗时
func HeroAnalysis(PlayerID uint64, fv int) (result [][28]float64, total uint64) {
	var players []db.Player
	db.SqlDB.Model(&db.Player{}).Where("player_id = ? and fv >= ?", PlayerID, fv).Find(&players)
	total = 0

	addValue := func(v [28]float64, me db.Player) [28]float64 {
		if me.FV < fv {
			return v
		}
		minute := float64(me.UsedTime) / 60
		v[1]++
		total++
		if me.Result == 1 || me.Result == 3 {
			v[2]++
		}
		v[3] += float64(me.KillUnit)
		v[4] += float64(me.KillUnit) / minute
		v[5] += float64(me.KillPlayer)
		v[6] += float64(me.KillPlayer) / minute
		v[7] += float64(me.Death)
		v[8] += float64(me.Death) / minute
		v[9] += float64(me.Assist)
		v[10] += float64(me.Assist) / minute
		v[11] += float64(me.DestoryTower)
		v[12] += float64(me.DestoryTower) / minute
		v[13] += float64(me.PutEyes)
		v[14] += float64(me.PutEyes) / minute
		v[15] += float64(me.DestoryEyes)
		v[16] += float64(me.DestoryEyes) / minute
		v[17] += float64(me.TotalMoney)
		v[18] += float64(me.TotalMoney) / minute
		v[19] += me.TotalMoneyPercent
		v[20] += float64(me.MakeDamageSide) * me.MakeDamagePercent
		v[21] += float64(me.MakeDamageSide) * me.MakeDamagePercent / minute
		v[22] += me.MakeDamagePercent
		v[23] += float64(me.TakeDamageSide) * me.TakeDamagePercent
		v[24] += float64(me.TakeDamageSide) * me.TakeDamagePercent / minute
		v[25] += me.TakeDamagePercent
		v[26] += float64(me.MakeDamageSide) * me.MakeDamagePercent / float64(me.TotalMoney) * 100
		v[27] += float64(me.UsedTime)
		return v
	}

	dataMap := make(map[int][28]float64)
	for _, me := range players {
		if v, ok := dataMap[me.HeroID]; ok {
			dataMap[me.HeroID] = addValue(v, me)
		} else {
			v := [28]float64{}
			dataMap[me.HeroID] = addValue(v, me)
		}
	}

	for heroID, data := range dataMap {
		data[0] = float64(heroID)
		for i := 3; i <= 27; i++ {
			data[i] = data[i] / data[1]
		}
		result = append(result, data)
	}

	sort.Slice(result, func(i, j int) bool { return result[i][1] >= result[j][1] })
	return
}

// TeamAnalysis 开黑分析
//
//	sortedAllies[*][0]: 玩家id
//	sortedAllies[*][1]: 出现次数
func TeamAnalysis(PlayerID uint64) (sortedAllies [][2]uint64, sortedEnermies [][2]uint64, total int) {
	matchIds, sides := getMatchIdsAndSides(PlayerID)

	total = len(matchIds)
	allyCount := map[uint64]int{}
	enermyCount := map[uint64]int{}
	for i := range matchIds {
		// 找到这局的玩家
		var localPlayers []db.Player
		db.SqlDB.Model(&db.Player{}).Where("match_id = ?", matchIds[i]).Find(&localPlayers)

		for j := range localPlayers {
			if localPlayers[j].PlayerID == PlayerID {
				continue
			}
			if localPlayers[j].Side == sides[i] {
				allyCount[localPlayers[j].PlayerID]++
			} else {
				enermyCount[localPlayers[j].PlayerID]++
			}
		}
	}

	// map to array
	for k, v := range allyCount {
		sortedAllies = append(sortedAllies, [2]uint64{k, uint64(v)})
	}

	for k, v := range enermyCount {
		sortedEnermies = append(sortedEnermies, [2]uint64{k, uint64(v)})
	}

	sort.Slice(sortedAllies, func(i int, j int) bool { return sortedAllies[i][1] > sortedAllies[j][1] })
	sort.Slice(sortedEnermies, func(i int, j int) bool { return sortedEnermies[i][1] > sortedEnermies[j][1] })
	return
}

func getMatchIdsAndSides(PlayerID uint64) (matchIds []string, sides []int) {
	var players []db.Player
	db.SqlDB.Model(&db.Player{}).Where("player_id = ?", PlayerID).Find(&players)
	sides = make([]int, 0, len(players))

	matchIds = make([]string, 0, len(players))
	for i := range players {
		matchIds = append(matchIds, players[i].MatchID)
		sides = append(sides, players[i].Side)
	}
	return
}

func GlobalHeroAnalysis(HeroName string) (players []db.Player, err error) {
	if id, ok := db.HeroNameToID[HeroName]; ok {
		err = db.SqlDB.Model(db.Player{}).Where("hero_id = ?", id).Find(&players).Error
		return
	}
	return nil, fmt.Errorf("不存在 %s 该英雄", HeroName)
}
