package service

import (
	"errors"
	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"strings"
	"weixin_LLM/dao"
	"weixin_LLM/init/config"
	"weixin_LLM/init/log"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
	"weixin_LLM/service/wx_cron"
	"weixin_LLM/service/wx_llm"
)

type WxService struct {
	*openwechat.Bot
	*logrus.Logger
	*wx_llm.WxLLMService
	*wx_cron.WxCronService
	wxDao      *dao.WxDao
	sourceDao  *dao.SourceDao
	groupSend  []func(*openwechat.Message) error
	friendSend []func(*openwechat.Message) error
	groups     openwechat.Groups
}

func NewWxService() *WxService {
	ws := &WxService{
		Bot:       openwechat.DefaultBot(openwechat.Desktop),
		Logger:    log.Logger,
		wxDao:     dao.NewWxDao(),
		sourceDao: dao.NewSourceDao(),
	}
	ws.groupSend = []func(*openwechat.Message) error{ws.groupTextMsg, ws.groupImgMsg}
	ws.friendSend = []func(*openwechat.Message) error{ws.friendTextMsg, ws.friendImgMsg}
	return ws
}

func (ws *WxService) groupSender(msg *openwechat.Message) error {
	if !msg.IsSendByGroup() {
		return errors.New("not send by group")
	}
	if msg.IsSendBySelf() {
		return errors.New("msg send by self")
	}
	for _, f := range ws.groupSend {
		if f(msg) == nil {
			return nil
		}
	}
	return errors.New("no such group request")
}

func (ws *WxService) friendSender(msg *openwechat.Message) error {
	if !msg.IsSendByFriend() {
		return errors.New("not send by friend")
	}
	if msg.IsSendBySelf() {
		return errors.New("msg send by self")
	}
	for _, f := range ws.friendSend {
		if f(msg) == nil {
			return nil
		}
	}
	return errors.New("no such friend request")
}

func (ws *WxService) groupTextMsg(msg *openwechat.Message) error {
	if !msg.IsText() {
		return errors.New("not text")
	}
	msg.Content = strings.TrimSpace(msg.Content)
	if strings.HasPrefix(msg.Content, constant.LlmKeyWord) {
		msg.Content = strings.TrimSpace(msg.Content[len(constant.LlmKeyWord):])
	} else if strings.HasSuffix(msg.Content, constant.LlmKeyWord) {
		msg.Content = strings.TrimSpace(msg.Content[:len(msg.Content)-len(constant.LlmKeyWord)])
	} else {
		return errors.New("not groupTextMsg")
	}
	//把小写全部转为大写
	msg.Content = lib.ProcessingCommands(msg.Content)
	for _, f := range ws.WxLLMService.GetGroupTextProducer() {
		if f(msg) == nil {
			return nil
		}
	}
	return errors.New("no such group text request")
}

func (ws *WxService) groupImgMsg(msg *openwechat.Message) error {
	if !msg.IsPicture() {
		return errors.New("not pic")
	}
	for _, f := range ws.WxLLMService.GetGroupImgProducer() {
		if f(msg) == nil {
			return nil
		}
	}
	return errors.New("no such group img request")
}

func (ws *WxService) friendTextMsg(msg *openwechat.Message) error {
	if msg.IsSendBySelf() {
		return errors.New("msg send by self")
	}
	if !msg.IsText() {
		return errors.New("not text")
	}
	msg.Content = strings.TrimSpace(msg.Content)
	//把小写全部转为大写
	msg.Content = lib.ProcessingCommands(msg.Content)
	for _, f := range ws.WxLLMService.GetFriendTextProducer() {
		if f(msg) == nil {
			return nil
		}
	}
	return errors.New("no such group img request")
}

func (ws *WxService) friendImgMsg(msg *openwechat.Message) error {
	if !msg.IsPicture() {
		return errors.New("not pic")
	}
	for _, f := range ws.WxLLMService.GetFriendImgProducer() {
		if f(msg) == nil {
			return nil
		}
	}
	return errors.New("no such user img request")
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

	//初始化WxLLM
	ws.WxLLMService = wx_llm.NewWxLLMService(wx_llm.SetSelf(self), wx_llm.SetGroups(groups), wx_llm.SetFriends(friends),
		wx_llm.SetLog(ws.Logger), wx_llm.SetWxDao(ws.wxDao))
	//注册消息处理函数
	ws.MessageHandler = func(msg *openwechat.Message) {
		user, err := msg.Sender()
		if err != nil {
			ws.Logln(logrus.ErrorLevel, err.Error())
			return
		}
		ws.Logln(logrus.InfoLevel, "user:", user.NickName, " msgContent:", msg.Content)
		//对于不同的消息进行不同的处理
		err = ws.friendSender(msg)
		if err == nil {
			return
		}
		err = ws.groupSender(msg)

	}
	//初始化WxCron
	ws.WxCronService = wx_cron.NewWxCronService(wx_cron.SetWxCronServiceWxDao(ws.wxDao), wx_cron.SetWxCronServiceSourceDao(ws.sourceDao),
		wx_cron.SetSelf(self), wx_cron.SetBot(ws.Bot), wx_cron.SetWxCronGroups(groups), wx_cron.SetWxCronServiceLog(ws.Logger))
	//初始化,批量更新userID
	err = ws.ReloadAndUpdateUserName()
	if err != nil {
		ws.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	// llm功能
	ws.Process()
	ws.Reply()
	ws.MessageUpdateUserName()

	// cron功能
	c := cron.New()
	err = c.AddFunc(config.Config.HolidayTips, ws.SendHolidayTips)
	if err != nil {
		return err
	}
	err = c.AddFunc(config.Config.NewsTips, ws.SendNews)
	if err != nil {
		return err
	}
	err = c.AddFunc(config.Config.RegularUpdate, ws.RegularUpdate)
	if err != nil {
		return err
	}
	c.Start()
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
	return nil
}

func (service *WxService) getGroupUserNameToUserIDMap() (map[string]map[string]string, error) {
	//群 群员
	usersMap := make(map[string]map[string]string)
	for _, g := range service.groups {
		if g.NickName == "" {
			continue
		}
		usersMap[g.NickName] = make(map[string]string)
		member, err := g.Members()
		if err != nil {
			return nil, err
		}
		for _, u := range member {
			if u.DisplayName == "" {
				continue
			}
			usersMap[g.NickName][u.DisplayName] = u.UserName
		}
	}
	return usersMap, nil
}
