package client

import (
	"strconv"
	"testing"
	"weixin_LLM/init/config"
)

func TestMeiTuan(t *testing.T) {
	err := initTest("../../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	client := NewMeiTuanUnionClient(SetMeiTuanUnionAppKey(config.Config.MeiTuanConfigure.AppKey), SetMeiTuanUnionApiToken(config.Config.MeiTuanConfigure.ApiToken))
	if err != nil {
		t.Error(err)
		panic(err)
	}
	resp, err := client.GetTodayProfit()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	cnt := 0.0
	for _, v := range resp.DataList {
		f, err := strconv.ParseFloat(v.Profit, 64)
		if err != nil {
			panic(err)
		}
		cnt += f
	}
	t.Log(cnt)
}
