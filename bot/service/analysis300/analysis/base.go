package analysis

import (
	"eebot/bot/service/analysis300/db"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"sort"
	"time"
)

var PlayerListKey = "300analysis:player_list"

var ExpiryDate int64 = 30 * 24 * 60 * 60

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
	start := time.Now().Unix() - ExpiryDate
	db.SqlDB.Model(&db.Player{}).Where("player_id = ? and fv >= ? and create_time > ?", PlayerID, fv, start).Find(&players)
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

func GetMatchAndMyPlays(PlayerID uint64, span time.Duration) (matches []db.Match, myPlays []db.Player) {
	if span > 0 {
		db.SqlDB.Model(&db.Player{}).Where("player_id = ? and create_time > ?", PlayerID, time.Now().Add(-span).Unix()).Order("create_time asc").Find(&myPlays)
	} else {
		db.SqlDB.Model(&db.Player{}).Where("player_id = ?", PlayerID).Order("create_time asc").Find(&myPlays)
	}

	matchIds := make([]string, 0, len(myPlays))
	for i := range myPlays {
		matchIds = append(matchIds, myPlays[i].MatchID)
	}
	err := db.SqlDB.Model(db.Match{}).Preload("Players").Where("match_id in ?", matchIds).Order("create_time asc").Find(&matches).Error
	if err != nil || len(matches) == 0 {
		return
	}
	return
}

func GlobalHeroAnalysis(HeroName string) (players []db.Player, err error) {
	if id, ok := db.HeroNameToID[HeroName]; ok {
		start := time.Now().Unix() - ExpiryDate
		err = db.SqlDB.Model(db.Player{}).Where("create_time > ? and hero_id = ?", start, id).Find(&players).Error
		return
	}
	return nil, fmt.Errorf("不存在 %s 该英雄", HeroName)
}

func JJLWithTeamAnalysis(PlayerID uint64) (timeRange []string, jjl []uint64, team [][4]uint64) {
	timeToData := map[int64][6]uint64{} // 当天最晚时间戳 - [单排次数, 双排次数, 三黑次数, 四黑次数, jjl, jjl对应时间戳]
	_, myPlays, marks, _, _ := MarkTeam(PlayerID)
	for i := range myPlays {
		thatTime := time.Unix(int64(myPlays[i].CreateTime), 0)
		latest := time.Date(thatTime.Year(), thatTime.Month(), thatTime.Day(), 23, 59, 59, 0, time.Local).Unix()
		if v, ok := timeToData[latest]; ok {
			v[marks[i]]++
			if myPlays[i].FV > int(v[4]) {
				v[4] = uint64(myPlays[i].FV)
				v[5] = myPlays[i].CreateTime
			}
			timeToData[latest] = v
		} else {
			v := [6]uint64{}
			v[marks[i]]++
			if myPlays[i].FV > int(v[4]) {
				v[4] = uint64(myPlays[i].FV)
				v[5] = myPlays[i].CreateTime
			}
			timeToData[latest] = v
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
	return
}

// JJLCompositionAnalysis jjl成分分析
//
//	gangUp:
//		gangUp[*][0] - (*+1)黑上分
//		gangUp[*][1] - (*+1)黑掉分
//		gangUp[*][2] - (*+1)黑场次
//	allies: 队友情况
//	enermies: 敌人分段情况
//	scope: 比赛分段情况
//	hero: 英雄情况
//	total: 总计
func JJLCompositionAnalysis(PlayerID uint64, span time.Duration) (gangUp [4][3]float64, allies map[uint64][3]float64, enermies [][3]float64, scope [][3]float64, hero map[int][3]float64, total int) {
	matches, myPlays, marks, _, top10Allies := MarkTeam(PlayerID)
	total = len(matches)
	hero = map[int][3]float64{}
	allies = map[uint64][3]float64{}
	enermies = make([][3]float64, len(DefaultJJLCategoryKeys))
	scope = make([][3]float64, len(DefaultJJLCategoryKeys))
	for i := range top10Allies {
		allies[top10Allies[i][0]] = [3]float64{}
	}
	ddl := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Add(-span).Unix()
	for i := range matches {
		if int64(matches[i].CreateTime) < ddl{
			continue
		}
		var offset float64
		if i == 0 {
			offset = 0
		} else {
			offset = float64(myPlays[i].FV - myPlays[i-1].FV)
		}
		tr := func(fv float64) int {
			if fv >= 0 {
				return 0
			}
			return 1
		}
		// 开黑情况
		gangUp[marks[i]][tr(offset)] += offset
		gangUp[marks[i]][2]++
		avg := 0
		// 敌人分段情况
		for _, player := range matches[i].Players {
			avg += player.FV
			if player.Side == myPlays[i].Side {
				// 队友情况
				if v, ok := allies[player.PlayerID]; ok {
					v[tr(offset)] += offset / float64(marks[i])
					v[2]++
					allies[player.PlayerID] = v
				}
				continue
			}
			enermies[DefaultJJLCategoryKeys.Index(float64(player.FV))][2]++
			enermies[DefaultJJLCategoryKeys.Index(float64(player.FV))][tr(offset)] += offset / 7
		}
		// 比赛分段情况
		scope[DefaultJJLCategoryKeys.Index(float64(avg)/14)][2]++
		scope[DefaultJJLCategoryKeys.Index(float64(avg)/14)][tr(offset)] += offset
		// 英雄情况
		if v, ok := hero[myPlays[i].HeroID]; ok {
			v[2]++
			v[tr(offset)] += offset
			hero[myPlays[i].HeroID] = v
		} else {
			tmp := [3]float64{0, 0, 1}
			tmp[tr(offset)] += offset
			hero[myPlays[i].HeroID] = tmp
		}
	}
	return
}

// 安定段位分析
func StableJJLLAnalysis(PlayerID uint64) (stable int) {
	matches, myPlays, _, _, _ := MarkTeam(PlayerID)
	if len(matches) == 0 {
		return 0
	}
	// fv-offset 分数得分情况
	m := map[int]int{}
	for i := range matches {
		if i == 0 {
			continue
		}
		offset := myPlays[i].FV - myPlays[i-1].FV
		for _, localPlayer := range matches[i].Players {
			if localPlayer.Side != myPlays[i].Side {
				m[localPlayer.FV] += offset
			}
		}
	}
	sortedOffset := [][2]int{}
	for fv, offset := range m {
		sortedOffset = append(sortedOffset, [2]int{fv, offset})
	}
	slices.SortFunc[[][2]int](sortedOffset, func(a, b [2]int) int { return a[0] - b[0] })
	// 初始化区间
	stables := []int{}
	for span := 30; span <= 100; span += 10 {
		cluster := make([][2]int, 2500/span)
		// 双指针
		q := 0
		for p := 0; p < len(cluster); p++ {
			scope := [2]int{(p+1)*span - span/2, (p+1)*span + span/2}
			for q < len(sortedOffset) {
				if scope[0] <= sortedOffset[q][0] && sortedOffset[q][0] < scope[1] {
					cluster[p][0] += sortedOffset[q][1]
					cluster[p][1] += 1
					q++
					continue
				} else {
					break
				}
			}
		}
		for i := 0; i < len(cluster)-1; i++ {
			if cluster[i][1] != 0 && (i+1)*span-span/2 > 1500 {
				if cluster[i][0] >= 0 && cluster[i+1][0] <= 0 {
					stables = append(stables, (i+1)*span+span/2)
					break
				}
			}
		}
	}
	slices.Sort[[]int](stables)
	sum := 0
	for _, v := range stables[1 : len(stables)-1] {
		sum += v
	}
	sum /= len(stables[1 : len(stables)-1])
	return sum
}

// 标记开黑情况
//
//	matches: 按照时间戳排序的比赛记录
//	myPlays: 比赛对应的自己游玩记录
//	marks: 比赛对应的开黑，0单排，3四黑
//	marksDetail: 比赛对应开黑的队友
//	allies:
//		allies[*][0]: id
//		allies[*][1]: 出现次数
//		allies[*][2]: 胜利次数
func MarkTeam(PlayerID uint64) (matches []db.Match, myPlays []db.Player, marks []int, marksDetail [][]string, allies [][4]uint64) {
	// 获取全部数据来计算开黑情况
	matches, myPlays = GetMatchAndMyPlays(PlayerID, 0)
	if len(matches) == 0 {
		return
	}
	allyInfo := map[uint64][2]int{}
	enermyInfo := map[uint64][2]int{}
	for i, match := range matches {
		// 先找到自己
		me := myPlays[i]
		// 初始化开黑情况
		for _, player := range match.Players {
			if player.PlayerID == PlayerID {
				continue
			}
			if player.Side == me.Side {
				if v, ok := allyInfo[player.PlayerID]; ok {
					v[0]++
					if me.Result == 1 || me.Result == 3 {
						v[1]++
					}
					allyInfo[player.PlayerID] = v
				} else {
					v := [2]int{1, 0}
					if me.Result == 1 || me.Result == 3 {
						v[1]++
					}
					allyInfo[player.PlayerID] = v
				}
			} else {
				if v, ok := enermyInfo[player.PlayerID]; ok {
					v[0]++
					if me.Result == 1 || me.Result == 3 {
						v[1]++
					}
					enermyInfo[player.PlayerID] = v
				} else {
					v := [2]int{1, 0}
					if me.Result == 1 || me.Result == 3 {
						v[1]++
					}
					enermyInfo[player.PlayerID] = v
				}
			}
		}
	}
	sortedAllies := [][4]uint64{}
	for k, v := range allyInfo {
		if v[0] > 3 {
			cnt := max(v[0]-enermyInfo[k][0], 0)
			sortedAllies = append(sortedAllies, [4]uint64{k, uint64(v[0]), uint64(v[1]), uint64(cnt)})
		}
	}
	sort.Slice(sortedAllies, func(i int, j int) bool { return sortedAllies[i][3] > sortedAllies[j][3] })
	top10Allies := [][4]uint64{}
	for i := range sortedAllies {
		if i >= 10 || sortedAllies[i][3] <= 3 {
			top10Allies = sortedAllies[:i]
			break
		}
	}
	isGangUp := func(arr [][4]uint64, id uint64) bool {
		for i := range arr {
			if arr[i][0] == id && arr[i][3] > 3 {
				return true
			}
		}
		return false
	}
	for i := range matches {
		cnt := 0
		tmp := []string{}
		for _, player := range matches[i].Players {
			if player.Side == myPlays[i].Side && isGangUp(top10Allies, player.PlayerID) {
				cnt++
				tmp = append(tmp, player.Name)
			}
		}
		cnt = min(cnt, 3)
		marks = append(marks, cnt)
		marksDetail = append(marksDetail, tmp)
	}
	return matches, myPlays, marks, marksDetail, top10Allies
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

func PKAnalysis(PlayerID uint64, HeroID int) (selfData [14]float64, otherData [14]float64, err error) {
	// 获取自己数据
	players, _ := GetRankFromPlayers(HeroID, 0, []uint64{PlayerID})
	me, ok := players[PlayerID]
	if !ok {
		return selfData, otherData, errors.New("玩家无此数据")
	}
	selfData[0] = math.Ceil(me.WinRate*100*100) / 100.0
	selfData[1] = math.Round(me.AvgUsedTime/60*100) / 100
	selfData[2] = math.Round(me.AvgHitPerMinite*100) / 100
	selfData[3] = math.Round(me.AvgKillPerMinite*100) / 100
	selfData[4] = math.Round(me.AvgDeathPerMinite*100) / 100
	selfData[5] = math.Round(me.AvgAssistPerMinite*100) / 100
	selfData[6] = math.Round(me.AvgTower*100) / 100
	selfData[7] = math.Round(me.AvgPutEye*100) / 100
	selfData[8] = math.Round(me.AvgDestryEye*100) / 100
	selfData[9] = math.Round(me.AvgMoneyPerMinite*100) / 100
	selfData[10] = math.Round(me.AvgMakeDamagePerMinite*100) / 100
	selfData[11] = math.Round(me.AvgTakeDamagePerMinite*100) / 100
	selfData[12] = math.Round(me.AvgMoneyConversionRate*100) / 100
	selfData[13] = math.Round(me.Score*100) / 100

	// 获取top1数据
	top, _ := GetRankFromTop(HeroID, 0, 1)
	top1 := top[0]
	otherData[0] = math.Ceil(top1.WinRate*100*100) / 100.0
	otherData[1] = math.Round(top1.AvgUsedTime/60*100) / 100
	otherData[2] = math.Round(top1.AvgHitPerMinite*100) / 100
	otherData[3] = math.Round(top1.AvgKillPerMinite*100) / 100
	otherData[4] = math.Round(top1.AvgDeathPerMinite*100) / 100
	otherData[5] = math.Round(top1.AvgAssistPerMinite*100) / 100
	otherData[6] = math.Round(top1.AvgTower*100) / 100
	otherData[7] = math.Round(top1.AvgPutEye*100) / 100
	otherData[8] = math.Round(top1.AvgDestryEye*100) / 100
	otherData[9] = math.Round(top1.AvgMoneyPerMinite*100) / 100
	otherData[10] = math.Round(top1.AvgMakeDamagePerMinite*100) / 100
	otherData[11] = math.Round(top1.AvgTakeDamagePerMinite*100) / 100
	otherData[12] = math.Round(top1.AvgMoneyConversionRate*100) / 100
	otherData[13] = math.Round(top1.Score*100) / 100
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
