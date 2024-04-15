package client

import (
	"testing"
)

func TestTanshu(t *testing.T) {
	err := initTest("../../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	client := NewTanshuClient()
	goldPrice, err := client.GetGoldPrice()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	t.Log(goldPrice)
	news, err := client.GetNews()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	t.Log(news)
}
