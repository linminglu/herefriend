package handlers

import (
	"sync"
	"time"

	"herefriend/common"
)

type personInfo struct {
	common.PersonInfo
}

type reviewInfo struct {
	FoceShowReviewAlert bool
	ShowReviewAlert     bool
	ReviewAlertMsg      string
	ReviewAlertCancel   string
	ReviewAlertGo       string
}

type registerInfo struct {
	Id              int        //用户ID
	PassWord        string     //用户密码
	ShowGift        bool       //是否显示礼物
	ClientVersion   int        //客户端版本
	Member          personInfo //用户详细信息
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

type loveshowcomment struct {
	Id        int       //用户id
	Age       int       //年龄
	Name      string    //姓名
	District  string    //地区
	Education string    //教育程度
	Text      string    //评论内容
	TimeUTC   time.Time //评论时间(标准时间)
}

type loverInfo struct {
	Id       int    //用户id
	Age      int    //年龄
	Name     string //姓名
	Imgurl   string //小头像地址
	District string //地区
}

type loveShow struct {
	Id           int               //爱情秀编号
	Girl         loverInfo         //女生信息
	Guy          loverInfo         //男生信息
	Daysfalllove int               //见面多少天之后就恋爱了
	Blessnum     int               //收到的祝福数量
	Lovestatus   string            //目前的恋爱状态
	Lovetitle    string            //爱情秀主题
	Lovestory    string            //爱情故事
	TimeUTC      time.Time         //故事时间(标准时间)
	ShowPics     []string          //恋爱秀照片
	Comments     []loveshowcomment //评论
}

type loveShowList struct {
	loveShow
	Hide bool
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
	MsgId     int         //消息Id
	MsgText   string      `json:"MsgText,omitempty"` //消息内容, 无内容时此字段会自动隐藏
	UserId    int         //用户Id
	UserInfo  *personInfo `json:"UserInfo,omitempty"` //用户信息
	Direction int         //消息方向, 0: UserId发送给我的消息, 1: 我发送给UserId的消息
	Readed    bool        //客户端是否显示为已读
	TimeUTC   time.Time   //标准时间,用来参考转换为本地时间
}

type allMessageInfo struct {
	RecommendArray []messageInfo
	VisitArray     []messageInfo
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

type PushMessageInfo struct {
	/*
	 * 根据类型不同，消息实体的结构体不同，如下为具体对应关系:
	 * ------------------------------------------------------
	 * | Type值 |         Value对应的数据结构               |
	 * ------------------------------------------------------
	 * |    1   |          PushMsgUnread                    |
	 * ------------------------------------------------------
	 * |    2   |          PushMsgEvaluation                |
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
