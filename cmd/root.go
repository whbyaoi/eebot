package cmd

import (
	"eebot/g"

	"github.com/spf13/cobra"
)

var (
	CfgFlag string
)

var rootCmd = &cobra.Command{
	Use:   "eebot",
	Short: "run bot",
}

func init() {
	cobra.OnInitialize(InitConfig)

	rootCmd.PersistentFlags().StringVarP(&CfgFlag, "config", "c", "", "config file (required)")
	rootCmd.MarkPersistentFlagRequired("config")

	rootCmd.AddCommand(Analysis300Cmd)
	rootCmd.AddCommand(CollectDataCmd)
	rootCmd.AddCommand(AddTimestampCmd)
	rootCmd.AddCommand(UpdatePlayerSetCmd)
	rootCmd.AddCommand(HttpCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func InitConfig() {
	g.Config.SetConfigFile(CfgFlag)
	g.Config.SetConfigType("yaml")
	if err := g.Config.ReadInConfig(); err != nil {
		g.Logger.Fatalf("读取配置文件错误：%s", err.Error())
	}
}
