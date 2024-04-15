package wx_llm

import (
	"encoding/json"
	"errors"
	"github.com/eatmoreapple/openwechat"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strings"
	"weixin_LLM/dto/chat"
	"weixin_LLM/dto/reply"
	"weixin_LLM/dto/response"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) GetLlmReq(content, key string) ([]*chat.ChatForm, error) {
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
	//脏话封禁
	if strings.Contains(resp.Result, "粗鲁") || strings.Contains(resp.Result, "不礼貌") || strings.Contains(resp.Result, "侮辱") {
		err := service.wxDao.SetString(key, msg.Content, constant.ForbidForProfanity)
		if err != nil {
			return false, err
		}
		reply := &reply.Reply{
			Content: "善言结善缘,恶语伤人心。你这一句话我需要花60秒来治愈自己😭😭",
			Message: msg,
		}
		service.replyTextChan <- reply
		return true, nil
	}
	return false, nil
}
func (service *WxLLMService) llmChatProducer(msg *openwechat.Message) error {
	if !msg.IsText() {
		return errors.New("not chat req")

	}
	content, err := service.MessagePreprocessing(msg)
	if err != nil {
		return err
	}
	msg.Content = content
	service.llmChan <- msg
	service.Logln(logrus.InfoLevel, "send to llmChan")
	return nil
}
func (service *WxLLMService) llmChatProcess(msg *openwechat.Message) error {
	user, err := msg.Sender()
	if err != nil {
		return err
	}
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
	service.Logln(logrus.InfoLevel, user.NickName, " chat")
	key := constant.ChatMark + user.UserName
	llmReq, err := service.GetLlmReq(constant.PreContent+msg.Content, key)
	if err != nil {
		return err
	}
	resp, err := service.Chat(llmReq)
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
	err = service.StoreChat(key, resp.Result, llmReq)
	if err != nil {
		return err
	}
	for len([]rune(resp.Result)) > constant.MaxAnswerLen {
		llmReq, err = service.GetLlmReq(constant.Short+resp.Result, key)
		if err != nil {
			return err
		}
		resp, err = service.Chat(llmReq)
		if err != nil {
			return err
		}
		// 存入redis
		err = service.StoreChat(key, resp.Result, llmReq)
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
