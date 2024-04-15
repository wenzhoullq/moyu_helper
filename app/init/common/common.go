package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
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

func InitTool() error {
	ToolMap = map[string]string{
		"PDF": constant.PDF,
	}
	ToolReplySuf = map[string]string{
		"PDF": constant.PDFSuf,
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
