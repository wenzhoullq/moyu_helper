package wx_cron

import (
	"testing"
	"weixin_LLM/init/common"
	"weixin_LLM/init/config"
)

func TestWxCron(t *testing.T) {
	err := initConfig()
	if err != nil {
		panic(err)
	}
	service := NewWxCronService()
	service.SendHolidayTips()
}
func initConfig() error {
	err := config.ConfigInit("../../../config/configTest.toml")
	if err != nil {
		return err
	}
	err = common.InitHoliday("../../../file/holiday2024.json")
	if err != nil {
		return err
	}
	return nil
}

func TestWorkDay(t *testing.T) {
	err := initConfig()
	if err != nil {
		panic(err)
	}
	service := NewWxCronService()
	t.Log(service.isWorkDay())
}

func TestZhiHuTopic(t *testing.T) {
	err := initConfig()
	if err != nil {
		panic(err)
	}
	service := NewWxCronService()
	service.SendNews()
}
