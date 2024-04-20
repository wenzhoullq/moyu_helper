package wx_llm

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
	reply2 "weixin_LLM/dto/reply"
	"weixin_LLM/init/config"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) friendTextToImg(msg *openwechat.Message) (bool, error) {
	if !strings.HasPrefix(msg.Content, constant.TextToImgKeyWord) {
		return false, nil
	}
	user, err := msg.Sender()
	if err != nil {
		return true, err
	}
	err = service.DailyFreeTimeCheck(user, msg)
	if err != nil {
		return true, err
	}
	//回复正在生成中
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: constant.ImgReplyFriend,
	}
	service.friendTextToImgChan <- msg
	return true, nil
}

func (service *WxLLMService) groupTextToImg(msg *openwechat.Message) (bool, error) {
	if !strings.HasPrefix(msg.Content, constant.TextToImgKeyWord) {
		return false, nil
	}
	user, err := msg.SenderInGroup()
	if err != nil {
		return true, err
	}
	group, err := msg.Sender()
	if err != nil {
		return true, err
	}
	//查看余额
	err = service.checkGold(msg, user, group)
	if err != nil {
		return true, err
	}
	//回复正在生成中
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.ImgReplyGroup, constant.ImgGoldConsume),
	}
	service.groupTextToImgChan <- msg
	return true, nil
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
			return gorm.ErrRecordNotFound
		}
		return err
	}
	if u.Reward < constant.ImgGoldConsume {
		reply := &reply2.Reply{
			Message: msg,
			Content: constant.TransToImgApplicationFail,
		}
		service.replyTextChan <- reply
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (service *WxLLMService) getImgToImgRedisKey(userName string) string {
	return fmt.Sprintf("%s%s", constant.ImgToImgMark, userName)
}

func (service *WxLLMService) groupImgToImgMark(msg *openwechat.Message) (bool, error) {
	if !strings.HasPrefix(msg.Content, constant.ImgToImgKeyWord) {
		return false, nil
	}
	user, err := msg.SenderInGroup()
	if err != nil {
		return true, err
	}
	group, err := msg.Sender()
	if err != nil {
		return true, err
	}
	//查看余额
	err = service.checkGold(msg, user, group)
	if err != nil {
		return true, err
	}
	//打标记
	err = service.imgToImgMark(user, msg)
	if err != nil {
		return true, err
	}
	service.Logln(logrus.InfoLevel, "user:", user.DisplayName, " request user img to img")
	return true, nil
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

func (service *WxLLMService) redisKeyFriendImgToImgMark(user *openwechat.User) string {
	return fmt.Sprintf("%s%s", constant.FriendImgToImgMark, user.UserName)
}

func (service *WxLLMService) DailyFreeTimeCheck(user *openwechat.User, msg *openwechat.Message) error {
	//查看免费额度使
	key := service.redisKeyFriendImgToImgMark(user)
	times, err := service.wxDao.IncrKey(key)
	if err != nil && err != redis.Nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	_, err = service.wxDao.Expire(key, lib.SecondsUntilMidnight())
	if err != nil && err != redis.Nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	//如果大于最大次数,则退出
	if times > constant.DailyMAXFreeImgTransTime {
		service.replyTextChan <- &reply2.Reply{
			Message: msg,
			Content: constant.ExDailyMAXFreeImgTransTimeReply,
		}
		return errors.New("transToImg times more than daily")
	}
	return nil
}

func (service *WxLLMService) friendImgToImgMark(msg *openwechat.Message) (bool, error) {
	if !strings.HasPrefix(msg.Content, constant.ImgToImgKeyWord) {
		return false, nil
	}
	user, err := msg.Sender()
	if err != nil {
		return true, err
	}
	// 查看免费额度
	err = service.DailyFreeTimeCheck(user, msg)
	if err != nil {
		return true, err
	}
	//打标记
	err = service.imgToImgMark(user, msg)
	if err != nil {
		return true, err
	}
	service.Logln(logrus.InfoLevel, "user:", user.DisplayName, " request user img to img")
	return true, err
}

func (service *WxLLMService) imgToImg(msg *openwechat.Message, key, value string) error {
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
	_, err = service.wxDao.DelString(key)
	if err != nil {
		return err
	}
	msg.Content = resp.Response.ResultImage
	if err != nil {
		return err
	}
	return nil
}

func (service *WxLLMService) friendImgToImgProducer(msg *openwechat.Message) (bool, error) {
	user, err := msg.Sender()
	if err != nil {
		return true, err
	}
	key := service.getImgToImgRedisKey(user.UserName)
	value, err := service.wxDao.GetString(key)
	if err != nil {
		if err == redis.Nil {
			return true, nil
		}
		return true, err
	}
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: constant.ImgReplyFriend,
	}
	err = service.imgToImg(msg, key, value)
	if err != nil {
		return true, err
	}
	//处理图片
	service.friendImgToImgChan <- msg
	return true, err
}

func (service *WxLLMService) groupImgToImgProducer(msg *openwechat.Message) (bool, error) {
	user, err := msg.SenderInGroup()
	if err != nil {
		return true, err
	}
	key := service.getImgToImgRedisKey(user.UserName)
	value, err := service.wxDao.GetString(key)
	if err != nil {
		return true, err
	}
	//发送通知
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.ImgReplyGroup, constant.ImgGoldConsume),
	}
	//保存图片
	err = service.imgToImg(msg, key, value)
	if err != nil {
		return true, err
	}
	//进入下一段处理
	service.groupImgToImgChan <- msg
	return true, nil
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

func (service *WxLLMService) imgProcess(msg *openwechat.Message) error {
	var fileName string
	fileName = fmt.Sprintf("%d.jpg", time.Now().Unix())
	data, err := base64.StdEncoding.DecodeString(msg.Content)
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

func (service *WxLLMService) textToImgProcess(msg *openwechat.Message) error {
	resp, err := service.TxCloudClient.PostTextToImg(msg.Content)
	if err != nil {
		return err
	}
	msg.Content = resp.Response.ResultImage
	err = service.imgProcess(msg)
	if err != nil {
		return err
	}
	return nil
}

func (service *WxLLMService) groupTextToImgProcess(msg *openwechat.Message) error {
	err := service.textToImgProcess(msg)
	if err != nil {
		return err
	}
	//扣除金币
	user, err := msg.SenderInGroup()
	if err != nil {
		return err
	}
	err = service.deductionGold(msg, user)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	return nil
}

func (service *WxLLMService) friendTextToImgProcess(msg *openwechat.Message) error {
	err := service.textToImgProcess(msg)
	if err != nil {
		return err
	}
	return nil
}

func (service *WxLLMService) friendImgToImgProcess(msg *openwechat.Message) error {
	err := service.imgProcess(msg)
	if err != nil {
		return err
	}
	return nil
}
func (service *WxLLMService) groupImgToImgProcess(msg *openwechat.Message) error {
	err := service.imgProcess(msg)
	if err != nil {
		return err
	}
	user, err := msg.SenderInGroup()
	if err != nil {
		return err
	}
	//扣除金币
	err = service.deductionGold(msg, user)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	//发送
	return nil
}
