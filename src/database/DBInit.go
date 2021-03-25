package database

import (
	"SHUCTES/src/config"
	. "SHUCTES/src/log"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB

func InitDatabase() {
    //[username[:password]@][protocol[(address)]]/dbname
	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?interpolateParams=%t&parseTime=%t",
		config.Conf.Drv.Username, config.Conf.Drv.Password, config.Conf.Drv.Network, config.Conf.Drv.Server, config.Conf.Drv.Port, config.Conf.Drv.Database, config.Conf.Drv.InterpolateParams, config.Conf.Drv.ParseTime)
	var err error
	DB, err = sql.Open("mysql", conn)
	if err != nil {
		Logger.Panicf("Failed to connect to mysql: %s", err)
	} else {
		Logger.Infof("Connect to database successful")
	}

	DB.SetConnMaxLifetime(time.Minute * 3)  	//最大连接周期
	DB.SetMaxOpenConns(10)              		//最大连接数
	DB.SetMaxIdleConns(10)                   //最大闲置连接数

	Logger.Info("Database init successfully")
}





