package client

import (
	"fmt"
	"testing"
	"weixin_LLM/init/config"
)

func initTest(confAddress string) error {
	err := config.ConfigInit(confAddress)
	if err != nil {
		return err
	}
	return nil
}

func TestQianFanGetToken(t *testing.T) {
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
	fmt.Println(token)
}
