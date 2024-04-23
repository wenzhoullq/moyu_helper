package wx_llm

import (
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
	"weixin_LLM/dto/reply"
	uu "weixin_LLM/dto/user"
	"weixin_LLM/init/config"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) game(msg *openwechat.Message) (bool, error) {
	for _, f := range service.groupGameProducer {
		ok, err := f(msg)
		if err != nil {
			return true, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func (service *WxLLMService) sign(msg *openwechat.Message) (bool, error) {
	if msg.Content != constant.SignKeyWord {
		return false, nil
	}
	service.signChan <- msg
	return true, nil
}

func (service *WxLLMService) signProcess(msg *openwechat.Message) error {
	if !msg.IsSendByGroup() {
		return errors.New("is not group message")
	}
	if msg.IsSendBySelf() {
		return errors.New("msg send by self")
	}
	service.signLock.Lock()
	defer service.signLock.Unlock()
	group, err := msg.Sender()
	if err != nil {
		return err
	}
	user, err := msg.SenderInGroup()
	if err != nil {
		return err
	}
	if user.DisplayName == "" {
		reply := &reply.Reply{
			Content: constant.SignFailReply,
			Message: msg,
		}
		service.replyTextChan <- reply
		service.updateChan <- struct{}{}
		return nil
	}
	// 查看是否签到
	res, err := service.wxDao.AddBit(user.UserName+":"+group.NickName+":"+lib.GetCurYearAndMonth(), lib.GetCurDay())
	if err != nil {
		return err
	}
	if res == constant.SignFail {
		reply := &reply.Reply{
			Content: fmt.Sprintf(constant.RepeatSignReply, user.DisplayName),
			Message: msg,
		}
		service.replyTextChan <- reply
		return nil
	}
	times, err := service.wxDao.IncrKey(constant.SignTime)
	service.Logln(logrus.InfoLevel, fmt.Sprintf("累计签到次数:%d次", times))
	if err != nil {
		return err
	}
	rank, err := service.wxDao.IncrKey(config.Config.SignMark + ":" + group.NickName + ":" + lib.GetCurYearAndMonthAndDay())
	if err != nil {
		return err
	}
	_, err = service.wxDao.Expire(config.Config.SignMark+":"+group.NickName+":"+lib.GetCurYearAndMonthAndDay(), lib.SecondsUntilMidnight())
	if err != nil {
		return err
	}
	//奖励
	var gold int
	switch rank {
	case 1:
		gold = config.Config.SignRewardFirst
		break
	case 2:
		gold = config.Config.SignRewardSecond
		break
	case 3:
		gold = config.Config.SignRewardThird
		break
	default:
		gold = config.Config.SignRewardElse
	}
	u, err := service.wxDao.GetUserByUserNameAndGroupNameAndUserId(user.DisplayName, group.NickName, user.UserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u = &uu.User{
				UserName:  user.DisplayName,
				UserId:    user.UserName,
				GroupName: group.NickName,
			}
			err = service.wxDao.AddUser(u)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	u.Reward = u.Reward + gold
	if err = service.wxDao.UpdateUser(u); err != nil {
		return err
	}
	reply := &reply.Reply{
		Content: fmt.Sprintf(constant.SignSuccessReply,
			user.DisplayName, rank, gold, u.Reward, lib.GetCurTime()),
		Message: msg,
	}
	service.replyTextChan <- reply
	service.Logln(logrus.InfoLevel, user.NickName, " signProcess success")
	return nil
}

func (service *WxLLMService) MessageUpdateUserName() {
	go func() {
		for {
			select {
			case <-service.updateChan:
				//60s内不断更新
				for i := 0; i < 60; i++ {
					time.Sleep(time.Second)
					err := service.self.UpdateMembersDetail()
					if err != nil {
						service.Logln(logrus.ErrorLevel, err.Error())
						continue
					}
				}
			}
		}
	}()
}
