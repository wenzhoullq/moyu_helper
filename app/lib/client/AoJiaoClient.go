package client

import (
	"github.com/baidubce/app-builder/go/appbuilder"
	"weixin_LLM/init/config"
)

type AoJiaoClient struct {
	*appbuilder.AgentBuilder
}

func NewAoJiaoClient() *AoJiaoClient {
	conf, _ := appbuilder.NewSDKConfig("", config.Config.AoJiaoKey)
	agent, _ := appbuilder.NewAgentBuilder(config.Config.AojiaoAppId, conf)
	client := &AoJiaoClient{
		AgentBuilder: agent,
	}
	return client
}

func (client *AoJiaoClient) Chat(query string) (string, error) {
	conversationID, err := client.AgentBuilder.CreateConversation()
	if err != nil {
		return "", err
	}
	i, err := client.Run(conversationID, query, nil, true)
	if err != nil {
		return "", err
	}
	totalAnswer := ""
	for answer, err := i.Next(); err == nil; answer, err = i.Next() {
		totalAnswer = totalAnswer + answer.Answer
	}
	return totalAnswer, nil
}
