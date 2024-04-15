package wx_llm

import (
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
	"weixin_LLM/dao"
	"weixin_LLM/dto/reply"
	"weixin_LLM/init/common"
	"weixin_LLM/lib/client"
	"weixin_LLM/lib/constant"
)

type WxLLMService struct {
	*client.Ernie8KClient
	*client.TxCloudClient
	*logrus.Logger
	signChan      chan *openwechat.Message
	llmChan       chan *openwechat.Message
	imgChan       chan *openwechat.Message
	replyTextChan chan *reply.Reply
	replyImgChan  chan *reply.ImgReply
	updateChan    chan struct{}
	//生产者回调函数
	textProducer []func(*openwechat.Message) error
	imgProducer  []func(*openwechat.Message) error
	wxDao        *dao.WxDao
	self         *openwechat.Self
	groups       openwechat.Groups
	friends      openwechat.Friends
	signLock     *sync.Mutex
}

func NewWxLLMService(ops ...func(c *WxLLMService)) *WxLLMService {
	service := &WxLLMService{
		Ernie8KClient: client.NewErnie8KClient(client.SetToken(common.Token)),
		TxCloudClient: client.NewTxCloudClient(),
		signChan:      make(chan *openwechat.Message, constant.SignMaxNum),
		llmChan:       make(chan *openwechat.Message, constant.LlmMaxNum),
		imgChan:       make(chan *openwechat.Message, constant.ReplyPicMaxNum),
		replyTextChan: make(chan *reply.Reply, constant.ReplyMaxNum),
		replyImgChan:  make(chan *reply.ImgReply, constant.ReplyMaxNum),
		updateChan:    make(chan struct{}, constant.UpdateMaxNum),
		signLock:      &sync.Mutex{},
	}
	service.imgProducer = []func(*openwechat.Message) error{service.imgToImgProducer}
	service.textProducer = []func(*openwechat.Message) error{service.signProducer, service.toolsProducer, service.llmImgToImgReqProducer, service.llmTextToImgReqProducer, service.textToImgProducer, service.llmChatProducer}

	for _, op := range ops {
		op(service)
	}
	return service
}

func (service *WxLLMService) MessagePreprocessing(msg *openwechat.Message) (string, error) {
	if msg.IsSendBySelf() {
		return "", errors.New("msg send by self")
	}
	content := msg.Content
	if msg.IsSendByGroup() {
		if !strings.HasPrefix(msg.Content, constant.LlmKeyWord) {
			return "", errors.New("not llmOperate")
		}
		content = msg.Content[len(constant.LlmKeyWord):]
	}
	return strings.TrimSpace(content), nil
}

func SetWxDao(wxDao *dao.WxDao) func(ws *WxLLMService) {
	return func(ws *WxLLMService) {
		ws.wxDao = wxDao
	}
}

func (service *WxLLMService) GetTextProducer() []func(*openwechat.Message) error {
	return service.textProducer
}

func (service *WxLLMService) GetImgProducer() []func(*openwechat.Message) error {
	return service.imgProducer
}

func SetLog(log *logrus.Logger) func(*WxLLMService) {
	return func(wls *WxLLMService) {
		wls.Logger = log
	}
}

func SetSelf(self *openwechat.Self) func(*WxLLMService) {
	return func(wls *WxLLMService) {
		wls.self = self
	}
}

func SetFriends(friends openwechat.Friends) func(*WxLLMService) {
	return func(wls *WxLLMService) {
		wls.friends = friends
	}
}

func SetGroups(groups openwechat.Groups) func(*WxLLMService) {
	return func(wls *WxLLMService) {
		wls.groups = groups
	}
}

func (service *WxLLMService) OperateMsgWorker() {
	go func() {
		for {
			select {
			case msg := <-service.llmChan:
				err := service.llmChatProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
				}
			case msg := <-service.signChan:
				err := service.signProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
				}
			case msg := <-service.imgChan:
				err := service.transToImgProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
				}
			}
		}
	}()
}

func (service *WxLLMService) OperateReplyWorker() {
	go func() {
		for {
			select {
			case reply := <-service.replyTextChan:
				_, err := reply.ReplyText(reply.Content)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
				user, err := reply.Message.Sender()
				if reply.Message.IsSendByGroup() {
					user, err = reply.Message.SenderInGroup()
					if err != nil {
						continue
					}
				}
				service.Logln(logrus.InfoLevel, user.NickName, reply.Content)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
				}
			case imgReply := <-service.replyImgChan:
				file, err := os.Open(imgReply.Path)
				if err != nil {
					return
				}
				_, err = imgReply.ReplyImage(file)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
				user, err := imgReply.Message.Sender()
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
				if imgReply.Message.IsSendByGroup() {
					//去除redis标记
					_, err = service.wxDao.DelString(imgReply.Key)
					if err != nil {
						service.Logln(logrus.ErrorLevel, err.Error())
						continue
					}
					//金币扣除
					user, err = imgReply.Message.SenderInGroup()
					if err != nil {
						service.Logln(logrus.ErrorLevel, err.Error())
						continue
					}
					group, err := imgReply.Message.Sender()
					if err != nil {
						service.Logln(logrus.ErrorLevel, err.Error())
						continue
					}
					u, err := service.wxDao.GetUserByUserNameAndGroupName(user.DisplayName, group.NickName)
					if err != nil {
						service.Logln(logrus.ErrorLevel, err.Error())
						continue
					}
					u.Reward = u.Reward - constant.ImgGoldConsume
					err = service.wxDao.UpdateUserByUserNameAndGroupName(u)
					if err != nil {
						service.Logln(logrus.ErrorLevel, err.Error())
						continue
					}
					//发送金币扣除通知
					service.replyTextChan <- &reply.Reply{
						Message: imgReply.Message,
						Content: fmt.Sprintf(constant.ImgGoldConsumeReply, u.Reward),
					}
				}
				service.Logln(logrus.InfoLevel, user.NickName)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
				}
			}
		}
	}()
}
