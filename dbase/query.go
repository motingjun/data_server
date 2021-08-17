package dbase

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Mission struct {
	Mid int
	Url string
}

type EpisodeQuery struct {
	Id int
}

func QueryByCpid(db *sqlx.DB, cpID string) (int, string) {
	sqlStr := "SELECT mid, url FROM mission WHERE site_id = 12 AND cp_id = ?"
	var m Mission
	err := db.Get(&m, sqlStr, cpID)
	if err != nil {
		return 0, ""
	}
	return m.Mid, m.Url
}

func QueryEpisodeByTitle(db *sqlx.DB, title string) int {
	tables := [3]string{"episode", "episode_suc", "episode_suc_01"}
	var e EpisodeQuery
	for _, table := range tables {
		sqlStr := "SELECT id FROM %s WHERE title = ? and ctime >= DATE_SUB(CURDATE(), INTERVAL ? DAY)"
		sqlStr = fmt.Sprintf(sqlStr, table)
		err := db.Get(&e, sqlStr, title, 30)
		if err != nil {
			return 0
		}
		if e.Id != 0 {
			return e.Id
		}
	}
	return 0
}

func QueryEpisodeById(db *sqlx.DB, cont_id string, site_id int) int {
	tables := [3]string{"episode", "episode_suc", "episode_suc_01"}
	var e EpisodeQuery
	for _, table := range tables {
		sqlStr := "SELECT id FROM %s WHERE cont_id = ? AND site_id = ?"
		sqlStr = fmt.Sprintf(sqlStr, table)
		err := db.Get(&e, sqlStr, cont_id, site_id)
		if err != nil {
			return 0
		}
		if e.Id != 0 {
			return e.Id
		}
	}
	return 0
}
