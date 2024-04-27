package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"sort"
	"strings"
	"weixin_LLM/dto/response"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
)

type MeiTuanUnionClient struct {
	*resty.Client
	appKey   string
	apiToken string
}

func NewMeiTuanUnionClient(ops ...func(c *MeiTuanUnionClient)) *MeiTuanUnionClient {
	client := &MeiTuanUnionClient{
		Client: resty.New(),
	}
	for _, op := range ops {
		op(client)
	}
	return client
}

func SetMeiTuanUnionAppKey(appKey string) func(client *MeiTuanUnionClient) {
	return func(client *MeiTuanUnionClient) {
		client.appKey = appKey
	}
}
func SetMeiTuanUnionApiToken(apiToken string) func(client *MeiTuanUnionClient) {
	return func(client *MeiTuanUnionClient) {
		client.apiToken = apiToken
	}
}

func (client *MeiTuanUnionClient) GetTodayProfit() (*response.MeiTuanUnionOrderResp, error) {
	url := "https://openapi.meituan.com/api/orderList"
	queryMap := map[string]string{
		"appkey":    client.appKey,
		"ts":        fmt.Sprintf("%d", lib.GetUnix(0, 0, 0)),
		"actId":     constant.MeiTuanActID,
		"startTime": fmt.Sprintf("%d", lib.GetUnix(0, 0, -1)),
		"endTime":   fmt.Sprintf("%d", lib.GetUnix(0, 0, 0)),
		"page":      constant.MeiTuanUnionPage,
		"limit":     constant.MeiTuanUnionLimit,
	}
	sign := client.genSign(queryMap)
	queryMap["sign"] = sign
	resp, err := client.Client.R().
		SetQueryParams(queryMap).
		Get(url)
	if err != nil {
		return nil, err
	}
	meiTuanResp := &response.MeiTuanUnionOrderResp{}
	err = json.Unmarshal(resp.Body(), meiTuanResp)
	if err != nil {
		return nil, err
	}
	return meiTuanResp, nil
}
func (client *MeiTuanUnionClient) genSign(params map[string]string) string {
	// 按键排序并拼接字符串
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var stringBuilder strings.Builder
	stringBuilder.WriteString(client.apiToken)
	for _, key := range keys {
		stringBuilder.WriteString(key)
		stringBuilder.WriteString(params[key])
	}
	stringBuilder.WriteString(client.apiToken)
	// 计算MD5
	return md5String(stringBuilder.String())
}

func md5String(source string) string {
	hasher := md5.New()
	hasher.Write([]byte(source))
	return hex.EncodeToString(hasher.Sum(nil))
}
