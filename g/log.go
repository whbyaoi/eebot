package g

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

var CrawlLogger = logrus.New()

func init() {
	file, err := os.OpenFile("./crawl.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Logger.Info("创建crawl日志文件失败")
	}
	CrawlLogger.SetOutput(file)
}
