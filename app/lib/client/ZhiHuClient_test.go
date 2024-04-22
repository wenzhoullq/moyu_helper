package client

import (
	"testing"
)

func TestHotTopic(t *testing.T) {
	err := initTest("../../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	client := NewZhiHuClient()
	topic, err := client.GetHotTopic()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, v := range topic.Data {
		t.Log(v)
	}
}
