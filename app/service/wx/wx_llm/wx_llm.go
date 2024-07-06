package wx_llm

import (
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"os"
	"weixin_LLM/dao"
	"weixin_LLM/dto/reply"
	"weixin_LLM/init/common"
	"weixin_LLM/init/config"
	"weixin_LLM/lib/client"
	"weixin_LLM/lib/constant"
)

type WxLLMService struct {
	*client.Ernie8KClient
	*client.TxCloudClient
	*client.AoJiaoClient
	*logrus.Logger
	signChan           chan *openwechat.Message
	FriendDrawLotsChan chan *openwechat.Message
	GroupDrawLotsChan  chan *openwechat.Message
	//unDrawLotsChan      chan *openwechat.Message
	zodiacBlindBoxChan  chan *openwechat.Message
	friendTextToImgChan chan *openwechat.Message
	groupTextToImgChan  chan *openwechat.Message
	friendImgToImgChan  chan *openwechat.Message
	groupImgToImgChan   chan *openwechat.Message
	upgradeChan         chan *openwechat.Message
	replyTextChan       chan *reply.Reply
	replyImgChan        chan *reply.ImgReply
	updateChan          chan struct{}
	//生产者回调函数
	groupTextProducer  []func(*openwechat.Message) (bool, error)
	groupImgProducer   []func(*openwechat.Message) (bool, error)
	groupGameProducer  []func(*openwechat.Message) (bool, error)
	groupMarkProducer  []func(*openwechat.Message) (bool, error)
	friendTextProducer []func(*openwechat.Message) (bool, error)
	friendImgProducer  []func(*openwechat.Message) (bool, error)
	//对话模式
	GroupChatModel map[string]func(*openwechat.Message, *openwechat.User) error
	ForbidChat     map[string]string
	SignReward     map[int64]int
	wxDao          *dao.WxDao
	self           *openwechat.Self
	groups         openwechat.Groups
	friends        openwechat.Friends
	//signLock       *sync.Mutex
}

func NewWxLLMService(ops ...func(c *WxLLMService)) *WxLLMService {
	service := &WxLLMService{
		Ernie8KClient:      client.NewErnie8KClient(client.SetToken(common.Token)),
		TxCloudClient:      client.NewTxCloudClient(),
		AoJiaoClient:       client.NewAoJiaoClient(),
		signChan:           make(chan *openwechat.Message, constant.GameMaxNum),
		FriendDrawLotsChan: make(chan *openwechat.Message, constant.GameMaxNum),
		GroupDrawLotsChan:  make(chan *openwechat.Message, constant.GameMaxNum),
		//unDrawLotsChan:      make(chan *openwechat.Message, constant.GameMaxNum),
		zodiacBlindBoxChan:  make(chan *openwechat.Message, constant.GameMaxNum),
		upgradeChan:         make(chan *openwechat.Message, constant.GameMaxNum),
		friendTextToImgChan: make(chan *openwechat.Message, constant.ReplyPicMaxNum),
		groupTextToImgChan:  make(chan *openwechat.Message, constant.ReplyPicMaxNum),
		groupImgToImgChan:   make(chan *openwechat.Message, constant.ReplyPicMaxNum),
		friendImgToImgChan:  make(chan *openwechat.Message, constant.ReplyPicMaxNum),
		replyTextChan:       make(chan *reply.Reply, constant.ReplyMaxNum),
		replyImgChan:        make(chan *reply.ImgReply, constant.ReplyMaxNum),
		updateChan:          make(chan struct{}, constant.UpdateMaxNum),
		//signLock:            &sync.Mutex{},
	}
	//service.friendTextProducer = []func(*openwechat.Message) (bool, error){service.tools, service.friendImgToImgMark, service.friendTextToImg, service.friendDrawLots,service.friendAbilities, service.friendChat}
	service.friendTextProducer = []func(*openwechat.Message) (bool, error){service.tools, service.friendImgToImgMark, service.friendTextToImg, service.friendDrawLots, service.friendAbilities}
	service.friendImgProducer = []func(*openwechat.Message) (bool, error){service.friendImgToImg}
	//service.groupTextProducer = []func(*openwechat.Message) (bool, error){service.game, service.tools, service.groupMark, service.groupTextToImg, service.groupChat}
	service.groupTextProducer = []func(*openwechat.Message) (bool, error){service.game, service.tools, service.groupMark, service.groupTextToImg, service.GroupAbilities}
	service.groupImgProducer = []func(*openwechat.Message) (bool, error){service.groupImgToImg}
	service.groupGameProducer = []func(message *openwechat.Message) (bool, error){service.sign, service.upgrade, service.groupDrawLots, service.zodiacBlindBox}
	//service.groupGameProducer = []func(message *openwechat.Message) (bool, error){service.sign, service.upgrade, service.drawLots, service.unDrawLots, service.zodiacBlindBox}
	service.groupMarkProducer = []func(message *openwechat.Message) (bool, error){service.groupImgToImgMark, service.ModeChangeMark}
	service.GroupChatModel = map[string]func(*openwechat.Message, *openwechat.User) error{
		constant.NorMalModeChat: service.NormalChatProcess,
		constant.AoJiaoModeChat: service.AoJiaoChatProcess,
		constant.DoctorModeChat: service.DoctorChatProcess,
	}
	service.ForbidChat = map[string]string{
		constant.NorMalModeChat: constant.NorMalModeForbidDirty,
		constant.AoJiaoModeChat: constant.AoJiaoModelForbidDirty,
	}
	service.SignReward = map[int64]int{
		constant.SignFirst:  config.Config.SignRewardFirst,
		constant.SignSecond: config.Config.SignRewardSecond,
		constant.SignThird:  config.Config.SignRewardThird,
	}
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

func (service *WxLLMService) GetGroupTextProducer() []func(*openwechat.Message) (bool, error) {
	return service.groupTextProducer
}

func (service *WxLLMService) GetGroupImgProducer() []func(*openwechat.Message) (bool, error) {
	return service.groupImgProducer
}

func (service *WxLLMService) GetFriendTextProducer() []func(*openwechat.Message) (bool, error) {
	return service.friendTextProducer
}

func (service *WxLLMService) GetFriendImgProducer() []func(*openwechat.Message) (bool, error) {
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
					continue
				}
			case msg := <-service.upgradeChan:
				err := service.upgradeProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
			case msg := <-service.GroupDrawLotsChan:
				err := service.GroupDrawLotsProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
			//case msg := <-service.unDrawLotsChan:
			//	err := service.unDrawLotsProcess(msg)
			//	if err != nil {
			//		service.Logln(logrus.ErrorLevel, err.Error())
			//		continue
			//	}
			case msg := <-service.FriendDrawLotsChan:
				err := service.FriendDrawLotsProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
			case msg := <-service.zodiacBlindBoxChan:
				err := service.zodiacBlindBoxProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
			case msg := <-service.friendTextToImgChan:
				err := service.friendTextToImgProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}

			case msg := <-service.groupTextToImgChan:
				err := service.groupTextToImgProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}

			case msg := <-service.friendImgToImgChan:
				err := service.friendImgToImgProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
			case msg := <-service.groupImgToImgChan:
				err := service.groupImgToImgProcess(msg)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
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
				//打日志
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
					service.Logln(logrus.ErrorLevel, err.Error())
					return
				}
				_, err = imgReply.ReplyImage(file)
				if err != nil {
					service.Logln(logrus.ErrorLevel, err.Error())
					continue
				}
			}
		}
	}()
}

func (service *WxLLMService) checkGold(msg *openwechat.Message, goldConsume int) (bool, error) {
	user, err := msg.SenderInGroup()
	if err != nil {
		return true, err
	}
	group, err := msg.Sender()
	if err != nil {
		return true, err
	}
	u, err := service.wxDao.GetUserByUserNameAndGroupNameAndUserId(user.DisplayName, group.NickName, user.UserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if u.Reward < goldConsume {
		return false, nil
	}
	return true, nil
}
func (service *WxLLMService) deductionGold(msg *openwechat.Message, goldConsume int) error {
	//金币扣除
	user, err := msg.SenderInGroup()
	if err != nil {
		return err
	}
	group, err := msg.Sender()
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	u, err := service.wxDao.GetUserByUserNameAndGroupNameAndUserId(user.DisplayName, group.NickName, user.UserName)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	u.Reward = u.Reward - goldConsume
	err = service.wxDao.UpdateUserReward(u)
	if err != nil {
		service.Logln(logrus.ErrorLevel, err.Error())
		return err
	}
	//发送金币扣除通知
	service.replyTextChan <- &reply.Reply{
		Message: msg,
		Content: fmt.Sprintf(constant.GoldConsumeReply, goldConsume, u.Reward),
	}
	return nil
}
