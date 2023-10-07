package controller

import (
	"eebot/bot/service"
	"eebot/bot/service/analysis300"
	"eebot/bot/service/analysis300/collect"
	"eebot/g"
	"errors"
	"fmt"
	"strconv"
)

// AnalysisHub
//
//	rawMessageSlice[0]: 300
//	rawMessageSlice[1]: 指令缩写
//	rawMessageSlice[2:]: 参数(顺序:昵称 英雄名 团分下限)
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

	if svc != "h" && svc != "help" && svc != "帮助" && svc != "gh" {
		go service.Reply("别急，查询角色中", prefix, targetID)
		err = collect.CrawlPlayerByName(name)
		if err != nil {
			err = service.Reply(err.Error(), prefix, targetID)
			return
		}
	}
	switch svc {
	case "t": // 开黑
		suffix, err = analysis300.ExportTeamAnalysisAdvanced(name)
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
				break
			}
		}
		if assgin == "" {
			err = errors.New("该指令必须指定英雄")
		} else {
			suffix, err = analysis300.ExportAssignHeroAnalysisAdvanced(name, assgin, fv)
		}
	case "gh": // 全局英雄
		suffix, err = analysis300.ExportGlobalHeroAnalysis(name, 0)
	case "l": // 常用
		suffix, err = analysis300.ExportLikeAnalysis(name)
	case "all": // 全部
		if !HasAuth(sourceID) {
			err = errors.New("无权使用该指令")
		}
		tmp, err := analysis300.ExportTeamAnalysisAdvanced(name)
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
	case "help": // 帮助
		suffix += "所有指令一览：\n"
		suffix += "n --- 胜负分析\n"
		suffix += "t --- 开黑分析\n"
		suffix += "s  --- 洗牌分析\n"
		suffix += "l  --- 常用分析\n"
		suffix += "h 英雄名称 [可选]团分下限 - 英雄分析\n"
		suffix += "gh 英雄名称 - 全局英雄分析\n"
		suffix += "本项目github: github.com/whbyaoi/eebot"
	default:
		suffix = "未知指令"
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
