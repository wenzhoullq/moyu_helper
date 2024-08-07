package wx

import (
	"encoding/json"
	"errors"
	"strings"
	"weixin_LLM/dao"
	group2 "weixin_LLM/dto/group"
	"weixin_LLM/init/config"
	"weixin_LLM/init/log"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
	"weixin_LLM/service/wx/wx_cron"
	"weixin_LLM/service/wx/wx_llm"

	"github.com/eatmoreapple/openwechat"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type WxService struct {
	*openwechat.Bot
	*logrus.Logger
	*wx_llm.WxLLMService
	*wx_cron.WxCronService
	wxDao      *dao.WxDao
	sourceDao  *dao.SourceDao
	groupSend  []func(*openwechat.Message) (bool, error)
	friendSend []func(*openwechat.Message) (bool, error)
	groups     openwechat.Groups
	friends    openwechat.Friends
	corn       *cron.Cron
}

func NewWxService() *WxService {
	ws := &WxService{
		Bot:       openwechat.DefaultBot(openwechat.Desktop),
		Logger:    log.Logger,
		wxDao:     dao.NewWxDao(),
		sourceDao: dao.NewSourceDao(),
		corn:      cron.New(),
	}
	ws.groupSend = []func(*openwechat.Message) (bool, error){ws.groupTextMsg, ws.groupImgMsg}
	ws.friendSend = []func(*openwechat.Message) (bool, error){ws.friendTextMsg, ws.friendImgMsg}
	return ws
}

func (ws *WxService) groupSender(msg *openwechat.Message) (bool, error) {
	if !msg.IsSendByGroup() {
		return false, nil
	}
	if msg.IsSendBySelf() {
		return true, nil
	}
	for _, f := range ws.groupSend {
		ok, err := f(msg)
		if err != nil {
			return true, err
		}
		if ok {
			return true, nil
		}
	}
	return true, nil
}

func (ws *WxService) friendSender(msg *openwechat.Message) (bool, error) {
	if !msg.IsSendByFriend() {
		return false, nil
	}
	if msg.IsSendBySelf() {
		return true, nil
	}
	for _, f := range ws.friendSend {
		ok, err := f(msg)
		if err != nil {
			return true, err
		}
		if ok {
			return true, nil
		}
	}
	return true, nil
}

func (ws *WxService) groupTextMsg(msg *openwechat.Message) (bool, error) {
	if !msg.IsText() {
		return false, nil
	}
	msg.Content = strings.TrimSpace(msg.Content)
	if strings.HasPrefix(msg.Content, constant.LlmKeyWord) {
		msg.Content = strings.TrimSpace(msg.Content[len(constant.LlmKeyWord):])
	} else if strings.HasSuffix(msg.Content, constant.LlmKeyWord) {
		msg.Content = strings.TrimSpace(msg.Content[:len(msg.Content)-len(constant.LlmKeyWord)])
	} else {
		return false, nil
	}
	//把小写全部转为大写
	msg.Content = lib.ProcessingCommands(msg.Content)
	for _, f := range ws.WxLLMService.GetGroupTextProducer() {
		ok, err := f(msg)
		if err != nil {
			return true, err
		}
		if ok {
			return true, nil
		}
	}
	return true, errors.New("no such group text request")
}

func (ws *WxService) groupImgMsg(msg *openwechat.Message) (bool, error) {
	if !msg.IsPicture() {
		return false, nil
	}
	for _, f := range ws.WxLLMService.GetGroupImgProducer() {
		ok, err := f(msg)
		if err != nil {
			return true, err
		}
		if ok {
			return true, nil
		}
	}
	return true, errors.New("no such group img request")
}

func (ws *WxService) friendTextMsg(msg *openwechat.Message) (bool, error) {
	if msg.IsSendBySelf() {
		return false, errors.New("msg send by self")
	}
	if !msg.IsText() {
		return false, nil
	}
	msg.Content = strings.TrimSpace(msg.Content)
	//把小写全部转为大写
	msg.Content = lib.ProcessingCommands(msg.Content)
	for _, f := range ws.WxLLMService.GetFriendTextProducer() {
		ok, err := f(msg)
		if err != nil {
			return true, err
		}
		if ok {
			return true, nil
		}
	}
	return true, nil
}

func (ws *WxService) friendImgMsg(msg *openwechat.Message) (bool, error) {
	if !msg.IsPicture() {
		return false, nil
	}
	for _, f := range ws.WxLLMService.GetFriendImgProducer() {
		ok, err := f(msg)
		if err != nil {
			return true, err
		}
		if ok {
			return true, nil
		}
	}
	return true, errors.New("no such user img request")
}

func (ws *WxService) InitGroupScribe() error {
	for _, v := range ws.groups {
		//如果没有该组,则添加订阅
		if _, err := ws.wxDao.GetGroupByName(v.NickName); err != nil {
			if gorm.IsRecordNotFoundError(err) {
				scribe := &group2.Subscribe{
					News: true,
					Tips: true,
				}
				scribeStr, err := json.Marshal(&scribe)
				if err != nil {
					return err
				}
				group := group2.Groups{
					GroupName: v.NickName,
					Subscribe: string(scribeStr),
				}
				if err = ws.wxDao.CreateGroup(&group); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (ws *WxService) InitWxRobot() error {
	// 注册登陆二维码回调
	ws.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// 登陆
	if err := ws.Login(); err != nil {
		return err
	}
	// 获取登陆的用户
	self, err := ws.GetCurrentUser()
	if err != nil {
		ws.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	friends, err := self.Friends()
	if err != nil {
		ws.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	ws.Logln(logrus.InfoLevel, friends)
	groups, err := self.Groups()
	if err != nil {
		ws.Logln(logrus.ErrorLevel, err.Error())
		return err
	}

	ws.groups = groups
	ws.Logln(logrus.InfoLevel, groups)
	// 初始化group订阅
	if err = ws.InitGroupScribe(); err != nil {
		ws.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	ws.friends = friends
	//初始化WxLLM
	ws.WxLLMService = wx_llm.NewWxLLMService(wx_llm.SetSelf(self), wx_llm.SetGroups(groups), wx_llm.SetFriends(friends),
		wx_llm.SetLog(ws.Logger), wx_llm.SetWxDao(ws.wxDao))
	//注册消息处理函数
	ws.MessageHandler = func(msg *openwechat.Message) {
		//对于不同的消息进行不同的处理
		ok, err := ws.friendSender(msg)
		if err != nil {
			ws.Logln(logrus.ErrorLevel, err.Error())
			return
		}
		if ok {
			return
		}
		ok, err = ws.groupSender(msg)
		if err != nil {
			ws.Logln(logrus.ErrorLevel, err.Error())
			return
		}
		if ok {
			return
		}
	}
	//初始化WxCron
	ws.WxCronService = wx_cron.NewWxCronService(wx_cron.SetWxCronServiceWxDao(ws.wxDao), wx_cron.SetWxCronServiceSourceDao(ws.sourceDao),
		wx_cron.SetSelf(self), wx_cron.SetBot(ws.Bot), wx_cron.SetWxCronGroups(groups), wx_cron.SetWxCronFriends(friends), wx_cron.SetWxCronServiceLog(ws.Logger))
	//初始化,批量更新userID
	//err = ws.ReloadAndUpdateUserName()
	//if err != nil {
	//	ws.Logln(logrus.ErrorLevel, err.Error())
	//	return err
	//}
	// llm功能
	ws.Process()
	ws.Reply()
	//ws.MessageUpdateUserName()
	err = ws.corn.AddFunc(config.Config.HolidayTips, ws.SendHolidayTips)
	if err != nil {
		return err
	}
	err = ws.corn.AddFunc(config.Config.NewsTips, ws.SendNews)
	if err != nil {
		return err
	}
	//err = c.AddFunc(config.Config.RegularUpdate, ws.RegularUpdate)
	//if err != nil {
	//	return err
	//}
	err = ws.corn.AddFunc(config.Config.RegularSendDailyProfit, ws.RegularSendDailyProfit)
	if err != nil {
		return err
	}
	ws.corn.Start()
	return nil
}

func (service *WxService) ReloadAndUpdateUserName() error {
	groups, err := service.getGroupUserNameToUserIDMap()
	if err != nil {
		return err
	}
	for k, v := range groups {
		users, err := service.wxDao.GetUsersByGroupName(k)
		if err != nil {
			return err
		}
		for _, user := range users {
			user.UserId = v[user.UserName]
			err = service.wxDao.UpdateUserID(user)
			if err != nil {
				return err
			}
		}
	}
	service.Logln(logrus.InfoLevel, "init Success")
	return nil
}

func (service *WxService) getGroupUserNameToUserIDMap() (map[string]map[string]string, error) {
	//群 群员
	usersMap := make(map[string]map[string]string)
	for _, g := range service.groups {
		nickName := strings.TrimSpace(g.NickName)
		if nickName == "" {
			continue
		}
		usersMap[nickName] = make(map[string]string)
		member, err := g.Members()
		if err != nil {
			return nil, err
		}
		for _, u := range member {
			disPlayName := strings.TrimSpace(u.DisplayName)
			if disPlayName == "" {
				continue
			}
			usersMap[nickName][disPlayName] = u.UserName
		}
	}
	return usersMap, nil
}
