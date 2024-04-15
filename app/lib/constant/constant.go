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
	ImgToImgKeyWord  = "图片风格转换"
	TextToImgKeyWord = "文字生成图片"
)

const (
	PreContent                  = "用50个字以内回答:\n"
	Short                       = "把如下内容缩短至50个字以内:\n"
	SignFailReply               = "╭┈┈⚠️失败⚠️┈┈╮\n❗ 请设置群昵称后签到\n╰┈┈⚠️失败⚠️┈┈╯"
	SignSuccessReply            = "╭┈┈┈🏡签到🏡┈┈┈╮\n🌼昵称：%s\n🐾排名：第%d名\n💰奖励：%d金币\n📆累积：%d金币\n⏰签到时间: %s\n╰┈┈┈┈┈┈┈┈┈┈┈╯"
	RepeatSignReply             = "╭┈┈┈🏡签到🏡┈┈┈╮\n🌼昵称：%s\n🐾状态：今日已签到\n╰┈┈┈🏡签到🏡┈┈┈╯"
	System                      = "你是汤鸽科技集团有限公司研发的摸鱼小助手,旨在各位摸鱼人进行摸鱼和薅羊毛,@摸鱼小助手并输入薅羊毛获得各种福利"
	ImgToImgApplicationSuccess  = "好的,请发出图片,我将图片进行风格转换"
	TextToImgApplicationSuccess = "好的,请用文字叙述你想要的图片,我将按您的描述生成相应的图片"
	TransToImgApplicationFail   = "抱歉,您的金币不足15枚,无法进行AI转换。\n若今天无签到,则可通过签到的方式获得金币"
	ImgReply                    = "图片正在生成中,请稍等...\n本次预计消耗金币%d枚"
	ImgGoldConsumeReply         = "图片转换成功,剩余金币数量:%d枚"
	ReplyPre                    = "你想要的是否是:\n"
	PDFSuf                      = "\n全能免费PDF工具，无广告，随时处理PDF。"
)
const (
	SignSuccess = 0
	SignFail    = 1
)

const (
	MaxAnswerLen = 50
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
	ImgExp             = 60
	ImgGoldConsume     = 15
	RegularUpdate      = 60 * 60
)

const (
	PDF = "https://tools.pdf24.org/zh/"
)

const (
	TxSecretId    = "TENCENTCLOUD_SECRET_ID"
	TxSecretKey   = "TENCENTCLOUD_SECRET_KEY"
	RegionShanhai = "ap-shanghai"
)
