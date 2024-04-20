package wx_llm

import (
	"github.com/eatmoreapple/openwechat"
	reply2 "weixin_LLM/dto/reply"
	"weixin_LLM/init/common"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) toolsProcess(msg *openwechat.Message) (bool, error) {
	if _, ok := common.ToolMap[msg.Content]; !ok {
		return false, nil
	}
	reply := &reply2.Reply{
		Message: msg,
		Content: constant.ReplyPre + common.ToolMap[msg.Content] + common.ToolReplySuf[msg.Content],
	}
	service.replyTextChan <- reply
	return true, nil
}
