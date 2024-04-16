package wx_llm

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"os"
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
	imgChan       chan *openwechat.Message
	replyTextChan chan *reply.Reply
	replyImgChan  chan *reply.ImgReply
	updateChan    chan struct{}
	//生产者回调函数
	groupTextProducer  []func(*openwechat.Message) error
	groupImgProducer   []func(*openwechat.Message) error
	friendTextProducer []func(*openwechat.Message) error
	friendImgProducer  []func(*openwechat.Message) error

	wxDao    *dao.WxDao
	self     *openwechat.Self
	groups   openwechat.Groups
	friends  openwechat.Friends
	signLock *sync.Mutex
}

func NewWxLLMService(ops ...func(c *WxLLMService)) *WxLLMService {
	service := &WxLLMService{
		Ernie8KClient: client.NewErnie8KClient(client.SetToken(common.Token)),
		TxCloudClient: client.NewTxCloudClient(),
		signChan:      make(chan *openwechat.Message, constant.SignMaxNum),
		imgChan:       make(chan *openwechat.Message, constant.ReplyPicMaxNum),
		replyTextChan: make(chan *reply.Reply, constant.ReplyMaxNum),
		replyImgChan:  make(chan *reply.ImgReply, constant.ReplyMaxNum),
		updateChan:    make(chan struct{}, constant.UpdateMaxNum),
		signLock:      &sync.Mutex{},
	}
	service.friendTextProducer = []func(*openwechat.Message) error{service.toolsProcess, service.friendImgToImgMark, service.friendTextToImg, service.friendChatProcess}
	service.friendImgProducer = []func(*openwechat.Message) error{service.friendImgToImgProducer}
	service.groupTextProducer = []func(*openwechat.Message) error{service.signProducer, service.toolsProcess, service.groupImgToImgMark, service.groupTextToImg, service.groupChatProcess}
	service.groupImgProducer = []func(*openwechat.Message) error{service.groupImgToImgProducer}
	for _, op := range ops {
		op(service)
	}
	return service
}

func SetWxDao(wxDao *dao.WxDao) func(ws *WxLLMService) {
	return func(ws *WxLLMService) {
		ws.wxDao = wxDao
	}
}

func (service *WxLLMService) GetGroupTextProducer() []func(*openwechat.Message) error {
	return service.groupTextProducer
}

func (service *WxLLMService) GetGroupImgProducer() []func(*openwechat.Message) error {
	return service.groupImgProducer
}

func (service *WxLLMService) GetFriendTextProducer() []func(*openwechat.Message) error {
	return service.friendTextProducer
}

func (service *WxLLMService) GetFriendImgProducer() []func(*openwechat.Message) error {
	return service.friendImgProducer
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

func (service *WxLLMService) Process() {
	go func() {
		for {
			select {
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

func (service *WxLLMService) Reply() {
	go func() {
		for {
			select {
			case textReply := <-service.replyTextChan:
				_, err := textReply.ReplyText(textReply.Content)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
				user, err := textReply.Message.Sender()
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
				if textReply.Message.IsSendByGroup() {
					user, err = textReply.Message.SenderInGroup()
					if err != nil {
						continue
					}
				}
				service.Logln(logrus.InfoLevel, user.NickName, textReply.Content)
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
					user, err = imgReply.Message.SenderInGroup()
					if err != nil {
						continue
					}
				}
				service.Logln(logrus.InfoLevel, user.NickName)
			}
		}
	}()
}
