package fsync

import (
	"data_server/handler"
	"data_server/logger"
	"data_server/reader"

	"github.com/jmoiron/sqlx"
)

func SyncReadFile(buffered chan string, db *sqlx.DB, log *logger.Log) {
	filename, ok := <-buffered
	if !ok {
		// 这意味着通道已经为空了，并且已被关闭
		log.Logger.Infoln("All files have readed, chan has closed!")
		return
	}

	// 读取xlsx文件数据
	data := reader.SyncReadXlsx(filename, log)
	// 同步数据到mysql
	for len(data) > 0 {
		episodes := make([][]string, 0)
		if len(data) >= 100 {
			episodes = append(episodes, data[:100]...)
			data = data[100:]
		} else {
			episodes = append(episodes, data...)
			data = nil
		}
		handler.DataSave(db, episodes, log)
	}
}
