package client

import (
	"fmt"
	"testing"
	"weixin_LLM/dto/chat"
)

func TestDoctorChat(t *testing.T) {
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
	DoctorClient := NewDoctorClient(SetDoctorClientToken(token))
	msgConetent := append(make([]*chat.ChatForm, 0), &chat.ChatForm{
		Role:    "user",
		Content: "你是傻逼",
	})
	res, err := DoctorClient.Chat(msgConetent)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	t.Log(res)
	fmt.Sprintf("%#v", res)
}
