package wx_llm

import (
	"encoding/json"
	"github.com/eatmoreapple/openwechat"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strings"
	"weixin_LLM/dto/chat"
	"weixin_LLM/dto/reply"
	"weixin_LLM/dto/response"
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
func (service *WxLLMService) Forbid(resp *response.Ernie8kResponse, key string, msg *openwechat.Message) (bool, error) {
	//内容有违法,不接受该用户半小时的发言
	if resp.Flag != constant.RESP_NORMAL || resp.NeedClearHistory {
		err := service.wxDao.SetString(key, msg.Content, constant.ForbidForPolitics)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	dirtyWords := []string{"粗鲁", "不礼貌", "侮辱"}
	//脏话封禁
	for _, word := range dirtyWords {
		if !strings.Contains(resp.Result, word) {
			continue
		}
		err := service.wxDao.SetString(key, msg.Content, constant.ForbidForProfanity)
		if err != nil {
			return false, err
		}
		reply := &reply.Reply{
			Content: constant.ForbidDirty,
			Message: msg,
		}
		service.replyTextChan <- reply
		return true, nil
	}
	return false, nil
}

func (service *WxLLMService) friendChatProcess(msg *openwechat.Message) error {
	user, err := msg.Sender()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	err = service.chatProcess(msg, user)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	service.Logln(logrus.InfoLevel, user.NickName, " chat")
	return nil
}

func (service *WxLLMService) groupChatProcess(msg *openwechat.Message) error {
	user, err := msg.SenderInGroup()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	err = service.chatProcess(msg, user)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	service.Logln(logrus.InfoLevel, user.NickName, " chat")
	return nil
}

func (service *WxLLMService) chatProcess(msg *openwechat.Message, user *openwechat.User) error {
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
	chatReq, err := service.GetChatReq(constant.PreContent+msg.Content, key)
	if err != nil {
		return err
	}
	resp, err := service.Chat(chatReq)
	if err != nil {
		return err
	}
	//封禁
	forbid, err := service.Forbid(resp, forbidKey, msg)
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
		chatReq, err = service.GetChatReq(constant.Short+resp.Result, key)
		if err != nil {
			return err
		}
		resp, err = service.Chat(chatReq)
		if err != nil {
			return err
		}
		// 存入redis
		err = service.StoreChat(key, resp.Result, chatReq)
		if err != nil {
			return err
		}
	}
	reply := &reply.Reply{
		Content: resp.Result,
		Message: msg,
	}
	service.replyTextChan <- reply
	service.Logln(logrus.InfoLevel, user.NickName, "llmProcess success")
	return nil
}
