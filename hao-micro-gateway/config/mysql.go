package config

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DbConfig struct {
	DbName       string `json:"dbName"`
	DriverName   string `json:"driverName"`
	Dsn          string `json:"dsn"`
	ShowSql      bool   `json:"showSql"`
	ShowExecTime bool   `json:"showExecTime"`
	MaxIdle      int    `json:"maxIdle"`
	MaxOpen      int    `json:"maxOpen"`
}

// var Db = map[string]DbConfig{
// 	"db1": {
// 		DriverName:   "mysql",
// 		Dsn:          "root:My.931206@tcp(101.43.7.108:3306)/hao-micro?charset=utf8mb4&parseTime=true&loc=Local",
// 		ShowSql:      true,
// 		ShowExecTime: false,
// 		MaxIdle:      10,
// 		MaxOpen:      200,
// 	},
// }

// var Db = make(map[string]DbConfig)

var MySQLEngine *xorm.Engine

type DbHandler interface {
	AddMySqlDB()
}

func (db DbConfig) AddMySqlDB() {
	// Db[db.DbName] = db
	if MySQLEngine == nil {
		var err error
		MySQLEngine, err = xorm.NewEngine(db.DriverName, db.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		MySQLEngine.SetMaxIdleConns(db.MaxIdle) //空闲连接
		MySQLEngine.SetMaxOpenConns(db.MaxOpen) //最大连接数
		MySQLEngine.ShowSQL(db.ShowSql)
		MySQLEngine.ShowExecTime(db.ShowExecTime)
	}
}
