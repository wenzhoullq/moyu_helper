package constant

const (
	RESP_NORMAL = 0
)

const (
	SignMaxNum     = 100
	ReplyMaxNum    = 300
	ReplyPicMaxNum = 3
	UpdateMaxNum   = 300
)

const (
	LlmKeyWord        = "@摸鱼小助手"
	SignKeyWord       = "签到"
	ImgToImgKeyWord   = "风格转换"
	TextToImgKeyWord  = "画"
	ModeChangeKeyWord = "切换"
)

const (
	PreContent                      = "用50个字以内回答:\n"
	Short                           = "把如下内容缩短至50个字以内:\n"
	SignFailReply                   = "╭┈┈⚠️失败⚠️┈┈╮\n❗ 请设置群昵称后签到\n╰┈┈⚠️失败⚠️┈┈╯"
	SignSuccessReply                = "╭┈┈┈🏡签到🏡┈┈┈╮\n🌼昵称：%s\n🐾排名：第%d名\n💰奖励：%d金币\n📆累积：%d金币\n⏰签到时间: %s\n╰┈┈┈┈┈┈┈┈┈┈┈╯"
	RepeatSignReply                 = "╭┈┈┈🏡签到🏡┈┈┈╮\n🌼昵称：%s\n🐾状态：今日已签到\n╰┈┈┈🏡签到🏡┈┈┈╯"
	System                          = "你是汤鸽科技集团有限公司研发的摸鱼小助手,旨在各位摸鱼人进行摸鱼和薅羊毛,@摸鱼小助手并输入薅羊毛获得各种福利"
	ImgToImgApplicationSuccess      = "好的,请在60秒内请发出图片,我将图片进行风格转换"
	TransToImgApplicationFail       = "抱歉,您的金币不足15枚,无法进行AI绘画。\n若今天无签到,则可通过签到的方式获得金币。"
	ImgReplyGroup                   = "图片正在生成中,请稍等...\n本次预计消耗金币%d枚"
	ImgReplyFriend                  = "图片正在生成中,请稍等..."
	ImgGoldConsumeReply             = "图片生成成功!\n剩余金币数量:%d枚"
	ReplyPre                        = "你想要的是否是:\n"
	NewsSuf                         = "【摸鱼小助手】提醒您:三点几了饮茶先啦🥤。\n这里是今天的摸鱼小新闻,祝各位摸鱼人摸鱼愉快！\n"
	GoldPriceNews                   = "今日黄金价格:%s元/克"
	NorMalModeForbidDirty           = "善言结善缘,恶语伤人心。你这一句话我需要花60秒来治愈自己😭😭"
	AoJiaoModelForbidDirty          = "你说话太没有礼貌了,我不想跟你说话了!"
	ExDailyMAXFreeImgTransTimeReply = "今日免费生图功能次数已用完,请明日再来使用"
	HolidayTip                      = "【摸鱼小助手】提醒您:各位摸鱼人上午好🌹！\n工作再累，一定不要忘记摸🐟！有事没事起身去茶水间、去厕所、去廊道走走，别老在工位上坐着，💴是老板的，但命是自己的！\n"
	WednesdayAd                     = "今日领取美团红包必得9元红包。\n"
	ModeChatSet                     = "已切换为%s模式"
	ModeChatSetFail                 = "无该模式"
)
const (
	SignSuccess = 0
	SignFail    = 1
)

const (
	NorMalModeChat = "正常"
	AoJiaoModeChat = "傲娇"
)

const (
	MaxAnswerLen = 50
)

const (
	MaxNewsNum = 5
)

const (
	SignTime                 = "signTimesCnt"
	ChatMark                 = "chat:"
	ChatExp                  = 60 * 2
	Forbid                   = "forbid:"
	ForbidForPolitics        = 60 * 30
	ForbidForProfanity       = 60
	ImgToImgMark             = "imgToImg:"
	FriendImgToImgMark       = "friendImgToImg:"
	ChatMode                 = "ChatMode:"
	ChatModeExp              = 60 * 5
	ImgExp                   = 60 * 2
	ImgGoldConsume           = 15
	DailyMAXFreeImgTransTime = 3 //每日生图功能额度
)

const (
	TxSecretId     = "TENCENTCLOUD_SECRET_ID"
	TxSecretKey    = "TENCENTCLOUD_SECRET_KEY"
	RegionShanghai = "ap-shanghai"
)
const (
	SourceNorMal = 1 //正常
	SourceExp    = 2 //已经过期

	PublicSource     = 1 //公共资源
	CommissionSource = 2 //返佣资源
)
const (
	NeverExp = "2050-1-1"
)

const (
	Success = iota
	ParamErr
	ServerErr
	DBErr
	ClientErr
)
