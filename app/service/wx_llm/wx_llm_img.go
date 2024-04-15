package wx_llm

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
	reply2 "weixin_LLM/dto/reply"
	"weixin_LLM/init/config"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) llmTextToImgReqProducer(msg *openwechat.Message) error {
	content, err := service.MessagePreprocessing(msg)
	if err != nil {
		return err
	}
	if content != constant.TextToImgKeyWord {
		return errors.New("not imgReq")
	}
	user, err := msg.Sender()
	if err != nil {
		return err
	}
	if msg.IsSendByGroup() {
		user, err = msg.SenderInGroup()
		if err != nil {
			service.Logln(logrus.ErrorLevel, err.Error())
			return err
		}
		group, err := msg.Sender()
		if err != nil {
			service.Logln(logrus.ErrorLevel, err.Error())
			return err
		}
		u, err := service.wxDao.GetUserByUserNameAndGroupName(user.DisplayName, group.NickName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				reply := &reply2.Reply{
					Message: msg,
					Content: constant.TransToImgApplicationFail,
				}
				service.replyTextChan <- reply
				return nil
			}
			service.Logln(logrus.ErrorLevel, err.Error())
			return err
		}
		if u.Reward < constant.ImgGoldConsume {
			reply := &reply2.Reply{
				Message: msg,
				Content: constant.TransToImgApplicationFail,
			}
			service.replyTextChan <- reply
			return nil
		}
	}
	//存入redis标记
	service.wxDao.SetString(fmt.Sprintf("%s%s", constant.TextToImgMark, user.UserName), nil, constant.ImgExp)
	reply := &reply2.Reply{
		Message: msg,
		Content: constant.TextToImgApplicationSuccess,
	}
	service.replyTextChan <- reply
	service.Logln(logrus.InfoLevel, "user:", "", "req pic to pic")
	return nil
}

func (service *WxLLMService) textToImgProducer(msg *openwechat.Message) error {
	content, err := service.MessagePreprocessing(msg)
	if err != nil {
		return err
	}
	user, err := msg.Sender()
	if err != nil {
		return err
	}
	if msg.IsSendByGroup() {
		user, err = msg.SenderInGroup()
		if err != nil {
			return err
		}
	}
	_, err = service.wxDao.GetString(fmt.Sprintf("%s%s", constant.TextToImgMark, user.UserName))
	if err != nil {
		return err
	}
	resp, err := service.TxCloudClient.PostTextToImg(content)
	if err != nil {
		return err
	}
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.ImgReply, constant.ImgGoldConsume),
	}
	msg.Content = resp.Response.ResultImage
	service.imgChan <- msg
	return nil
}

func (service *WxLLMService) llmImgToImgReqProducer(msg *openwechat.Message) error {
	content, err := service.MessagePreprocessing(msg)
	if err != nil {
		return err
	}
	if content != constant.ImgToImgKeyWord {
		return errors.New("not imgReq")
	}
	user, err := msg.Sender()
	if err != nil {
		return err
	}
	if msg.IsSendByGroup() {
		user, err = msg.SenderInGroup()
		if err != nil {
			service.Logln(logrus.ErrorLevel, err.Error())
			return err
		}
		group, err := msg.Sender()
		if err != nil {
			service.Logln(logrus.ErrorLevel, err.Error())
			return err
		}
		u, err := service.wxDao.GetUserByUserNameAndGroupName(user.DisplayName, group.NickName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				reply := &reply2.Reply{
					Message: msg,
					Content: constant.TransToImgApplicationFail,
				}
				service.replyTextChan <- reply
				return nil
			}
			service.Logln(logrus.ErrorLevel, err.Error())
			return err
		}
		if u.Reward < constant.ImgGoldConsume {
			reply := &reply2.Reply{
				Message: msg,
				Content: constant.TransToImgApplicationFail,
			}
			service.replyTextChan <- reply
			return nil
		}
	}
	//存入redis标记
	service.wxDao.SetString(fmt.Sprintf("%s%s", constant.ImgToImgMark, user.UserName), nil, constant.ImgExp)
	reply := &reply2.Reply{
		Message: msg,
		Content: constant.ImgToImgApplicationSuccess,
	}
	service.replyTextChan <- reply
	service.Logln(logrus.InfoLevel, "user:", "", "req pic to pic")
	return nil
}
func (service *WxLLMService) imgToImgProducer(msg *openwechat.Message) error {
	user, err := msg.Sender()
	if err != nil {
		return err
	}
	if msg.IsSendByGroup() {
		user, err = msg.SenderInGroup()
		if err != nil {
			return err
		}
	}
	_, err = service.wxDao.GetString(fmt.Sprintf("%s%s", constant.ImgToImgMark, user.UserName))
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
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.ImgReply, constant.ImgGoldConsume),
	}
	resp, err := service.TxCloudClient.PostImgToImg(picBase64)
	if err != nil {
		return err
	}
	msg.Content = resp.Response.ResultImage
	service.imgChan <- msg
	return nil
}

func (service *WxLLMService) transToImgProcess(msg *openwechat.Message) error {
	var fileName string
	if msg.IsSendByGroup() {
		user, err := msg.SenderInGroup()
		if err != nil {
			return err
		}
		group, err := msg.Sender()
		if err != nil {
			return err
		}
		fileName = fmt.Sprintf("%s-%s-%d.jpg", group.NickName, user.NickName, time.Now().Unix())
	} else if msg.IsSendByFriend() {
		user, err := msg.Sender()
		if err != nil {
			return err
		}
		fileName = fmt.Sprintf("%s.jpg", user.NickName)
	} else {
		return errors.New("not send by friend or group")
	}
	data, err := base64.StdEncoding.DecodeString(msg.Content)
	if err != nil {
		return err
	}
	path := config.Config.FileConfigure.ImgFile + fileName
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	// 如果是
	user, err := msg.Sender()
	if msg.IsSendByGroup() {
		user, err = msg.SenderInGroup()
		if err != nil {
			service.Logln(logrus.ErrorLevel, err.Error())
			return err
		}
	}
	key := fmt.Sprintf("%s%s", constant.TextToImgMark, user.UserName)
	if msg.IsPicture() {
		key = fmt.Sprintf("%s%s", constant.ImgToImgMark, user.UserName)
	}
	service.replyImgChan <- &reply2.ImgReply{
		Message: msg,
		Path:    path,
		Key:     key,
	}
	return nil
}
