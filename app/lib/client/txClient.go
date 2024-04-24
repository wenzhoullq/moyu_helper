package client

import (
	"encoding/json"
	aiart "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/aiart/v20221229"
	tx "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	r1 "weixin_LLM/dto/response"
	"weixin_LLM/init/config"
	"weixin_LLM/lib/constant"
)

type TxCloudClient struct {
	*aiart.Client
}

func NewTxCloudClient() *TxCloudClient {
	credential := tx.NewCredential(
		config.Config.TxConfigure.SecretId,
		config.Config.TxConfigure.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "aiart.tencentcloudapi.com"
	client, _ := aiart.NewClient(credential, constant.RegionShanghai, cpf)
	txClient := &TxCloudClient{
		Client: client,
	}
	return txClient
}

func (tc *TxCloudClient) PostImgToImg(base64 string, value string) (*r1.TxImgToImgResp, error) {
	request := aiart.NewImageToImageRequest()
	request.InputImage = &base64
	var logoAdd int64 = 0
	request.LogoAdd = &logoAdd
	request.Prompt = &value
	response, err := tc.Client.ImageToImage(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	txImgToImgResp := &r1.TxImgToImgResp{}
	err = json.Unmarshal([]byte(response.ToJsonString()), txImgToImgResp)
	if err != nil {
		return nil, err
	}
	return txImgToImgResp, nil
}

func (tc *TxCloudClient) PostTextToImg(text, style string) (*r1.TxImgToImgResp, error) {
	request := aiart.NewTextToImageRequest()
	request.Prompt = &text
	var logoAdd int64 = 0
	request.LogoAdd = &logoAdd
	styles := []*string{&style}
	request.Styles = styles
	response, err := tc.Client.TextToImage(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	txImgToImgResp := &r1.TxImgToImgResp{}
	err = json.Unmarshal([]byte(response.ToJsonString()), txImgToImgResp)
	if err != nil {
		return nil, err
	}
	return txImgToImgResp, nil
}
