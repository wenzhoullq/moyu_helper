package client

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"weixin_LLM/dto/response"
)

type ZhiHuClient struct {
	client *resty.Client
}

func NewZhiHuClient(ops ...func(model *ZhiHuClient)) *ZhiHuClient {
	client := &ZhiHuClient{
		client: resty.New(),
	}
	for _, op := range ops {
		op(client)
	}
	return client
}

func (client *ZhiHuClient) GetHotTopic() (*response.ZhiHuTopicResponse, error) {
	url := "https://www.zhihu.com/api/v3/feed/topstory/hot-lists/total"
	resp, err := client.client.R().
		Get(url)
	if err != nil {
		return nil, err
	}
	result := &response.ZhiHuTopicResponse{}
	err = json.Unmarshal(resp.Body(), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
