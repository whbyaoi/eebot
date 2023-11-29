package collect

import (
	"eebot/bot/service/analysis300/db"
	"eebot/g"
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

func InitCrawler() {
	Crawler.Auto = make(chan Request)
	Crawler.User = make(chan Request)
	go Crawler.Run()
	if g.Config.GetBool("auto-collect.run") {
		go Crawler.AutoIncrementalCrawl()
	}
}

var Crawler = new(crawler)

// 返回数据结构体
type Response struct {
	Data *db.Match
	Err  error
}

// 请求结构体
type Request struct {
	MatchID string

	RespChan chan Response
}

type crawler struct {
	Auto chan Request

	User chan Request
}

func (c *crawler) Run() {
	for {
		// 先处理User的战绩爬取请求，再处理Auto爬取请求
		select {
		case req := <-c.User:
			data, err := c.getMatchDetail(req.MatchID)
			req.RespChan <- Response{Data: data, Err: err}
		case autoReq := <-c.Auto:
			if g.Config.GetBool("auto-collect.run") {
			Auto:
				// 检查有没有User请求，如果有，先处理干净User请求
				for {
					select {
					case req := <-c.Auto:
						data, err := c.getMatchDetail(req.MatchID)
						req.RespChan <- Response{Data: data, Err: err}
					default:
						break Auto
					}
				}
				// 最终执行Auto请求
				data, err := c.getMatchDetail(autoReq.MatchID)
				autoReq.RespChan <- Response{Data: data, Err: err}
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (c *crawler) getMatchDetail(MatchID string) (match *db.Match, err error) {
	// 是否为新战绩
	var count int64
	db.SqlDB.Model(&db.Match{}).Where("match_id = ?", MatchID).Count(&count)
	if count > 0 {
		return nil, nil
	}
	// 查询比赛详情并保存
	tmp, err := SearchMatchInfo(MatchID)
	if err != nil {
		return nil, err
	}
	if len(tmp.Players) < 14 || tmp.MID != 254 {
		return nil, errors.New("wrong type of match detail")
	}
	match = &tmp
	match.MatchID = MatchID
	var totalMoney1 int
	var MakeDamage1 int
	var TakeDamage1 int
	var totalMoney2 int
	var MakeDamage2 int
	var TakeDamage2 int
	for j := range match.Players {
		// 赋值游戏时间
		match.Players[j].MatchID = match.MatchID
		match.Players[j].UsedTime = match.UsedTime
		match.Players[j].CreateTime = match.CreateTime

		if match.Players[j].Side == 1 {
			totalMoney1 += match.Players[j].TotalMoney
			MakeDamage1 += match.Players[j].MD[len(match.Players[j].MD)-1]
			TakeDamage1 += match.Players[j].TD[len(match.Players[j].TD)-1]
		} else {
			totalMoney2 += match.Players[j].TotalMoney
			MakeDamage2 += match.Players[j].MD[len(match.Players[j].MD)-1]
			TakeDamage2 += match.Players[j].TD[len(match.Players[j].TD)-1]
		}
	}
	// 自行计算缺省值
	for j := range match.Players {
		if match.Players[j].Side == 1 {
			match.Players[j].TotalMoneySide = totalMoney1
			match.Players[j].TotalMoneyPercent = float64(match.Players[j].TotalMoney) / float64(totalMoney1)
			match.Players[j].MakeDamageSide = MakeDamage1
			if MakeDamage1 != 0 {
				match.Players[j].MakeDamagePercent = float64(match.Players[j].MD[len(match.Players[j].MD)-1]) / float64(MakeDamage1)
			} else {
				match.Players[j].MakeDamagePercent = 0
			}
			match.Players[j].TakeDamageSide = TakeDamage1
			if TakeDamage1 != 0 {
				match.Players[j].TakeDamagePercent = float64(match.Players[j].TD[len(match.Players[j].TD)-1]) / float64(TakeDamage1)
			} else {
				match.Players[j].TakeDamagePercent = 0
			}
		} else {
			match.Players[j].TotalMoneySide = totalMoney2
			match.Players[j].TotalMoneyPercent = float64(match.Players[j].TotalMoney) / float64(totalMoney2)
			match.Players[j].MakeDamageSide = MakeDamage2
			if MakeDamage2 != 0 {
				match.Players[j].MakeDamagePercent = float64(match.Players[j].MD[len(match.Players[j].MD)-1]) / float64(MakeDamage2)
			} else {
				match.Players[j].MakeDamagePercent = 0
			}
			match.Players[j].TakeDamageSide = TakeDamage2
			if TakeDamage2 != 0 {
				match.Players[j].TakeDamagePercent = float64(match.Players[j].TD[len(match.Players[j].TD)-1]) / float64(TakeDamage2)
			} else {
				match.Players[j].TakeDamagePercent = 0
			}
		}
	}

	return
}

// CrawlAllAndSave 爬取所有战绩详情并保存，返回涉及的其他人id
//
//	source: 0-用户，1-自动爬取
func (c *crawler) CrawlAllAndSave(PlayerID uint64, source int) (ids []interface{}) {
	// 通过网页api查询战绩列表
	name := SearchName(PlayerID)
	total := 0
	t0 := time.Now()
	idMap := map[uint64]struct{}{}
	defer func() {
		if r := recover(); r != nil {
			g.CrawlLogger.Errorf("爬取玩家 %s(%d) 战绩时发生错误，%s", name, PlayerID, string(debug.Stack()))
		}
	}()
	g.CrawlLogger.Infof("开始爬取玩家 %s(%d) 战绩", name, PlayerID)
	for page := 1; ; page++ {
		matches, err := SearchMatches(PlayerID, 1, page)
		if err != nil || len(matches) == 0 {
			g.CrawlLogger.Infof("页面 %d，玩家 %s(%d) 无新战绩", page, name, PlayerID)
			break
		}
		// 缓存最新id-name
		if page == 1 {
			db.RDB.Set(Ctx, fmt.Sprintf("%s_%d", PlayerIDToNameKey, PlayerID), matches[0].Players[0].Name, Expiration)
			name = matches[0].Players[0].Name
		}
		var saveMatches []db.Match
		var wg sync.WaitGroup
		wg.Add(len(matches))
		for i := range matches {
			go func(i int) {
				defer wg.Done()
				req := Request{
					MatchID:  matches[i].MatchID,
					RespChan: make(chan Response),
				}
				// 发送请求
				if source == 0 {
					c.User <- req
				} else {
					c.Auto <- req
				}
				resp := <-req.RespChan
				if resp.Err != nil {
					g.CrawlLogger.Errorf("爬取 %s(%d) 玩家 %s 战绩时错误：%s", name, PlayerID, matches[i].MatchID, resp.Err.Error())
					return
				}
				if resp.Data != nil {
					// 保存比赛和玩家记录
					saveMatches = append(saveMatches, *resp.Data)
				}
			}(i)
		}
		wg.Wait()
		if len(saveMatches) == 0 {
			g.CrawlLogger.Infof("页面 %d，玩家 %s(%d) 无新战绩", page, name, PlayerID)
			break
		}
		// 保存战绩
		err = db.SqlDB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&saveMatches).Error
		if err != nil {
			g.CrawlLogger.Errorf("保存玩家 %s(%d) 战绩时错误：%s", name, PlayerID, err.Error())
			continue
		}
		// 保存分区战绩
		partitionPlayers := []db.PlayerPartition{}
		for i := range saveMatches {
			for _, play := range saveMatches[i].Players {
				partitionPlayers = append(partitionPlayers, db.ToPartition(play))
			}
		}
		err = db.SqlDB.Create(&partitionPlayers).Error
		if err != nil {
			g.CrawlLogger.Errorf("保存玩家 %s(%d) 分区战绩时错误：%s", name, PlayerID, err.Error())
			continue
		}
		total += len(saveMatches)
		for i := range saveMatches {
			for j := range saveMatches[i].Players {
				if saveMatches[i].Players[j].PlayerID == PlayerID {
					continue
				}
				idMap[saveMatches[i].Players[j].PlayerID] = struct{}{}
			}
		}
	}
	g.CrawlLogger.Infof("玩家 %s(%d) 战绩查询完毕，总计 %d 条新战绩，用时 %v", name, PlayerID, total, time.Since(t0))
	for id := range idMap {
		ids = append(ids, id)
	}
	return
}

// AutoIncrementalCrawl 增量更新
func (c *crawler) AutoIncrementalCrawl() {
	for {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), g.Config.GetInt("auto-collect.start"), 0, 0, 0, time.Local)
		end := time.Date(now.Year(), now.Month(), now.Day(), g.Config.GetInt("auto-collect.end"), 0, 0, 0, time.Local)
		if !(now.After(start) && now.Before(end)) {
			var wait time.Duration
			if now.After(start) {
				wait = start.Add((24*60*60 + 10) * time.Second).Sub(now)
			} else {
				wait = start.Sub(now)
			}
			g.CrawlLogger.Infof("增量更新停止: 当前时间%s不在更新时间段%s-%s内，等待%v", now.Format(time.TimeOnly), start.Format(time.TimeOnly), end.Format(time.TimeOnly), wait)
			time.Sleep(wait)
			g.CrawlLogger.Infof("增量更新开始")
		}
		playerID, err := db.RDB.SPop(Ctx, PlayerKeySet).Result()
		if err != nil {
			g.CrawlLogger.Error("增量更新错误: redis set pop, ", err.Error())
			break
		}
		id, err := strconv.ParseUint(playerID, 10, 64)
		if err != nil {
			g.CrawlLogger.Error("增量更新错误: wrong type of player id, ", err.Error())
			break
		}
		ids := c.CrawlAllAndSave(id, 1)
		if len(ids) > 0 {
			_, err = db.RDB.SAdd(Ctx, PlayerKeySet, ids...).Result()
			if err != nil {
				g.CrawlLogger.Error("增量更新错误: redis set add, ", err.Error())
			}
		}
	}
}

func (c *crawler) ManualIncrementalCrawl() {
	for {
		g.CrawlLogger.Infof("增量更新开始")
		playerID, err := db.RDB.SPop(Ctx, PlayerKeySet).Result()
		if err != nil {
			g.CrawlLogger.Error("增量更新错误: redis set pop, ", err.Error())
			break
		}
		id, err := strconv.ParseUint(playerID, 10, 64)
		if err != nil {
			g.CrawlLogger.Error("增量更新错误: wrong type of player id, ", err.Error())
			break
		}
		ids := c.CrawlAllAndSave(id, 1)
		if len(ids) > 0 {
			_, err = db.RDB.SAdd(Ctx, PlayerKeySet, ids...).Result()
			if err != nil {
				g.CrawlLogger.Error("增量更新错误: redis set add, ", err.Error())
			}
		}
	}
}

func (c *crawler) UpdatePlayerSet() {
	type result struct {
		PlayerID string
	}
	var results []result
	db.SqlDB.Model(db.Player{}).Distinct("player_id").Select("player_id").Scan(&results)
	ids := []interface{}{}
	for i := range results {
		ids = append(ids, results[i].PlayerID)
	}
	db.RDB.SAdd(Ctx, PlayerKeySet, ids...).Result()
}
