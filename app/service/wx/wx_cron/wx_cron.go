package wx_cron

import (
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
	"weixin_LLM/dao"
	"weixin_LLM/dto/group"
	holiday2 "weixin_LLM/dto/holiday"
	"weixin_LLM/init/common"
	"weixin_LLM/init/config"
	"weixin_LLM/lib"
	"weixin_LLM/lib/client"
	"weixin_LLM/lib/constant"
)

type WxCronService struct {
	*openwechat.Bot
	*logrus.Logger
	*client.ZhiHuClient
	*client.EleUnionClient
	*client.MeiTuanUnionClient
	*client.DidiUnionClient
	groups    openwechat.Groups
	wxDao     *dao.WxDao
	sourceDao *dao.SourceDao
	self      *openwechat.Self
	friends   openwechat.Friends
}

func NewWxCronService(ops ...func(*WxCronService)) *WxCronService {
	service := &WxCronService{
		ZhiHuClient:        client.NewZhiHuClient(),
		EleUnionClient:     client.NewEleUnionClient(client.SetEleUnionAppKey(config.Config.EleConfigure.AppKey), client.SetEleUnionSecret(config.Config.EleConfigure.Secret)),
		MeiTuanUnionClient: client.NewMeiTuanUnionClient(client.SetMeiTuanUnionAppKey(config.Config.MeiTuanConfigure.AppKey), client.SetMeiTuanUnionApiToken(config.Config.MeiTuanConfigure.ApiToken)),
		DidiUnionClient:    client.NewDidiUnionClient(client.SetDidiUnionAppKey(config.Config.DiDiConfigure.AppKey), client.SetDidiUnionAccessKey(config.Config.DiDiConfigure.AccessKey)),
	}
	for _, op := range ops {
		op(service)
	}
	return service
}

func SetWxCronFriends(friends openwechat.Friends) func(*WxCronService) {
	return func(wls *WxCronService) {
		wls.friends = friends
	}
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
	//节假日
	if len(common.Holidays) > 0 && today == common.Holidays[0].Date && common.Holidays[0].IsOffDay {
		return false
	}
	// 周六周末
	if weekday == time.Saturday || weekday == time.Sunday {
		if len(common.Holidays) > 0 && today == common.Holidays[0].Date && !common.Holidays[0].IsOffDay {
			return true
		}
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
		return
	}
	//查看是否订阅HolidayTips

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
	ad := common.AdMap[time.Now().Weekday()]
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
	//发送前更新
	err := service.self.UpdateMembersDetail()
	if err != nil {
		service.Log(logrus.ErrorLevel, err.Error())
		return
	}
	groups, err := service.wxDao.GetGroupList()
	if err != nil {
		service.Log(logrus.ErrorLevel, err)
		return
	}
	groupMap := make(map[string]*group.Groups)
	for _, v := range groups {
		scribe := &group.Subscribe{}
		err = json.Unmarshal([]byte(v.Subscribe), scribe)
		if err != nil {
			service.Log(logrus.ErrorLevel, err)
			continue
		}
		if !scribe.Tips {
			continue
		}
		groupMap[v.GroupName] = v
	}
	for _, group := range service.groups {
		//?到底是什么name
		if _, ok := groupMap[group.UserName]; !ok {
			//dd
			continue
		}
		_, err = group.SendText(holidayTip)
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
		newsSuf += fmt.Sprintf("%d.%s\n %s \n", i+1, v.Target.Title, v.Target.URL)
	}
	ad := common.AdMap[time.Now().Weekday()]
	news := newsPre + ad + newsSuf
	//发送前更新
	//err = service.self.UpdateMembersDetail()
	//if err != nil {
	//	service.Log(logrus.ErrorLevel, err.Error())
	//	return
	//}
	groups, err := service.wxDao.GetGroupList()
	if err != nil {
		service.Log(logrus.ErrorLevel, err)
		return
	}
	groupMap := make(map[string]*group.Groups)
	for _, v := range groups {
		scribe := &group.Subscribe{}
		err = json.Unmarshal([]byte(v.Subscribe), scribe)
		if err != nil {
			service.Log(logrus.ErrorLevel, err)
			continue
		}
		if !scribe.News {
			continue
		}
		groupMap[v.GroupName] = v
	}
	for _, group := range service.groups {
		//NickName:群昵称
		if _, ok := groupMap[group.NickName]; !ok {
			continue
		}
		_, err = group.SendText(news)
		if err != nil {
			service.Log(logrus.ErrorLevel, err.Error())
			return
		}
	}
	return
}

func (service *WxCronService) RegularUpdate() {
	//err := service.RegularUpdateUserName()
	//if err != nil {
	//	service.Logln(logrus.ErrorLevel, err.Error())
	//}
	err := service.RegularSource()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
	}
	err = service.RegularUpdateGroup()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
	}

}

// 更新订阅情况
func (service *WxCronService) RegularUpdateGroup() error {
	groups, err := service.self.Groups()
	if err != nil {
		return err
	}
	service.groups = groups
	for _, v := range groups {
		_, err = service.wxDao.GetGroupByName(v.NickName)
		if gorm.IsRecordNotFoundError(err) {
			//不存在则新增
			scribe := &group.Subscribe{
				News: true,
				Tips: true,
			}
			scribeStr, err := json.Marshal(&scribe)
			if err != nil {
				return err
			}
			group := &group.Groups{
				GroupName:  v.NickName,
				Subscribe:  string(scribeStr),
				CreateTime: time.Now(),
			}
			err = service.wxDao.CreateGroup(group)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (service *WxCronService) RegularSource() error {
	sources, err := service.sourceDao.GetNotExpSources()
	if err != nil {
		return err
	}
	for i := range sources {
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

func (service *WxCronService) RegularSendDailyProfit() {
	//计算didi收益
	didiResp, err := service.DidiUnionClient.GetTodayProfit()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return
	}
	var didiProfit float64
	for _, v := range didiResp.Data.OrderList {
		//单位为分
		didiProfit += float64(v.CPAProfit)/100 + float64(v.CPSProfit)/100
	}
	//计算美团收益
	var meiTuanProfit = 0.0
	meiTuanResp, err := service.MeiTuanUnionClient.GetTodayProfit()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return
	}
	for _, v := range meiTuanResp.DataList {
		f, err := strconv.ParseFloat(v.Profit, 64)
		if err != nil {
			service.Logln(logrus.ErrorLevel, err.Error())
			continue
		}
		meiTuanProfit += f
	}
	totalProfit := meiTuanProfit + didiProfit
	//发送每日收益
	user := service.friends.GetByNickName(constant.SendDailyProfitUser)
	_, err = user.SendText(fmt.Sprintf(constant.DailyProfit, meiTuanProfit, didiProfit, totalProfit))
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return
	}
	return
}
