package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"weixin_LLM/dto/response"
	"weixin_LLM/init/config"
	"weixin_LLM/lib/constant"
)

type TanshuClient struct {
	client *resty.Client
	key    string
}

func NewTanshuClient(ops ...func(model *TanshuClient)) *TanshuClient {
	client := &TanshuClient{
		client: resty.New(),
		key:    config.Config.TanshuConfigure.Key,
	}
	for _, op := range ops {
		op(client)
	}
	return client
}

func (client *TanshuClient) GetGoldPrice() (string, error) {
	url := "http://api.tanshuapi.com/api/gold/v1/shgold"
	queryMap := map[string]string{
		"key": client.key,
	}
	resp, err := client.client.R().
		SetQueryParams(queryMap).
		Get(url)
	if err != nil {
		return "", err
	}
	goldResp := &response.GoldResponse{}
	err = json.Unmarshal(resp.Body(), goldResp)
	if err != nil {
		return "", err
	}
	if goldResp.Code != constant.TanshuSuccess {
		return "", errors.New(fmt.Sprintf("code:%d;codeMsg:%s", goldResp.Code, goldResp.Msg))
	}
	var huGoldPricce string
	for _, v := range goldResp.Data.List {
		if v.Typename != "AU9999" {
			continue
		}
		huGoldPricce = v.Price
		break
	}
	return huGoldPricce, nil
}

func (client *TanshuClient) GetNews() ([]string, error) {
	url := "http://api.tanshuapi.com/api/toutiao/v1/index"
	queryMap := map[string]string{
		"key": client.key,
		"num": constant.NewsNum,
	}
	resp, err := client.client.R().
		SetQueryParams(queryMap).
		Get(url)
	if err != nil {
		return nil, err
	}
	newsResp := &response.NewsResp{}
	err = json.Unmarshal(resp.Body(), newsResp)
	if err != nil {
		return nil, err
	}
	if newsResp.Code != constant.TanshuSuccess {
		return nil, errors.New(fmt.Sprintf("code:%d;codeMsg:%s", newsResp.Code, newsResp.Msg))
	}
	news := make([]string, 0)
	for _, v := range newsResp.Data.List {
		news = append(news, v.Title)
	}
	return news, nil
}
