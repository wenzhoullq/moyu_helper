package constant

const (
	RESP_NORMAL = 0
)

const (
	GameMaxNum     = 100
	ReplyMaxNum    = 300
	ReplyPicMaxNum = 3
	UpdateMaxNum   = 300
)

const (
	LlmKeyWord       = "@摸鱼小助手"
	SignKeyWord      = "签到"
	ImgToImgKeyWord  = "风格转换"
	TextToImgKeyWord = "画"
	EleKeyWord       = "饿了么"
	Upgrade          = "升级"
	DrawLots         = "抽签"
	UnDrawLots       = "解签"
	ZodiacBlindBox   = "生肖盲盒"
)

const (
	PreContent         = "请结合上下文并用50个字以内回答:\n"
	Short              = "把如下内容缩短至50个字以内:\n"
	SignFailReply      = "╭┈┈⚠️失败⚠️┈┈╮\n❗ 请设置群昵称后签到\n╰┈┈⚠️失败⚠️┈┈╯"
	SignSuccessReply   = "╭┈┈┈🏡签到🏡┈┈┈╮\n🌼昵称：%s\n🐾排名：第%d名\n📆奖励：%d枚金币\n💰财富：%d枚金币\n📝签到次数：%d次\n⏰签到时间: %s\n╰┈┈┈┈┈┈┈┈┈┈┈╯"
	StatusSuccessReply = "╭┈┈┈🥷人物🥷┈┈┈╮\n🌼昵称：%s\n💯等级：LV%d\n👊武力值:%d\n😺幸运值:%d\n💰财富：%d枚金币\n╰┈┈┈┈┈┈┈┈┈┈┈╯"
	RepeatSignReply    = "╭┈┈┈🏡签到🏡┈┈┈╮\n🌼昵称：%s\n🐾状态：今日已签到\n╰┈┈┈🏡签到🏡┈┈┈╯"
	SystemNormal       = "你是汤鸽科技集团有限公司研发的摸鱼小助手,旨在各位摸鱼人进行摸鱼和薅羊毛,@摸鱼小助手并输入薅羊毛获得各种福利"
	//SystemDoctor                    = "你是汤鸽科技集团有限公司研发的摸鱼小助手医疗藐视,旨在各位摸鱼人进行各种医疗相关问题咨询,但是实际情况还得线下医生问诊"
	ImgToImgApplicationSuccess      = "好的,请在60秒内请发出图片,我将图片进行风格转换"
	TransToImgApplicationFail       = "抱歉,您的金币不足%d枚,无法进行AI绘画。%s"
	UpgradeApplicationFail          = "抱歉,您的金币不足%d枚,无法进行升级。%s"
	GoldGetTip                      = "\n可通过签到的方式获得金币。"
	ImgReplyGroup                   = "图片正在生成中,请稍等...\n本次预计消耗金币%d枚"
	ImgReplyFriend                  = "图片正在生成中,请稍等..."
	GoldConsumeReply                = "金币扣除%d枚!\n剩余金币数量:%d枚"
	EmptyReply                      = "你好,我是摸鱼小助手,有什么能帮助到您？"
	ReplyPre                        = "你想要的是否是:\n"
	NewsSuf                         = "【摸鱼小助手】这里是今天的摸鱼小新闻,祝各位摸鱼人摸鱼愉快！\n"
	NorMalModeForbidDirty           = "善言结善缘,恶语伤人心。你这一句话我需要花60秒来治愈自己😭😭"
	AoJiaoModelForbidDirty          = "你说话太没有礼貌了,我不想跟你说话了!"
	ExDailyMAXFreeImgTransTimeReply = "今日免费生图功能次数已用完,请明日再来使用"
	HolidayTip                      = "【摸鱼小助手】提醒您:各位摸鱼人上午好🌹！\n工作再累，一定不要忘记摸🐟！有事没事起身去茶水间、去厕所、去廊道走走，别老在工位上坐着，💴是老板的，但命是自己的！\n"
	WednesdayAd                     = "周三神券节,@摸鱼小助手发送 美团外卖 必得9元红包。\n"
	DailyAd                         = "@摸鱼小助手发送美团/饿了么领取8元红包，发送滴滴打车领取优惠券。\n"
	WeekendAd                       = "@摸鱼小助手发送美团/饿了么领取8元红包，发送滴滴打车领取五折优惠券。\n"
	ModeChatSet                     = "已切换为%s对话模式"
	ModeImgSet                      = "已切换为%s生图模式"
	DailyProfit                     = "昨日收益:\n美团:%.2f元。\n滴滴:%.2f元。\n总收益:%.2f元。"
	SendDailyProfitUser             = "哥哥"
	UnDrawLotsSuf                   = "解@%s的签:\n%s"
	ZodiacBlindBoxSuf               = "恭喜您抽中了%s"
	HasDraw                         = "今日已抽过签,请明天再来~"
	HasZodiacBlindBox               = "今日已抽过生肖盲盒,请明天再来~"
	UnHasDrawLots                   = "今日还未抽签,请抽签后再进行解签~"
)
const (
	SignSuccess = 0
	SignFail    = 1
)

const (
	NorMalModeChat = "正常模式"
	AoJiaoModeChat = "傲娇模式"
	DoctorModeChat = "医疗模式"
)

const (
	SignFirst  = 1
	SignSecond = 2
	SignThird  = 3
)

const (
	Unlimited              = "不限定风格"
	InkWash                = "水墨画"
	ConceptualArt          = "概念艺术"
	OilPainting1           = "油画1"
	OilPainting2           = "油画2(梵高)"
	Watercolor             = "水彩画"
	PixelPainting          = "像素画"
	ThickCoating           = "厚涂风格"
	Illustration           = "插图"
	PaperCuttings          = "剪纸"
	Impressionism1         = "印象派1(莫奈)"
	Impressionism2         = "印象派2"
	D25                    = "2.5D"
	D3                     = "3D"
	ClassicalPortrait      = "古典肖像画"
	BlackAndWhiteSketching = "黑白素描画"
	Cyberpunk              = "赛博朋克"
	ScienceFiction         = "科幻风格"
	Dark                   = "暗黑风格"
	SteamWave              = "蒸汽波"
	JapaneseAnime          = "日系动漫"
	Monster                = "怪兽风格"
	BeautifulAncient       = "唯美古风"
	RetroAnime             = "复古动漫"
	GameCartoon            = "游戏卡通手绘"
	Universal              = "通用写实风格"
)

const (
	UnlimitedMark              = "000"
	InkWashMark                = "101"
	ConceptualArtMark          = "102"
	OilPainting1Mark           = "103"
	OilPainting2Mark           = "118"
	WatercolorMark             = "104"
	PixelPaintingMark          = "105"
	ThickCoatingMark           = "106"
	IllustrationMark           = "107"
	PaperCuttingsMark          = "108"
	Impressionism1Mark         = "109"
	Impressionism2Mark         = "119"
	D25Mark                    = "110"
	ClassicalPortraitMark      = "111"
	BlackAndWhiteSketchingMark = "112"
	CyberpunkMark              = "113"
	ScienceFictionMark         = "114"
	DarkMark                   = "115"
	D3Mark                     = "116"
	SteamWaveMark              = "117"
	JapaneseAnimeMark          = "201"
	MonsterMark                = "202"
	BeautifulAncientMark       = "203"
	RetroAnimeMark             = "204"
	GameCartoonMark            = "301"
	UniversalMark              = "401"
)

const (
	MaxAnswerLen = 50
)

const (
	MaxNewsNum = 3
)

const (
	ShortTime                = "shortTimeCnt"
	SignTime                 = "signTimesCnt"
	SignMark                 = "sign:"
	ChatMark                 = "chat:"
	ChatExp                  = 60 * 2
	Forbid                   = "forbid:"
	ForbidForPolitics        = 60 * 30
	ForbidForProfanity       = 60
	ImgToImgMark             = "imgToImg:"
	FriendImgToImgMark       = "friendImgToImg:"
	GroupDrawLotsMark        = "groupDrawLots:%s-%s"
	GroupZodiacBlindBoxMark  = "groupZodiacBlindBox:%s-%s"
	ChatMode                 = "ChatMode:"
	ImgMode                  = "ImgMode:"
	ModeExp                  = -1
	ImgExp                   = 60 * 2
	ImgGoldConsume           = 25
	LvUpConsume              = 100
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
const (
	MeiTuanUnionLimit = "100"
	MeiTuanUnionPage  = "1"
	MeiTuanActID      = "33"
)
