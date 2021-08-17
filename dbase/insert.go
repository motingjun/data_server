package dbase

import (
	"fmt"
	"strings"

	"data_server/logger"

	"github.com/jmoiron/sqlx"
)

type Episode struct {
	ContId   string `db:"cont_id"`
	SiteId   int    `db:"site_id"`
	Title    string `db:"title"`
	Url      string `db:"url"`
	ThumbUrl string `db:"thumb_url"`
	Duration int    `db:"duration"`
	CpName   string `db:"cp_name"`
	Played   string `db:"played"`
	Liked    string `db:"liked"`
	Comment  string `db:"comment"`
	Utime    string `db:"utime"`
	Mid      string `db:"mid"`
}

// 批量插入
func InsertRow(db *sqlx.DB, episodes []*Episode) error {
	// 存放 (?, ?) 的slice
	valueStrings := make([]string, 0, len(episodes))
	// 存放values的slice
	valueArgs := make([]interface{}, 0, len(episodes)*2)
	// 遍历episodes准备相关数据
	for _, epi := range episodes {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, epi.ContId)
		valueArgs = append(valueArgs, epi.SiteId)
		valueArgs = append(valueArgs, epi.Title)
		valueArgs = append(valueArgs, epi.Url)
		valueArgs = append(valueArgs, epi.ThumbUrl)
		valueArgs = append(valueArgs, epi.Duration)
		valueArgs = append(valueArgs, epi.CpName)
		valueArgs = append(valueArgs, epi.Played)
		valueArgs = append(valueArgs, epi.Liked)
		valueArgs = append(valueArgs, epi.Comment)
		valueArgs = append(valueArgs, epi.Utime)
		valueArgs = append(valueArgs, epi.Mid)
	}
	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO episode (cont_id, site_id, title, url, thumb_url, duration, cp_name, played, liked, comment, utime, mid, ctime, mtime) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}

func InsertRowSingle(db *sqlx.DB, episode []string, mid int, contID string, log *logger.Log) {
	sqlStr := "INSERT INTO episode (cont_id, site_id, title, url, thumb_url, duration, cp_name, played, liked, comment, utime, mid, ctime, mtime) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,now(),now())"
	var (
		site_id  int
		duration int
	)
	duration = 20
	site_id = 12
	if len(episode) > 0 {
		ret, err := db.Exec(sqlStr, contID, site_id, episode[1], episode[2], episode[3], duration, episode[6], episode[7], episode[9], episode[8], episode[4], mid)
		if err != nil {
			log.Logger.Errorln(fmt.Sprintf("insert episode failed, err: %s", err))
			return
		}
		theID, err := ret.LastInsertId() // 新插入数据的id
		if err != nil {
			log.Logger.Errorln(fmt.Sprintf("get lastinsert ID failed, err: %s", err))
			return
		}
		log.Logger.Infoln(fmt.Sprintf("insert episode success, the id is : %d.", theID))
	}

}
