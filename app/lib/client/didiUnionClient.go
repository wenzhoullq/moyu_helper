package client

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/url"
	"sort"
	"strings"
	"weixin_LLM/dto/response"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
)

type DidiUnionClient struct {
	*resty.Client
	appKey    string
	accessKey string
}

func NewDidiUnionClient(ops ...func(c *DidiUnionClient)) *DidiUnionClient {
	client := &DidiUnionClient{
		Client: resty.New(),
	}
	for _, op := range ops {
		op(client)
	}
	return client
}

func SetDidiUnionAccessKey(accessKey string) func(client *DidiUnionClient) {
	return func(client *DidiUnionClient) {
		client.accessKey = accessKey
	}
}
func SetDidiUnionAppKey(appKey string) func(client *DidiUnionClient) {
	return func(client *DidiUnionClient) {
		client.appKey = appKey
	}
}
func (client *DidiUnionClient) GetTodayProfit() (*response.DidiUnionResp, error) {
	url := "https://union.didi.cn/openapi/v1.0/order/list"
	queryMap := map[string]string{
		"pay_start_time": fmt.Sprint(lib.GetUnix(0, 0, -1)),
		"pay_end_time":   fmt.Sprint(lib.GetUnix(0, 0, 0)),
	}
	headerMap := map[string]string{
		"App-Key":   client.appKey,
		"Timestamp": fmt.Sprint(lib.GetUnix(0, 0, 0)),
	}
	signMap := map[string]string{
		"App-Key":        client.appKey,
		"Timestamp":      fmt.Sprint(lib.GetUnix(0, 0, 0)),
		"pay_start_time": fmt.Sprint(lib.GetUnix(0, 0, -1)),
		"pay_end_time":   fmt.Sprint(lib.GetUnix(0, 0, 0)),
	}
	sign := client.getSign(signMap, client.accessKey)
	headerMap["Sign"] = sign
	resp, err := client.Client.R().
		SetQueryParams(queryMap).
		SetHeaders(headerMap).
		Get(url)
	if err != nil {
		return nil, err
	}
	didiResp := &response.DidiUnionResp{}
	err = json.Unmarshal(resp.Body(), didiResp)
	if err != nil {
		return nil, err
	}
	if didiResp.Errno != constant.Success {
		return nil, errors.New(didiResp.ErrMsg)
	}
	return didiResp, nil
}
func (client *DidiUnionClient) getSign(params map[string]string, accessKey string) string {
	// key排序
	arr := sort.StringSlice{}
	for k := range params {
		if k != "sign" {
			arr = append(arr, k)
		}
	}
	arr.Sort()
	// 参数拼接
	var build strings.Builder
	for idx, k := range arr {
		if idx != 0 {
			build.WriteString("&")
		}
		build.WriteString(fmt.Sprintf("%s=%v", k, params[k]))
	}
	build.WriteString(accessKey)
	// URL encode
	sourceStr := url.QueryEscape(build.String())
	// sha1加密
	h := sha1.New()
	_, _ = io.WriteString(h, sourceStr)
	shaStr := hex.EncodeToString(h.Sum([]byte("")))
	// 返回base64字符串
	b64Str := base64.StdEncoding.EncodeToString([]byte(shaStr))
	// base64字符串含有=和/，再一次URL encode
	return url.QueryEscape(b64Str)
}
