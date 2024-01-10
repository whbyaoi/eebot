package analysis300

import (
	"context"
	"eebot/bot/service/analysis300/analysis"
	"eebot/bot/service/analysis300/collect"
	"eebot/bot/service/analysis300/db"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func ExportRelatedAnalysis(name string, hero string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	plays := []db.Player{}
	if heroID, ok := db.HeroNameToID[hero]; ok {
		err = db.SqlDB.Model(&db.Player{}).Where("player_id = ? and hero_id = ?", PlayerID, heroID).Order("create_time desc").Find(&plays).Error
		msg += fmt.Sprintf("玩家：%s，最近十场JJC %s 战绩\n", name, hero)
	} else {
		if hero == "" {
			err = db.SqlDB.Model(&db.Player{}).Where("player_id = ?", PlayerID).Order("create_time desc").Find(&plays).Error
			msg += fmt.Sprintf("玩家：%s，最近十场JJC战绩\n", name)
		} else {
			return "", fmt.Errorf("不存在 %s 该英雄", hero)
		}
	}
	if err != nil {
		return
	}
	plays = plays[:10]
	matches := []db.Match{}
	for i := range plays {
		tmp := db.Match{}
		db.SqlDB.Model(&db.Match{}).Preload("Players").Where("match_id = ?", plays[i].MatchID).Take(&tmp)
		matches = append(matches, tmp)
	}
	allMathes, _, marks, _, _ := analysis.MarkTeam(PlayerID)
	matchIdToIndex := map[string]int{}
	for i := range allMathes {
		matchIdToIndex[allMathes[i].MatchID] = i
	}
	marksTran := []string{"单排", "双排", "三黑", "四黑"}
	msg += "时间 输赢 英雄 k/d/a/补 我方/敌方均分 [开黑] [评价]\n"
	for i := range matches {
		sum1 := 0
		sum2 := 0
		for j := range matches[i].Players {
			if matches[i].Players[j].Side == plays[i].Side {
				sum1 += matches[i].Players[j].FV
			} else {
				sum2 += matches[i].Players[j].FV
			}
		}
		sum1 /= 7
		sum2 /= 7
		result := "输"
		if plays[i].Result == 1 || plays[i].Result == 3 {
			result = "赢"
		}
		msg += fmt.Sprintf("%s %s %s %d/%d/%d/%d %d/%d [%s] [%s]\n",
			time.Unix(int64(plays[i].CreateTime), 0).Format(time.DateOnly)[2:],
			result,
			db.HeroIDToName[plays[i].HeroID],
			plays[i].KillPlayer,
			plays[i].Death,
			plays[i].Assist,
			plays[i].KillUnit,
			sum1,
			sum2,
			marksTran[marks[matchIdToIndex[matches[i].MatchID]]],
			analysis.GetAppraise(plays[i]))
	}
	return
}

func ExportTeamAnalysis(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	matches, myPlays, marks, marksDetail, allies := analysis.MarkTeam(PlayerID)
	gangUp, alliesDetail, _, _, _, _ := analysis.JJLCompositionAnalysis(PlayerID, 0)
	msg += fmt.Sprintf("昵称：%s，记录场次：%d\n", name, len(matches))
	msg += fmt.Sprintln("队友情况：")
	for i := range allies {
		if i >= 7 {
			break
		}
		if allies[i][1] > 0 {
			msg += fmt.Sprintf("%d、%s，%d局(%.1f%%)，胜率：%.1f%%，净上分：%d\n",
				i+1,
				collect.SearchName(allies[i][0]),
				allies[i][1],
				analysis.Divide(allies[i][2], uint64(len(myPlays)))*100,
				analysis.Divide(allies[i][2], allies[i][1])*100,
				int64(alliesDetail[allies[i][0]][0]+alliesDetail[allies[i][0]][1]),
			)
		}
	}
	m := [4]map[string]struct{}{{}, {}, {}, {}}
	gangUpDetail := [4][2]int{}
	gangUpAllies := [4][]string{}
	for i := range myPlays {
		for j := range marksDetail[i] {
			m[marks[i]][marksDetail[i][j]] = struct{}{}
		}
		gangUpDetail[marks[i]][0]++
		if myPlays[i].Result == 1 || myPlays[i].Result == 3 {
			gangUpDetail[marks[i]][1]++
		}
	}
	for i := range m {
		for name := range m[i] {
			gangUpAllies[i] = append(gangUpAllies[i], name)
		}
	}
	msg += fmt.Sprintln("开黑情况：")
	for i := range gangUpDetail {
		if gangUpDetail[i][0] > 0 {
			msg += fmt.Sprintf("%s%d局(%.1f%%)，胜率%.1f%%，净上分：%d，队友：%s\n",
				analysis.GangUpCategoryKeys[i],
				gangUpDetail[i][0],
				analysis.Divide[int](gangUpDetail[i][0], len(myPlays))*100,
				analysis.Divide[int](gangUpDetail[i][1], gangUpDetail[i][0])*100,
				int64(gangUp[i][0]+gangUp[i][1]),
				strings.Join(gangUpAllies[i], "，"),
			)
		}
	}
	return
}

func ExportWinOrLoseAnalysisAdvanced(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}

	rs, diff, fixDiff, fixCount, svd, fvNow, timeRange := analysis.WinOrLoseAnalysisAdvanced(PlayerID)
	_, _, _, scope, _, _ := analysis.JJLCompositionAnalysis(PlayerID, 0)
	if len(rs) == 0 {
		return "", errors.New("查询不到任何战绩")
	}
	win := 0
	lose := 0
	cnt1 := 0
	cnt2 := 0
	diff2 := 0
	fvRange := [2]int{2500, 0}
	for i := range rs {
		fvRange[0] = min(int(rs[i][3]), fvRange[0])
		fvRange[1] = max(int(rs[i][3]), fvRange[1])
		if rs[i][2] == 1 {
			win++
			if rs[i][0] > rs[i][1] {
				cnt1++
			}
		} else {
			lose++
			if rs[i][0] < rs[i][1] {
				cnt2++
			}
		}
		diff2 += int(rs[i][0] - rs[i][1])
	}
	diff2 /= len(rs)
	msg += fmt.Sprintf("昵称：%s，记录场次：%d，竞技力跨度：%d - %d，时间跨度：%s - %s\n", name, len(rs), fvRange[0], fvRange[1], time.Unix(int64(timeRange[0]), 0).Format("20060102"), time.Unix(int64(timeRange[1]), 0).Format("20060102"))
	msg += fmt.Sprintf("当前竞技力：%d，胜率：%.1f%%\n", fvNow, float32(win)/float32(win+lose)*100)
	msg += fmt.Sprintf("玩家分相对场均分水平：%d\n", diff)
	msg += fmt.Sprintf("总记录 %d 局中有 %d 局(%.1f%%) 己方均分高于对面\n",
		len(rs),
		cnt1+lose-cnt2,
		float32(cnt1+lose-cnt2)/float32(len(rs))*100)
	msg += fmt.Sprintf("除开自身与对位敌人有 %d 局(%.1f%%) 己方六人竞技力高于对面六人\n", fixCount, float64(fixCount)/float64(len(rs))*100)
	msg += fmt.Sprintf("己方均分相对敌方均分水平：%d\n", diff2)
	msg += fmt.Sprintf("己方其他六人均分相对敌方六人均分水平：%d\n", fixDiff)
	msg += fmt.Sprintf("己方竞技力离散度相对敌方竞技力离散度水平：%d\n", svd)

	stages := make([][2]int, len(analysis.DefaultJJLCategoryKeys))
	for i := range rs {
		stages[analysis.DefaultJJLCategoryKeys.Index(rs[i][5])][0]++
		if rs[i][2] == 1 {
			stages[analysis.DefaultJJLCategoryKeys.Index(rs[i][5])][1]++
		}
	}
	suffix := "(占比 / 胜率 / 场次 / 净上分，下同)\n"
	msg += suffix
	for i := range stages {
		if stages[i][0] > 0 {
			msg += fmt.Sprintf("%d-%d(%.1f%%)：%.1f%% / %d / %d\n",
				int(analysis.DefaultJJLCategoryKeys.GetRange(i)[0]),
				int(analysis.DefaultJJLCategoryKeys.GetRange(i)[1]),
				analysis.Divide[int](stages[i][0], len(rs))*100,
				analysis.Divide[int](stages[i][1], stages[i][0])*100,
				stages[i][0],
				int64(scope[i][0]+scope[i][1]),
			)
		}
	}

	var a1, a2, a3 uint64 = 0, 0, 0
	var w1, w2, w3 uint64 = 0, 0, 0
	for i := range rs {
		if rs[i][4] == 0 {
			a1++
			if rs[i][2] == 1 {
				w1++
			}
		} else if rs[i][4] == 1 {
			a2++
			if rs[i][2] == 1 {
				w2++
			}
		} else if rs[i][4] == 2 {
			a3++
			if rs[i][2] == 1 {
				w3++
			}
		}
	}
	msg += fmt.Sprintf("进入杀鸡局场次(%.1f%%)：%.1f%% / %d\n", analysis.Divide(a1, uint64(len(rs)))*100, analysis.Divide(w1, a1)*100, a1)
	msg += fmt.Sprintf("进入本地局场次(%.1f%%)：%.1f%% / %d\n", analysis.Divide(a2, uint64(len(rs)))*100, analysis.Divide(w2, a2)*100, a2)
	msg += fmt.Sprintf("进入壮丁局场次(%.1f%%)：%.1f%% / %d\n", analysis.Divide(a3, uint64(len(rs)))*100, analysis.Divide(w3, a3)*100, a3)

	return
}

func ExportShuffleAnalysis(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	avg, than10min, total := analysis.ShuffleAnalysis(PlayerID)
	// rank, total2 := analysis.GetMatchInterval(PlayerID)

	msg += fmt.Sprintf("洗牌分析，昵称：%s\n", name)
	msg += fmt.Sprintf("有效间隔数：%d，平均间隔：%d秒，超过十分钟的间隔数：%d (占比%.1f%%)", total, avg, than10min, analysis.Divide(uint64(than10min), uint64(total))*100)
	return
}

func ExportShuffleAnalysisAdvanced(name string) (msg string, err error) {
	playerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	avg, than10min, total := analysis.ShuffleAnalysis(playerID)

	msg += fmt.Sprintf("洗牌分析，昵称：%s\n", name)
	msg += fmt.Sprintf("本人有效间隔数：%d，平均间隔：%d秒，超过十分钟的间隔数：%d (占比%.1f%%)\n", total, avg, than10min, analysis.Divide(uint64(than10min), uint64(total))*100)

	matches, _, _, _, allies := analysis.MarkTeam(playerID)
	for i, ally := range allies {
		if i >= 7 {
			break
		}
		if analysis.Divide(ally[1], uint64(len(matches)))*100 < 2.0 || ally[1] <= 2 {
			continue
		}
		collect.CrawlPlayerByID(ally[0])
		avg, than10min, total := analysis.ShuffleAnalysis(ally[0])
		msg += fmt.Sprintf("开黑队友 %s 有效间隔数：%d，平均间隔：%d秒，超过十分钟的间隔数：%d (%.1f%%)\n", collect.SearchName(ally[0]), total, avg, than10min, analysis.Divide(uint64(than10min), uint64(total))*100)
	}
	return
}

// ExportAssignHeroAnalysisAdvancedV2 分析fv竞技力上的玩家的英雄数据
func ExportAssignHeroAnalysisAdvancedV2(name string, hero string, fv int) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	if _, ok := db.HeroNameToID[hero]; !ok {
		return "未知英雄 " + hero, nil
	}
	if _, ok := analysis.HeroFactor[db.HeroIDToName[db.HeroNameToID[hero]]]; !ok {
		return "", fmt.Errorf("英雄 %s 的计算权重不存在", db.HeroIDToName[db.HeroNameToID[hero]])
	}
	heroDataSlice, total := analysis.GetRankFromPlayers(db.HeroNameToID[hero], fv, []uint64{PlayerID})
	if _, ok := heroDataSlice[PlayerID]; !ok {
		return fmt.Sprintf("%s 最近30天无 %s 战绩", name, hero), nil
	}
	heroData := heroDataSlice[PlayerID]
	_, _, _, _, jjl, _ := analysis.JJLCompositionAnalysis(PlayerID, 24*30*time.Hour)
	_, plays, marks, _, _ := analysis.MarkTeam(PlayerID)
	marksMap := map[string]int{}
	for i := range plays {
		marksMap[plays[i].MatchID] = marks[i]
	}
	cnt := 0.0
	for _, matchID := range heroData.MatchIDs {
		if v, ok := marksMap[matchID]; ok && v == 0 {
			cnt++
		}
	}
	msg += fmt.Sprintf("昵称：%s(只会计算近30天战绩)\n", name)
	msg += fmt.Sprintf("英雄：%s\n", hero)
	msg += fmt.Sprintf("有 %d 名玩家记录场次超过了 %d 次，竞技力下限：%d\n", total, int(analysis.ValidTimes), fv)
	msg += fmt.Sprintf("实际场次：%d，参与计算场次：%d\n", uint64(heroData.ActualTotal), uint64(heroData.Total))
	msg += fmt.Sprintf("单排率：%1.f%%\n", cnt/heroData.Total*100)
	msg += fmt.Sprintf("净上分：%d，上分：%d，掉分：%d\n", int(jjl[db.HeroNameToID[hero]][0]+jjl[db.HeroNameToID[hero]][1]),
		int(jjl[db.HeroNameToID[hero]][0]), int(jjl[db.HeroNameToID[hero]][1]))
	msg += fmt.Sprintf("胜率：%.1f%% (超越%.1f%%的玩家，下同)\n", heroData.WinRate*100, heroData.Rank.WinRate)
	msg += fmt.Sprintf("玩家均分：%d (%.1f%%)\n", int64(heroData.AvgJJL), heroData.Rank.AvgJJL)
	msg += fmt.Sprintf("场均耗时：%.1f (%.1f%%) \n", heroData.AvgUsedTime/60, heroData.Rank.AvgUsedTime)
	msg += fmt.Sprintf("场均补刀：%.1f (%.1f%%)\n", heroData.AvgHit, heroData.Rank.AvgHit)
	msg += fmt.Sprintf("场均每分补刀：%.2f (%.1f%%)\n", heroData.AvgHitPerMinite, heroData.Rank.AvgHitPerMinite)
	msg += fmt.Sprintf("场均kda：%.1f (%.1f%%) / %.1f (%.1f%%) / %.1f (%.1f%%)\n",
		heroData.AvgKill, heroData.Rank.AvgKill, heroData.AvgDeath, heroData.Rank.AvgDeath, heroData.AvgAssist, heroData.Rank.AvgAssist)
	msg += fmt.Sprintf("场均每分kda：%.2f (%.1f%%) / %.2f (%.1f%%) / %.2f (%.1f%%)\n",
		heroData.AvgKillPerMinite, heroData.Rank.AvgKillPerMinite,
		heroData.AvgDeathPerMinite, heroData.Rank.AvgDeathPerMinite,
		heroData.AvgAssistPerMinite, heroData.Rank.AvgAssistPerMinite)
	msg += fmt.Sprintf("场均推塔：%.1f (%.1f%%)\n", heroData.AvgTower, heroData.Rank.AvgTower)
	msg += fmt.Sprintf("场均每分推塔：%.2f (%.1f%%)\n", heroData.AvgTowerPerMinite, heroData.Rank.AvgTowerPerMinite)
	msg += fmt.Sprintf("场均插/排眼：%.2f (%.1f%%) / %.2f (%.1f%%)\n",
		heroData.AvgPutEye, heroData.Rank.AvgPutEye,
		heroData.AvgDestryEye, heroData.Rank.AvgDestryEye)
	msg += fmt.Sprintf("场均每分插/排眼：%.3f (%.1f%%) / %.3f (%.1f%%)\n",
		heroData.AvgPutEyePerMinite, heroData.Rank.AvgPutEyePerMinite,
		heroData.AvgDestryEyePerMinite, heroData.Rank.AvgDestryEyePerMinite)
	msg += fmt.Sprintf("场均经济：%.1f (%.1f%%)\n", heroData.AvgMoney, heroData.Rank.AvgMoney)
	msg += fmt.Sprintf("场均每分经济：%.1f (%.1f%%)\n", heroData.AvgMoneyPerMinite, heroData.Rank.AvgMoneyPerMinite)
	msg += fmt.Sprintf("场均输出：%.1f (%.1f%%)\n", heroData.AvgMakeDamage, heroData.Rank.AvgMakeDamage)
	msg += fmt.Sprintf("场均每分输出：%.1f (%.1f%%)\n", heroData.AvgMakeDamagePerMinite, heroData.Rank.AvgMakeDamagePerMinite)
	msg += fmt.Sprintf("场均承伤：%.1f (%.1f%%)\n", heroData.AvgTakeDamage, heroData.Rank.AvgTakeDamage)
	msg += fmt.Sprintf("场均每分承伤：%.1f (%.1f%%)\n", heroData.AvgTakeDamagePerMinite, heroData.Rank.AvgTakeDamagePerMinite)
	msg += fmt.Sprintf("场均经济转换率：%.1f%% (%.1f%%)\n", heroData.AvgMoneyConversionRate, heroData.Rank.AvgMoneyConversionRate)
	msg += fmt.Sprintf("综合评分：%d (%.1f%%)\n", uint64(heroData.Score), heroData.Rank.Score)
	return
}

func ExportLikeAnalysis(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	rs, total := analysis.HeroAnalysis(PlayerID, 0)
	_, _, _, _, heroDetail, _ := analysis.JJLCompositionAnalysis(PlayerID, 30*24*time.Hour)
	msg += fmt.Sprintf("英雄分析，昵称：%s，总场次：%d(只会计算近30天战绩)\n", name, total)
	for i := range rs {
		if i >= 5 {
			break
		}
		heroData, _ := analysis.GetRankFromPlayers(int(rs[i][0]), 0, []uint64{PlayerID})
		if _, ok := heroData[PlayerID]; !ok {
			return "", errors.New("异常错误")
		}
		msg += fmt.Sprintf("%d、英雄：%s，场次：%d，胜率：%.1f%%，净上分：%d，评分：%.1f(%.1f%%)\n",
			i+1,
			db.HeroIDToName[int(rs[i][0])],
			int(heroData[PlayerID].ActualTotal),
			heroData[PlayerID].ActualWin/heroData[PlayerID].ActualTotal*100,
			int(heroDetail[int(rs[i][0])][1]+heroDetail[int(rs[i][0])][0]),
			heroData[PlayerID].Score,
			heroData[PlayerID].Rank.Score,
		)
	}
	return
}

func ExportFindPlayer(args string, page int) (msg string, err error) {
	argArr := strings.Split(args, "/")
	if len(argArr) < 4 {
		return "", errors.New("填入的参数不足")
	}
	if _, ok := db.HeroNameToID[argArr[0]]; !ok {
		return "", fmt.Errorf("不存在 %s 英雄", argArr[0])
	}
	argArr[0] = fmt.Sprint(db.HeroNameToID[argArr[0]])
	keyArr := []string{"hero_id", "kill_player", "death", "assist", "kill_unit", "fv"}
	if db.HasPartition() {
		plays := []db.PlayerPartition{}
		tmp := db.SqlDB.Model(db.PlayerPartition{})
		for i := range argArr {
			if argArr[i] != "" {
				tmp = tmp.Where(fmt.Sprintf("%s = %v", keyArr[i], argArr[i]))
			}
		}
		tmp.Order("create_time desc").Find(&plays)
		msg += fmt.Sprintf("%d条匹配的游戏记录：\n", len(plays))
		msg += "名称 k/d/a/补 时间\n"
		if len(plays) != 0 && (page-1)*10 > len(plays)-1 {
			return "页码超出范围", nil
		}
		for i := range plays[(page-1)*10 : min(page*10, len(plays))] {
			msg += fmt.Sprintf("%s %d/%d/%d/%d %s\n", collect.SearchName(plays[i].PlayerID), plays[i].KillPlayer, plays[i].Death, plays[i].Assist, plays[i].KillUnit, time.Unix(int64(plays[i].CreateTime), 0).Format(time.DateTime))
		}
		total := len(plays) / 10
		if len(plays)%10 > 0 {
			total++
		}
		msg += fmt.Sprintf("页码：%d/%d", page, total)
	} else {
		plays := []db.Player{}
		tmp := db.SqlDB.Model(db.Player{})
		for i := range argArr {
			if argArr[i] != "" {
				tmp = tmp.Where(fmt.Sprintf("%s = %v", keyArr[i], argArr[i]))
			}
		}
		tmp.Order("create_time desc").Find(&plays)
		msg += fmt.Sprintf("%d条匹配的游戏记录：\n", len(plays))
		msg += "名称 k/d/a/补 时间\n"
		if len(plays) != 0 && (page-1)*10 > len(plays)-1 {
			return "页码超出范围", nil
		}
		for i := range plays[(page-1)*10 : min(page*10, len(plays))] {
			msg += fmt.Sprintf("%s %d/%d/%d/%d %s\n", collect.SearchName(plays[i].PlayerID), plays[i].KillPlayer, plays[i].Death, plays[i].Assist, plays[i].KillUnit, time.Unix(int64(plays[i].CreateTime), 0).Format(time.DateTime))
		}
		total := len(plays) / 10
		if len(plays)%10 > 0 {
			total++
		}
		msg += fmt.Sprintf("页码：%d/%d", page, total)
	}
	return msg, err
}

func ExportGlobalHeroAnalysis(HeroName string) (msg string, err error) {
	ps, err := analysis.GlobalHeroAnalysis(HeroName)
	if err != nil {
		return
	}
	MatchIDToPlayers := map[string][]db.PlayerPartition{}
	for i := range ps {
		MatchIDToPlayers[ps[i].MatchID] = append(MatchIDToPlayers[ps[i].MatchID], ps[i])
	}
	// [win, total]
	count := make([][2]float64, len(analysis.DefaultJJLCategoryKeys))
	all := 0.0
	win := 0.0
	for _, players := range MatchIDToPlayers {
		// 双方都有该英雄
		if len(players) >= 2 {
			continue
		}
		// 场次+1
		all += 1
		count[analysis.DefaultJJLCategoryKeys.Index(float64(players[0].FV))][1]++
		if players[0].Result == 1 || players[0].Result == 3 {
			win++
			count[analysis.DefaultJJLCategoryKeys.Index(float64(players[0].FV))][0]++
		}
	}

	msg += fmt.Sprintf("英雄：%s，单方面出现场次：%d(最近30天)\n", HeroName, int(all))
	msg += fmt.Sprintf("全局单方面胜率：%.1f%%\n", win/all*100)
	msg += "[使用者分段(占比): 胜率, 场次]\n"
	for i := range count{
		if count[i][1] > 0{
			msg += fmt.Sprintf("%s(%.1f%%): %.1f%%, %d\n", analysis.DefaultJJLCategoryKeys[i], count[i][1]/all, count[i][0]/count[i][1]*100, int(count[i][1]))
		}
	}
	return
}

func ExportGlobalHeroAnalysis2(HeroName string) (msg string, err error) {
	if _, ok := db.HeroNameToID[HeroName]; !ok {
		return "未知英雄 " + HeroName, nil
	}
	stages, total, timestamp, err := analysis.GetHeroWinRate(db.HeroNameToID[HeroName])
	if err != nil {
		return "", err
	}
	msg += fmt.Sprintf("英雄：%s，单方面出现场次：%d(最近30天)\n", HeroName, int(total[1]))
	msg += fmt.Sprintf("全局单方面胜率：%.1f%%\n", total[0]/total[1]*100)
	msg += "[对局分段(占比): 胜率, 场次]\n"
	for i := range stages {
		if stages[i][1] > 0 {
			msg += fmt.Sprintf("%s(%.1f%%): %.1f%% / %d\n", analysis.DefaultJJLCategoryKeys[i], stages[i][1]/total[1]*100, stages[i][0]/stages[i][1]*100, int(stages[i][1]))
		}
	}
	msg += "缓存时间：" + time.Unix(int64(timestamp), 0).Format(time.DateTime)
	return
}

func ExportTopAnalysis(HeroName string, fv int) (msg string, err error) {
	if _, ok := db.HeroNameToID[HeroName]; !ok {
		if HeroName[:3] == "jjl" {
			var page int64 = 1
			if len(HeroName) != 3 {
				page, err = strconv.ParseInt(HeroName[3:], 10, 64)
				if err != nil {
					return "", errors.New("错误页码")
				}
			}
			now := time.Now()
			t0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Add(-7 * 24 * 60 * 60 * time.Second)
			matches := []db.Match{}
			db.SqlDB.Model(db.Match{}).Where("create_time > ?", t0.Unix()).Find(&matches)
			matchIDs := []string{}
			for i := range matches {
				matchIDs = append(matchIDs, matches[i].MatchID)
			}
			type result struct {
				PlayerID uint64
				FV       int
			}
			allPlays := map[uint64]int{}
			step := 1000
			for start := 0; start < len(matchIDs); start += step {
				plays := []result{}
				end := min(start+step, len(matchIDs))
				db.SqlDB.Raw("select player_id, max(fv) fv from players where match_id in ? group by player_id", matchIDs[start:end]).Scan(&plays)
				for i := range plays {
					allPlays[plays[i].PlayerID] = max(allPlays[plays[i].PlayerID], plays[i].FV)
				}
			}
			plays := []result{}
			for id, fv := range allPlays {
				plays = append(plays, result{id, fv})
			}
			sort.Slice(plays, func(i, j int) bool {
				return plays[i].FV > plays[j].FV
			})
			msg += "竞技力排行榜(只会计算近30天战绩)：\n"
			for i := (page - 1) * 10; i < int64(len(plays)); i++ {
				if i >= page*10 {
					break
				}

				msg += fmt.Sprintf("%d、%s，%d\n", i+1, collect.SearchName(plays[i].PlayerID), plays[i].FV)
			}
			total := len(plays) / 10
			if len(plays)%10 > 0 {
				total++
			}
			msg += fmt.Sprintf("总计人数：%d\n页码：%d/%d\n", len(plays), page, total)
			return
		}
		return "", fmt.Errorf("不存在 %s 英雄", HeroName)
	}
	if _, ok := analysis.HeroFactor[db.HeroIDToName[db.HeroNameToID[HeroName]]]; !ok {
		return "", fmt.Errorf("英雄 %s 的计算权重不存在", db.HeroIDToName[db.HeroNameToID[HeroName]])
	}
	data, total := analysis.GetRankFromTop(db.HeroNameToID[HeroName], fv, 10)
	msg += fmt.Sprintf("英雄：%s，玩家竞技力下限：%d，总计人数：%d(只会计算近30天游玩次数至少5次、至多50次的战绩)\n", HeroName, fv, total)
	for i := range data {
		_, plays, marks, _, _ := analysis.MarkTeam(data[i].PlayerID)
		marksMap := map[string]int{}
		for i := range plays {
			marksMap[plays[i].MatchID] = marks[i]
		}
		cnt := 0.0
		for _, matchID := range data[i].MatchIDs {
			if v, ok := marksMap[matchID]; ok && v == 0 {
				cnt++
			}
		}
		msg += fmt.Sprintf("%d、%s，评分：%.1f\n均分：%d，单排率：%.1f%%\n", i+1, collect.SearchName(data[i].PlayerID), data[i].Score, int(data[i].AvgJJL), cnt/data[i].Total*100)
	}
	return
}

func ExportWinRateAnalysis(pageStr, scope string) (msg string, err error) {
	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
	}
	if err != nil || page < 1 {
		return msg, errors.New("错误页码")
	}
	dataMap := map[string][2]int{}
	if scope == "" {
		scope = "total"
	}
	rs, _ := db.RDB.ZRangeWithScores(context.Background(), analysis.HeroPlayCountKeyPrefix+scope, 0, -1).Result()
	for i := range rs {
		dataMap[rs[i].Member.(string)] = [2]int{0, int(rs[i].Score)}
	}
	rs, _ = db.RDB.ZRangeWithScores(context.Background(), analysis.HeroWinCountKeyPrefix+scope, 0, -1).Result()
	for i := range rs {
		if v, ok := dataMap[rs[i].Member.(string)]; ok {
			v[0] = int(rs[i].Score)
			dataMap[rs[i].Member.(string)] = v
		}
	}
	data := [][3]int{}
	for name, arr := range dataMap {
		data = append(data, [3]int{db.HeroNameToID[name], arr[0], arr[1]})
	}
	if err != nil {
		return "", err
	}
	sort.Slice(data, func(i, j int) bool {
		return float64(data[i][1])/float64(data[i][2]) > float64(data[j][1])/float64(data[j][2])
	})
	if scope == "" {
		scope = "全局"
	}
	msg += fmt.Sprintf("近30天英雄胜率排行，范围：%s\n", scope)
	msg += "英雄名称 胜率 场次\n"
	for i := (page - 1) * 10; i < min(len(data), page*10); i++ {
		msg += fmt.Sprintf("%d、%s，%.1f%%，%d\n", i+1, db.HeroIDToName[data[i][0]], float64(data[i][1])/float64(data[i][2])*100, data[i][2])
	}
	total := len(data) / 10
	if len(data)%10 > 0 {
		total++
	}
	msg += fmt.Sprintf("页码：%d/%d", page, total)
	return
}

func ExportTopWithDetailAnalysis(HeroName string, fv int) (msg string, err error) {
	if _, ok := db.HeroNameToID[HeroName]; !ok {
		return "", fmt.Errorf("不存在 %s 英雄", HeroName)
	}
	if _, ok := analysis.HeroFactor[db.HeroIDToName[db.HeroNameToID[HeroName]]]; !ok {
		return "", fmt.Errorf("英雄 %s 的计算权重不存在", db.HeroIDToName[db.HeroNameToID[HeroName]])
	}

	data, total := analysis.GetRankFromTop(db.HeroNameToID[HeroName], fv, 10)

	msg += fmt.Sprintf("英雄：%s，玩家竞技力下限：%d，总计人数：%d(只会计算近30天游玩次数超过五次的战绩)\n", HeroName, fv, total)
	for i := range data {
		msg += fmt.Sprintf("%d、%s，评分：%.1f\n", i+1, collect.SearchName(data[i].PlayerID), data[i].Score)
	}
	msg += fmt.Sprintln("英雄详情：")
	for i, factor := range analysis.HeroFactor[db.HeroIDToName[db.HeroNameToID[HeroName]]] {
		if factor > 0.0 {
			msg += fmt.Sprintf("%.2f * %s，", factor, analysis.HeroTypeSlice[i])
		}
	}
	msg += "\n计算详情：\n"
	factorSlice := [][2]float64{}
	for i, factor := range analysis.MergeImportance(analysis.HeroFactor[db.HeroIDToName[db.HeroNameToID[HeroName]]]) {
		factorSlice = append(factorSlice, [2]float64{float64(i), factor})
	}
	sort.Slice(factorSlice, func(i, j int) bool {
		return factorSlice[i][0] < factorSlice[j][0]
	})
	for i := range factorSlice {
		msg += fmt.Sprintf("%s：%.2f\n", analysis.ImportanceMap[int(factorSlice[i][0])], factorSlice[i][1])
	}
	return
}

func ExportFlushTop(HeroName string) (msg string, err error) {
	if _, ok := db.HeroNameToID[HeroName]; !ok {
		return "", fmt.Errorf("不存在 %s 英雄", HeroName)
	}
	data, _ := analysis.GetRankFromTop(db.HeroNameToID[HeroName], 0, 10)
	for i := range data {
		err = collect.CrawlPlayerByName(fmt.Sprintf("id:%d", data[i].PlayerID))
		if err != nil {
			return
		}
	}
	return fmt.Sprintf("刷新 %s 月榜完毕", HeroName), nil
}

func ExportJJLWithTeamAnalysis(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}

	timeRange, jjl, team := analysis.JJLWithTeamAnalysis(PlayerID)
	if len(timeRange) == 0 {
		return "", errors.New("分析异常")
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "竞技力与开黑情况",
			Subtitle: "玩家：" + name,
			Left:     "10%",
		}),
		charts.WithSingleAxisOpts(opts.SingleAxis{
			Type: "time",
		}),
		charts.WithYAxisOpts(opts.YAxis{Name: "场次", NameLocation: "start", Show: true, AxisLabel: &opts.AxisLabel{Show: true}}),
	)
	generateData := func(index int) (arr []opts.BarData) {
		var sum uint64
		for i := range team {
			sum += team[i][index]
			arr = append(arr, opts.BarData{Value: sum})
		}
		return
	}
	data0 := generateData(0)
	data1 := generateData(1)
	data2 := generateData(2)
	data3 := generateData(3)
	bar.SetXAxis(timeRange)
	bar.AddSeries("单排", data0, charts.WithBarChartOpts(opts.BarChart{YAxisIndex: 0})).
		AddSeries("双排", data1, charts.WithBarChartOpts(opts.BarChart{YAxisIndex: 0})).
		AddSeries("三黑", data2, charts.WithBarChartOpts(opts.BarChart{YAxisIndex: 0})).
		AddSeries("四黑", data3, charts.WithBarChartOpts(opts.BarChart{YAxisIndex: 0})).
		SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
			Stack:          "stackA",
			BarCategoryGap: "0%",
		}))
	bar.ExtendYAxis(opts.YAxis{Name: "竞技力", NameLocation: "start", Show: true, AxisLabel: &opts.AxisLabel{Show: true}, Min: "dataMin", Max: "dataMax"})

	data4 := []opts.LineData{}
	for i := range jjl {
		data4 = append(data4, opts.LineData{Value: jjl[i]})
	}
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithAnimation(),
		charts.WithSingleAxisOpts(opts.SingleAxis{
			Type:   "time",
			Bottom: "10%",
		}),
	)
	line.SetXAxis(timeRange).AddSeries("竞技力", data4, charts.WithLineChartOpts(opts.LineChart{YAxisIndex: 1})).
		SetSeriesOptions(charts.WithMarkPointNameTypeItemOpts(
			opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
			opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
		))

	bar.Overlap(line)
	f, err := os.Create(fmt.Sprintf("./files/%d.html", PlayerID))
	if err != nil {
		return "", err
	}
	err = bar.Render(f)
	if err != nil {
		return "", err
	}
	err = SavePNG(fmt.Sprintf("./files/%d", PlayerID))
	if err != nil {
		return "", err
	}
	abs, _ := filepath.Abs(fmt.Sprintf("./files/%d", PlayerID))
	if runtime.GOOS == "windows" {
		abs = "/" + abs
	}
	return fmt.Sprintf("[CQ:image,file=file://%s.png]", abs), nil
}

func ExportJJLCompositionAnalysis(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	team, _, enermies, _, hero, total := analysis.JJLCompositionAnalysis(PlayerID, 0)
	heroArr := [][4]float64{}
	for heroID, data := range hero {
		heroArr = append(heroArr, [4]float64{data[0], data[1], data[2], float64(heroID)})
	}
	sort.Slice(heroArr, func(i int, j int) bool { return heroArr[i][2] > heroArr[j][2] })
	msg += fmt.Sprintf("玩家昵称：%s，总场次：%d\n", name, total)
	msg += "jjl来自开黑情况：\n"
	ranges := []string{"单排", "双排", "三黑", "四黑"}
	for i := range ranges {
		if team[i][2] > 0 {
			msg += fmt.Sprintf("%s%d场，占比%.1f%%，净上分：%d\n", ranges[i], int(team[i][2]), float64(team[i][2])/float64(total)*100, int(team[i][0])+int(team[i][1]))
		}
	}
	ranges = analysis.DefaultJJLCategoryKeys
	msg += "jjl来自对手玩家的分段情况：\n"
	for i := range ranges {
		if enermies[i][2] > 0 {
			msg += fmt.Sprintf("分段%s，%d人，净上分：%d\n", ranges[i], int(enermies[i][2]), int(enermies[i][0])+int(enermies[i][1]))
		}
	}
	other := [3]int{}
	if len(heroArr) > 5 {
		for i := 5; i < len(heroArr); i++ {
			other[0] += int(heroArr[i][0])
			other[1] += int(heroArr[i][1])
			other[2] += int(heroArr[i][2])
		}
	}
	msg += "jjl来自英雄情况：\n"
	for i := range heroArr {
		if i > 5 {
			break
		}
		msg += fmt.Sprintf("%d、%s，%d场，占比%.1f%%，净上分：%d\n", i+1, db.HeroIDToName[int(heroArr[i][3])], int(heroArr[i][2]), heroArr[i][2]/float64(total)*100, int(heroArr[i][0])+int(heroArr[i][1]))
	}
	if other[2] > 0 {
		msg += fmt.Sprintf("其他，%d场，占比%.1f%%，净上分：%d", int(other[2]), float64(other[2])/float64(total)*100, int(other[0])+int(other[1]))
	}
	return
}

func ExportPKAnalysis(name string, hero string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	if _, ok := db.HeroNameToID[hero]; !ok {
		return "未知英雄 " + name, nil
	}
	if _, ok := analysis.HeroFactor[db.HeroIDToName[db.HeroNameToID[hero]]]; !ok {
		return "", fmt.Errorf("英雄 %s 的计算权重不存在", db.HeroIDToName[db.HeroNameToID[hero]])
	}
	you, top1, err := analysis.PKAnalysis(PlayerID, db.HeroNameToID[hero])
	if err != nil {
		return
	}

	var indicators = []*opts.Indicator{
		{Name: "胜率", Max: float32(max(you[0], top1[0]) * 1.2)},
		{Name: "场均耗时", Max: float32(max(you[1], top1[1]) * 1.2)},
		{Name: "场均每分补刀", Max: float32(max(you[2], top1[2]) * 1.2)},
		{Name: "场均每分击杀", Max: float32(max(you[3], top1[3]) * 1.2)},
		{Name: "场均每分死亡", Max: float32(max(you[4], top1[4]) * 1.2)},
		{Name: "场均每分助攻", Max: float32(max(you[5], top1[5]) * 1.2)},
		{Name: "场均推塔", Max: float32(max(you[6], top1[6]) * 1.2)},
		{Name: "场均插眼", Max: float32(max(you[7], top1[7]) * 1.2)},
		{Name: "场均排眼", Max: float32(max(you[8], top1[8]) * 1.2)},
		{Name: "场均每分经济", Max: float32(max(you[9], top1[9]) * 1.2)},
		{Name: "场均每分输出", Max: float32(max(you[10], top1[10]) * 1.2)},
		{Name: "场均每分承伤", Max: float32(max(you[11], top1[11]) * 1.2)},
		{Name: "场均经济转换率", Max: float32(max(you[12], top1[12]) * 1.2)},
		{Name: "综合评分", Max: float32(max(you[13], top1[13]) * 1.2), Min: -100},
	}

	radar := charts.NewRadar()
	radar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "英雄场均数据及榜一对比",
			Subtitle: "英雄：" + hero,
			Left:     "20%",
			TitleStyle: &opts.TextStyle{
				Color: "#eee",
			},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			BackgroundColor: "#161627",
		}),
		charts.WithRadarComponentOpts(opts.RadarComponent{
			Indicator:   indicators,
			Shape:       "polygon",
			SplitNumber: 5,
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Opacity: 0.1,
				},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{
			Show:   true,
			Bottom: "5px",
			TextStyle: &opts.TextStyle{
				Color: "#eee",
			},
		}),
	)

	radar.AddSeries("排行第一", []opts.RadarData{{Value: top1}},
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "#F9713C"}),
		charts.WithLabelOpts(opts.Label{
			Show:     true,
			Position: "top",
			Color:    "#F9713C",
		})).
		AddSeries(name, []opts.RadarData{{Value: you}},
			charts.WithItemStyleOpts(opts.ItemStyle{Color: "#B3E4A1"}),
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "bottom",
				Color:    "#B3E4A1",
			})).
		SetSeriesOptions(
			charts.WithLineStyleOpts(opts.LineStyle{
				Width:   1,
				Opacity: 0.5,
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: 0.1,
			}),
		)
	f, err := os.Create(fmt.Sprintf("./files/%d_pk.html", PlayerID))
	if err != nil {
		return "", err
	}
	err = radar.Render(f)
	if err != nil {
		return "", err
	}
	err = SavePNG(fmt.Sprintf("./files/%d_pk", PlayerID))
	if err != nil {
		return "", err
	}
	abs, _ := filepath.Abs(fmt.Sprintf("./files/%d_pk", PlayerID))
	if runtime.GOOS == "windows" {
		abs = "/" + abs
	}
	return fmt.Sprintf("[CQ:image,file=file://%s.png]", abs), nil
}

func ExportActiveAnalysis() (msg string, err error) {
	now := time.Now()
	t0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Add(-7 * 24 * 60 * 60 * time.Second)
	matches := []db.Match{}
	db.SqlDB.Model(db.Match{}).Where("create_time > ?", t0.Unix()).Find(&matches)
	matchIDs := []string{}
	for i := range matches {
		matchIDs = append(matchIDs, matches[i].MatchID)
	}
	type result struct {
		PlayerID uint64
		FV       int
	}
	allPlays := map[uint64]int{}
	step := 1000
	for start := 0; start < len(matchIDs); start += step {
		plays := []result{}
		end := min(start+step, len(matchIDs))
		db.SqlDB.Raw("select player_id, max(fv) fv from players where match_id in ? group by player_id", matchIDs[start:end]).Scan(&plays)
		for i := range plays {
			allPlays[plays[i].PlayerID] = max(allPlays[plays[i].PlayerID], plays[i].FV)
		}
	}
	count := make([]int, len(analysis.DefaultJJLCategoryKeys))
	for _, fv := range allPlays {
		count[analysis.DefaultJJLCategoryKeys.Index(float64(fv))]++
	}

	items := make([]opts.PieData, 0)
	for i := range count {
		name := fmt.Sprintf("%s(%.2f%%)", analysis.DefaultJJLCategoryKeys[i], float64(count[i])/float64(len(allPlays))*100)
		items = append(items, opts.PieData{Name: name, Value: count[i]})
	}
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "过去七天活跃玩家数量及分布(仅供参考)",
			Subtitle: fmt.Sprintf("玩家总数：%d", len(allPlays)),
			Top:      "0%",
			Left:     "10%",
		}),
		charts.WithLegendOpts(
			opts.Legend{
				Show: false,
			},
		),
	)
	pie.AddSeries("pie", items).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
			charts.WithPieChartOpts(
				opts.PieChart{
					RoseType: "radius",
					Radius:   []string{"30%", "75%"},
				},
			),
		)

	f, err := os.Create("./files/active.html")
	if err != nil {
		return "", err
	}
	err = pie.Render(f)
	if err != nil {
		return "", err
	}
	err = SavePNG("./files/active")
	if err != nil {
		return "", err
	}
	abs, _ := filepath.Abs("./files/active")
	if runtime.GOOS == "windows" {
		abs = "/" + abs
	}
	return fmt.Sprintf("[CQ:image,file=file://%s.png]", abs), nil
}

func SavePNG(file string) (err error) {
	abs, _ := filepath.Abs(file)

	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var buf []byte

	task := chromedp.Tasks{
		chromedp.Navigate("file://" + abs + ".html"),
		// chromedp.EmulateViewport(1200, 600),
		chromedp.FullScreenshot(&buf, 100),
	}

	if err = chromedp.Run(ctx, task); err != nil {
		return
	}
	if err = os.WriteFile(abs+".png", buf, 0644); err != nil {
		return
	}
	return
}
