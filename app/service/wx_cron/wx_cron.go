package wx_cron

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"time"
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
	*client.TanshuClient
}

func NewWxCronService(ops ...func(*WxCronService)) *WxCronService {
	service := &WxCronService{
		TanshuClient: client.NewTanshuClient(),
	}
	for _, op := range ops {
		op(service)
	}
	return service
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
	holidayTipPre := "【摸鱼小助手】提醒您:各位摸鱼人上午好🌹！\n工作再累，一定不要忘记摸🐟哦！有事没事起身去茶水间、去厕所、去廊道走走，别老在工位上坐着，💴是老板的，但命是自己的！\n"
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
	holidayTip := holidayTipPre + holidayTipSuf
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

	todayNews, err := service.GetNews()
	if err != nil {
		service.Log(logrus.ErrorLevel, err)
		return
	}
	//添加金价新闻
	goldPrice, err := service.GetGoldPrice()
	if err != nil {
		service.Log(logrus.ErrorLevel, err)
		return
	}
	goldNews := fmt.Sprintf(constant.GoldPriceNews, goldPrice)
	todayNews = append([]string{goldNews}, todayNews...)
	for i, v := range todayNews {
		newsSuf += fmt.Sprintf("%d.%s。\n", i+1, v)
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
