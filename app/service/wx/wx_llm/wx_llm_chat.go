package wx_llm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strings"
	"weixin_LLM/dto/chat"
	group2 "weixin_LLM/dto/group"
	"weixin_LLM/dto/reply"
	"weixin_LLM/init/config"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) GetChatReq(content, key string) ([]*chat.ChatForm, error) {
	llmReq := make([]*chat.ChatForm, 0)
	chats, err := service.wxDao.GetString(key)
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}
	if chats != "" {
		tmp := make([]chat.ChatForm, 0)
		err = json.Unmarshal([]byte(chats), &tmp)
		if err != nil {
			return nil, err
		}
		for i := range tmp {
			llmReq = append(llmReq, &tmp[i])
		}
	}
	llmReq = append(llmReq, &chat.ChatForm{
		Role:    "user",
		Content: content,
	})
	return llmReq, nil
}

func (service *WxLLMService) StoreChat(key, resp string, llmReq []*chat.ChatForm) error {
	llmReq = append(llmReq, &chat.ChatForm{
		Role:    "assistant",
		Content: resp,
	})
	//最多只保存上一轮对话
	llmReq = llmReq[lib.Max(len(llmReq)-4, 0):]
	res, err := json.Marshal(&llmReq)
	if err != nil {
		return err
	}
	err = service.wxDao.SetString(key, string(res), constant.ChatExp)
	if err != nil {
		return err
	}
	return nil
}
func (service *WxLLMService) Forbid(content, modeType, key string, msg *openwechat.Message) (bool, error) {
	//dirtyWords := []string{"粗鲁", "不礼貌", "侮辱"}
	////脏话封禁
	//for _, word := range dirtyWords {
	//	if !strings.Contains(content, word) {
	//		continue
	//	}
	//	err := service.wxDao.SetString(key, msg.Content, constant.ForbidForProfanity)
	//	if err != nil {
	//		return false, err
	//	}
	//	reply := &reply.Reply{
	//		Content: service.ForbidChat[modeType],
	//		Message: msg,
	//	}
	//	service.replyTextChan <- reply
	//	return true, nil
	//}
	return false, nil
}

func (service *WxLLMService) AbilityIntroduce(msg *openwechat.Message) {
	content := constant.AbilitiesIntroduce
	for i, v := range config.Config.AbilityConfigure.Abilities {
		content = fmt.Sprintf("%s%d:%s\n", content, i+1, v)
	}
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: content,
	}
}

func (service *WxLLMService) SubscribeNews(msg *openwechat.Message) error {
	group, err := msg.Sender()
	if err != nil {
		return err
	}
	g, err := service.wxDao.GetGroupByName(group.NickName)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
		subscribeContent := &group2.Subscribe{
			News: true,
		}
		subscribe, err := json.Marshal(subscribeContent)
		if err != nil {
			return err
		}
		g = &group2.Groups{
			GroupName: group.UserName,
			Subscribe: string(subscribe),
		}
		err = service.wxDao.CreateGroup(g)
		if err != nil {
			return err
		}
	}
	subscribeContent := &group2.Subscribe{}
	err = json.Unmarshal([]byte(g.Subscribe), subscribeContent)
	if err != nil {
		return err
	}
	subscribeContent.News = true
	subscribe, err := json.Marshal(subscribeContent)
	g.Subscribe = string(subscribe)
	err = service.wxDao.UpdateGroup(g)
	if err != nil {
		return err
	}
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: constant.SubscribeNewsSuccess,
	}
	return nil
}
func (service *WxLLMService) UnSubscribeNews(msg *openwechat.Message) error {
	group, err := msg.Sender()
	if err != nil {
		return err
	}
	g, err := service.wxDao.GetGroupByName(group.NickName)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
		subscribeContent := &group2.Subscribe{
			News: false,
		}
		subscribe, err := json.Marshal(subscribeContent)
		if err != nil {
			return err
		}
		g = &group2.Groups{
			GroupName: group.UserName,
			Subscribe: string(subscribe),
		}
		err = service.wxDao.CreateGroup(g)
		if err != nil {
			return err
		}
	}
	subscribeContent := &group2.Subscribe{}
	err = json.Unmarshal([]byte(g.Subscribe), subscribeContent)
	if err != nil {
		return err
	}
	subscribeContent.News = false
	subscribe, err := json.Marshal(subscribeContent)
	g.Subscribe = string(subscribe)
	err = service.wxDao.UpdateGroup(g)
	if err != nil {
		return err
	}
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: constant.UnSubscribeNewsSuccess,
	}
	return nil
}
func (service *WxLLMService) SubscribeTips(msg *openwechat.Message) error {
	group, err := msg.Sender()
	if err != nil {
		return err
	}
	g, err := service.wxDao.GetGroupByName(group.NickName)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
		subscribeContent := &group2.Subscribe{
			Tips: true,
		}
		subscribe, err := json.Marshal(subscribeContent)
		if err != nil {
			return err
		}
		g = &group2.Groups{
			GroupName: group.UserName,
			Subscribe: string(subscribe),
		}
		err = service.wxDao.CreateGroup(g)
		if err != nil {
			return err
		}
	}
	subscribeContent := &group2.Subscribe{}
	err = json.Unmarshal([]byte(g.Subscribe), subscribeContent)
	if err != nil {
		return err
	}
	subscribeContent.Tips = true
	subscribe, err := json.Marshal(subscribeContent)
	g.Subscribe = string(subscribe)
	err = service.wxDao.UpdateGroup(g)
	if err != nil {
		return err
	}
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: constant.SubscribeMoyuTipsSuccess,
	}
	return nil
}
func (service *WxLLMService) UnSubscribeTips(msg *openwechat.Message) error {
	group, err := msg.Sender()
	if err != nil {
		return err
	}
	g, err := service.wxDao.GetGroupByName(group.NickName)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
		subscribeContent := &group2.Subscribe{
			Tips: false,
		}
		subscribe, err := json.Marshal(subscribeContent)
		if err != nil {
			return err
		}
		g = &group2.Groups{
			GroupName: group.UserName,
			Subscribe: string(subscribe),
		}
		err = service.wxDao.CreateGroup(g)
		if err != nil {
			return err
		}
	}
	subscribeContent := &group2.Subscribe{}
	err = json.Unmarshal([]byte(g.Subscribe), subscribeContent)
	if err != nil {
		return err
	}
	subscribeContent.Tips = false
	subscribe, err := json.Marshal(subscribeContent)
	g.Subscribe = string(subscribe)
	err = service.wxDao.UpdateGroup(g)
	if err != nil {
		return err
	}
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: constant.UnSubscribeMoyuTipsSuccess,
	}
	return nil
}
func (service *WxLLMService) groupSubscribeNews(msg *openwechat.Message) (bool, error) {
	f, ok := service.SubscribeNewsMap[msg.Content]
	if !ok {
		return false, nil
	}
	err := f(msg)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (service *WxLLMService) groupSubscribeTips(msg *openwechat.Message) (bool, error) {
	f, ok := service.SubscribeTipsMap[msg.Content]
	if !ok {
		return false, nil
	}
	err := f(msg)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (service *WxLLMService) redisKeyFriendDrawLotsMark(user *openwechat.User) string {
	return fmt.Sprintf(constant.FriendDrawLotsMark, user.UserName)
}

func (service *WxLLMService) friendAbilities(msg *openwechat.Message) (bool, error) {
	if msg.Content != "" && !strings.Contains(msg.Content, "菜单") && !strings.Contains(msg.Content, "功能") && !strings.HasSuffix(msg.Content, "我是") {
		return false, nil
	}
	service.AbilityIntroduce(msg)
	user, err := msg.Sender()
	if err != nil {
		return true, err
	}
	service.Logln(logrus.InfoLevel, user.NickName, "Abilities")
	return true, nil
}

func (service *WxLLMService) GroupAbilities(msg *openwechat.Message) (bool, error) {
	if msg.Content != "" && !strings.Contains(msg.Content, "菜单") && !strings.Contains(msg.Content, "功能") {
		return false, nil
	}
	service.AbilityIntroduce(msg)
	user, err := msg.SenderInGroup()
	if err != nil {
		return true, err
	}
	group, err := msg.Sender()
	if err != nil {
		return true, err
	}
	service.Logln(logrus.InfoLevel, group.UserName, user.NickName, "Abilities")
	return true, nil
}

func (service *WxLLMService) friendChat(msg *openwechat.Message) (bool, error) {
	user, err := msg.Sender()
	if err != nil {
		return true, err
	}
	err = service.NormalChatProcess(msg, user)
	if err != nil {
		return true, err
	}
	service.Logln(logrus.InfoLevel, user.NickName, " chat")
	return true, nil
}

func (service *WxLLMService) groupChat(msg *openwechat.Message) (bool, error) {
	user, err := msg.SenderInGroup()
	if err != nil {
		return true, err
	}
	//根据不同模式进行不同的对话
	key := service.getChatModeKey(user)
	value, err := service.wxDao.GetString(key)
	if err != nil {
		if err != redis.Nil {
			return true, err
		}
		value = constant.NorMalModeChat
	}
	err = service.GroupChatModel[value](msg, user)
	if err != nil {
		return true, err
	}
	service.Logln(logrus.InfoLevel, user.NickName, " chat")
	return true, nil
}

func (service *WxLLMService) ChatProcess(msg *openwechat.Message, user *openwechat.User) error {
	forbidKey := constant.Forbid + user.UserName
	value, err := service.wxDao.GetString(forbidKey)
	if err != nil {
		if err != redis.Nil {
			return err
		}
	}
	// 该用户还在被封禁
	if value != "" {
		return nil
	}
	return nil
}

func (service *WxLLMService) DoctorChatProcess(msg *openwechat.Message, user *openwechat.User) error {
	//forbidKey := constant.Forbid + user.UserName
	//value, err := service.wxDao.GetString(forbidKey)
	//if err != nil {
	//	if err != redis.Nil {
	//		return err
	//	}
	//}
	//// 该用户还在被封禁
	//if value != "" {
	//	return nil
	//}
	resp, err := service.AoJiaoClient.Chat(msg.Content)
	if err != nil {
		return err
	}
	//封禁
	//forbid, err := service.Forbid(resp, constant.AoJiaoModeChat, forbidKey, msg)
	//if err != nil {
	//	return err
	//}
	//if forbid {
	//	return nil
	//}
	service.replyTextChan <- &reply.Reply{
		Content: resp,
		Message: msg,
	}
	service.Logln(logrus.InfoLevel, user.NickName, "llmProcess success")
	return nil
}

func (service *WxLLMService) AoJiaoChatProcess(msg *openwechat.Message, user *openwechat.User) error {
	forbidKey := constant.Forbid + user.UserName
	value, err := service.wxDao.GetString(forbidKey)
	if err != nil {
		if err != redis.Nil {
			return err
		}
	}
	// 该用户还在被封禁
	if value != "" {
		return nil
	}
	resp, err := service.AoJiaoClient.Chat(msg.Content)
	if err != nil {
		return err
	}
	//封禁
	forbid, err := service.Forbid(resp, constant.AoJiaoModeChat, forbidKey, msg)
	if err != nil {
		return err
	}
	if forbid {
		return nil
	}
	service.replyTextChan <- &reply.Reply{
		Content: resp,
		Message: msg,
	}
	service.Logln(logrus.InfoLevel, user.NickName, "llmProcess success")
	return nil
}

func (service *WxLLMService) NormalChatProcess(msg *openwechat.Message, user *openwechat.User) error {
	forbidKey := constant.Forbid + user.UserName
	value, err := service.wxDao.GetString(forbidKey)
	if err != nil {
		if err != redis.Nil {
			return err
		}
	}
	// 该用户还在被封禁
	if value != "" {
		return nil
	}
	key := constant.ChatMark + user.UserName
	chatReq, err := service.GetChatReq(msg.Content, key)
	if err != nil {
		return err
	}
	resp, err := service.Ernie8KClient.Chat(chatReq)
	if err != nil {
		return err
	}
	//内容有违法,不接受该用户半小时的发言
	if resp.Flag != constant.RESP_NORMAL || resp.NeedClearHistory {
		err := service.wxDao.SetString(key, msg.Content, constant.ForbidForPolitics)
		if err != nil {
			return err
		}
		return errors.New("forbid for politics")
	}
	//封禁
	forbid, err := service.Forbid(resp.Result, constant.NorMalModeChat, forbidKey, msg)
	if err != nil {
		return err
	}
	if forbid {
		return nil
	}
	// 存入redis
	err = service.StoreChat(key, resp.Result, chatReq)
	if err != nil {
		return err
	}
	for len([]rune(resp.Result)) > constant.MaxAnswerLen {
		service.Logln(logrus.InfoLevel, user.NickName, resp.Result)
		_, err := service.wxDao.IncrKey(constant.ShortTime)
		if err != nil {
			return err
		}
		chatReq, err = service.GetChatReq(constant.Short+resp.Result, key)
		if err != nil {
			return err
		}
		resp, err = service.Ernie8KClient.Chat(chatReq)
		if err != nil {
			return err
		}
		// 存入redis
		err = service.StoreChat(key, resp.Result, chatReq)
		if err != nil {
			return err
		}
	}
	service.replyTextChan <- &reply.Reply{
		Content: resp.Result,
		Message: msg,
	}
	service.Logln(logrus.InfoLevel, user.NickName, "llmProcess success")
	return nil
}
