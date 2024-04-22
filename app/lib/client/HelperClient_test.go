package client

import (
	"testing"
)

func TestNorMalChat(t *testing.T) {
	err := initTest("../../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	client := NewHelperClient()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	ans, err := client.Chat("金币有什么用")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	t.Log(ans)
}
