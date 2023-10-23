package analysis

import (
	"eebot/bot/service/analysis300/db"
	"math"
	"sort"
	"time"
)

// ShuffleAnalysis 洗牌分析
//
//	两种情况：
//	1、自己秒的情况
//	2、开黑小号秒的情况
//	return: 平均游戏开启时间间隔(仅仅计算2个小时内的时间间隔)
func ShuffleAnalysis(PlayerID uint64) (avgInterval int, than10min int, validCnt int) {
	matchIds, _ := getMatchIdsAndSides(PlayerID)

	times := make([][2]uint64, 0, len(matchIds)) // [开始时间, 结束时间]

	for i := range matchIds {
		var match db.Match
		db.SqlDB.Model(&db.Match{}).Where("match_id = ?", matchIds[i]).First(&match)

		times = append(times, [2]uint64{match.CreateTime - match.UsedTime, match.CreateTime})
	}

	sort.Slice(times, func(i, j int) bool { return times[i][0] < times[j][0] })

	intervalSum := 0
	for i := 0; i < len(times)-1; i++ {
		interval := times[i+1][0] - times[i][1]
		if interval <= 2*60*60 {
			if interval > 60*10 {
				than10min += 1
			}

			validCnt++
			intervalSum += int(interval)
		}
	}
	avgInterval = intervalSum / validCnt
	return
}

// TeamAnalysisAdvanced 进阶开黑分析(开黑胜率，开几黑)
//
//	return: sortedAllies [][3]uint64
//			- sortedAllies[*][0] -- PlayerID
//			- sortedAllies[*][1] -- 胜场
//			- sortedAllies[*][2] -- 场次
//
//			teams [4][2]uint64
//			- teams[*][0] -- 开(*+1)黑胜场
//			- teams[*][1] -- 开(*+1)黑场次
func TeamAnalysisAdvanced(PlayerID uint64) (sortedAllies [][3]uint64, sortedEnermies [][3]uint64, teams [4][2]uint64, teamAllies map[int]map[string]struct{}, total int) {
	matchIds, sides := getMatchIdsAndSides(PlayerID)

	total = len(matchIds)
	allyInfo := map[uint64][2]int{} // playerID -- [win, total]
	enermyInfo := map[uint64][2]int{}
	matchToPlayers := map[string][]db.Player{}
	for i := range matchIds {
		// 找到这局的玩家
		var localPlayers []db.Player
		db.SqlDB.Model(&db.Player{}).Where("match_id = ?", matchIds[i]).Find(&localPlayers)
		matchToPlayers[matchIds[i]] = localPlayers

		for j := range localPlayers {
			if localPlayers[j].PlayerID == PlayerID {
				continue
			}
			// 队友还是敌人
			if localPlayers[j].Side == sides[i] {
				// 胜还是负
				if v, ok := allyInfo[localPlayers[j].PlayerID]; ok {
					v[1]++
					if localPlayers[j].Result == 1 || localPlayers[j].Result == 3 {
						v[0]++
					}
					allyInfo[localPlayers[j].PlayerID] = v
				} else {
					if localPlayers[j].Result == 1 || localPlayers[j].Result == 3 {
						allyInfo[localPlayers[j].PlayerID] = [2]int{1, 1}
					} else {
						allyInfo[localPlayers[j].PlayerID] = [2]int{0, 1}
					}
				}
			} else {
				if v, ok := enermyInfo[localPlayers[j].PlayerID]; ok {
					v[1]++
					if localPlayers[j].Result == 2 || localPlayers[j].Result == 4 {
						v[0]++
					}
					enermyInfo[localPlayers[j].PlayerID] = v
				} else {
					if localPlayers[j].Result == 2 || localPlayers[j].Result == 4 {
						enermyInfo[localPlayers[j].PlayerID] = [2]int{1, 1}
					} else {
						enermyInfo[localPlayers[j].PlayerID] = [2]int{0, 1}
					}
				}
			}
		}
	}

	// map to array
	for k, v := range allyInfo {
		sortedAllies = append(sortedAllies, [3]uint64{k, uint64(v[0]), uint64(v[1])})
	}
	for k, v := range enermyInfo {
		sortedEnermies = append(sortedEnermies, [3]uint64{k, uint64(v[0]), uint64(v[1])})
	}

	sort.Slice(sortedAllies, func(i int, j int) bool { return sortedAllies[i][2] > sortedAllies[j][2] })
	sort.Slice(sortedEnermies, func(i int, j int) bool { return sortedEnermies[i][2] > sortedEnermies[j][2] })

	// 进阶分析
	top10Allies := sortedAllies[:10]
	contain := func(slice [][3]uint64, id uint64) bool {
		for i := range slice {
			if slice[i][0] == id {
				return true
			}
		}
		return false
	}

	teamAllies = make(map[int]map[string]struct{})
	for _, localPlayers := range matchToPlayers {

		cnt := 0
		var me db.Player
		names := make([]string, 0, 4)
		// 是否包含在高频队友内
		for j := range localPlayers {
			if localPlayers[j].PlayerID == PlayerID {
				me = localPlayers[j]
				continue
			}
			if contain(top10Allies, localPlayers[j].PlayerID) {
				names = append(names, localPlayers[j].Name)
				cnt++
			}
		}
		// 对应开黑+1
		if cnt >= 4 {
			cnt = 3
		}
		if me.Result == 1 || me.Result == 3 {
			teams[cnt][0]++
		}
		teams[cnt][1]++
		for _, name := range names {
			if teamAllies[cnt] == nil {
				teamAllies[cnt] = make(map[string]struct{})
			}
			teamAllies[cnt][name] = struct{}{}
		}
	}

	return
}

// WinOrLoseAnalysisAdvanced 进阶输赢分析(是否匹配当前分数段)
//
//	result[0]: 己方均分
//	result[1]: 敌方均分
//	result[2]: 输赢 1-赢，2-输
//	result[3]: 自己团分
//	result[4]: 局类型 0-杀鸡，1-本地，2-壮丁
//	diff: 均分差
//	svd: 离散度差
//	fvRange: 团分范围
//	fvNow: 目前团分
//	timeRange: 时间范围
func WinOrLoseAnalysisAdvanced(PlayerID uint64) (result [][5]int, diff int, svd int, fvRange [2]int, fvNow int, timeRange [2]uint64) {
	matchIds, sides := getMatchIdsAndSides(PlayerID)

	fvRange[0] = 2500
	fvRange[1] = 0
	timeRange[0] = uint64(time.Now().Unix())
	var maxTimestamps uint64
	for i := range matchIds {
		var selfFV int
		var localPlayers []db.Player
		db.SqlDB.Model(&db.Player{}).Where("match_id = ?", matchIds[i]).Find(&localPlayers)

		var tmp [5]int
		fvSum1 := 0 // 己方团分
		fvSum2 := 0 // 对面团分
		for j := range localPlayers {
			if localPlayers[j].PlayerID == PlayerID {
				tmp[2] = localPlayers[j].Result
				selfFV = localPlayers[j].FV
				fvRange[0] = min(localPlayers[j].FV, fvRange[0])
				fvRange[1] = max(localPlayers[j].FV, fvRange[1])
				timeRange[0] = min(localPlayers[j].CreateTime, timeRange[0])
				timeRange[1] = max(localPlayers[j].CreateTime, timeRange[1])
				if localPlayers[j].CreateTime >= maxTimestamps {
					maxTimestamps = localPlayers[j].CreateTime
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
		tmp[3] = selfFV

		diff += selfFV - (tmp[0]+tmp[1])/2

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
		svd += _svd1 - _svd2

		// 判断局类型
		tmp[4] = 1
		avg := (tmp[0] + tmp[1]) / 2
		// 只要自己团分比均分低100，直接判断为壮丁局
		if selfFV-avg < -100 {
			tmp[4] = 2
		}
		// 自己团分比均分高100
		// 找到对面有没有比自己高或者和自己团分相似的人，若没有，则是杀鸡
		if selfFV-avg > 100 {
			flag := false
			for j := range localPlayers {
				if localPlayers[j].PlayerID == PlayerID {
					continue
				}
				if localPlayers[j].Side == sides[i] {
					continue
				} else {
					if localPlayers[j].FV >= selfFV || IsSimilarFV(selfFV, localPlayers[j].FV) {
						flag = true
						break
					}
				}
			}
			if !flag {
				tmp[4] = 0
			}
		}

		result = append(result, tmp)
	}
	if len(matchIds) != 0 {
		diff /= len(matchIds)
		svd /= len(matchIds)
	}
	return
}

// IsSimilarFV 比较两个团分是否相似
func IsSimilarFV(fv1, fv2 int) bool {
	if fv1 >= 2100 && fv2 >= 2100 {
		return true
	} else if fv1-fv2 <= 50 && fv1-fv2 >= -50 {
		return true
	}
	return false
}
