package cmd

import (
	"context"
	"eebot/bot/router"
	"eebot/bot/service/analysis300/analysis"
	"eebot/bot/service/analysis300/collect"
	"eebot/bot/service/analysis300/db"
	"eebot/ws"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var Analysis300Cmd = &cobra.Command{
	Use:   "300",
	Short: "300 bot",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()
		collect.InitCrawler()
		for {
			if err := ws.InitWebsocket(); err != nil {
				continue
			}
			if err := ws.Read(router.WsMessageHandler); err != nil {
				continue
			}
			ws.WsClient.Close()
			time.Sleep(time.Minute)
		}
	},
}

var CollectDataCmd = &cobra.Command{
	Use:   "300-collect",
	Short: "collect 300 data via redis keys (may block 300 bot service)",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()
		collect.Crawler.IncrementalCrawl()
	},
}

var RefreshIntervalCmd = &cobra.Command{
	Use:   "300-refresh-interval",
	Short: "refresh intervals of players for shuffle anslysis (may block 300 bot service)",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()
		analysis.InitMatchInterval()
	},
}

var AddTimestampCmd = &cobra.Command{
	Use:   "300-add-timestamp",
	Short: "add timestamp to table players (may block 300 bot service)",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()

		matches := []db.Match{}
		db.SqlDB.Model(db.Match{}).Find(&matches)
		fmt.Printf("total: %d\n", len(matches))
		for i := range matches {
			db.SqlDB.Model(db.Player{}).Where("match_id = ?", matches[i].MatchID).Update("create_time", matches[i].CreateTime)
			if (i+1)%100 == 0 {
				fmt.Printf("%d over\n", i+1)
			}
		}
	},
}

var TranTimestampCmd = &cobra.Command{
	Use:   "300-tran",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()

		playerIDs, err := db.RDB.LRange(context.Background(), "300analysis:player_list", 0, -1).Result()
		if err != nil {
			return
		}
		var data []interface{}
		for i := range playerIDs {
			data = append(data, playerIDs[i])
		}
		_ = db.RDB.SAdd(context.Background(), "300analysis:player_set", data...)
	},
}
