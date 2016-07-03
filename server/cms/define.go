package cms

import "time"

type cmsSystemSummary struct {
	OSDescribe  string  // 操作系统信息
	CpuDescribe string  // CPU信息
	MemTotal    uint64  // 内存总大小(MB)
	MemUsed     uint64  // 内存使用大小(MB)
	MemUsage    float64 // 内存使用率
	HDTotal     uint64  // HD总大小(GB)
	HDUsed      uint64  // HD使用大小(GB)
	HDUsage     float64 // HD使用率
}

type cmsCpuInfo struct {
	CpuUsage float64 // CPU使用率
}

type cmsSystemUsersSummary struct {
	GirlsNum  int // 女生总数
	GuysNum   int // 男生总数
	ActiveNum int // 实时在线
	OnlineNum int // 在线人数
	RegistNum int // 新注册人数
}

type cmsCommentSummary struct {
	TalkNum   int // 聊天消息个数
	PushNum   int // 个推发送消息个数
	BuyVIPNum int // 购买VIP人数
}

type cmsCommentInfo struct {
	MsgId     int
	From      string
	FromPic   string
	FromId    int
	To        string
	ToPic     string
	ToId      int
	TimeUTC   time.Time
	Direction int //消息方向, 0: 收到消息, 1: 发送消息
	MsgType   int
	MsgText   string
}

type cmsTalkCommentInfo struct {
	MsgId   int
	FromId  int
	ToId    int
	TimeUTC time.Time
	MsgText string
}

type cmsTalkHistoryInfo struct {
	UserName   string
	TalkerName string
	UserPic    string
	TalkerPic  string
	Comments   []cmsTalkCommentInfo
}

type cmsMessageTempalte struct {
	Id       int
	Template string
}

type cmsUserInfo struct {
	Id               int
	Age              int
	Name             string
	Img              string
	Province         string
	Selected         bool
	Usertype         int
	VipLevel         int
	VipSetAppVersion string
}

type cmsSearchInfo struct {
	Count int
	Users []cmsUserInfo
}

type cmsImageInfo struct {
	filename string
	tag      int
}
