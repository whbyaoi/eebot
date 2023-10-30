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
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func ExportTeamAnalysisAdvanced(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	allies, _, teams, teamAllies, total := analysis.TeamAnalysisAdvanced(PlayerID)
	msg += fmt.Sprintf("昵称：%s，记录场次：%d\n", name, total)
	msg += fmt.Sprintln("队友情况：")
	for i := range allies {
		if i >= 10 {
			break
		}
		msg += fmt.Sprintf("%d、%s，作为队友场次：%d (%.1f%%)，胜率：%.1f%%\n", i+1, "昵称："+collect.SearchName(allies[i][0]), allies[i][2], analysis.Divide(allies[i][2], uint64(total))*100, analysis.Divide(allies[i][1], allies[i][2])*100)
	}
	// msg += fmt.Sprintln("对手情况：")
	// for i := range enermies {
	// 	if i >= 5 {
	// 		break
	// 	}
	// 	msg += fmt.Sprintf("%d、%s，作为对手场次：%d (%.1f%%)，胜率：%.1f%%\n", i+1, "昵称："+collect.SearchName(enermies[i][0]), enermies[i][2], analysis.Divide(enermies[i][2], uint64(total))*100, analysis.Divide(enermies[i][1], enermies[i][2])*100)
	// }

	var arr [4][]string
	for k, v := range teamAllies {
		for name := range v {
			arr[k] = append(arr[k], name)
		}
	}
	msg += fmt.Sprintln("开黑情况(仅供参考)：")
	msg += fmt.Sprintf("单排 %d(%.1f%%) 局，胜率 %.1f%%\n", teams[0][1], analysis.Divide(teams[0][1], uint64(total))*100, analysis.Divide(teams[0][0], teams[0][1])*100)
	msg += fmt.Sprintf("双排 %d(%.1f%%) 局，胜率 %.1f%%，黑车队友：%s\n", teams[1][1], analysis.Divide(teams[1][1], uint64(total))*100, analysis.Divide(teams[1][0], teams[1][1])*100, strings.Join(arr[1], "，"))
	msg += fmt.Sprintf("三黑 %d(%.1f%%) 局，胜率 %.1f%%，黑车队友：%s\n", teams[2][1], analysis.Divide(teams[2][1], uint64(total))*100, analysis.Divide(teams[2][0], teams[2][1])*100, strings.Join(arr[2], "，"))
	msg += fmt.Sprintf("四黑 %d(%.1f%%) 局，胜率 %.1f%%，黑车队友：%s\n", teams[3][1], analysis.Divide(teams[3][1], uint64(total))*100, analysis.Divide(teams[3][0], teams[3][1])*100, strings.Join(arr[3], "，"))
	return
}

func ExportWinOrLoseAnalysisAdvanced(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}

	rs, diff, svd, fvRange, fvNow, timeRange := analysis.WinOrLoseAnalysisAdvanced(PlayerID)
	if len(rs) == 0 {
		return "", errors.New("查询不到任何战绩")
	}
	win := 0
	lose := 0
	cnt1 := 0
	cnt2 := 0
	diff2 := 0
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
		diff2 += rs[i][0] - rs[i][1]
	}
	diff2 /= len(rs)

	tran := func(stage [2]int) string {
		return fmt.Sprintf("%.1f%% / %d", float32(stage[0])/float32(stage[0]+stage[1])*100, stage[0]+stage[1])
	}

	tran2 := func(stage [2]int) string {
		return fmt.Sprintf("%.1f%%", float32(stage[0]+stage[1])/float32(len(rs))*100)
	}
	msg += fmt.Sprintf("昵称：%s，记录场次：%d，团分跨度：%d - %d，时间跨度：%s - %s\n", name, len(rs), fvRange[0], fvRange[1], time.Unix(int64(timeRange[0]), 0).Format("20060102"), time.Unix(int64(timeRange[1]), 0).Format("20060102"))
	msg += fmt.Sprintf("当前团分：%d，胜率：%.1f%%\n", fvNow, float32(win)/float32(win+lose)*100)
	msg += fmt.Sprintf("总记录 %d 局中有 %d 局(%.1f%%) 己方均分高于对面\n", len(rs), cnt1+lose-cnt2, float32(cnt1+lose-cnt2)/float32(len(rs))*100)
	msg += fmt.Sprintf("%d 胜场中有 %d 场(%.1f%%)的己方均分高于对面\n", win, cnt1, float32(cnt1)/float32(win)*100)
	msg += fmt.Sprintf("%d 负场中有 %d 场(%.1f%%)的己方均分低于对面\n", lose, cnt2, float32(cnt2)/float32(lose)*100)
	msg += fmt.Sprintf("玩家分相对场均分水平：%d\n", diff)
	msg += fmt.Sprintf("己方均分相对敌方均分水平：%d\n", diff2)
	msg += fmt.Sprintf("己方团分离散度相对敌方团分离散度水平：%d\n", svd)

	stage1 := analysis.ExtractByFVAdvanced(1000, 1500, rs)
	stage2 := analysis.ExtractByFVAdvanced(1500, 1700, rs)
	stage3 := analysis.ExtractByFVAdvanced(1700, 1800, rs)
	stage4 := analysis.ExtractByFVAdvanced(1800, 1900, rs)
	stage5 := analysis.ExtractByFVAdvanced(1900, 2000, rs)
	stage6 := analysis.ExtractByFVAdvanced(2000, 2500, rs)

	msg += fmt.Sprintf("分段%s(%s)：%s (占比 / 胜率 / 场次，下同)\n", "1000-1500", tran2(stage1), tran(stage1))
	msg += fmt.Sprintf("分段%s(%s)：%s\n", "1500-1700", tran2(stage2), tran(stage2))
	msg += fmt.Sprintf("分段%s(%s)：%s\n", "1700-1800", tran2(stage3), tran(stage3))
	msg += fmt.Sprintf("分段%s(%s)：%s\n", "1800-1900", tran2(stage4), tran(stage4))
	msg += fmt.Sprintf("分段%s(%s)：%s\n", "1900-2000", tran2(stage5), tran(stage5))
	msg += fmt.Sprintf("分段%s(%s)：%s\n", "2000-2500", tran2(stage6), tran(stage6))

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

	allies, _, sum := analysis.TeamAnalysis(playerID)
	for i, ally := range allies {
		if i >= 7 {
			break
		}
		if analysis.Divide(ally[1], uint64(sum))*100 < 2.0 || ally[1] <= 2 {
			continue
		}
		collect.CrawlPlayerByID(ally[0])
		avg, than10min, total := analysis.ShuffleAnalysis(ally[0])
		msg += fmt.Sprintf("开黑队友 %s 有效间隔数：%d，平均间隔：%d秒，超过十分钟的间隔数：%d (%.1f%%)\n", collect.SearchName(ally[0]), total, avg, than10min, analysis.Divide(uint64(than10min), uint64(total))*100)
	}
	return
}

// ExportAssignHeroAnalysisAdvanced 分析fv团分上的玩家的英雄数据
func ExportAssignHeroAnalysisAdvanced(name string, hero string, fv int) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	if _, ok := db.HeroNameToID[hero]; !ok {
		return "未知英雄 " + name, nil
	}
	analysis.UpdateHeroOfPlayerRank(db.HeroNameToID[hero], fv)
	rs, total := analysis.HeroAnalysis(PlayerID, fv)
	_, rank, overallScore, total2 := analysis.GetHeroOfPlayerRank(db.HeroNameToID[hero], PlayerID, fv)
	msg += fmt.Sprintf("昵称：%s，总场次：%d\n", name, total)
	for i := range rs {
		if db.HeroIDToName[int(rs[i][0])] != hero {
			continue
		}
		msg += fmt.Sprintf("英雄：%s\n", hero)
		msg += fmt.Sprintf("有 %d 名玩家记录场次超过了 %d 次，团分下限：%d\n", total2, int(analysis.ValidTimes), fv)
		msg += fmt.Sprintf("实际场次：%d，参与计算场次：%d\n", uint64(rs[i][28]), uint64(rs[i][1]))
		msg += fmt.Sprintf("胜率：%.1f%% (超越%.1f%%的玩家，下同)\n", rs[i][2]/rs[i][1]*100, rank[2])
		msg += fmt.Sprintf("场均耗时：%.1f (%.1f%%) \n", rs[i][27]/60, rank[27])
		msg += fmt.Sprintf("场均补刀：%.1f (%.1f%%)\n", rs[i][3], rank[3])
		msg += fmt.Sprintf("场均每分补刀：%.2f (%.1f%%)\n", rs[i][4], rank[4])
		msg += fmt.Sprintf("场均kda：%.1f (%.1f%%) / %.1f (%.1f%%) / %.1f (%.1f%%)\n", rs[i][5], rank[5], rs[i][7], rank[7], rs[i][9], rank[9])
		msg += fmt.Sprintf("场均每分kda：%.2f (%.1f%%) / %.2f (%.1f%%) / %.2f (%.1f%%)\n", rs[i][6], rank[6], rs[i][8], rank[8], rs[i][10], rank[10])
		msg += fmt.Sprintf("场均推塔：%.1f (%.1f%%)\n", rs[i][11], rank[11])
		msg += fmt.Sprintf("场均每分推塔：%.2f (%.1f%%)\n", rs[i][12], rank[12])
		msg += fmt.Sprintf("场均插/排眼：%.2f (%.1f%%) / %.2f (%.1f%%)\n", rs[i][13], rank[13], rs[i][15], rank[15])
		msg += fmt.Sprintf("场均每分插/排眼：%.3f (%.1f%%) / %.3f (%.1f%%)\n", rs[i][14], rank[14], rs[i][16], rank[16])
		msg += fmt.Sprintf("场均经济：%.1f (%.1f%%)\n", rs[i][17], rank[17])
		msg += fmt.Sprintf("场均每分经济：%.1f (%.1f%%)\n", rs[i][18], rank[18])
		msg += fmt.Sprintf("场均经济占比：%.1f%% (%.1f%%)\n", rs[i][19]*100, rank[19])
		msg += fmt.Sprintf("场均输出：%.1f (%.1f%%)\n", rs[i][20], rank[20])
		msg += fmt.Sprintf("场均每分输出：%.1f (%.1f%%)\n", rs[i][21], rank[21])
		msg += fmt.Sprintf("场均输出占比：%.1f%% (%.1f%%)\n", rs[i][22]*100, rank[22])
		msg += fmt.Sprintf("场均承伤：%.1f (%.1f%%)\n", rs[i][23], rank[23])
		msg += fmt.Sprintf("场均每分承伤：%.1f (%.1f%%)\n", rs[i][24], rank[24])
		msg += fmt.Sprintf("场均承伤占比：%.1f%% (%.1f%%)\n", rs[i][25]*100, rank[25])
		msg += fmt.Sprintf("场均经济转换率：%.1f%% (%.1f%%)\n", rs[i][26], rank[26])
		msg += fmt.Sprintf("综合评分：%d (%.1f%%)\n", overallScore, rank[28])
	}
	return
}

func ExportLikeAnalysis(name string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	rs, total := analysis.HeroAnalysis(PlayerID, 0)
	msg += fmt.Sprintf("英雄分析，昵称：%s，总场次：%d\n", name, total)
	for i := range rs {
		if i >= 5 {
			break
		}
		analysis.UpdateHeroOfPlayerRank(int(rs[i][0]), 0)
		_, rank, overallScore, _ := analysis.GetHeroOfPlayerRank(int(rs[i][0]), PlayerID, 0)
		players := []db.Player{}
		db.SqlDB.Model(db.Player{}).Where("player_id = ? and hero_id = ?", PlayerID, rs[i][0]).Find(&players)
		win := 0
		for j := range players {
			if players[j].Result == 1 || players[j].Result == 3 {
				win++
			}
		}
		msg += fmt.Sprintf("%d、英雄：%s，场次：%d，胜率：%.1f%%，评分：%d(%.1f%%)\n", i+1, db.HeroIDToName[int(rs[i][0])], len(players), float64(win)/float64(len(players))*100, overallScore, rank[28])
	}
	return
}

func ExportGlobalHeroAnalysis(HeroName string, fv int) (msg string, err error) {
	ps, err := analysis.GlobalHeroAnalysis(HeroName)
	if err != nil {
		return
	}
	MatchIDToPlayers := map[string][]db.Player{}
	for i := range ps {
		MatchIDToPlayers[ps[i].MatchID] = append(MatchIDToPlayers[ps[i].MatchID], ps[i])
	}

	all := 0
	win := 0
	range0 := 0
	win0 := 0
	range1 := 0
	win1 := 0
	range2 := 0
	win2 := 0
	range3 := 0
	win3 := 0
	range4 := 0
	win4 := 0
	range5 := 0
	win5 := 0
	for _, players := range MatchIDToPlayers {
		// 双方都有该英雄
		if len(players) >= 2 {
			continue
		}
		// 玩家团分低于下限
		for i := range players {
			if players[i].FV < fv {
				continue
			}
		}
		// 场次+1
		all += 1
		for i := range players {
			tmp := 0
			// 赢了
			if players[i].Result == 1 || players[i].Result == 3 {
				win++
				tmp = 1
			}
			if players[i].FV < 1500 {
				range0++
				win0 += tmp
			} else if players[i].FV < 1700 && players[i].FV >= 1500 {
				range1++
				win1 += tmp
			} else if players[i].FV < 1800 && players[i].FV >= 1700 {
				range2++
				win2 += tmp
			} else if players[i].FV < 1900 && players[i].FV >= 1800 {
				range3++
				win3 += tmp
			} else if players[i].FV < 2000 && players[i].FV >= 1900 {
				range4++
				win4 += tmp
			} else if players[i].FV < 3000 && players[i].FV >= 2000 {
				range5++
				win5 += tmp
			}
		}
	}

	msg += fmt.Sprintf("英雄：%s，玩家团分下限：%d, 出现次数：%d\n", HeroName, fv, all)
	msg += fmt.Sprintf("全局单方面胜率：%.1f%%\n", analysis.Divide(uint64(win), uint64(all))*100)
	msg += fmt.Sprintf("分段%s(%.1f%%)：%d, %.1f%% (占比 / 场次 / 胜率，下同)\n", "1000-1500", analysis.Divide(uint64(range0), uint64(all))*100, range0, analysis.Divide(uint64(win0), uint64(range0))*100)
	msg += fmt.Sprintf("分段%s(%.1f%%)：%d, %.1f%%\n", "1500-1700", analysis.Divide(uint64(range1), uint64(all))*100, range1, analysis.Divide(uint64(win1), uint64(range1))*100)
	msg += fmt.Sprintf("分段%s(%.1f%%)：%d, %.1f%%\n", "1700-1800", analysis.Divide(uint64(range2), uint64(all))*100, range2, analysis.Divide(uint64(win2), uint64(range2))*100)
	msg += fmt.Sprintf("分段%s(%.1f%%)：%d, %.1f%%\n", "1800-1900", analysis.Divide(uint64(range3), uint64(all))*100, range3, analysis.Divide(uint64(win3), uint64(range3))*100)
	msg += fmt.Sprintf("分段%s(%.1f%%)：%d, %.1f%%\n", "1900-2000", analysis.Divide(uint64(range4), uint64(all))*100, range4, analysis.Divide(uint64(win4), uint64(range4))*100)
	msg += fmt.Sprintf("分段%s(%.1f%%)：%d, %.1f%%\n", "2000-2500", analysis.Divide(uint64(range5), uint64(all))*100, range5, analysis.Divide(uint64(win5), uint64(range5))*100)
	return
}

func ExportTopAnalysis(HeroName string, fv int) (msg string, err error) {
	if _, ok := db.HeroNameToID[HeroName]; !ok {
		return "", fmt.Errorf("不存在 %s 英雄", HeroName)
	}

	analysis.UpdateHeroOfPlayerRank(db.HeroNameToID[HeroName], fv)
	result, total, err := analysis.GetTopRank(db.HeroNameToID[HeroName], fv)
	if err != nil {
		return
	}

	msg += fmt.Sprintf("英雄：%s，玩家团分下限：%d，总计人数：%d\n", HeroName, fv, total)
	for i := range result {
		idStr := result[i].Member.(string)
		id, _ := strconv.ParseUint(idStr, 10, 64)
		msg += fmt.Sprintf("%d、昵称：%s，评分：%.1f\n", i+1, collect.SearchName(id), result[i].Score)
	}
	return
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
	return fmt.Sprintf("[CQ:image,file=file://%s.png]", abs), nil
}

func ExportPKAnalysis(name string, hero string) (msg string, err error) {
	PlayerID, err := collect.SearchRoleID(name)
	if err != nil {
		return
	}
	if _, ok := db.HeroNameToID[hero]; !ok {
		return "未知英雄 " + name, nil
	}
	you, top1 := analysis.PKAnalysis(PlayerID, db.HeroNameToID[hero])

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
