package print

import (
	"eebot/bot/service/analysis300/analysis"
	"eebot/bot/service/analysis300/collect"
	"eebot/bot/service/analysis300/db"
	"fmt"
)

func PrintTeamAnalysis(playerID uint64, limit int) {
	allies, enermies, total := analysis.TeamAnalysis(playerID)
	fmt.Printf("昵称：%s，记录场次：%d\n", collect.SearchName(playerID), total)
	fmt.Println("队友情况：")
	for i := range allies {
		if i >= limit {
			break
		}
		fmt.Printf("%d、%s，作为队友场次：%d，%.1f%%\n", i+1, "昵称："+collect.SearchName(allies[i][0]), allies[i][1], Divide(allies[i][1], uint64(total))*100)
	}
	fmt.Println("对手情况：")
	for i := range enermies {
		if i >= limit {
			break
		}
		fmt.Printf("%d、%s，作为对手场次：%d，%.1f%%\n", i+1, "昵称："+collect.SearchName(enermies[i][0]), enermies[i][1], Divide(enermies[i][1], uint64(total))*100)
	}
}

func PrintWinOrLoseAnalysis(PlayerID uint64) {
	rs, _, fvRange, fvNow := analysis.WinOrLoseAnalysis(PlayerID)

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

	fmt.Printf("昵称：%s，记录场次：%d，团分跨度：%d - %d，当前团分：%d\n", collect.SearchName(PlayerID), len(rs), fvRange[0], fvRange[1], fvNow)
	fmt.Printf("胜率：%.1f%%\n", float32(win)/float32(win+lose)*100)
	fmt.Printf("总记录 %d 局中有 %d 局(%.1f%%) 己方均分高于对面\n", len(rs), cnt1+lose-cnt2, float32(cnt1+lose-cnt2)/float32(len(rs))*100)
	fmt.Printf("%d 胜场中有 %d 场(%.1f%%)的己方均分高于对面\n", win, cnt1, float32(cnt1)/float32(win)*100)
	fmt.Printf("%d 负场中有 %d 场(%.1f%%)的己方均分低于对面\n", lose, cnt2, float32(cnt2)/float32(lose)*100)

	stage1 := ExtractByFV(1000, 1500, rs)
	stage2 := ExtractByFV(1500, 1700, rs)
	stage3 := ExtractByFV(1700, 1800, rs)
	stage4 := ExtractByFV(1800, 1900, rs)
	stage5 := ExtractByFV(1900, 2000, rs)
	stage6 := ExtractByFV(2000, 2500, rs)

	fmt.Printf("%-20s %-20s %-20s %-20s %-20s %-20s %-20s\n", "range", "1000-1500("+tran2(stage1), "1500-1700("+tran2(stage2), "1700-1800("+tran2(stage3), "1800-1900("+tran2(stage4), "1900-2000("+tran2(stage5), "2000-2500("+tran2(stage6))
	fmt.Printf("%-20s %-20s %-20s %-20s %-20s %-20s %-20s\n", "win(%%) / total", tran(stage1), tran(stage2), tran(stage3), tran(stage4), tran(stage5), tran(stage6))
}

func Divide(a uint64, b uint64) float64 {
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

func ExtractByFVAdvanced(start, end int, result [][4]int) (cnt [2]int) {
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

func PrintHeroAnalysis(PlayerID uint64) {
	rs, total := analysis.HeroAnalysis(PlayerID, 0)
	fmt.Printf("英雄分析，昵称：%s，总场次：%d\n", collect.SearchName(PlayerID), total)
	for i := range rs {
		if i >= 1 {
			break
		}
		fmt.Printf("英雄：%s\n", db.HeroIDToName[int(rs[i][0])])
		fmt.Printf("场次：%d，%.1f%%\n", uint64(rs[i][1]), rs[i][1]/float64(total)*100)
		fmt.Printf("胜率：%.1f%%\n", rs[i][2]/rs[i][1]*100)
		fmt.Printf("场均耗时：%.1f 分\n", rs[i][27]/60)
		fmt.Printf("场均补刀：%.1f\n", rs[i][3])
		fmt.Printf("场均每分补刀：%.2f\n", rs[i][4])
		fmt.Printf("场均kda：%.1f / %.1f / %.1f\n", rs[i][5], rs[i][7], rs[i][9])
		fmt.Printf("场均每分kda：%.2f / %.2f / %.2f\n", rs[i][6], rs[i][8], rs[i][10])
		fmt.Printf("场均推塔：%.1f\n", rs[i][11])
		fmt.Printf("场均每分推塔：%.2f\n", rs[i][12])
		fmt.Printf("场均插/排眼：%.1f / %.1f\n", rs[i][13], rs[i][15])
		fmt.Printf("场均每分插/排眼：%.2f / %.2f\n", rs[i][14], rs[i][16])
		fmt.Printf("场均经济：%.1f\n", rs[i][17])
		fmt.Printf("场均每分经济：%.1f\n", rs[i][18])
		fmt.Printf("场均经济占比：%.1f%%\n", rs[i][19]*100)
		fmt.Printf("场均输出：%.1f\n", rs[i][20])
		fmt.Printf("场均每分输出：%.1f\n", rs[i][21])
		fmt.Printf("场均输出占比：%.1f%%\n", rs[i][22]*100)
		fmt.Printf("场均承伤：%.1f\n", rs[i][23])
		fmt.Printf("场均每分承伤：%.1f\n", rs[i][24])
		fmt.Printf("场均承伤占比：%.1f%%\n", rs[i][25]*100)
		fmt.Printf("场均经济转换率：%.1f%%\n", rs[i][26])
	}
}
