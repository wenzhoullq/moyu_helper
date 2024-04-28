package client

import (
	"testing"
	"weixin_LLM/init/config"
)

func TestEleClient(t *testing.T) {
	err := initTest("../../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	client := NewEleUnionClient(SetEleUnionAppKey(config.Config.EleConfigure.AppKey), SetEleUnionSecret(config.Config.EleConfigure.Secret))
	if err != nil {
		t.Error(err)
		panic(err)
	}
	err = client.GetTodayProfit()
	if err != nil {
		panic(err)
	}
	//resp, err := client.GetTodayProfit()
	//if err != nil {
	//	panic(err)
	//}
	//if err != nil {
	//	t.Error(err)
	//	panic(err)
	//}
	//cnt := 0
	//for _, v := range resp.Data.OrderList {
	//	cnt += v.CPAProfit + v.CPSProfit
	//}
	//t.Log(cnt)
}
