package constant

const (
	RESP_NORMAL = 0
)

const (
	SignMaxNum     = 100
	LlmMaxNum      = 100
	ReplyMaxNum    = 300
	ReplyPicMaxNum = 3
	UpdateMaxNum   = 300
)

const (
	LlmKeyWord       = "@摸鱼小助手"
	SignKeyWord      = "签到"
	ImgToImgKeyWord  = "风格转换"
	TextToImgKeyWord = "画"
)

const (
	PreContent                 = "用50个字以内回答:\n"
	Short                      = "把如下内容缩短至50个字以内:\n"
	SignFailReply              = "╭┈┈⚠️失败⚠️┈┈╮\n❗ 请设置群昵称后签到\n╰┈┈⚠️失败⚠️┈┈╯"
	SignSuccessReply           = "╭┈┈┈🏡签到🏡┈┈┈╮\n🌼昵称：%s\n🐾排名：第%d名\n💰奖励：%d金币\n📆累积：%d金币\n⏰签到时间: %s\n╰┈┈┈┈┈┈┈┈┈┈┈╯"
	RepeatSignReply            = "╭┈┈┈🏡签到🏡┈┈┈╮\n🌼昵称：%s\n🐾状态：今日已签到\n╰┈┈┈🏡签到🏡┈┈┈╯"
	System                     = "你是汤鸽科技集团有限公司研发的摸鱼小助手,旨在各位摸鱼人进行摸鱼和薅羊毛,@摸鱼小助手并输入薅羊毛获得各种福利"
	ImgToImgApplicationSuccess = "好的,请在60秒内请发出图片,我将图片进行风格转换"
	TransToImgApplicationFail  = "抱歉,您的金币不足15枚,无法进行AI绘画。\n若今天无签到,则可通过签到的方式获得金币。"
	ImgReplyGroup              = "图片正在生成中,请稍等...\n本次预计消耗金币%d枚"
	ImgReplyFriend             = "图片正在生成中,请稍等..."
	ImgGoldConsumeReply        = "图片生成成功!\n剩余金币数量:%d枚"
	ReplyPre                   = "你想要的是否是:\n"
	PDFSuf                     = "\n全能免费PDF工具，无广告，随时处理PDF。"
	FlowSuf                    = "\n2月1日-4月30日，移动用户登陆页面可领取2张2GB流量日包券，1张7天有效，1张30天有效，每个用户每月可领1次。。"
	NewsSuf                    = "【摸鱼小助手】提醒您:三点几了饮茶先啦🥤。\n这里是今天的摸鱼小新闻,祝各位摸鱼人摸鱼愉快！\n"
	GoldPriceNews              = "今日黄金价格:%s元/克"
	ForbidDirty                = "善言结善缘,恶语伤人心。你这一句话我需要花60秒来治愈自己😭😭"
)
const (
	SignSuccess = 0
	SignFail    = 1
)

const (
	MaxAnswerLen = 50
	MaxCharLen   = 20000
)

const (
	NewsNum       = "4"
	TanshuSuccess = 1
)

const (
	SignTime           = "signTimesCnt"
	ChatMark           = "chat:"
	ChatExp            = 60 * 2
	Forbid             = "forbid:"
	ForbidForPolitics  = 60 * 30
	ForbidForProfanity = 60
	ImgToImgMark       = "imgToImg:"
	TextToImgMark      = "textToImg:"
	ImgExp             = 60 * 2
	ImgGoldConsume     = 15
	RegularUpdate      = 60 * 60
)

const (
	PDF  = "https://tools.pdf24.org/zh/"
	Flow = "https://wx.10086.cn/qwhdhub/leadin/1024013102?A_C_CODE=Q0NXXsZMFT&channelId=P00000016916#/"
)

const (
	TxSecretId    = "TENCENTCLOUD_SECRET_ID"
	TxSecretKey   = "TENCENTCLOUD_SECRET_KEY"
	RegionShanhai = "ap-shanghai"
)
