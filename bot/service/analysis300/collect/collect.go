package collect

import (
	"context"
	"eebot/bot/service/analysis300/db"
	"eebot/g"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var PlayerListKey = "300analysis:player_list"

var PlayerIDToNameKey = "300analysis:player_id_to_name"

var Expiration = 24 * time.Hour

var Ctx = context.Background()

// BeginCrawl 以某某玩家为中心爬取战绩
//   - 从redis队列中取出一个玩家id
//   - 查询其战绩、对比数据库中的战绩是否更新
//   - 将新战绩中的其他玩家存入redis队列，将新战绩中的本人战绩存入mysql
func BeginCrawl() {
	for {
		playerID, err := RDB.LPop(Ctx, PlayerListKey).Result()
		if err != nil {
			g.Logger.Info("error: redis list pop, ", err.Error())
			break
		}
		id, err := strconv.ParseUint(playerID, 10, 64)
		if err != nil {
			g.Logger.Info("error: parse player id, ", err.Error())
			break
		}
		beginCrawl(id)
		// if len(ids) > 0 {
		// 	_, err = RDB.RPush(ctx, PlayerListKey, ids...).Result()
		// 	if err != nil {
		// 		g.Logger.Info("error: redis list push, ", err.Error())
		// 	}
		// }
	}
}

func beginCrawl(PlayerID uint64) (ids []interface{}) {
	cnt := 0
	// 通过网页api查询战绩列表
	for searchIndex := 1; ; searchIndex++ {
		matches, err := SearchMatches(PlayerID, 1, searchIndex)
		if err != nil || len(matches) == 0 {
			break
		}
		for i, match := range matches {
			// 是否为新战绩
			var count int64
			db.SqlDB.Model(&db.Match{}).Where("match_id = ?", match.MatchID).Count(&count)
			// 不为新战绩
			if count > 0 {
				continue
			}
			// 查询比赛详情并保存
			cnt += 1
			matchInfo, err := SearchMatchInfo(match.MatchID)
			matchInfo.MatchID = match.MatchID
			if err != nil {
				g.Logger.Errorf("index %d, %d match %s pass: %s", searchIndex, i, match.MatchID, err.Error())
				continue
			}
			if len(matchInfo.Players) == 0 || matchInfo.MID != 254 {
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
			db.SqlDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&matchInfo)
			
		}
	}
	return
}

func CrawlPlayerByName(name string) (err error) {
	playerID, err := SearchRoleID(name)
	if err != nil {
		g.Logger.Info("error: search player id, ", err.Error())
		return err
	}
	beginCrawl(playerID)
	return
}

func CrawlPlayerByID(PlayerID uint64) {
	beginCrawl(PlayerID)
}
