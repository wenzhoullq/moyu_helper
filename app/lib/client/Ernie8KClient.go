package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"weixin_LLM/dto/chat"
	"weixin_LLM/dto/response"
	"weixin_LLM/lib/constant"
)

type Ernie8KClient struct {
	client *resty.Client
	token  string
}

func NewErnie8KClient(ops ...func(c *Ernie8KClient)) *Ernie8KClient {
	client := &Ernie8KClient{
		client: resty.New(),
	}
	for _, op := range ops {
		op(client)
	}
	return client
}

func SetToken(token string) func(client *Ernie8KClient) {
	return func(client *Ernie8KClient) {
		client.token = token
	}
}

func (client *Ernie8KClient) Chat(content []*chat.ChatForm) (*response.Ernie8kResponse, error) {
	url := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions"
	headerMap := map[string]string{
		"Content-Type": "application/json",
	}
	queryMap := map[string]string{
		"access_token": client.token,
	}
	bodyMap := map[string]interface{}{
		"messages": content,
		"system":   constant.SystemNormal,
	}
	resp, err := client.client.R().
		SetBody(bodyMap).
		SetHeaders(headerMap).
		SetQueryParams(queryMap).
		Post(url)
	if err != nil {
		return nil, err
	}
	ernie8kResponse := &response.Ernie8kResponse{}
	err = json.Unmarshal(resp.Body(), ernie8kResponse)
	if err != nil {
		return nil, err
	}
	if ernie8kResponse.Flag != constant.RESP_NORMAL {
		return nil, errors.New(fmt.Sprintf("errNo:%d;errMsg:%s", ernie8kResponse.ErrorCode, ernie8kResponse.ErrorMsg))
	}
	return ernie8kResponse, nil
}
