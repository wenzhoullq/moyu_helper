package main

import (
	l "log"
	"net/http"
	"os"
	"weixin_LLM/init/common"
	"weixin_LLM/init/config"
	"weixin_LLM/init/db"
	"weixin_LLM/init/log"
	"weixin_LLM/init/redis"
	"weixin_LLM/init/route"
	"weixin_LLM/service"
)

func main() {
	var err error
	switch os.Getenv("ENV") {
	case "test":
		err = config.ConfigInit("../config/configTest.toml")
		break
	case "dev":
		err = config.ConfigInit("../config/configDev.toml")
		break
	default:
		l.Panicln("Env is wrong and env is " + os.Getenv("ENV"))
		return
	}
	if err != nil {
		panic(err)
	}
	err = log.InitLog()
	if err != nil {
		panic(err)
	}
	err = db.InitDB()
	if err != nil {
		panic(err)
	}
	err = redis.InitRedis()
	if err != nil {
		panic(err)
	}
	err = common.InitCommon()
	if err != nil {
		panic(err)
	}
	ws := service.NewWxService()
	err = ws.InitWxRobot()
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Addr:    config.Config.ServerAddr,
		Handler: route.RouteInit(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
