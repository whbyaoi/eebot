package cmd

import (
	"eebot/bot/router"
	"eebot/bot/service/analysis300/collect"
	"eebot/bot/service/analysis300/db"
	"eebot/bot/service/http"
	"eebot/g"
	"eebot/ws"
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/spf13/cobra"
)

var Analysis300Cmd = &cobra.Command{
	Use:   "300",
	Short: "300 bot",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()
		g.InitLog()
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
		g.InitLog()
		collect.Crawler.AutoIncrementalCrawl()
	},
}

var UpdatePlayerSetCmd = &cobra.Command{
	Use:   "300-update-player-set",
	Short: "Update the player set of redis for 300-collect command",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()
		g.InitLog()
		collect.Crawler.UpdatePlayerSet()
	},
}

var AddTimestampCmd = &cobra.Command{
	Use:   "300-add-timestamp",
	Short: "add timestamp to table players (may block 300 bot service)",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()
		g.InitLog()
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

var HttpCmd = &cobra.Command{
	Use:   "300-http",
	Short: "export service by http",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitRedis()
		db.InitMysql()
		g.InitLog()
		r := http.New()
		r.Run(":8090")
	},
}
