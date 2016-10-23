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
	ID              int               `json:"Id"` //用户ID
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
	ProductID    string `json:"ProductId"`
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
				ProductID:    "12xx",
				IabItemType:  "subs"},
			{
				Days:         93,
				Month:        "3个月",
				OldPrice:     0,
				Price:        98,
				Discount:     "6折优惠",
				DailyAverage: "1.0",
				ProductID:    "3xx",
				IabItemType:  "inapp"},
			{
				Days:         31,
				Month:        "1个月",
				OldPrice:     0,
				Price:        50,
				Discount:     "",
				DailyAverage: "1.6",
				ProductID:    "1xx",
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
				ProductID:    "12zs",
				IabItemType:  "inapp"},
			{
				Days:         93,
				Month:        "3个月",
				OldPrice:     0,
				Price:        60,
				Discount:     "6折优惠",
				DailyAverage: "0.6",
				ProductID:    "3zs",
				IabItemType:  "inapp"},
			{
				Days:         31,
				Month:        "1个月",
				OldPrice:     0,
				Price:        30,
				Discount:     "",
				DailyAverage: "1.0",
				ProductID:    "1zs",
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
				ProductID:    "12zz",
				IabItemType:  "inapp"},
			{
				Days:         93,
				Month:        "3个月",
				OldPrice:     0,
				Price:        188,
				Discount:     "",
				DailyAverage: "2.0",
				ProductID:    "3zz",
				IabItemType:  "inapp"},
		},
	},
}

type goldBeansPrice struct {
	Price     int    //价格
	Count     int    //普通会员购买数量
	Song      int    //赠送
	ProductID string `json:"ProductId"` //产品ID
}

var gGoldBeansPrices = []goldBeansPrice{
	{
		Price:     6,
		Count:     60,
		Song:      0,
		ProductID: "6yuan1"},
	{
		Price:     30,
		Count:     300,
		Song:      0,
		ProductID: "30yuan1"},
	{
		Price:     98,
		Count:     980,
		Song:      118,
		ProductID: "98yuan1"},
	{
		Price:     298,
		Count:     2980,
		Song:      388,
		ProductID: "298yuan1"},
	{
		Price:     588,
		Count:     5880,
		Song:      638,
		ProductID: "588yuan1"},
	{
		Price:     998,
		Count:     9980,
		Song:      1388,
		ProductID: "998yuan1"},
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
	ID        int   `json:"Id"` //用户id
	Blacklist []int //用户id的黑名单
}

const (
	// MessageDirectionToMe .
	MessageDirectionToMe = 0
	// MessageDirectionFromMe .
	MessageDirectionFromMe = 1
)

type messageInfo struct {
	MsgID     int                `json:"MsgId"`              //消息Id
	MsgText   string             `json:"MsgText,omitempty"`  //消息内容, 无内容时此字段会自动隐藏
	UserID    int                `json:"UserId"`             //用户Id
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

// PushMsgUnread .
type PushMsgUnread struct {
	UnreadRecommend int //未读的聊天消息
	UnreadVisit     int //未读的访问消息
	Badge           int //badge: the icon number of app
}

// PushMsgEvaluation .
type PushMsgEvaluation struct {
	Enable      bool   //是否要弹出评价对话框
	ShowMessage string //弹出对话框显示的信息
}

// PushMsgRefreshVIP .
type PushMsgRefreshVIP struct {
	ShowMessage string //弹出对话框显示的信息
}

// PushMsgRecvGift .
type PushMsgRecvGift struct {
	SenderID    int    `json:"SenderId"` //赠送者ID
	GiftID      int    `json:"GiftId"`   //礼物ID
	GiftNum     int    //礼物数量
	GiftName    string //礼物名称
	ShowMessage string //弹出对话框显示的信息
}

// PushMessageInfo .
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

// GiftInfo .
type GiftInfo struct {
	ID int `json:"Id"` //礼物固定id
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
	ImageURL            string `json:"ImageUrl"` //礼物图片URL
	Effect              int    //礼物特效，需要客户端支持
	Price               int    //价格(beans)
	OriginPrice         int    //原价(beans)，对于折扣礼物和Price不同
	DiscountDescription string //折扣描述信息，对折扣作说明
}

type presentGiftInfo struct {
	UserInfo    common.PersonInfo //个人信息
	WhoRecvGift common.PersonInfo //收到礼物的人的信息
}

// GiftListVerbose 礼物列表详情
type GiftListVerbose struct {
	ID      int               `json:"Id"` //数据唯一性标识
	Person  common.PersonInfo //赠送礼物或者收到礼物的用户信息
	GiftID  int               `json:"GiftId"` //礼物ID
	GiftNum int               //礼物数量
	Message string            //礼物留言
	TimeUTC time.Time         //送礼物的时间
}

type giftRecvListInfo struct {
	toid    int
	giftid  int
	giftnum int
}

// AppConfig .
type AppConfig struct {
	Person      common.PersonInfo
	StartupView struct {
		ImageURL   string `json:"ImageUrl"` //图片地址
		Duration   int    //图片显示时间
		LinkEnable bool   //链接是否使能
		LinkURL    string `json:"LinkUrl"` //链接跳转地址
	}
	VersionInfo struct {
		VersionStr string //版本信息
	}
}
