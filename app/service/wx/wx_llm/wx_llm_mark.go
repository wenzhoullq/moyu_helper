package wx_llm

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"weixin_LLM/dto/reply"
	"weixin_LLM/init/common"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) getChatModeKey(user *openwechat.User) string {
	return fmt.Sprintf("%s:%s", constant.ChatMode, user.UserName)
}
func (service *WxLLMService) getImgModeKey(user *openwechat.User) string {
	return fmt.Sprintf("%s:%s", constant.ImgMode, user.UserName)
}
func (service *WxLLMService) ModeChangeMark(msg *openwechat.Message) (bool, error) {
	user, err := msg.SenderInGroup()
	if err != nil {
		return true, err
	}
	//对话模型
	if _, ok := common.ChatModeMap[msg.Content]; ok {
		key := service.getChatModeKey(user)
		//将标记插入redis
		err = service.wxDao.SetString(key, msg.Content, constant.ModeExp)
		if err != nil {
			return true, err
		}
		//发送切换成功
		service.replyTextChan <- &reply.Reply{
			Message: msg,
			Content: fmt.Sprintf(constant.ModeChatSet, msg.Content),
		}
		return true, nil
	}
	//生图模型
	if v, ok := common.ImgModeMap[msg.Content]; ok {
		key := service.getImgModeKey(user)
		//将标记插入redis
		err = service.wxDao.SetString(key, v, constant.ModeExp)
		if err != nil {
			return true, err
		}
		//发送切换成功
		service.replyTextChan <- &reply.Reply{
			Message: msg,
			Content: fmt.Sprintf(constant.ModeImgSet, msg.Content),
		}
		return true, nil
	}
	return false, nil
}
