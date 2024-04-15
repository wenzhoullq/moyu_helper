package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"weixin_LLM/dto/user"
	"weixin_LLM/init/config"
	"weixin_LLM/init/log"
)

var DB *gorm.DB

func InitDB() (err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local&timeout=%s", config.Config.UserName, config.Config.Pw, config.Config.MysqlConfigure.Host, config.Config.MysqlConfigure.Port, config.Config.DbName, config.Config.TimeOut)
	DB, err = gorm.Open(config.Config.Driver, dns)
	if err != nil {
		return err
	}
	err = DB.DB().Ping()
	if err != nil {
		return err
	}
	DB.AutoMigrate(&user.User{})
	DB.SetLogger(log.Logger)
	DB.LogMode(true)
	return nil
}
