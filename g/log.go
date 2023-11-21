package g

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

var CrawlLogger = logrus.New()

func InitLog() {
	file, err := os.OpenFile(Config.GetString("auto-collect.log"), os.O_RDWR|os.O_TRUNC|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		Logger.Info("创建crawl日志文件失败")
	}
	CrawlLogger.SetOutput(file)
}
