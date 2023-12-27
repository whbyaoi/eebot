package analysis

import (
	"math"
	"slices"
	"sort"
)

// ShuffleAnalysis 洗牌分析
//
//	两种情况：
//	1、自己秒的情况
//	2、开黑小号秒的情况
//	return: 平均游戏开启时间间隔(仅仅计算2个小时内的时间间隔)
func ShuffleAnalysis(PlayerID uint64) (avgInterval int, than10min int, validCnt int) {
	matches, _ := GetMatchAndMyPlays(PlayerID, 0)

	times := make([][2]uint64, 0, len(matches)) // [开始时间, 结束时间]

	for i := range matches {
		times = append(times, [2]uint64{matches[i].CreateTime - matches[i].UsedTime, matches[i].CreateTime})
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

// WinOrLoseAnalysisAdvanced 进阶输赢分析(是否匹配当前分数段)
//
//	result[0]: 己方均分
//	result[1]: 敌方均分
//	result[2]: 输赢 1-赢，2-输
//	result[3]: 自己竞技力
//	result[4]: 局类型 0-杀鸡，1-本地，2-壮丁
//	result[5]: 均分
//	diff: 玩家分相对场均分差
//	svd: 离散度差
//	fixDiff: 修正双方均分差
//	fvNow: 目前竞技力
//	timeRange: 时间范围
func WinOrLoseAnalysisAdvanced(PlayerID uint64) (result [][6]float64, diff, fixDiff, fixCount, svd, fvNow int, timeRange [2]uint64) {
	matches, myPlays := GetMatchAndMyPlays(PlayerID, 0)
	if len(matches) == 0 {
		return
	}
	timeRange[0] = matches[0].CreateTime
	timeRange[1] = matches[len(matches)-1].CreateTime
	fvNow = myPlays[len(myPlays)-1].FV
	for i := range matches {
		var tmp [6]float64
		fvSum1 := 0 // 己方竞技力
		fvSum2 := 0 // 对面竞技力
		fvArr1 := []int{}
		fvArr2 := []int{}
		for j := range matches[i].Players {
			tmp[5] += float64(matches[i].Players[j].FV)
			if matches[i].Players[j].PlayerID == PlayerID {
				tmp[2] = float64(matches[i].Players[j].Result)
			}
			if matches[i].Players[j].Side == myPlays[i].Side {
				fvSum1 += matches[i].Players[j].FV
				fvArr1 = append(fvArr1, matches[i].Players[j].FV)
			} else {
				fvSum2 += matches[i].Players[j].FV
				fvArr2 = append(fvArr2, matches[i].Players[j].FV)
			}
		}
		tmp[0] = float64(fvSum1 / 7)
		tmp[1] = float64(fvSum2 / 7)
		tmp[3] = float64(myPlays[i].FV)
		tmp[5] /= 14

		diff += myPlays[i].FV - int((tmp[0]+tmp[1])/2)

		// 计算标准差
		_svd1 := 0
		_svd2 := 0
		for j := range matches[i].Players {
			if matches[i].Players[j].Side == myPlays[i].Side {
				_svd1 += (matches[i].Players[j].FV - int(tmp[0])) * (matches[i].Players[j].FV - int(tmp[0]))
			} else {
				_svd2 += (matches[i].Players[j].FV - int(tmp[1])) * (matches[i].Players[j].FV - int(tmp[1]))
			}
		}
		_svd1 = int(math.Sqrt(float64(_svd1) / 6))
		_svd2 = int(math.Sqrt(float64(_svd2) / 6))
		svd += _svd1 - _svd2

		// 判断局类型
		tmp[4] = 1
		avg := (tmp[0] + tmp[1]) / 2
		// 只要自己竞技力比均分低100，直接判断为壮丁局
		if tmp[3]-avg < -100 {
			tmp[4] = 2
		}
		// 自己竞技力比均分高100
		// 找到对面有没有比自己高或者和自己竞技力相似的人，若没有，则是杀鸡
		if tmp[3]-avg > 100 {
			flag := false
			for j := range matches[i].Players {
				if matches[i].Players[j].PlayerID == PlayerID {
					continue
				}
				if matches[i].Players[j].Side == myPlays[i].Side {
					continue
				} else {
					if float64(matches[i].Players[j].FV) >= tmp[3] || IsSimilarFV(int(tmp[3]), matches[i].Players[j].FV) {
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

		sort.Slice(fvArr1, func(i, j int) bool { return fvArr1[i] < fvArr1[j] })
		sort.Slice(fvArr2, func(i, j int) bool { return fvArr2[i] < fvArr2[j] })
		index, _ := slices.BinarySearch[[]int](fvArr1, myPlays[i].FV)
		fvSum1 -= fvArr1[index]
		fvSum2 -= fvArr2[index]
		if fvSum1 > fvSum2 {
			fixCount++
		}
		fixDiff += (fvSum1 - fvSum2) / 6
	}
	if len(matches) != 0 {
		diff /= len(matches)
		fixDiff /= len(matches)
		svd /= len(matches)
	}
	return
}

// IsSimilarFV 比较两个竞技力是否相似
func IsSimilarFV(fv1, fv2 int) bool {
	if fv1 >= 2100 && fv2 >= 2100 {
		return true
	} else if fv1-fv2 <= 50 && fv1-fv2 >= -50 {
		return true
	}
	return false
}
