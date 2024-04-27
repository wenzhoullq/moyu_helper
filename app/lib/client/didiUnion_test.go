package client

import (
	"testing"
	"weixin_LLM/init/config"
)

func TestDidi(t *testing.T) {
	err := initTest("../../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	client := NewDidiUnionClient(SetDidiUnionAppKey(config.Config.DiDiConfigure.AppKey), SetDidiUnionAccessKey(config.Config.DiDiConfigure.AccessKey))
	if err != nil {
		t.Error(err)
		panic(err)
	}
	resp, err := client.GetTodayProfit()
	if err != nil {
		panic(err)
	}
	if err != nil {
		t.Error(err)
		panic(err)
	}
	cnt := 0
	for _, v := range resp.Data.OrderList {
		cnt += v.CPAProfit + v.CPSProfit
	}
	t.Log(cnt)
}
