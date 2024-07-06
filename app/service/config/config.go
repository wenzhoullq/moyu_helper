package config

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"os"
	"weixin_LLM/init/config"
	"weixin_LLM/init/log"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
)

type ConfigService struct {
	*logrus.Logger
	*common.Client
}

func NewConfigService() *ConfigService {
	cs := &ConfigService{
		Logger: log.Logger,
	}
	return cs
}

func (cs *ConfigService) LoadConfig(c context.Context) (*lib.Response, error) {
	resp := lib.NewResponse()
	var err error
	switch os.Getenv("ENV") {
	case "test":
		err = config.ConfigInit("../config/configTest.toml")
		if err != nil {
			lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.ParamErr))
			return resp, err
		}
		break
	case "dev":
		err = config.ConfigInit("../config/configDev.toml")
		if err != nil {
			lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.ParamErr))
			return resp, err
		}
		break
	default:
		return resp, errors.New("Env is wrong and env is " + os.Getenv("ENV"))
	}
	cs.Logln(logrus.InfoLevel, "load Config success")
	return resp, nil
}
