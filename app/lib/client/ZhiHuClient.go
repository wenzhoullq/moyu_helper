package client

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
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
	for i := range result.Data {
		result.Data[i].Target.URL = fmt.Sprintf("https://www.zhihu.com/question/%d", result.Data[i].Target.ID)
		titleSub := strings.Split(result.Data[i].Target.Title, "？")
		if len(titleSub) > 0 {
			result.Data[i].Target.Title = titleSub[0] + "?"
		}
	}
	return result, nil
}
