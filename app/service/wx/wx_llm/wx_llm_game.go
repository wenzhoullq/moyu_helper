package wx_llm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"strings"
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

func (service *WxLLMService) unDrawLots(msg *openwechat.Message) (bool, error) {
	if msg.Content != constant.UnDrawLots {
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
	key := service.redisKeyGroupDrawLotsMark(user, group)
	_, err = service.wxDao.GetString(key)
	if err != nil {
		if err == redis.Nil {
			service.replyTextChan <- &reply.Reply{
				Message: msg,
				Content: constant.UnHasDrawLots,
			}
			return true, nil
		}
		return true, err
	}
	service.unDrawLotsChan <- msg
	return true, nil
}

func (service *WxLLMService) drawLots(msg *openwechat.Message) (bool, error) {
	if msg.Content != constant.DrawLots {
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
	key := service.redisKeyGroupDrawLotsMark(user, group)
	_, err = service.wxDao.GetString(key)
	if err != nil {
		if err == redis.Nil {
			service.drawLotsChan <- msg
			return true, nil
		}
		return true, err
	}
	//已经抽过签了
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: constant.HasDraw,
	}
	return true, nil
}

func (service *WxLLMService) redisKeyGroupZodiacBlindBox(user *openwechat.User, group *openwechat.User) string {
	return fmt.Sprintf(constant.GroupZodiacBlindBoxMark, group.UserName, user.UserName)
}

func (service *WxLLMService) zodiacBlindBox(msg *openwechat.Message) (bool, error) {
	if msg.Content != constant.ZodiacBlindBox {
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
	key := service.redisKeyGroupZodiacBlindBox(user, group)
	_, err = service.wxDao.GetString(key)
	if err != nil {
		if err == redis.Nil {
			service.zodiacBlindBoxChan <- msg
			return true, nil
		}
		return true, err
	}
	//已经抽过签了
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: constant.HasZodiacBlindBox,
	}
	return true, nil
}

func (service *WxLLMService) upgrade(msg *openwechat.Message) (bool, error) {
	if msg.Content != constant.Upgrade {
		return false, nil
	}
	ok, err := service.checkGold(msg, constant.LvUpConsume)
	if err != nil {
		return true, err
	}
	if !ok {
		reply := &reply.Reply{
			Message: msg,
			Content: fmt.Sprintf(constant.UpgradeApplicationFail, constant.LvUpConsume, constant.GoldGetTip),
		}
		service.replyTextChan <- reply
		return true, nil
	}
	service.upgradeChan <- msg
	return true, nil

}

func (service *WxLLMService) redisKeyGroupDrawLotsMark(user *openwechat.User, group *openwechat.User) string {
	return fmt.Sprintf(constant.GroupDrawLotsMark, group.UserName, user.UserName)
}
func (service *WxLLMService) unDrawLotsProcess(msg *openwechat.Message) error {
	user, err := msg.SenderInGroup()
	if err != nil {
		return err
	}
	group, err := msg.Sender()
	key := service.redisKeyGroupDrawLotsMark(user, group)
	value, err := service.wxDao.GetString(key)
	if err != nil {
		return err
	}
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.UnDrawLotsSuf, user.DisplayName, value),
	}
	return nil
}

func (service *WxLLMService) randomSend(path string, msg *openwechat.Message) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	// 如果文件夹为空，提前返回
	if len(files) == 0 {
		return "", err
	}
	// 随机选择一个文件
	randomIndex := rand.Intn(len(files))
	randomFile := files[randomIndex]
	fileName := randomFile.Name()
	// 输出随机选择的文件名
	service.replyImgChan <- &reply.ImgReply{
		Message: msg,
		Path:    path + fileName,
	}
	texts := strings.Split(fileName, ".jpg")
	return texts[0], nil
}

func (service *WxLLMService) zodiacBlindBoxProcess(msg *openwechat.Message) error {
	text, err := service.randomSend(config.Config.FileConfigure.ZodiacBlindBoxFile, msg)
	if err != nil {
		return err
	}
	user, err := msg.SenderInGroup()
	if err != nil {
		return err
	}
	group, err := msg.Sender()
	if err != nil {
		return err
	}
	//记录抽生肖盲盒
	key := service.redisKeyGroupZodiacBlindBox(user, group)
	err = service.wxDao.SetString(key, text, lib.SecondsUntilMidnight())
	if err != nil {
		return err
	}
	//发送抽中的盲盒信息
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.ZodiacBlindBoxSuf, text),
	}
	return nil
}

func (service *WxLLMService) drawLotsProcess(msg *openwechat.Message) error {
	//发送图片
	text, err := service.randomSend(config.Config.FileConfigure.DrawLotsFile, msg)
	if err != nil {
		return err
	}
	//记录抽签
	user, err := msg.SenderInGroup()
	if err != nil {
		return err
	}
	group, err := msg.Sender()
	if err != nil {
		return err
	}
	key := service.redisKeyGroupDrawLotsMark(user, group)
	err = service.wxDao.SetString(key, text, lib.SecondsUntilMidnight())
	if err != nil {
		return err
	}
	return nil
}

func (service *WxLLMService) upgradeProcess(msg *openwechat.Message) error {
	//service.signLock.Lock()
	//defer service.signLock.Unlock()
	group, err := msg.Sender()
	if err != nil {
		return err
	}
	user, err := msg.SenderInGroup()
	u, err := service.wxDao.GetUserByUserNameAndGroupNameAndUserId(user.DisplayName, group.NickName, user.UserName)
	if err != nil {
		return err
	}
	extra := &uu.Extra{}
	err = json.Unmarshal([]byte(u.Extra), extra)
	//等级 武力 幸运值均升级
	extra.Lv = extra.Lv + 1
	extra.Force = extra.Force + 1
	extra.Luck = extra.Luck + 1
	extraJson, err := json.Marshal(extra)
	if err != nil {
		return err
	}
	u.Extra = string(extraJson)
	if err = service.wxDao.UpdateUserExtra(u); err != nil {
		return err
	}
	if err = service.deductionGold(msg, constant.LvUpConsume); err != nil {
		return err
	}
	reply := &reply.Reply{
		Content: fmt.Sprintf(constant.StatusSuccessReply,
			user.DisplayName, extra.Lv, extra.Force, extra.Luck, u.Reward-constant.LvUpConsume),
		Message: msg,
	}
	service.replyTextChan <- reply
	service.Logln(logrus.InfoLevel, user.NickName, " upgradeProcess success")
	return nil
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
	//service.signLock.Lock()
	//defer service.signLock.Unlock()
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
	rank, err := service.wxDao.IncrKey(constant.SignMark + group.NickName + ":" + lib.GetCurYearAndMonthAndDay())
	if err != nil {
		return err
	}
	_, err = service.wxDao.Expire(constant.SignMark+group.NickName+":"+lib.GetCurYearAndMonthAndDay(), lib.SecondsUntilMidnight())
	if err != nil {
		return err
	}
	//签到奖励
	reward := config.Config.SignRewardElse
	if v, ok := service.SignReward[rank]; ok {
		reward = v
	}
	u, err := service.wxDao.GetUserByUserNameAndGroupNameAndUserId(user.DisplayName, group.NickName, user.UserName)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		extra := &uu.Extra{}
		extra.Lv = 1
		extra.SignTime = 1
		extraJson, err := json.Marshal(extra)
		if err != nil {
			return err
		}
		u = &uu.User{
			UserName:  user.DisplayName,
			UserId:    user.UserName,
			GroupName: group.NickName,
			Reward:    reward,
			Extra:     string(extraJson),
		}
		err = service.wxDao.AddUser(u)
		if err != nil {
			return err
		}
		//回复信息
		service.replyTextChan <- &reply.Reply{
			Content: fmt.Sprintf(constant.SignSuccessReply,
				user.DisplayName, rank, reward, u.Reward, extra.SignTime, lib.GetCurTime()),
			Message: msg,
		}
		service.Logln(logrus.InfoLevel, user.NickName, " signProcess success")
		return nil
	}
	extra := &uu.Extra{}
	err = json.Unmarshal([]byte(u.Extra), extra)
	extra.SignTime = extra.SignTime + 1
	extraJson, err := json.Marshal(extra)
	if err != nil {
		return err
	}
	u.Reward = u.Reward + reward
	u.Extra = string(extraJson)
	if err = service.wxDao.UpdateUserReward(u); err != nil {
		return err
	}
	if err = service.wxDao.UpdateUserExtra(u); err != nil {
		return err
	}
	//回复信息
	service.replyTextChan <- &reply.Reply{
		Content: fmt.Sprintf(constant.SignSuccessReply,
			user.DisplayName, rank, reward, u.Reward, extra.SignTime, lib.GetCurTime()),
		Message: msg,
	}
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
