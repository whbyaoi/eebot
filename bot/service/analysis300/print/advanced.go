package print

import (
	"eebot/bot/service/analysis300/analysis"
	"eebot/bot/service/analysis300/collect"
	"eebot/bot/service/analysis300/db"
	"fmt"
	"strings"
)

func PrintShuffleAnalysis(playerID uint64) {
	avg, than10min, total := analysis.ShuffleAnalysis(playerID)

	fmt.Printf("洗牌分析，昵称：%s\n", collect.SearchName(playerID))
	fmt.Printf("有效间隔数：%d，平均间隔：%d秒，超过十分钟的间隔数：%d", total, avg, than10min)
}

func PrintShuffleAnalysisAdvanced(playerID uint64) {
	avg, than10min, total := analysis.ShuffleAnalysis(playerID)

	fmt.Printf("洗牌分析，昵称：%s\n", collect.SearchName(playerID))
	fmt.Printf("本人有效间隔数：%d，平均间隔：%d秒，超过十分钟的间隔数：%d\n", total, avg, than10min)

	allies, _, sum := analysis.TeamAnalysis(playerID)
	for i, ally := range allies {
		if i >= 10 {
			break
		}
		if Divide(ally[1], uint64(sum)) < 2.0 || ally[1] <= 2 {
			continue
		}
		collect.CrawlPlayerByID(ally[0])
		avg, than10min, total := analysis.ShuffleAnalysis(ally[0])
		fmt.Printf("开黑队友 %s 有效间隔数：%d，平均间隔：%d秒，超过十分钟的间隔数：%d\n", collect.SearchName(ally[0]), total, avg, than10min)
	}
}

func PrintTeamAnalysisAdvanced(PlayerID uint64) {
	allies, enermies, teams, teamAllies, total := analysis.TeamAnalysisAdvanced(PlayerID)
	fmt.Printf("昵称：%s，记录场次：%d\n", collect.SearchName(PlayerID), total)
	fmt.Println("队友情况：")
	for i := range allies {
		if i >= 10 {
			break
		}
		fmt.Printf("%d、%s，作为队友场次：%d (%.1f%%)，胜率：%.1f%%\n", i+1, "昵称："+collect.SearchName(allies[i][0]), allies[i][2], Divide(allies[i][2], uint64(total))*100, Divide(allies[i][1], allies[i][2])*100)
	}
	fmt.Println("对手情况：")
	for i := range enermies {
		if i >= 5 {
			break
		}
		fmt.Printf("%d、%s，作为对手场次：%d (%.1f%%)，胜率：%.1f%%\n", i+1, "昵称："+collect.SearchName(enermies[i][0]), enermies[i][2], Divide(enermies[i][2], uint64(total))*100, Divide(enermies[i][1], enermies[i][2])*100)
	}

	var arr [4][]string
	for k, v := range teamAllies {
		for name := range v {
			arr[k] = append(arr[k], name)
		}
	}
	fmt.Println("开黑情况(仅供参考)：")
	fmt.Printf("单排 %d 局，胜率 %.1f%%\n", teams[0][1], Divide(teams[0][0], teams[0][1])*100)
	fmt.Printf("双排 %d 局，胜率 %.1f%%，黑车队友：%s\n", teams[1][1], Divide(teams[1][0], teams[1][1])*100, strings.Join(arr[1], "，"))
	fmt.Printf("三黑 %d 局，胜率 %.1f%%，黑车队友：%s\n", teams[2][1], Divide(teams[2][0], teams[2][1])*100, strings.Join(arr[2], "，"))
	fmt.Printf("四黑 %d 局，胜率 %.1f%%，黑车队友：%s\n", teams[3][1], Divide(teams[3][0], teams[3][1])*100, strings.Join(arr[3], "，"))
}

func PrintWinOrLoseAnalysisAdvanced(PlayerID uint64) {
	rs, diff, fvRange, fvNow, _ := analysis.WinOrLoseAnalysisAdvanced(PlayerID)

	win := 0
	lose := 0
	cnt1 := 0
	cnt2 := 0
	for i := range rs {
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
	}

	tran := func(stage [2]int) string {
		return fmt.Sprintf("%.1f%% / %d", float32(stage[0])/float32(stage[0]+stage[1])*100, stage[0]+stage[1])
	}

	tran2 := func(stage [2]int) string {
		return fmt.Sprintf("%.1f%%", float32(stage[0]+stage[1])/float32(len(rs))*100)
	}

	fmt.Printf("基本信息：\n")
	fmt.Printf("昵称：%s，记录场次：%d，团分跨度：%d - %d，当前团分：%d\n", collect.SearchName(PlayerID), len(rs), fvRange[0], fvRange[1], fvNow)
	fmt.Printf("胜率：%.1f%%\n", float32(win)/float32(win+lose)*100)
	fmt.Printf("均分分析：\n")
	fmt.Printf("总记录 %d 局中有 %d 局(%.1f%%) 己方均分高于对面\n", len(rs), cnt1+lose-cnt2, float32(cnt1+lose-cnt2)/float32(len(rs))*100)
	fmt.Printf("%d 胜场中有 %d 场(%.1f%%)的己方均分高于对面\n", win, cnt1, float32(cnt1)/float32(win)*100)
	fmt.Printf("%d 负场中有 %d 场(%.1f%%)的己方均分低于对面\n", lose, cnt2, float32(cnt2)/float32(lose)*100)
	fmt.Printf("总记录 %d 局平均水平来看，玩家分相对场均分（正数为高，负数为低） %d\n", len(rs), diff)

	stage1 := ExtractByFVAdvanced(1000, 1500, rs)
	stage2 := ExtractByFVAdvanced(1500, 1700, rs)
	stage3 := ExtractByFVAdvanced(1700, 1800, rs)
	stage4 := ExtractByFVAdvanced(1800, 1900, rs)
	stage5 := ExtractByFVAdvanced(1900, 2000, rs)
	stage6 := ExtractByFVAdvanced(2000, 2500, rs)
	fmt.Printf("匹配分析：\n")
	fmt.Printf("%-20s %-20s %-20s %-20s %-20s %-20s %-20s\n", "range", "1000-1500("+tran2(stage1), "1500-1700("+tran2(stage2), "1700-1800("+tran2(stage3), "1800-1900("+tran2(stage4), "1900-2000("+tran2(stage5), "2000-2500("+tran2(stage6))
	fmt.Printf("%-20s %-20s %-20s %-20s %-20s %-20s %-20s\n", "win(%%) / total", tran(stage1), tran(stage2), tran(stage3), tran(stage4), tran(stage5), tran(stage6))

	var a1, a2, a3 uint64 = 0, 0, 0
	for i := range rs {
		avg := (rs[i][0] + rs[i][1]) / 2
		flag := IsSameRange(avg, rs[i][3])
		if flag == 0 {
			a1++
		} else if flag == 1 {
			a2++
		} else if flag == 2 {
			a3++
		}

	}
	fmt.Printf("进入低分局场次：%d(%.1f%%)\n", a1, Divide(a1, uint64(len(rs)))*100)
	fmt.Printf("进入本地局场次：%d(%.1f%%)\n", a2, Divide(a2, uint64(len(rs)))*100)
	fmt.Printf("进入高分局场次：%d(%.1f%%)\n", a3, Divide(a3, uint64(len(rs)))*100)
}

// IsSameRange
// return: 0-低分 1-本地 2-高分
func IsSameRange(avg, fv int) int {
	diff := fv - avg
	if diff < -100 {
		return 2
	} else if diff >= 100 {
		return 0
	}
	return 1
}

func PrintAssignHeroAnalysisAdvanced(PlayerID uint64, name string, fv int) {
	rs, total := analysis.HeroAnalysis(PlayerID, fv)
	analysis.UpdateHeroOfPlayerRank(db.HeroNameToID[name], fv)
	rank, total2 := analysis.GetHeroOfPlayerRank(db.HeroNameToID[name], PlayerID, fv)
	fmt.Printf("英雄分析，昵称：%s，总场次：%d\n", collect.SearchName(PlayerID), total)
	for i := range rs {
		if db.HeroIDToName[int(rs[i][0])] != name {
			continue
		}
		fmt.Printf("英雄：%s\n", name)
		fmt.Printf("场次：%d，占比%.1f%% (超越%d人数百分比：%.1f%%，下同)\n", uint64(rs[i][1]), rs[i][1]/float64(total)*100, total2, rank[1])
		fmt.Printf("胜率：%.1f%% (%.1f%%，仅该数据仅供参考)\n", rs[i][2]/rs[i][1]*100, rank[2])
		fmt.Printf("场均耗时：%.1f (%.1f%%) 分 \n", rs[i][27]/60, rank[27])
		fmt.Printf("场均补刀：%.1f (%.1f%%)\n", rs[i][3], rank[3])
		fmt.Printf("场均每分补刀：%.2f (%.1f%%)\n", rs[i][4], rank[4])
		fmt.Printf("场均kda：%.1f (%.1f%%) / %.1f (%.1f%%) / %.1f (%.1f%%)\n", rs[i][5], rank[5], rs[i][7], rank[7], rs[i][9], rank[9])
		fmt.Printf("场均每分kda：%.2f (%.1f%%) / %.2f (%.1f%%) / %.2f (%.1f%%)\n", rs[i][6], rank[6], rs[i][8], rank[8], rs[i][10], rank[10])
		fmt.Printf("场均推塔：%.1f (%.1f%%)\n", rs[i][11], rank[11])
		fmt.Printf("场均每分推塔：%.2f (%.1f%%)\n", rs[i][12], rank[12])
		fmt.Printf("场均插/排眼：%.1f (%.1f%%) / %.1f (%.1f%%)\n", rs[i][13], rank[13], rs[i][15], rank[15])
		fmt.Printf("场均每分插/排眼：%.2f (%.1f%%) / %.2f (%.1f%%)\n", rs[i][14], rank[14], rs[i][16], rank[16])
		fmt.Printf("场均经济：%.1f (%.1f%%)\n", rs[i][17], rank[17])
		fmt.Printf("场均每分经济：%.1f (%.1f%%)\n", rs[i][18], rank[18])
		fmt.Printf("场均经济占比：%.1f%% (%.1f%%)\n", rs[i][19]*100, rank[19])
		fmt.Printf("场均输出：%.1f (%.1f%%)\n", rs[i][20], rank[20])
		fmt.Printf("场均每分输出：%.1f (%.1f%%)\n", rs[i][21], rank[21])
		fmt.Printf("场均输出占比：%.1f%% (%.1f%%)\n", rs[i][22]*100, rank[22])
		fmt.Printf("场均承伤：%.1f (%.1f%%)\n", rs[i][23], rank[23])
		fmt.Printf("场均每分承伤：%.1f (%.1f%%)\n", rs[i][24], rank[24])
		fmt.Printf("场均承伤占比：%.1f%% (%.1f%%)\n", rs[i][25]*100, rank[25])
		fmt.Printf("场均经济转换率：%.1f%% (%.1f%%)\n", rs[i][26], rank[26])
	}
}
