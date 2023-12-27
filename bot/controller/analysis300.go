package controller

import (
	"eebot/bot/service"
	"eebot/bot/service/analysis300"
	"eebot/bot/service/analysis300/analysis"
	"eebot/bot/service/analysis300/collect"
	"eebot/g"
	"errors"
	"fmt"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
)

var NoWait = []string{"help", "菜单", "g", "g1", "f", "g2", "win", "top", "topa", "active", "flush", "", "test"}

var mutexes map[string]*sync.Mutex = map[string]*sync.Mutex{}

// AnalysisHub
//
//	rawMessageSlice[0]: 300
//	rawMessageSlice[1]: 指令缩写
//	rawMessageSlice[2:]: 参数(顺序:昵称 英雄名 竞技力下限)
func AnalysisHub(rawMessageSlice []string, isGroup bool, sourceID int64, targetID int64) (err error) {
	var svc string
	var name string
	if len(rawMessageSlice) > 1 {
		svc = rawMessageSlice[1]
	}
	if len(rawMessageSlice) > 2 {
		name = rawMessageSlice[2]
	}

	var suffix string
	var prefix string
	if isGroup {
		prefix = fmt.Sprintf("[CQ:at,qq=%d] \n", sourceID)
	}
	defer func() {
		if r := recover(); r != nil {
			service.Reply("未知错误", prefix, targetID)
			err = fmt.Errorf("%+v", r)
			debug.PrintStack()
		}
	}()

	if _, ok := mutexes[strings.Join(rawMessageSlice, "")]; !ok {
		mutexes[strings.Join(rawMessageSlice, "")] = new(sync.Mutex)
	}
	excceed := mutexes[strings.Join(rawMessageSlice, "")].TryLock()
	if !excceed {
		service.Reply("已有重复的查询正在进行", prefix, targetID)
		return
	}
	defer mutexes[strings.Join(rawMessageSlice, "")].Unlock()

	if !slices.Contains[[]string](NoWait, svc) {
		go service.Reply("别急，查询角色中(第一次查询会较慢)", prefix, targetID)
		err = collect.CrawlPlayerByName(name)
		if err != nil {
			err = service.Reply(err.Error(), prefix, targetID)
			return
		}
	}
	switch svc {
	case "r": // 最近
		assgin := ""
		if len(rawMessageSlice) > 3 {
			assgin = rawMessageSlice[3]
		}
		suffix, err = analysis300.ExportRelatedAnalysis(name, assgin)
	case "t": // 开黑
		suffix, err = analysis300.ExportTeamAnalysis(name)
	case "n": // 常规
		suffix, err = analysis300.ExportWinOrLoseAnalysisAdvanced(name)
	case "s": // 洗牌
		suffix, err = analysis300.ExportShuffleAnalysis(name)
	case "as": // 进阶洗牌
		if HasAuth(sourceID) {
			suffix, err = analysis300.ExportShuffleAnalysisAdvanced(name)
		} else {
			err = errors.New("无权使用该命令")
		}
	case "h": // 英雄
		assgin := ""
		if len(rawMessageSlice) > 3 {
			assgin = rawMessageSlice[3]
		}
		var fv int
		if len(rawMessageSlice) > 4 {
			fv, err = strconv.Atoi(rawMessageSlice[4])
			if err != nil {
				fv = 0
			}
		}
		if assgin == "" {
			err = errors.New("该指令必须指定英雄")
		} else {
			suffix, err = analysis300.ExportAssignHeroAnalysisAdvancedV2(name, assgin, fv)
		}
	case "f":
		var page int64 = 1
		if len(rawMessageSlice) > 3 {
			page, err = strconv.ParseInt(rawMessageSlice[3], 10, 64)
			if err != nil || page <= 0 {
				err = errors.New("错误页码")
			}
		}
		suffix, err = analysis300.ExportFindPlayer(name, int(page))
	case "g", "g1": // 全局英雄
		if name == "" {
			err = errors.New("该指令必须指定英雄")
		} else {
			suffix, err = analysis300.ExportGlobalHeroAnalysis(name)
		}
	case "g2":
		if name == "" {
			err = errors.New("该指令必须指定英雄")
		} else {
			suffix, err = analysis300.ExportGlobalHeroAnalysis2(name)
		}
	case "l": // 常用
		suffix, err = analysis300.ExportLikeAnalysis(name)
	case "jjl": // 竞技力
		suffix, err = analysis300.ExportJJLWithTeamAnalysis(name)
	case "jjl2":
		suffix, err = analysis300.ExportJJLCompositionAnalysis(name)
	case "pk":
		assgin := ""
		if len(rawMessageSlice) > 3 {
			assgin = rawMessageSlice[3]
		}
		if assgin == "" {
			err = errors.New("该指令必须指定英雄")
		} else {
			suffix, err = analysis300.ExportPKAnalysis(name, assgin)
		}
	case "active":
		suffix, err = analysis300.ExportActiveAnalysis()
	case "flush":
		if name == "" {
			err = errors.New("该指令必须指定英雄")
		} else {
			suffix, err = analysis300.ExportFlushTop(name)
		}
	case "top": // top10
		var fv int
		if len(rawMessageSlice) > 3 {
			fv, err = strconv.Atoi(rawMessageSlice[3])
			if err != nil {
				fv = 0
			}
		}
		if name == "" {
			err = errors.New("该指令必须包含参数")
		} else {
			suffix, err = analysis300.ExportTopAnalysis(name, fv)
		}
	case "win":
		scope := ""
		if len(rawMessageSlice) > 3 {
			for _, tmp := range analysis.DefaultJJLCategoryKeys {
				if tmp == rawMessageSlice[3] {
					scope = tmp
				}
			}
			if scope == "" {
				suffix += fmt.Sprintf("未知竞技力范围：%s\n支持的竞技力范围：", rawMessageSlice[3])
				for _, tmp := range analysis.DefaultJJLCategoryKeys {
					suffix += tmp + "\n"
				}
				break
			}
		}
		suffix, err = analysis300.ExportWinRateAnalysis(name, scope)
	case "topa":
		var fv int
		if len(rawMessageSlice) > 3 {
			fv, err = strconv.Atoi(rawMessageSlice[3])
			if err != nil {
				fv = 0
			}
		}
		if name == "" {
			err = errors.New("该指令必须指定英雄")
		} else {
			suffix, err = analysis300.ExportTopWithDetailAnalysis(name, fv)
		}
	case "all": // 全部
		if !HasAuth(sourceID) {
			err = errors.New("无权使用该指令")
			break
		}
		tmp, err := analysis300.ExportTeamAnalysis(name)
		if err != nil {
			break
		}
		suffix += tmp
		tmp, err = analysis300.ExportWinOrLoseAnalysisAdvanced(name)
		if err != nil {
			break
		}
		suffix += tmp
		tmp, err = analysis300.ExportShuffleAnalysis(name)
		if err != nil {
			break
		}
		suffix += tmp
	case "help", "菜单": // 帮助
		suffix += "指令一览(任何指令需要前缀300，详见个人空间或github项目首页)：\n"
		suffix += "n 玩家 --- 胜负分析\n"
		suffix += "t 玩家 --- 开黑分析\n"
		suffix += "s 玩家 --- 洗牌分析\n"
		suffix += "l 玩家 --- 常用分析\n"
		suffix += "jjl 玩家 - 竞技力与开黑分析\n"
		suffix += "jjl2 玩家 - 竞技力成分分析\n"
		suffix += "pk 玩家 英雄名称 - 与榜一比较\n"
		suffix += "r 玩家 [可选]英雄名称 - 近10场jjc数据\n"
		suffix += "h 玩家 英雄名称 [可选]竞技力下限 - 英雄分析\n"
		suffix += "f 英雄/k/d/a/补刀/经济/竞技力 - 战绩找人，其中英雄/k/d/a为必填，其余选填可空\n"
		suffix += "g 英雄名称 - 全局英雄分析(使用者各分段的出场及胜率情况)\n"
		suffix += "g2 英雄名称 - 全局英雄分析(场均各分段的出场及胜率情况)\n"
		suffix += "top 英雄名称 [可选]竞技力下限 - 月榜前10\n"
		suffix += "top jjl[可选]页码 - jjl月榜\n"
		suffix += "win 页码 竞技力范围 - 英雄胜率榜\n"
		suffix += "topa 英雄名称 [可选]竞技力下限 - 月榜前10(附带计算详情)\n"
		suffix += "flush 英雄名称 - 刷新月榜"
	case "test":
		if !HasAuth(sourceID) {
			err = errors.New("未知指令：" + svc + "\n查看可用指令请“有效”@机器人 + 300 + help\n如：@男神 300 help")
			break
		}
		id, _ := collect.SearchRoleID("晚约")
		analysis.StableJJLLAnalysis(id)
	default:
		suffix = "未知指令：" + svc + "\n查看可用指令请“有效”@机器人 + 300 + help\n如：@男神 300 help"
	}

	if err == nil {
		err = service.Reply(suffix, prefix, targetID)
	} else {
		err = service.Reply(err.Error(), prefix, targetID)
	}
	return
}

func HasAuth(sourceID int64) bool {
	return sourceID == int64(g.Config.GetInt("analysis.vip"))
}
