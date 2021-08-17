package handler

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"data_server/dbase"
	"data_server/logger"
)

func DataSave(mysqlDB *sqlx.DB, episodes [][]string, log *logger.Log) {
	var (
		site_id int
		cont_id string
		title   string
		url     string
	)

	site_id = 12
	for _, episode := range episodes {
		if len(episode) == 14 {
			title = strings.TrimSpace(episode[1])
			url = strings.TrimSpace(episode[2])
			if len(title) != 0 || len(url) != 0 {
				s := strings.Split(episode[2], "/")
				cont_id = s[len(s)-1]
				eid := dbase.QueryEpisodeByTitle(mysqlDB, title)
				if eid != 0 {
					log.Logger.Infoln(fmt.Sprintf("episode title exist, title : %s", title))
					fmt.Println(title, len(title))
					return
				}

				eid = dbase.QueryEpisodeById(mysqlDB, cont_id, site_id)
				if eid != 0 {
					log.Logger.Infoln(fmt.Sprintf("episode cont_id and site_id exist, id : %d", eid))
					return
				}

				mid, _ := dbase.QueryByCpid(mysqlDB, episode[13])
				if mid != 0 {
					// log.Logger.Infoln(fmt.Sprintf("mission is exist, cp_id: %s, mid: %d", episode[13], mid))
					dbase.InsertRowSingle(mysqlDB, episode, mid, cont_id, log)
				} else {
					log.Logger.Errorln(fmt.Sprintf("cp_id: %s is wrong, no mid return!", episode[13]))
				}
			}
		}
	}
}
