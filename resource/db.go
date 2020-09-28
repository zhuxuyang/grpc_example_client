package resource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	go_logger "github.com/phachon/go-logger"
	"github.com/spf13/viper"
	"log"
)

var GORM_REQUEST_ID = "X-Request-ID"
var db *gorm.DB

type DbLog struct {
	Logger *go_logger.Logger
}

func (l DbLog) Print(v ...interface{}) {
	if l.Logger != nil {
		l.Logger.Infof("%+v", v)
	}
}

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
	mdb.SetLogger(DbLog{})

	if err = db.DB().Ping(); err != nil {
		log.Println("db ping err", err)
	}
	AddGormCallbacks(db)
	return
}

// GetDB 获取数据库链接实例
func GetDB() *gorm.DB {
	return db
}

func GetDBWithEchoContext(c echo.Context) *gorm.DB {
	key := c.Response().Header().Get(echo.HeaderXRequestID)
	return db.Set(GORM_REQUEST_ID, key)
}

func SetKeyToGorm(c echo.Context, d *gorm.DB) *gorm.DB {
	key := c.Response().Header().Get(echo.HeaderXRequestID)
	return d.Set(GORM_REQUEST_ID, key)
}

// AddGormCallbacks adds callbacks for tracing, you should call SetSpanToGorm to make them work
func AddGormCallbacks(db *gorm.DB) {
	callbacks := newCallbacks()
	registerCallbacks(db, "create", callbacks)
	registerCallbacks(db, "query", callbacks)
	registerCallbacks(db, "update", callbacks)
	registerCallbacks(db, "delete", callbacks)
	registerCallbacks(db, "row_query", callbacks)
}

type callbacks struct{}

func newCallbacks() *callbacks {
	return &callbacks{}
}

func (c *callbacks) beforeCreate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterCreate(scope *gorm.Scope)    { c.after(scope, "INSERT") }
func (c *callbacks) beforeQuery(scope *gorm.Scope)    { c.before(scope) }
func (c *callbacks) afterQuery(scope *gorm.Scope)     { c.after(scope, "SELECT") }
func (c *callbacks) beforeUpdate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterUpdate(scope *gorm.Scope)    { c.after(scope, "UPDATE") }
func (c *callbacks) beforeDelete(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterDelete(scope *gorm.Scope)    { c.after(scope, "DELETE") }
func (c *callbacks) beforeRowQuery(scope *gorm.Scope) { c.before(scope) }
func (c *callbacks) afterRowQuery(scope *gorm.Scope)  { c.after(scope, "") }

func (c *callbacks) before(scope *gorm.Scope) {


}

func (c *callbacks) after(scope *gorm.Scope, operation string) {
	key, ok := scope.Get(GORM_REQUEST_ID)
	if !ok {
		return
	}
	info := fmt.Sprintf("%s: %v err: %v vars: %v RowsAffected %v  %v  %v",
		GORM_REQUEST_ID, key, scope.DB().Error, scope.SQLVars, scope.DB().RowsAffected, scope.SQL, scope.Value)
	if Logger != nil {
		Logger.Info(info)
	} else {
		log.Println(info)
	}
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	switch name {
	case "create":
		db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
	case "query":
		db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
	case "update":
		db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
	case "delete":
		db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
	case "row_query":
		db.Callback().RowQuery().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery)
		db.Callback().RowQuery().After(gormCallbackName).Register(afterName, c.afterRowQuery)
	}
}

func Close() {
	_ = db.Close()
}
