package client

import (
	"testing"
	"weixin_LLM/dto/chat"
)

func TestBaZongChat(t *testing.T) {
	err := initTest("../../../config/configTest.toml")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	client := NewQianFanClient()
	token, err := client.GetToken()
	if err != nil {
		t.Error(err)
		panic(err)
	}
	EClient := NewBaDaoClient(SetBaDaoToken(token))
	msgContent := append(make([]*chat.ChatForm, 0), &chat.ChatForm{
		Role:    "user",
		Content: "描述下霸道总裁吗",
	})
	res, err := EClient.ChatToBaZong(msgContent)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	t.Log(res)
}
