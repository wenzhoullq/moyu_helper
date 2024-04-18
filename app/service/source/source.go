package source

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"weixin_LLM/dao"
	"weixin_LLM/dto/request"
	source2 "weixin_LLM/dto/source"
	"weixin_LLM/init/common"
	"weixin_LLM/init/log"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
)

type SourceService struct {
	*dao.SourceDao
	*logrus.Logger
}

func NewSourceService() *SourceService {
	ss := &SourceService{
		SourceDao: dao.NewSourceDao(),
		Logger:    log.Logger,
	}
	return ss
}

func (ss *SourceService) UpdateSource(c context.Context, req *request.UpdateSourceRequest) (*lib.Response, error) {
	resp := lib.NewResponse()
	return resp, nil
}
func (ss *SourceService) checkCreateSource(req *request.CreateSourceRequest) error {
	if req.SourceName == "" || req.SourceDesc == "" || req.SourceLink == "" || req.SourceType == 0 {
		return errors.New("param error")
	}
	if req.SourceType == constant.CommissionSource && req.SourceExp == "" {
		return errors.New("param error")
	}
	return nil
}
func (ss *SourceService) CreateSource(c context.Context, req *request.CreateSourceRequest) (*lib.Response, error) {
	resp := lib.NewResponse()
	if err := ss.checkCreateSource(req); err != nil {
		lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.ParamErr))
		return nil, err
	}
	source := &source2.Source{
		Status:     constant.SourceNorMal,
		SourceName: req.SourceName,
		SourceDesc: req.SourceDesc,
		SourceLink: req.SourceLink,
		SourceType: req.SourceType,
		SourceExp:  req.SourceExp,
	}
	if req.SourceType == constant.PublicSource {
		source.SourceExp = constant.NeverExp
	}
	if err := ss.SourceDao.CreateSource(source); err != nil {
		lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.DBErr))
		return resp, err
	}
	if err := common.InitTool(); err != nil {
		lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.ServerErr))
		return resp, err
	}
	ss.Logln(logrus.InfoLevel, "create Source success")
	return resp, nil
}

func (ss *SourceService) GetSource(c context.Context) (*lib.Response, error) {
	resp := lib.NewResponse()
	return resp, nil
}
