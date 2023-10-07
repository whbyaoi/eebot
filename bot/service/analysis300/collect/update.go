package collect

import (
	"eebot/bot/service/analysis300/db"
	"eebot/g"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// UpdateMatchAndPlayer 以某某玩家为中心爬取战绩
//   - 将数据库中所有玩家的id存入redis
//   - 从redis队列中取出一个玩家id
//   - 查询其战绩、对比数据库中的战绩是否更新
//   - 将新战绩中所有玩家战绩存入mysql
func UpdateMatchAndPlayer() {
	for {
		playerID, err := RDB.LPop(Ctx, PlayerListKey).Result()
		if err != nil {
			g.Logger.Errorln("error: redis list pop, ", err.Error())
			break
		}
		id, err := strconv.ParseUint(playerID, 10, 64)
		if err != nil {
			g.Logger.Errorln("error: parse player id, ", err.Error())
			break
		}
		updateCrawl(id)
	}
	g.Logger.Info("func UpdateMatchAndPlayer() over")
}

func updateCrawl(PlayerID uint64) (ids []interface{}) {
	t0 := time.Now()
	// 通过网页api查询战绩列表
	total := 0
	for searchIndex := 1; ; searchIndex++ {
		cnt := 0
		matches, err := SearchMatches(PlayerID, 1, searchIndex)
		if err != nil || len(matches) == 0 {
			break
		}
		var saveMatches []db.Match
		for i, match := range matches {
			// 是否为新战绩
			var count int64
			db.SqlDB.Model(&db.Match{}).Where("match_id = ?", match.MatchID).Count(&count)
			// 不为新战绩
			if count > 0 {
				continue
			}
			// 查询比赛详情并保存
			matchInfo, err := SearchMatchInfo(match.MatchID)
			matchInfo.MatchID = match.MatchID
			if err != nil {
				log.Printf("index %d, %d match %s pass: %s", searchIndex, i, match.MatchID, err.Error())
				continue
			}
			if len(matchInfo.Players) < 14 && matchInfo.MID != 254 {
				continue
			}
			var totalMoney1 int
			var MakeDamage1 int
			var TakeDamage1 int
			var totalMoney2 int
			var MakeDamage2 int
			var TakeDamage2 int
			for j := range matchInfo.Players {
				// 赋值游戏时间
				matchInfo.Players[j].MatchID = matchInfo.MatchID
				matchInfo.Players[j].UsedTime = matchInfo.UsedTime

				if matchInfo.Players[j].Side == 1 {
					totalMoney1 += matchInfo.Players[j].TotalMoney
					MakeDamage1 += matchInfo.Players[j].MD[len(matchInfo.Players[j].MD)-1]
					TakeDamage1 += matchInfo.Players[j].TD[len(matchInfo.Players[j].TD)-1]
				} else {
					totalMoney2 += matchInfo.Players[j].TotalMoney
					MakeDamage2 += matchInfo.Players[j].MD[len(matchInfo.Players[j].MD)-1]
					TakeDamage2 += matchInfo.Players[j].TD[len(matchInfo.Players[j].TD)-1]
				}
				if matchInfo.Players[j].PlayerID != PlayerID {
					ids = append(ids, matchInfo.Players[j].PlayerID)
				} else {
					if searchIndex == 1 && i == 0 {
						// 缓存id-name
						go func(idx int) {
							RDB.Set(Ctx, fmt.Sprintf("%s_%d", PlayerIDToNameKey, matchInfo.Players[idx].PlayerID), matchInfo.Players[idx].Name, Expiration)
						}(j)
					}
				}
			}
			// 自行计算缺省值
			for j := range matchInfo.Players {
				if matchInfo.Players[j].Side == 1 {
					matchInfo.Players[j].TotalMoneySide = totalMoney1
					matchInfo.Players[j].TotalMoneyPercent = float64(matchInfo.Players[j].TotalMoney) / float64(totalMoney1)
					matchInfo.Players[j].MakeDamageSide = MakeDamage1
					matchInfo.Players[j].MakeDamagePercent = float64(matchInfo.Players[j].MD[len(matchInfo.Players[j].MD)-1]) / float64(MakeDamage1)
					matchInfo.Players[j].TakeDamageSide = TakeDamage1
					matchInfo.Players[j].TakeDamagePercent = float64(matchInfo.Players[j].TD[len(matchInfo.Players[j].TD)-1]) / float64(TakeDamage1)
				} else {
					matchInfo.Players[j].TotalMoneySide = totalMoney2
					matchInfo.Players[j].TotalMoneyPercent = float64(matchInfo.Players[j].TotalMoney) / float64(totalMoney2)
					matchInfo.Players[j].MakeDamageSide = MakeDamage2
					matchInfo.Players[j].MakeDamagePercent = float64(matchInfo.Players[j].MD[len(matchInfo.Players[j].MD)-1]) / float64(MakeDamage2)
					matchInfo.Players[j].TakeDamageSide = TakeDamage2
					matchInfo.Players[j].TakeDamagePercent = float64(matchInfo.Players[j].TD[len(matchInfo.Players[j].TD)-1]) / float64(TakeDamage2)
				}
			}
			// 保存比赛和玩家记录
			saveMatches = append(saveMatches, matchInfo)
			cnt += 1
		}
		total += cnt
		db.SqlDB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&saveMatches)
		if cnt == 0 {
			log.Printf("player %s, no new matches", SearchName(PlayerID))
			return nil
		}
	}
	log.Printf("player %s over, %d total new matches, %vs used", SearchName(PlayerID), total, time.Since(t0))
	return
}
