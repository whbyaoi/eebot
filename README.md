# eebot

一个与cqhttp交互的bot后端

需要自行配置golang(版本>1.21.0，推荐最新)，cqhttp，mysql及redis，配置文件参考/config/server.yaml

sql文件在/bot/service/analysis300/db下，其中db_partition.sql为可选文件（用于加速英雄相关的查询）

用法，在源路径下运行：go mod tidy & go run main.go 300 -c /path/to/config

/path/to/config替换成你的配置文件路径

~~已知问题，在windows环境下调用生成图片的命令会报错图片文件找不到，可以自行修改export.go文件夹内的寻找图片路径的代码中的file://为file:/// (即ExportActiveAnalysis，ExportPKAnalysis，ExportJJLWithTeamAnalysis函数的return代码)~~（已尝试修复）