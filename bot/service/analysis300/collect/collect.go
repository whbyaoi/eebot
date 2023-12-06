package collect

import (
	"context"
	"eebot/bot/service/analysis300/db"
	"eebot/g"
	"time"
)

var PlayerListKey = "300analysis:player_list"

var PlayerKeySet = "300analysis:player_set"

var PlayerIDToNameKey = "300analysis:player_id_to_name"

var Expiration = 7 * 24 * time.Hour

var Ctx = context.Background()

func CrawlPlayerByName(name string) (err error) {
	PlayerID, err := SearchRoleID(name)
	if err != nil {
		g.Logger.Info("error: search player id, ", err.Error())
		return err
	}
	ids := Crawler.CrawlAllAndSave(PlayerID, 0)
	if len(ids) > 0 {
		db.RDB.SAdd(Ctx, PlayerKeySet, ids...).Result()
	}
	return
}

func CrawlPlayerByID(PlayerID uint64) {
	ids := Crawler.CrawlAllAndSave(PlayerID, 0)
	if len(ids) > 0 {
		db.RDB.SAdd(Ctx, PlayerKeySet, ids...).Result()
	}
}
