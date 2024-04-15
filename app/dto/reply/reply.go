package reply

import "github.com/eatmoreapple/openwechat"

type Reply struct {
	*openwechat.Message
	Content string
}
type ImgReply struct {
	*openwechat.Message
	Path string
	Key  string
}
