package wx_llm

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
	reply2 "weixin_LLM/dto/reply"
	"weixin_LLM/init/config"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) friendTextToImg(msg *openwechat.Message) error {
	if !strings.HasPrefix(msg.Content, constant.TextToImgKeyWord) {
		return errors.New("not text To img")
	}
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: constant.ImgReplyFriend,
	}
	service.imgChan <- msg
	return nil
}

func (service *WxLLMService) groupTextToImg(msg *openwechat.Message) error {
	if !strings.HasPrefix(msg.Content, constant.TextToImgKeyWord) {
		return errors.New("not text To img")
	}
	user, err := msg.SenderInGroup()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	group, err := msg.Sender()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	//查看余额
	err = service.checkGold(msg, user, group)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.ImgReplyGroup, constant.ImgGoldConsume),
	}
	service.imgChan <- msg
	//扣除金币
	err = service.deductionGold(msg, user)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	return nil
}

func (service *WxLLMService) checkGold(msg *openwechat.Message, user *openwechat.User, group *openwechat.User) error {
	u, err := service.wxDao.GetUserByUserNameAndGroupNameAndUserId(user.DisplayName, group.NickName, user.UserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			reply := &reply2.Reply{
				Message: msg,
				Content: constant.TransToImgApplicationFail,
			}
			service.replyTextChan <- reply
			return nil
		}
		return err
	}
	if u.Reward < constant.ImgGoldConsume {
		reply := &reply2.Reply{
			Message: msg,
			Content: constant.TransToImgApplicationFail,
		}
		service.replyTextChan <- reply
		return errors.New("gold not enough")
	}
	return nil
}

func (service *WxLLMService) getImgToImgRedisKey(userName string) string {
	return fmt.Sprintf("%s%s", constant.ImgToImgMark, userName)
}

func (service *WxLLMService) groupImgToImgMark(msg *openwechat.Message) error {
	if !strings.HasPrefix(msg.Content, constant.ImgToImgKeyWord) {
		return errors.New("not imgTo img Req")
	}
	user, err := msg.SenderInGroup()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	group, err := msg.Sender()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	//查看余额
	err = service.checkGold(msg, user, group)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	err = service.imgToImgMark(user, msg)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	service.Logln(logrus.InfoLevel, "user:", user.DisplayName, " req user img to img")
	return nil
}

func (service *WxLLMService) imgToImgMark(user *openwechat.User, msg *openwechat.Message) error {
	content := strings.TrimSpace(msg.Content)
	value := ""
	contents := strings.Split(content, constant.ImgToImgKeyWord)
	if len(contents) > 0 {
		value = contents[len(contents)-1]
	}
	//存入redis标记
	service.wxDao.SetString(service.getImgToImgRedisKey(user.UserName), value, constant.ImgExp)
	reply := &reply2.Reply{
		Message: msg,
		Content: constant.ImgToImgApplicationSuccess,
	}
	service.replyTextChan <- reply
	return nil
}

func (service *WxLLMService) friendImgToImgMark(msg *openwechat.Message) error {
	if !strings.HasPrefix(msg.Content, constant.ImgToImgKeyWord) {
		return errors.New("not imgTo img Req")
	}
	user, err := msg.Sender()
	if err != nil {
		return err
	}
	err = service.imgToImgMark(user, msg)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	service.Logln(logrus.InfoLevel, "user:", user.DisplayName, " req user img to img")
	return nil
}

func (service *WxLLMService) imgToImgProducer(user *openwechat.User, msg *openwechat.Message) error {
	value, err := service.wxDao.GetString(service.getImgToImgRedisKey(user.UserName))
	if err != nil {
		return err
	}
	picHttp, err := msg.GetPicture()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(picHttp.Body)
	if err != nil {
		return err
	}
	picBase64 := base64.StdEncoding.EncodeToString(body)
	resp, err := service.TxCloudClient.PostImgToImg(picBase64, value)
	if err != nil {
		return err
	}
	//去除redis标记
	_, err = service.wxDao.DelString(service.getImgToImgRedisKey(user.UserName))
	if err != nil {
		return err
	}
	//回复图片
	msg.Content = resp.Response.ResultImage
	service.imgChan <- msg
	return nil
}

func (service *WxLLMService) friendImgToImgProducer(msg *openwechat.Message) error {
	user, err := msg.Sender()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	//回复正在生成图片的信息
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: constant.ImgReplyFriend,
	}
	err = service.imgToImgProducer(user, msg)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	return nil
}

func (service *WxLLMService) groupImgToImgProducer(msg *openwechat.Message) error {
	user, err := msg.SenderInGroup()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	//回复正在生成图片的信息
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.ImgReplyGroup, constant.ImgGoldConsume),
	}
	err = service.imgToImgProducer(user, msg)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	err = service.deductionGold(msg, user)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	return nil
}

func (service *WxLLMService) deductionGold(msg *openwechat.Message, user *openwechat.User) error {
	//金币扣除
	group, err := msg.Sender()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	u, err := service.wxDao.GetUserByUserNameAndGroupNameAndUserId(user.DisplayName, group.NickName, user.UserName)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	u.Reward = u.Reward - constant.ImgGoldConsume
	err = service.wxDao.UpdateUserByUserNameAndGroupNameAndUserId(u)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	//发送金币扣除通知
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.ImgGoldConsumeReply, u.Reward),
	}
	return nil
}

func (service *WxLLMService) transToImgProcess(msg *openwechat.Message) error {
	resp, err := service.TxCloudClient.PostTextToImg(msg.Content)
	if err != nil {
		return err
	}
	var fileName string
	fileName = fmt.Sprintf("%d.jpg", time.Now().Unix())
	data, err := base64.StdEncoding.DecodeString(resp.Response.ResultImage)
	if err != nil {
		return err
	}
	path := config.Config.FileConfigure.ImgFile + fileName
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	service.replyImgChan <- &reply2.ImgReply{
		Message: msg,
		Path:    path,
	}
	return nil
}
