package common

import (
	"time"
)

const (
	// UserTypeRobot .
	UserTypeRobot = 0
	// UserTypeUser .
	UserTypeUser = 1
	// UserTypeWeibo .
	UserTypeWeibo = 2
	// UserTypeZhenQing .
	UserTypeZhenQing = 3

	// ClientAgent ...
	ClientAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.99 Safari/537.36"
)

// GiftSendRecvInfo .
type GiftSendRecvInfo struct {
	GiftID int //礼物id
	Number int //礼物数量
}

// PersonInfo is infomation shows to the clients
type PersonInfo struct {
	ID              int                `json:"Id"`         //ID号
	Height          int                `json:",omitempty"` //身高
	Weight          int                `json:",omitempty"` //体重
	Age             int                `json:",omitempty"` //年龄
	Gender          int                //性别: 0(女) 1(男)
	OnlineStatus    int                //在线状态
	VipLevel        int                //Vip级别
	VipExpireTime   time.Time          `json:",omitempty"` //会员到期时间
	Name            string             //姓名
	Province        string             `json:",omitempty"` //所在省/直辖市/自治区
	District        string             `json:",omitempty"` //所在区域
	Native          string             `json:",omitempty"` //家乡
	LoveType        string             `json:",omitempty"` //恋爱类型
	BodyType        string             `json:",omitempty"` //体型
	BloodType       string             `json:",omitempty"` //体型
	Animal          string             `json:",omitempty"` //属相
	Constellation   string             `json:",omitempty"` //星座
	Lang            string             `json:",omitempty"` //语言
	Introduction    string             //自我介绍
	Selfjudge       string             `json:",omitempty"`                //自评
	Education       string             `json:",omitempty"`                //教育程度
	Income          string             `json:",omitempty"`                //收入情况
	IncomeMin       int                `json:",omitempty"`                //收入最低
	IncomeMax       int                `json:",omitempty"`                //收入最高
	School          string             `json:",omitempty"`                //毕业学校
	Occupation      string             `json:",omitempty"`                //职业
	Housing         string             `json:",omitempty"`                //购房情况
	Carstatus       string             `json:",omitempty"`                //购车情况
	Speciality      string             `json:",omitempty"`                //技能
	Marriage        string             `json:",omitempty"`                //婚姻状况
	Companytype     string             `json:",omitempty"`                //公司类型
	Companyindustry string             `json:",omitempty"`                //公司领域
	Nationnality    string             `json:",omitempty"`                //民族
	Religion        string             `json:",omitempty"`                //信仰
	Charactor       string             `json:",omitempty"`                //性格类型
	Hobbies         string             `json:",omitempty"`                //兴趣爱好
	CityLove        int                `json:",omitempty"`                //是否接受异地恋: 0(视情况而定) 1(接受) 2(不接受)
	Naken           int                `json:",omitempty"`                //是否接受婚前性行为: 0(视情况而定) 1(接受) 2(不接受)
	AllowAge        string             `json:"Allow_age,omitempty"`       //择偶条件:年龄
	AllowResidence  string             `json:"Allow_residence,omitempty"` //择偶条件:居住地
	AllowHeight     string             `json:"Allow_height,omitempty"`    //择偶条件:身高
	AllowMarriage   string             `json:"Allow_marriage,omitempty"`  //择偶条件:婚姻状况
	AllowEducation  string             `json:"Allow_education,omitempty"` //择偶条件:教育程度
	AllowHousing    string             `json:"Allow_housing,omitempty"`   //择偶条件:购房情况
	AllowIncome     string             `json:"Allow_income,omitempty"`    //择偶条件:收入
	AllowKidStatus  string             `json:"Allow_kidstatus,omitempty"` //择偶条件:子女情况
	IconURL         string             `json:"IconUrl,omitempty"`         //头像url
	Pics            []string           `json:",omitempty"`                //照片列表
	GoldBeans       int                `json:",omitempty"`                //用户的金币数量
	RecvGiftList    []GiftSendRecvInfo `json:",omitempty"`                //收到的礼物列表
	SendGiftList    []GiftSendRecvInfo `json:",omitempty"`                //送出的礼物列表
}

// UserCharmInfo charm information of user
type UserCharmInfo struct {
	Person      PersonInfo //用户信息
	GiftValue   int        //收到礼物的总价值
	AdmireCount int        //被心仪的数量
}

// UserCharmInfoList .
type UserCharmInfoList []UserCharmInfo

func (list UserCharmInfoList) Len() int {
	return len(list)
}

func (list UserCharmInfoList) Less(i, j int) bool {
	if list[i].GiftValue > list[j].GiftValue {
		return true
	}

	return false
}

func (list UserCharmInfoList) Swap(i, j int) {
	temp := list[i]
	list[i] = list[j]
	list[j] = temp
}

// UserWealthInfo is the wealth information of user
type UserWealthInfo struct {
	Person        PersonInfo //用户信息
	ConsumedBeans int        //花费金币的总数量
}

// UserWealthInfoList .
type UserWealthInfoList []UserWealthInfo

func (list UserWealthInfoList) Len() int {
	return len(list)
}

func (list UserWealthInfoList) Less(i, j int) bool {
	if list[i].ConsumedBeans > list[j].ConsumedBeans {
		return true
	}

	return false
}

func (list UserWealthInfoList) Swap(i, j int) {
	temp := list[i]
	list[i] = list[j]
	list[j] = temp
}
