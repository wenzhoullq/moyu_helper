package wx_llm

import (
	"errors"
	"github.com/eatmoreapple/openwechat"
	reply2 "weixin_LLM/dto/reply"
	"weixin_LLM/init/common"
	"weixin_LLM/lib"
	"weixin_LLM/lib/constant"
)

func (service *WxLLMService) toolsProducer(msg *openwechat.Message) error {
	content, err := service.MessagePreprocessing(msg)
	if err != nil {
		return err
	}
	content = lib.ProcessingCommands(content)
	if _, ok := common.ToolMap[content]; !ok {
		return errors.New("not toolsReq")
	}
	reply := &reply2.Reply{
		Message: msg,
		Content: constant.ReplyPre + common.ToolMap[content] + common.ToolReplySuf[content],
	}
	service.replyTextChan <- reply
	return nil
}
