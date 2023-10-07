package cmd

import (
	"eebot/bot/router"
	"eebot/bot/service/analysis300/collect"
	"eebot/bot/service/analysis300/db"
	"eebot/ws"
	"time"

	"github.com/spf13/cobra"
)

var Analysis300Cmd = &cobra.Command{
	Use:   "300",
	Short: "300 bot",
	Run: func(cmd *cobra.Command, args []string) {
		collect.InitRedis()
		db.InitMysql()
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
