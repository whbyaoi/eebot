package analysis

import (
	"eebot/bot/service/analysis300/collect"
	"eebot/bot/service/analysis300/db"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ValidTimes = 5.0

var ValidIntervalTimes = 100

// 索引 - 英雄数据水平
var HeroDataToName = map[int]string{
	1:  "total",
	2:  "win",
	3:  "avg_last_hit",
	4:  "avg_last_hit_per_minute",
	5:  "avg_kill",
	6:  "avg_kill_per_minute",
	7:  "avg_death",
	8:  "avg_death_per_minute",
	9:  "avg_assist",
	10: "avg_assist_per_minute",
	11: "avg_tower",
	12: "avg_tower_per_minute",
	13: "avg_put_eye",
	14: "avg_put_eye_per_minute",
	15: "avg_destroy_eye",
	16: "avg_destroy_eye_per_minute",
	17: "avg_money",
	18: "avg_money_per_minute",
	19: "avg_money_percent",
	20: "avg_make_damage",
	21: "avg_make_damage_per_minute",
	22: "avg_make_damage_percent",
	23: "avg_take_damage",
	24: "avg_take_damage_per_minute",
	25: "avg_take_damage_percent",
	26: "avg_money_conversion_rate",
	27: "avg_used_time",
	28: "overall_score",
}

// 存储按照玩家类别的玩家水平的zset key值
var HeroOfPlayerRankKey = "300analysis_hero_player_rank"

// 存储时间间隔水平zset key值
var MatchIntervalKey = "300analysis_shuffle"

// UpdateHeroOfPlayerRank 按照玩家类别更新某个英雄数据水平
func UpdateHeroOfPlayerRank(HeroID int, fv int) {
	var players []db.Player
	db.SqlDB.Model(db.Player{}).Where("hero_id = ?", HeroID).Find(&players)

	idToRecord := map[uint64][]db.Player{}
	for i := range players {
		idToRecord[players[i].PlayerID] = append(idToRecord[players[i].PlayerID], players[i])
	}
	data := getHeroOfPlayerData(idToRecord, fv)

	prefix := HeroOfPlayerRankKey + fmt.Sprintf("_%s_%d:", db.HeroIDToName[HeroID], fv)
	for id, detail := range data {
		for k, score := range detail {
			key := prefix + k
			collect.RDB.ZAdd(collect.Ctx, key, redis.Z{Score: score, Member: id})
		}
	}

	// 计算综合水平
	factors := MergeImportance(HeroNameToID[db.HeroIDToName[HeroID]])
	for id := range data {
		rank, _ := GetHeroOfPlayerRankWithoutOverallScore(HeroID, id, fv)
		overallScore := 0.0
		for i, factor := range factors {
			overallScore += rank[i] * factor
		}
		overallScore = overallScore * (0.9 + min(data[id]["total"]-ValidTimes, 10)/10*0.1) * (0.9 + rank[1]/100*0.1)
		key := prefix + HeroDataToName[28]
		collect.RDB.ZAdd(collect.Ctx, key, redis.Z{Score: overallScore, Member: id})
	}
}

// InitMatchInterval 初始化比赛间隔水平(慎用)
func InitMatchInterval() {
	var players []db.Player
	sub := db.SqlDB.Model(db.Player{}).Select("player_id").Group("player_id").Having("count(*) > ?", ValidIntervalTimes)
	db.SqlDB.Model(db.Player{}).Distinct("player_id").Where("player_id in (?)", sub).Find(&players)

	prefix := fmt.Sprintf("%s:", MatchIntervalKey)
	for i := range players {
		avg, than10min, total := ShuffleAnalysis(players[i].PlayerID)

		key := prefix + "avg"
		collect.RDB.ZAdd(collect.Ctx, key, redis.Z{Score: float64(avg), Member: players[i].PlayerID})

		key = prefix + "than_10_min"
		collect.RDB.ZAdd(collect.Ctx, key, redis.Z{Score: float64(than10min) / float64(total), Member: players[i].PlayerID})
	}
}

// InitMatchInterval 更新比赛间隔水平
func UpdateMatchInterval(PlayerID uint64) {

	prefix := fmt.Sprintf("%s:", MatchIntervalKey)
	avg, than10min, total := ShuffleAnalysis(PlayerID)

	key := prefix + "avg"
	collect.RDB.ZAdd(collect.Ctx, key, redis.Z{Score: float64(avg), Member: PlayerID})

	key = prefix + "than_10_min"
	collect.RDB.ZAdd(collect.Ctx, key, redis.Z{Score: float64(than10min) / float64(total), Member: PlayerID})
}

// getHeroOfPlayerData 计算玩家的英雄数据水平
//
//	data: playerID-HeroDataName-value
func getHeroOfPlayerData(raw map[uint64][]db.Player, fv int) (data map[uint64]map[string]float64) {

	addValue := func(v [28]float64, me db.Player) [28]float64 {
		// 过滤低于阈值的战绩
		if me.FV < fv {
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

	data = map[uint64]map[string]float64{}
	for playerID, heroData := range raw {
		tmp := [28]float64{}
		for _, p := range heroData {
			tmp = addValue(tmp, p)
		}
		// 过滤总场次不足 ValidTimes 的数据
		if tmp[1] < ValidTimes {
			continue
		}
		tmp[2] = tmp[2] / tmp[1]
		for i := 3; i <= 27; i++ {
			tmp[i] = tmp[i] / tmp[1]
		}
		data[playerID] = map[string]float64{}
		for index, name := range HeroDataToName {
			if index == 28 {
				continue
			}
			data[playerID][name] = tmp[index]
		}

	}
	return
}

// GetHeroOfPlayerRank 获取某位玩家的英雄水平(包含综合评分)
func GetHeroOfPlayerRank(HeroID int, PlayerID uint64, fv int) (rank [29]float64, overallScore uint64, total int64) {
	prefix := HeroOfPlayerRankKey + fmt.Sprintf("_%s_%d:", db.HeroIDToName[HeroID], fv)
	for index, name := range HeroDataToName {
		key := prefix + name
		total, _ = collect.RDB.ZCard(collect.Ctx, key).Result()
		pos, _ := collect.RDB.ZRank(collect.Ctx, key, fmt.Sprintf("%d", PlayerID)).Result()
		rank[index] = float64(pos) / float64(total-1) * 100
		if index == 28 {
			tmp, _ := collect.RDB.ZScore(collect.Ctx, key, fmt.Sprintf("%d", PlayerID)).Result()
			overallScore = uint64(tmp)
		}
	}
	return
}

// GetHeroOfPlayerRankWithoutOverallScore 获取某位玩家的英雄水平(不包含综合评分)
func GetHeroOfPlayerRankWithoutOverallScore(HeroID int, PlayerID uint64, fv int) (rank [28]float64, total int64) {
	prefix := HeroOfPlayerRankKey + fmt.Sprintf("_%s_%d:", db.HeroIDToName[HeroID], fv)
	for index, name := range HeroDataToName {
		if index == 28 {
			continue
		}
		key := prefix + name
		total, _ = collect.RDB.ZCard(collect.Ctx, key).Result()
		pos, _ := collect.RDB.ZRank(collect.Ctx, key, fmt.Sprintf("%d", PlayerID)).Result()
		rank[index] = float64(pos) / float64(total-1) * 100
	}
	return
}

// GetTopRank 获取某英雄综合评分前10
func GetTopRank(HeroID int, fv int) (result []redis.Z, total int64, err error) {
	prefix := HeroOfPlayerRankKey + fmt.Sprintf("_%s_%d:", db.HeroIDToName[HeroID], fv)
	key := prefix + HeroDataToName[28]
	result, err = collect.RDB.ZRevRangeWithScores(collect.Ctx, key, 0, 9).Result()
	total, _ = collect.RDB.ZCard(collect.Ctx, key).Result()
	return
}

// GetMatchInterval 获取某位玩家的更新比赛间隔水平
func GetMatchInterval(PlayerID uint64) (rank [2]float64, total [2]int64) {
	UpdateMatchInterval(PlayerID)
	prefix := fmt.Sprintf("%s:", MatchIntervalKey)

	key := prefix + "avg"
	total[0], _ = collect.RDB.ZCard(collect.Ctx, key).Result()
	pos, _ := collect.RDB.ZRank(collect.Ctx, key, fmt.Sprintf("%d", PlayerID)).Result()
	rank[0] = float64(pos) / float64(total[0]-1) * 100

	key = prefix + "than_10_min"
	total[1], _ = collect.RDB.ZCard(collect.Ctx, key).Result()
	pos, _ = collect.RDB.ZRank(collect.Ctx, key, fmt.Sprintf("%d", PlayerID)).Result()
	rank[1] = float64(pos) / float64(total[1]-1) * 100

	return
}
