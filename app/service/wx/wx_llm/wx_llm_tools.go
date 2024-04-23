package wx_llm

import (
	"github.com/eatmoreapple/openwechat"
	"strings"
	reply2 "weixin_LLM/dto/reply"
	"weixin_LLM/init/common"
	"weixin_LLM/init/config"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) tools(msg *openwechat.Message) (bool, error) {
	if _, ok := common.ToolMap[msg.Content]; !ok {
		return false, nil
	}
	service.replyTextChan <- &reply2.Reply{
		Message: msg,
		Content: constant.ReplyPre + common.ToolMap[msg.Content] + common.ToolReplySuf[msg.Content],
	}
	//发送饿了么城市卡
	if strings.Contains(msg.Content, constant.EleKeyWord) {
		service.replyImgChan <- &reply2.ImgReply{
			Message: msg,
			Path:    config.Config.ElePosterFile,
		}
	}
	return true, nil
}
