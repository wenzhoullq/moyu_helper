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

func (service *WxLLMService) signProducer(msg *openwechat.Message) error {
	if msg.Content != constant.SignKeyWord {
		return errors.New("not signOperator")
	}
	service.signChan <- msg
	return nil
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
	u, err := service.wxDao.GetUserByUserNameAndGroupName(user.DisplayName, group.NickName)
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

func (service *WxLLMService) RegularUpdateUserName() {
	go func() {
		for {
			time.Sleep(time.Second * constant.RegularUpdate)
			select {
			default:

				err := service.updateDBNickName()

				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
				}
			}
		}
	}()
}

func (service *WxLLMService) updateDBNickName() error {
	defer func() {
		if err := recover(); err != nil {
			service.Logln(logrus.PanicLevel, "panic:", err)
		}
	}()
	service.signLock.Lock()
	defer service.signLock.Unlock()
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
			user, err := service.wxDao.GetUserByUserNameAndGroupName(v, g.NickName)
			if err != nil {
				return err
			}
			user.UserName = newMap[g.NickName][k]
			err = service.wxDao.UpdateUserName(user)
			if err != nil {
				return err
			}
		}
	}
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

func (service *WxLLMService) getGroupUserIDToUserNameMap() (map[string]map[string]string, error) {
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

func (service *WxLLMService) getGroupUserMapFromDB() (map[string]map[string]string, error) {
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
