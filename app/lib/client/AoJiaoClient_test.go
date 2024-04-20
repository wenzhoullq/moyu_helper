package client

import (
	"testing"
)

func TestAoJiaoChat(t *testing.T) {
	err := initTest("../../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	client, err := NewBaiduClient()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	ans, err := client.Chat("你好蠢")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	t.Log(ans)
}
