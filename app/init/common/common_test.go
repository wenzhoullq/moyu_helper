package common

import (
	"testing"
	"weixin_LLM/init/config"
)

func TestCommon(t *testing.T) {
	err := initCommonConfig()
	if err != nil {
		panic(err)
	}
	err = InitHoliday("../../../file/holiday2024.json")
	if err != nil {
		panic(err)
	}
	//jsonStr, _ := json.Marshal(Holidays)
	for _, v := range Holidays {
		t.Log(v)
	}
	//t.Log(string(jsonStr))
}
func initCommonConfig() error {
	err := config.ConfigInit("../../../config/configTest.toml")
	if err != nil {
		return err
	}
	return nil
}
