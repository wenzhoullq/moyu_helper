package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	dao2 "weixin_LLM/dao"
	"weixin_LLM/dto/holiday"
	"weixin_LLM/init/config"
	"weixin_LLM/lib/client"
	"weixin_LLM/lib/constant"
)

var Token string
var Holidays []*holiday.Day
var ToolMap map[string]string
var ToolReplySuf map[string]string
var KeyMap map[string]string
var AdMap map[time.Weekday]string
var ChatModeMap map[string]struct{}
var ImgModeMap map[string]string

func InitCommon() error {
	client := client.NewQianFanClient()
	token, err := client.GetToken()
	if err != nil {
		return err
	}
	Token = token
	err = InitHoliday(fmt.Sprintf("%s%d.json", config.Config.HolidayFile, time.Now().Year()))
	if err != nil {
		return err
	}
	err = InitTool()
	if err != nil {
		return err
	}
	err = InitKeyMap()
	if err != nil {
		return err
	}
	err = InitAd()
	if err != nil {
		return err
	}
	err = InitModeMap()
	if err != nil {
		return err
	}
	return nil
}

func InitModeMap() error {
	ChatModeMap = map[string]struct{}{
		constant.NorMalModeChat: {},
		constant.AoJiaoModeChat: {},
	}
	ImgModeMap = map[string]string{
		constant.Unlimited:              constant.UnlimitedMark,
		constant.InkWash:                constant.InkWashMark,
		constant.ConceptualArt:          constant.ConceptualArtMark,
		constant.OilPainting1:           constant.OilPainting1Mark,
		constant.OilPainting2:           constant.OilPainting2Mark,
		constant.Watercolor:             constant.WatercolorMark,
		constant.PixelPainting:          constant.PixelPaintingMark,
		constant.ThickCoating:           constant.ThickCoatingMark,
		constant.Illustration:           constant.IllustrationMark,
		constant.PaperCuttings:          constant.PaperCuttingsMark,
		constant.Impressionism1:         constant.Impressionism1Mark,
		constant.Impressionism2:         constant.Impressionism2Mark,
		constant.D25:                    constant.D25Mark,
		constant.D3:                     constant.D3Mark,
		constant.ClassicalPortrait:      constant.ClassicalPortraitMark,
		constant.BlackAndWhiteSketching: constant.BlackAndWhiteSketchingMark,
		constant.Cyberpunk:              constant.CyberpunkMark,
		constant.ScienceFiction:         constant.ScienceFictionMark,
		constant.Dark:                   constant.DarkMark,
		constant.SteamWave:              constant.SteamWaveMark,
		constant.JapaneseAnime:          constant.JapaneseAnimeMark,
		constant.Monster:                constant.MonsterMark,
		constant.BeautifulAncient:       constant.BeautifulAncientMark,
		constant.RetroAnime:             constant.RetroAnimeMark,
		constant.GameCartoon:            constant.GameCartoonMark,
		constant.Universal:              constant.UniversalMark,
	}
	return nil
}

func InitKeyMap() error {
	KeyMap = map[string]string{
		constant.TxSecretId:  os.Getenv(constant.TxSecretId),
		constant.TxSecretKey: os.Getenv(constant.TxSecretKey),
	}
	return nil
}

func InitAd() error {
	AdMap = map[time.Weekday]string{
		time.Sunday:    constant.WeekendAd,
		time.Monday:    constant.DailyAd,
		time.Tuesday:   constant.DailyAd,
		time.Wednesday: constant.WednesdayAd,
		time.Thursday:  constant.DailyAd,
		time.Friday:    constant.WeekendAd,
		time.Saturday:  constant.WeekendAd,
	}
	return nil
}

func InitTool() error {
	dao := dao2.NewSourceDao()
	sources, err := dao.GetNotExpSources()
	if err != nil {
		return err
	}
	ToolMap = make(map[string]string)
	ToolReplySuf = make(map[string]string)
	for _, s := range sources {
		ToolMap[s.SourceName] = s.SourceLink
		ToolReplySuf[s.SourceName] = s.SourceDesc
	}
	return nil
}

func InitHoliday(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	jsonStr, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	holiday := &holiday.Holiday{}
	err = json.Unmarshal(jsonStr, holiday)
	if err != nil {
		return err
	}
	Holidays = holiday.Days
	cur := time.Now().Format("2006-01-02")
	for _, v := range holiday.Days {
		if v.Date < cur {
			Holidays = Holidays[1:]
			continue
		}
		break
	}
	return nil
}
