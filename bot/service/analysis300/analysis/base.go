package analysis

import (
	"eebot/bot/service/analysis300/db"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
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
//			diff2   相对均分离散
//			fvRange 竞技力范围
func WinOrLoseAnalysis(PlayerID uint64) (result [][3]int, diff int, diff2 int, fvRange [2]int, fvNow int) {
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

		// 计算标准差
		_svd1 := 0
		_svd2 := 0
		for j := range localPlayers {
			if localPlayers[j].Side == sides[i] {
				_svd1 += (localPlayers[j].FV - tmp[0]) * (localPlayers[j].FV - tmp[0])
			} else {
				_svd2 += (localPlayers[j].FV - tmp[1]) * (localPlayers[j].FV - tmp[1])
			}
		}
		_svd1 = int(math.Sqrt(float64(_svd1) / 6))
		_svd2 = int(math.Sqrt(float64(_svd2) / 6))
		diff2 += _svd1 - _svd2
	}
	if len(matchIds) != 0 {
		diff /= len(matchIds)
		diff2 /= len(matchIds)
	}
	return
}

// HeroAnalysis 英雄分析(常用英雄及其数据)
//
//	return:
//		result[*][0] -- 英雄id
//		result[*][1] -- 计算场次
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
//		result[*][28] -- 总场次
func HeroAnalysis(PlayerID uint64, fv int) (result [][29]float64, total uint64) {
	var players []db.Player
	db.SqlDB.Model(&db.Player{}).Where("player_id = ? and fv >= ?", PlayerID, fv).Find(&players)
	sort.Slice(players, func(i, j int) bool {
		return players[i].CreateTime >= players[j].CreateTime
	})
	total = uint64(len(players))

	addValue := func(v [29]float64, me db.Player) [29]float64 {
		v[28]++
		// 过滤低于阈值的战绩
		if me.FV < fv {
			return v
		}
		// 只计算最近50场
		if v[1] >= MaxPlayTimes {
			return v
		}
		minute := float64(me.UsedTime) / 60
		v[1]++
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

	dataMap := make(map[int][29]float64)
	for _, me := range players {
		if v, ok := dataMap[me.HeroID]; ok {
			dataMap[me.HeroID] = addValue(v, me)
		} else {
			v := [29]float64{}
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

	sort.Slice(result, func(i, j int) bool { return result[i][28] >= result[j][28] })
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

func JJLWithTeamAnalysis(PlayerID uint64) (timeRange []string, jjl []uint64, team [][4]uint64) {
	timeToData := map[int64][6]uint64{} // 时间戳 - [单排次数, 双排次数, 三黑次数, 四黑次数, jjl, jjl对应时间戳]

	// 先获得开黑情况
	matchIds, sides := getMatchIdsAndSides(PlayerID)
	if len(matchIds) == 0 {
		return
	}

	allyInfo := map[uint64]int{}               // id-cnt
	matchToPlayers := map[string][]db.Player{} // match_id - players
	for i := range matchIds {
		// 找到这局的玩家
		var localPlayers []db.Player
		db.SqlDB.Model(&db.Player{}).Where("match_id = ?", matchIds[i]).Find(&localPlayers)
		matchToPlayers[matchIds[i]] = localPlayers
		for j := range localPlayers {
			if localPlayers[j].PlayerID == PlayerID || localPlayers[j].Side != sides[i] {
				continue
			}
			allyInfo[localPlayers[j].PlayerID]++
		}
	}

	sortedAllies := [][2]uint64{}
	for k, v := range allyInfo {
		sortedAllies = append(sortedAllies, [2]uint64{k, uint64(v)})
	}
	sort.Slice(sortedAllies, func(i int, j int) bool { return sortedAllies[i][1] > sortedAllies[j][1] })
	top10Allies := sortedAllies[:10]
	contain := func(arr [][2]uint64, id uint64) bool {
		for i := range arr {
			// 包含在频次top10中
			// 并且频次超过3把或者场次占比超过2%
			// 则认为是开黑队友
			if arr[i][0] == id && (arr[i][1] > 3 || (float64(arr[i][1])/float64(len(matchIds))) >= 0.02) {
				return true
			}
		}
		return false
	}

	for _, players := range matchToPlayers {
		cnt := 0
		for i := range players {
			if players[i].PlayerID == PlayerID {
				continue
			}
			if contain(top10Allies, players[i].PlayerID) {
				cnt++
			}
		}
		cnt = min(cnt, 3)

		tmp := time.Unix(int64(players[0].CreateTime), 0)
		timestamp := time.Date(tmp.Year(), tmp.Month(), tmp.Day(), 23, 59, 59, 0, time.Local).Unix()
		if v, ok := timeToData[timestamp]; ok {
			v[cnt] += 1
			timeToData[timestamp] = v
		} else {
			tmp := [6]uint64{0, 0, 0, 0, 0, 0}
			tmp[cnt] += 1
			timeToData[timestamp] = tmp
		}
	}

	// 获得团分情况
	players := []db.Player{}
	db.SqlDB.Model(db.Player{}).Where("player_id = ?", PlayerID).Find(&players)
	for i := range players {
		tmp := time.Unix(int64(players[i].CreateTime), 0)
		timestamp := time.Date(tmp.Year(), tmp.Month(), tmp.Day(), 23, 59, 59, 0, time.Local).Unix()
		if v, ok := timeToData[timestamp]; ok {
			if uint64(players[i].CreateTime) >= v[5] {
				v[5] = uint64(players[i].CreateTime)
				v[4] = uint64(players[i].FV)
				timeToData[timestamp] = v
			}
		} else {
			tmp := [6]uint64{0, 0, 0, 0, uint64(players[i].FV), uint64(players[i].CreateTime)}
			timeToData[timestamp] = tmp
		}
	}

	// 按照时间戳排序
	sortedData := [][7]uint64{}
	for timestamp, data := range timeToData {
		tmp := [7]uint64{}
		tmp[0] = uint64(timestamp)
		tmp[1] = data[0]
		tmp[2] = data[1]
		tmp[3] = data[2]
		tmp[4] = data[3]
		tmp[5] = data[4]
		tmp[6] = data[5]
		sortedData = append(sortedData, tmp)
	}
	sort.Slice(sortedData, func(i, j int) bool { return sortedData[i][0] < sortedData[j][0] })

	// fmt.Printf("len(sortedData): %v\n", len(sortedData))
	// str := FormatJson(sortedData[:10], true)
	// fmt.Printf("str: %v\n", str)

	// 插入第一条
	timeRange = append(timeRange, time.Unix(int64(sortedData[0][0]), 0).Format(time.DateOnly))
	jjl = append(jjl, sortedData[0][5])
	team = append(team, [4]uint64{sortedData[0][1], sortedData[0][2], sortedData[0][3], sortedData[0][4]})
	for i := 1; i < len(sortedData); i++ {
		add := 0
		for (sortedData[i][0] - sortedData[i-1][0] - uint64(add*86400)) > 86400 {
			// 仿造一条数据插入
			add++
			timeRange = append(timeRange, time.Unix(int64(sortedData[i-1][0]+uint64(add*86400)), 0).Format(time.DateOnly))
			jjl = append(jjl, sortedData[i-1][5])
			team = append(team, [4]uint64{0, 0, 0, 0})
		}
		// 插入今天的数据
		timeRange = append(timeRange, time.Unix(int64(sortedData[i][0]), 0).Format(time.DateOnly))
		jjl = append(jjl, sortedData[i][5])
		team = append(team, [4]uint64{sortedData[i][1], sortedData[i][2], sortedData[i][3], sortedData[i][4]})
	}

	// fmt.Printf("timeRange: %v\n", timeRange)
	return
}

// FormatJson 格式化Json以便更容器查看, 如果m格式错误则返回空字符串
func FormatJson(m interface{}, indent bool) string {
	var b []byte
	var err error
	if !indent {
		b, err = json.Marshal(m)
	} else {
		b, err = json.MarshalIndent(m, "", "  ")
	}
	if err != nil {
		return ""
	}
	return string(b)
}

func PKAnalysis(PlayerID uint64, HeroID int) (selfData [14]float64, otherData [14]float64) {
	UpdateHeroOfPlayerRank(HeroID, 0)
	// 获取自己数据
	scores, _, overallScore, _ := GetHeroOfPlayerRank(HeroID, PlayerID, 0)
	selfData[0] = math.Round(scores[2]*100) / 100
	selfData[1] = math.Round(scores[27]/60*100) / 100
	selfData[2] = math.Round(scores[4]*100) / 100
	selfData[3] = math.Round(scores[6]*100) / 100
	selfData[4] = math.Round(scores[8]*100) / 100
	selfData[5] = math.Round(scores[10]*100) / 100
	selfData[6] = math.Round(scores[11]*100) / 100
	selfData[7] = math.Round(scores[13]*100) / 100
	selfData[8] = math.Round(scores[15]*100) / 100
	selfData[9] = math.Round(scores[18]*100) / 100
	selfData[10] = math.Round(scores[21]*100) / 100
	selfData[11] = math.Round(scores[24]*100) / 100
	selfData[12] = math.Round(scores[26]*100) / 100
	selfData[13] = math.Round(float64(overallScore)*100) / 100

	// 获取top1数据
	top10, _, _ := GetTopRank(HeroID, 0)
	top1IDStr := top10[0].Member.(string)
	top1ID, _ := strconv.ParseUint(top1IDStr, 10, 64)
	scores, _, overallScore, _ = GetHeroOfPlayerRank(HeroID, top1ID, 0)
	otherData[0] = math.Round(scores[2]*100) / 100
	otherData[1] = math.Round(scores[27]/60*100) / 100
	otherData[2] = math.Round(scores[4]*100) / 100
	otherData[3] = math.Round(scores[6]*100) / 100
	otherData[4] = math.Round(scores[8]*100) / 100
	otherData[5] = math.Round(scores[10]*100) / 100
	otherData[6] = math.Round(scores[11]*100) / 100
	otherData[7] = math.Round(scores[13]*100) / 100
	otherData[8] = math.Round(scores[15]*100) / 100
	otherData[9] = math.Round(scores[18]*100) / 100
	otherData[10] = math.Round(scores[21]*100) / 100
	otherData[11] = math.Round(scores[24]*100) / 100
	otherData[12] = math.Round(scores[26]*100) / 100
	otherData[13] = math.Round(float64(overallScore)*100) / 100
	return
}

func Divide[T uint64 | int64 | int](a T, b T) float64 {
	return float64(a) / float64(b)
}

func ExtractByFV(start, end int, result [][3]int) (cnt [2]int) {
	for i := range result {
		avg := (result[i][0] + result[i][1]) / 2
		if start <= avg && avg < end {
			if result[i][2] == 1 {
				cnt[0]++
			} else {
				cnt[1]++
			}
		}
	}
	return
}

func ExtractByFVAdvanced(start, end int, result [][5]int) (cnt [2]int) {
	for i := range result {
		avg := (result[i][0] + result[i][1]) / 2
		if start <= avg && avg < end {
			if result[i][2] == 1 {
				cnt[0]++
			} else {
				cnt[1]++
			}
		}
	}
	return
}

func Sum[T []any](s T, get func(e any) float64) float64 {
	rs := 0.0
	for i := range s {
		rs += get(s[i])
	}
	return rs
}
