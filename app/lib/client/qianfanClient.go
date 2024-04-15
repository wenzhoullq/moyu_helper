package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"weixin_LLM/dto/response"
	"weixin_LLM/init/config"
)

type QianFancClient struct {
	client    *resty.Client
	apiKey    string
	secretKey string
}

func NewQianFanClient(ops ...func(model *QianFancClient)) *QianFancClient {
	client := &QianFancClient{
		client:    resty.New(),
		apiKey:    config.Config.QianFanConfigure.ApiKey,
		secretKey: config.Config.QianFanConfigure.SecretKey,
	}
	for _, op := range ops {
		op(client)
	}
	return client
}

func (client *QianFancClient) GetToken() (string, error) {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	headerMap := map[string]string{
		"Content-Type": "application/json",
	}
	queryMap := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     client.apiKey,
		"client_secret": client.secretKey,
	}
	resp, err := client.client.R().
		SetHeaders(headerMap).
		SetQueryParams(queryMap).
		Post(url)
	if err != nil {
		return "", err
	}
	tokenResp := &response.TokenResponse{}
	err = json.Unmarshal(resp.Body(), tokenResp)
	if tokenResp.Error != "" {
		return "", errors.New(fmt.Sprintf("getToken fail,Error:%s,ErrorDescription:%s", tokenResp.Error, tokenResp.ErrorDescription))
	}
	return tokenResp.AccessToken, nil
}
