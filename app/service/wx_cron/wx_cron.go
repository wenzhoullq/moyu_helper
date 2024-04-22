package wx_cron

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"time"
	"weixin_LLM/dao"
	holiday2 "weixin_LLM/dto/holiday"
	"weixin_LLM/init/common"
	"weixin_LLM/lib"
	"weixin_LLM/lib/client"
	"weixin_LLM/lib/constant"
)

type WxCronService struct {
	*openwechat.Bot
	*logrus.Logger
	groups openwechat.Groups
	*client.ZhiHuClient
	wxDao     *dao.WxDao
	sourceDao *dao.SourceDao
	self      *openwechat.Self
}

func NewWxCronService(ops ...func(*WxCronService)) *WxCronService {
	service := &WxCronService{
		ZhiHuClient: client.NewZhiHuClient(),
	}
	for _, op := range ops {
		op(service)
	}
	return service
}

func SetSelf(self *openwechat.Self) func(*WxCronService) {
	return func(wls *WxCronService) {
		wls.self = self
	}
}

func SetWxCronServiceWxDao(wxDao *dao.WxDao) func(ws *WxCronService) {
	return func(ws *WxCronService) {
		ws.wxDao = wxDao
	}
}

func SetWxCronServiceSourceDao(sourceDao *dao.SourceDao) func(ws *WxCronService) {
	return func(ws *WxCronService) {
		ws.sourceDao = sourceDao
	}
}

func SetWxCronServiceLog(log *logrus.Logger) func(*WxCronService) {
	return func(ws *WxCronService) {
		ws.Logger = log
	}
}

func SetWxCronGroups(groups openwechat.Groups) func(*WxCronService) {
	return func(ws *WxCronService) {
		ws.groups = groups
	}
}

func SetBot(bot *openwechat.Bot) func(service *WxCronService) {
	return func(ws *WxCronService) {
		ws.Bot = bot
	}
}

func (service *WxCronService) isWorkDay() bool {
	today := time.Now().Format("2006-01-02")
	weekday := time.Now().Weekday()
	//今天是日历上的一天，但调休
	if len(common.Holidays) > 0 && today == common.Holidays[0].Date && common.Holidays[0].IsOffDay {
		return false
	}
	// 今天不是周六周末
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	return true
}

// 先判断最近的周六是不是在holiday之前
func (service *WxCronService) SendHolidayTips() {
	today := time.Now().Format("2006-01-02")
	//更新
	if len(common.Holidays) > 0 && common.Holidays[0].Date < today {
		common.Holidays = common.Holidays[1:]
	}
	// 判断是否是工作日
	if !service.isWorkDay() {
		service.Log(logrus.InfoLevel, "today is ", today, " ,this day is not workday")
		return
	}
	//下一个休息日数组
	nextRestDays := make([]*holiday2.Day, 0)
	// 获得最近的一个周六
	nextSatuday := lib.NextSaturday()
	//下一个周六小于下一个节假日
	if len(common.Holidays) > 0 && nextSatuday.Format("2006-01-02") < common.Holidays[0].Date {
		nextRestDays = append(nextRestDays, &holiday2.Day{
			Name: "周六",
			Date: nextSatuday.Format("2006-01-02"),
		})
	}
	holidaySet := make(map[string]struct{})
	for _, v := range common.Holidays {
		if !v.IsOffDay {
			continue
		}
		if _, ok := holidaySet[v.Name]; ok {
			continue
		}
		holidaySet[v.Name] = struct{}{}
		nextRestDays = append(nextRestDays, v)
	}
	holidayTipPre := constant.HolidayTip
	ad := common.AdMap[time.Wednesday]
	holidayTipSuf := ""
	for i, v := range nextRestDays {
		diffDay, err := lib.CalDays(today, v.Date)
		if err != nil {
			service.Log(logrus.ErrorLevel, err.Error())
			return
		}
		holidayTipSuf += fmt.Sprintf("离%s还有:%d天。", v.Name, diffDay)
		if i < len(nextRestDays)-1 {
			holidayTipSuf += "\n"
		}
	}
	holidayTip := holidayTipPre + ad + holidayTipSuf
	for _, group := range service.groups {
		_, err := group.SendText(holidayTip)
		if err != nil {
			service.Log(logrus.ErrorLevel, err.Error())
			return
		}
	}
	return
}

func (service *WxCronService) SendNews() {
	// 判断是否是工作日
	if !service.isWorkDay() {
		return
	}
	newsPre := constant.NewsSuf
	newsSuf := ""

	hotTopic, err := service.GetHotTopic()
	if err != nil {
		service.Log(logrus.ErrorLevel, err)
		return
	}
	for i, v := range hotTopic.Data {
		if i >= constant.MaxNewsNum {
			break
		}
		newsSuf += fmt.Sprintf("%d.%s   %s🔥。\n %s \n", i+1, v.Target.Title, v.DetailText, v.Target.URL)
	}
	news := newsPre + newsSuf
	for _, group := range service.groups {
		_, err := group.SendText(news)
		if err != nil {
			service.Log(logrus.ErrorLevel, err.Error())
			return
		}
	}
	return
}

func (service *WxCronService) RegularUpdate() {
	err := service.RegularUpdateUserName()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
	}
	err = service.RegularSource()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
	}
}
func (service *WxCronService) RegularSource() error {
	sources, err := service.sourceDao.GetNotExpSources()
	if err != nil {
		return err
	}
	for i, _ := range sources {
		s := sources[i]
		exp, err := lib.TimeHasExp(s.SourceExp)
		if err != nil {
			return err
		}
		if !exp {
			continue
		}
		//去除Map
		delete(common.ToolMap, s.SourceName)
		delete(common.ToolReplySuf, s.SourceName)
		//更新db
		s.Status = constant.SourceExp
		err = service.sourceDao.UpdateSource(s)
		if err != nil {
			return err
		}

	}
	return nil
}
func (service *WxCronService) RegularUpdateUserName() error {
	defer func() {
		if err := recover(); err != nil {
			service.Logln(logrus.PanicLevel, "panic:", err)
		}
	}()
	err := service.self.UpdateMembersDetail()
	if err != nil {
		return err
	}
	preMap, err := service.getGroupUserMapFromDB()
	if err != nil {
		return err
	}
	newMap, err := service.getGroupUserIDToUserNameMap()
	if err != nil {
		return err
	}
	//查看是否修改过昵称
	for _, g := range service.groups {
		for k, v := range preMap[g.NickName] {
			if v == newMap[g.NickName][k] {
				continue
			}
			// 修改昵称
			user, err := service.wxDao.GetUserByUserNameAndGroupNameAndUserId(v, g.NickName, k)
			if err != nil {
				service.Logln(logrus.ErrorLevel, err.Error())
				return err
			}
			user.UserName = newMap[g.NickName][k]
			err = service.wxDao.UpdateUserName(user)
			if err != nil {
				return err
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}
func (service *WxCronService) getGroupUserMapFromDB() (map[string]map[string]string, error) {
	//群 群员
	usersMap := make(map[string]map[string]string)
	for _, g := range service.groups {
		if g.NickName == "" {
			continue
		}
		usersMap[g.NickName] = make(map[string]string)
		users, err := service.wxDao.GetUsersByGroupName(g.NickName)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			usersMap[g.NickName][u.UserId] = u.UserName
		}
	}
	return usersMap, nil
}
func (service *WxCronService) getGroupUserIDToUserNameMap() (map[string]map[string]string, error) {
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
			usersMap[g.NickName][u.UserName] = u.DisplayName
		}
	}
	return usersMap, nil
}
