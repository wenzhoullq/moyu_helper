package source

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
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
func (ss *SourceService) checkUpdateSource(req *request.UpdateSourceRequest) error {
	if req.SourceName == "" {
		return errors.New("param error")
	}
	if req.SourceType == constant.CommissionSource && req.SourceExp == "" {
		return errors.New("param error")
	}
	return nil
}
func (ss *SourceService) UpdateSource(c context.Context, req *request.UpdateSourceRequest) (*lib.Response, error) {
	resp := lib.NewResponse()
	if err := ss.checkUpdateSource(req); err != nil {
		lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.ParamErr))
		return nil, err
	}
	source, err := ss.SourceDao.GetSourceByName(req.SourceName)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			lib.SetResponse(resp, lib.SetErrMsg("source not exist"), lib.SetErrNo(constant.DBErr))
			return resp, err
		}
		lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.DBErr))
		return resp, err
	}
	if len(req.SourceExp) > 0 {
		source.SourceExp = req.SourceExp
	}
	if len(req.SourceDesc) > 0 {
		source.SourceDesc = req.SourceDesc
	}
	if len(req.SourceLink) > 0 {
		source.SourceLink = req.SourceLink
	}
	if req.SourceType != 0 {
		source.SourceType = req.SourceType
	}
	if err := ss.SourceDao.UpdateSource(source); err != nil {
		lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.DBErr))
		return resp, err
	}
	if err := common.InitTool(); err != nil {
		lib.SetResponse(resp, lib.SetErrMsg(err.Error()), lib.SetErrNo(constant.ServerErr))
		return resp, err
	}
	ss.Logln(logrus.InfoLevel, "update Source success")
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
