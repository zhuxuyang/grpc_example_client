package resource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
)

var db *gorm.DB

// InitDB 初始化 MySQL 链接
func InitDB() {
	dbConf := viper.GetStringMapString("database")
	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["name"])
	mdb, err := gorm.Open("mysql", address)
	if err != nil {
		panic(err)
		return
	}
	if mdb == nil {
		panic("failed to connect database")
	}
	mdb.LogMode(true)
	log.Println("connected")
	db = mdb
	return
}

// GetDB 获取数据库链接实例
func GetDB() *gorm.DB {
	return db
}
