package collect

import (
	"context"
	"eebot/g"
	"time"
)

var PlayerListKey = "300analysis:player_list"

var PlayerIDToNameKey = "300analysis:player_id_to_name"

var Expiration = 24 * time.Hour

var Ctx = context.Background()

func CrawlPlayerByName(name string) (err error) {
	PlayerID, err := SearchRoleID(name)
	if err != nil {
		g.Logger.Info("error: search player id, ", err.Error())
		return err
	}
	Crawler.CrawlAllAndSave(PlayerID, 0)
	return
}

func CrawlPlayerByID(PlayerID uint64) {
	Crawler.CrawlAllAndSave(PlayerID, 0)
}
