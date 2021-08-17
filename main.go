package main

import (
	"fmt"
	"net/http"
	"path"
	"runtime"

	"data_server/dbase"
	fsync "data_server/fsync"
	"data_server/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 分配指定逻辑处理器给调度器使用
	runtime.GOMAXPROCS(viper.GetInt("NumCpu"))
	// fmt.Println("cpus:", runtime.NumCPU())
	// fmt.Println("goroot:", runtime.GOROOT())
	// fmt.Println("archive:", runtime.GOOS)

	// 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	errConfig := viper.ReadInConfig()
	if errConfig != nil {
		panic(fmt.Errorf("fatal error config file: %v", errConfig))
	}

	// 初始化日志
	log := logger.InitLog()

	// 创建mysql实例
	mysqlDB, _ := dbase.InitDB(log)
	// 尝试与数据库建立连接（校验dsn是否正确）
	err := mysqlDB.Ping()
	if err != nil {
		log.Logger.Errorln(err)
		mysqlDB, _ = dbase.InitDB(log)
	}

	// 创建有缓冲通道
	buffered := make(chan string, viper.GetInt("TaskLoad"))
	// 创建路由
	r := gin.Default()

	// 限制表单上传大小，默认为32MB
	r.MaxMultipartMemory = int64(viper.GetInt("MaxMultipartMemory"))
	r.LoadHTMLFiles("./template/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 上传文件
	r.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("get err %s", err.Error()),
			})
			log.Logger.Errorln(err.Error())
		}

		// 获取所有文件
		files := form.File["files"]
		// 遍历所有文件
		for _, file := range files {
			// 判断文件是否为xslx
			fileExt := path.Ext(file.Filename)
			if fileExt != ".xlsx" {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "上传失败：需上传 ‘.xlsx’ 格式文件",
				})
				log.Logger.Errorln("上传失败：需上传 ‘.xlsx’ 格式文件")
				return
			}
			buffered <- file.Filename // 添加文件名到通道
			filename := path.Join("./upload", file.Filename)
			// 上传文件到指定目录
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("上传文件到指定目录失败：%s", err.Error()),
				})
				log.Logger.Errorln(fmt.Sprintf("上传文件到指定目录失败：%s", err.Error()))
				return
			}
			// 异步执行数据同步
			go fsync.SyncReadFile(buffered, mysqlDB, log)
		}
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("上传成功 %d 个文件", len(files)),
		})
		log.Logger.Infoln(fmt.Sprintf("上传成功 %d 个文件", len(files)))
	})

	log.Logger.Infoln(fmt.Sprintf("data_server starting, ServerAddress : %s", viper.GetString("ServerAddress")))
	r.Run(viper.GetString("ServerAddress"))
}
