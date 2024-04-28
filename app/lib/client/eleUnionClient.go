package client

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
	"sort"
	"strings"
	"weixin_LLM/lib"
)

type EleUnionClient struct {
	*resty.Client
	appKey string
	secret string
}

func NewEleUnionClient(ops ...func(c *EleUnionClient)) *EleUnionClient {
	client := &EleUnionClient{
		Client: resty.New(),
	}
	for _, op := range ops {
		op(client)
	}
	return client
}
func SetEleUnionAppKey(appKey string) func(client *EleUnionClient) {
	return func(client *EleUnionClient) {
		client.appKey = appKey
	}
}
func SetEleUnionSecret(secret string) func(client *EleUnionClient) {
	return func(client *EleUnionClient) {
		client.secret = secret
	}
}
func (client *EleUnionClient) GetTodayProfit() error {
	url := "https://eco.taobao.com/router/rest"
	queryMap := map[string]string{
		"method":      "alibaba.alsc.union.kbcpx.positive.order.get",
		"app_key":     client.appKey,
		"timestamp":   lib.GetCurTimeDetail(0, 0, 0),
		"v":           "2.0",
		"startTime":   fmt.Sprintf("%d", lib.GetUnix(0, 0, -1)),
		"endTime":     fmt.Sprintf("%d", lib.GetUnix(0, 0, 0)),
		"sign_method": "md5",
		"format":      "json",
		"simplify":    "True",
		"date_type":   "1",
		"biz_unit":    "2",
		"page_size":   "50",
		"page_number": "1",
		"start_date":  lib.GetCurTimeDetail(0, 0, -1),
		"end_date":    lib.GetCurTimeDetail(0, 0, 0),
	}
	sign, err := client.signTopRequest(queryMap, client.secret)
	fmt.Println(sign)
	if err != nil {
		return err
	}
	queryMap["sign"] = sign
	resp, err := client.Client.R().
		SetQueryParams(queryMap).
		Get(url)
	lib.WriteToTxt(resp.Body())
	return nil
}
func (client *EleUnionClient) signTopRequest(params map[string]string, secret string) (string, error) {
	// 第一步：检查参数是否已经排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 第二步：把所有参数名和参数值串在一起
	var query bytes.Buffer
	query.WriteString(secret)
	for _, key := range keys {
		value := params[key]
		if key != "" && value != "" {
			query.WriteString(key)
			query.WriteString(value)
		}
	}
	// 第三步：使用MD5/HMAC加密
	var bytes []byte
	query.WriteString(secret)
	bytes, err := encryptMD5(query.String())
	if err != nil {
		return "", err
	}
	// 第四步：把二进制转化为大写的十六进制
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func encryptMD5(data string) ([]byte, error) {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hasher.Sum(nil), nil
}
