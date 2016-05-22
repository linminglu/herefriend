package handlers

import (
	"sync"
	"time"

	"herefriend/common"
)

type reviewInfo struct {
	FoceShowReviewAlert bool
	ShowReviewAlert     bool
	ReviewAlertMsg      string
	ReviewAlertCancel   string
	ReviewAlertGo       string
}

type registerInfo struct {
	Id              int               //用户ID
	PassWord        string            //用户密码
	ShowGift        bool              //是否显示礼物
	ClientVersion   int               //客户端版本
	Member          common.PersonInfo //用户详细信息
	ReviewAlertInfo reviewInfo
}

type vipUser struct {
	gender     int
	level      int   //VIP购买级别
	days       int   //VIP购买天数
	expiretime int64 //到期时间
}

type vipUsersInfo struct {
	users map[int]*vipUser
	lock  sync.RWMutex
}

type vipPriceInfo struct {
	Days         int
	Month        string
	OldPrice     int
	Price        int
	Discount     string
	DailyAverage string
	ProductId    string
	IabItemType  string
}

type vipInfo struct {
	VipLevel     int
	DiscountDesc string
	PriceList    []vipPriceInfo
}

var gVipLevels = []vipInfo{
	{
		VipLevel:     1,
		DiscountDesc: "限时特惠,优惠截止2016年7月12日",
		PriceList: []vipPriceInfo{
			{
				Days:         365,
				Month:        "12个月",
				OldPrice:     0,
				Price:        198,
				Discount:     "3折优惠",
				DailyAverage: "0.5",
				ProductId:    "12xx",
				IabItemType:  "subs"},
			{
				Days:         93,
				Month:        "3个月",
				OldPrice:     0,
				Price:        98,
				Discount:     "6折优惠",
				DailyAverage: "1.0",
				ProductId:    "3xx",
				IabItemType:  "inapp"},
			{
				Days:         31,
				Month:        "1个月",
				OldPrice:     0,
				Price:        50,
				Discount:     "",
				DailyAverage: "1.6",
				ProductId:    "1xx",
				IabItemType:  "subs"},
		},
	},
	{
		VipLevel:     2,
		DiscountDesc: "限时特惠,优惠截止2016年7月12日",
		PriceList: []vipPriceInfo{
			{
				Days:         365,
				Month:        "12个月",
				OldPrice:     0,
				Price:        98,
				Discount:     "3折优惠",
				DailyAverage: "0.3",
				ProductId:    "12zs",
				IabItemType:  "inapp"},
			{
				Days:         93,
				Month:        "3个月",
				OldPrice:     0,
				Price:        60,
				Discount:     "6折优惠",
				DailyAverage: "0.6",
				ProductId:    "3zs",
				IabItemType:  "inapp"},
			{
				Days:         31,
				Month:        "1个月",
				OldPrice:     0,
				Price:        30,
				Discount:     "",
				DailyAverage: "1.0",
				ProductId:    "1zs",
				IabItemType:  "inapp"},
		},
	},
	{
		VipLevel:     3,
		DiscountDesc: "限时特惠,优惠截止2016年7月12日",
		PriceList: []vipPriceInfo{
			{
				Days:         365,
				Month:        "12个月",
				OldPrice:     0,
				Price:        388,
				Discount:     "5折优惠",
				DailyAverage: "1.0",
				ProductId:    "12zz",
				IabItemType:  "inapp"},
			{
				Days:         93,
				Month:        "3个月",
				OldPrice:     0,
				Price:        188,
				Discount:     "",
				DailyAverage: "2.0",
				ProductId:    "3zz",
				IabItemType:  "inapp"},
		},
	},
}

type goldBeansPrice struct {
	Price     int    //价格
	Count     int    //普通会员购买数量
	Count_zs  int    //level 2会员购买数量
	Count_zz  int    //level 3会员购买数量
	ProductId string //产品ID
}

var gGoldBeansPrices = []goldBeansPrice{
	{
		Price:     6,
		Count:     30,
		Count_zs:  34,
		Count_zz:  50,
		ProductId: "6yuan"},
	{
		Price:     30,
		Count:     210,
		Count_zs:  234,
		Count_zz:  350,
		ProductId: "30yuan"},
	{
		Price:     98,
		Count:     784,
		Count_zs:  872,
		Count_zz:  1306,
		ProductId: "98yuan"},
	{
		Price:     298,
		Count:     2682,
		Count_zs:  2980,
		Count_zz:  4470,
		ProductId: "298yuan"},
	{
		Price:     588,
		Count:     5880,
		Count_zs:  6534,
		Count_zz:  9800,
		ProductId: "588yuan"},
	{
		Price:     998,
		Count:     14970,
		Count_zs:  16634,
		Count_zz:  24950,
		ProductId: "998yuan"},
}

type liveUser struct {
	gender   int
	status   int
	livetick int
	livetime int
}

type liveUsersInfo struct {
	users map[int]*liveUser
	lock  sync.RWMutex
}

type userBlacklist struct {
	Id        int   //用户id
	Blacklist []int //用户id的黑名单
}

const (
	MESSAGE_DIRECTION_TOME   = 0
	MESSAGE_DIRECTION_FROMME = 1
)

type messageInfo struct {
	MsgId     int                //消息Id
	MsgText   string             `json:"MsgText,omitempty"` //消息内容, 无内容时此字段会自动隐藏
	UserId    int                //用户Id
	UserInfo  *common.PersonInfo `json:"UserInfo,omitempty"` //用户信息
	Direction int                //消息方向, 0: UserId发送给我的消息, 1: 我发送给UserId的消息
	Readed    bool               //客户端是否显示为已读
	TimeUTC   time.Time          //标准时间,用来参考转换为本地时间
}

type allMessageInfo struct {
	RecommendArray []messageInfo
	VisitArray     []messageInfo
}

type unreadMessageInfo struct {
	UnreadRecommend int //未读的聊天消息
	UnreadVisit     int //未读的访问消息
	Badge           int //badge: the icon number of app
}

type PushMsgUnread struct {
	UnreadRecommend int //未读的聊天消息
	UnreadVisit     int //未读的访问消息
	Badge           int //badge: the icon number of app
}

type PushMsgEvaluation struct {
	Enable      bool   //是否要弹出评价对话框
	ShowMessage string //弹出对话框显示的信息
}

type PushMsgRefreshVIP struct {
	ShowMessage string //弹出对话框显示的信息
}

type PushMsgRecvGift struct {
	SenderId    int    //赠送者ID
	GiftId      int    //礼物ID
	GiftNum     int    //礼物数量
	GiftName    string //礼物名称
	ShowMessage string //弹出对话框显示的信息
}

type PushMessageInfo struct {
	/*
	 * 根据类型不同，消息实体的结构体不同，如下为具体对应关系:
	 * ------------------------------------------------------
	 * |   Type |				Value						|
	 * ------------------------------------------------------
	 * |    1   |          PushMsgUnread                    |
	 * ------------------------------------------------------
	 * |    2   |          PushMsgEvaluation                |
	 * ------------------------------------------------------
	 * |    3   |          PushMsgRecvGift                  |
	 * ------------------------------------------------------
	 */
	Type  int    //消息类型
	Value string //消息实体, 可以解析为对应的数据结构
}

type recommendQueueNode struct {
	timewait     int64
	fromid       int
	toid         int
	fromusertype int
	tousertype   int
	msgtype      int
	message      string
	timevalue    int64
}

var gHelloArray = [...]string{
	"你好呀",
	"在吗",
	"你好可以认识一下吗",
	"可以认识一下吗",
	"嘿嘿,在干嘛呢",
	"hi",
	"hello",
	"嗨",
	"打个招呼，嘿嘿",
}

var gRobotResponseCheckList = [...]string{
	"我叫",
	"小黄鸡",
	"小小鸡",
	"机器人",
	"我叫mimmimi",
	"simsimi",
	"主人",
	"我叫大脚",
	"2鸡",
	"对不起 是我不好",
	"器官",
}

type giftInfo struct {
	Id int //礼物固定id
	/* 礼物类型：
	 * 0 免费
	 * 1 普通礼物
	 * 2 折扣礼物
	 * 3 名人礼物
	 * ...
	 */
	Type                int
	Name                string //礼物名称
	ValidNum            int    //库存数量
	Description         string //礼物描述
	ImageUrl            string //礼物图片URL
	Effect              int    //礼物特效，需要客户端支持
	Price               int    //价格(beans)
	OriginPrice         int    //原价(beans)，对于折扣礼物和Price不同
	DiscountDescription string //折扣描述信息，对折扣作说明
}

type presentGiftInfo struct {
	UserInfo    common.PersonInfo //个人信息
	WhoRecvGift common.PersonInfo //收到礼物的人的信息
}

/*
 * 礼物列表详情
 */
type giftListVerbose struct {
	Person  common.PersonInfo //赠送礼物或者收到礼物的用户信息
	GiftId  int               //礼物ID
	GiftNum int               //礼物数量
	Message string            //礼物留言
	TimeUTC time.Time         //送礼物的时间
}

type giftRecvListInfo struct {
	toid    int
	giftid  int
	giftnum int
}

type userCharmInfo struct {
	Person      common.PersonInfo //用户信息
	GiftValue   int               //收到礼物的总价值
	AdmireCount int               //被心仪的数量
}
