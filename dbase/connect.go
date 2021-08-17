package dbase

import (
	"fmt"

	"data_server/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func InitDB(log *logger.Log) (db *sqlx.DB, err error) {
	dsn := "%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True"
	dsn = fmt.Sprintf(
		dsn,
		viper.GetString("DBUserPasswd"),
		viper.GetString("DBAddress"),
		viper.GetString("DBDatabase"))
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		// fmt.Printf("connect DB failed, err:%v\n", err)
		log.Logger.Errorln("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}
