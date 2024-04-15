package client

import (
	"fmt"
	"testing"
	"weixin_LLM/dto/chat"
)

func TestErnie8KChat(t *testing.T) {
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
	EClient := NewErnie8KClient(SetToken(token))
	mesConetent := append(make([]*chat.ChatForm, 0), &chat.ChatForm{
		Role:    "user",
		Content: "你是傻逼",
	})
	res, err := EClient.Chat(mesConetent)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(res)
}
