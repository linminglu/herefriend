package common

import (
	"time"
)

const (
	USERTYPE_BH    = 0
	USERTYPE_USER  = 1
	USERTYPE_WEIBO = 2
)

const (
	ClientAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.99 Safari/537.36"
)

/*
 * Infomation shows to the clients
 */
type PersonInfo struct {
	Id              int       //ID号
	Height          int       //身高
	Weight          int       //体重
	Age             int       //年龄
	Gender          int       //性别: 0(女) 1(男)
	OnlineStatus    int       //在线状态
	VipLevel        int       //Vip级别
	VipExpireTime   time.Time //会员到期时间
	Name            string    //姓名
	Province        string    //所在省/直辖市/自治区
	District        string    //所在区域
	Native          string    //家乡
	LoveType        string    //恋爱类型
	BodyType        string    //体型
	BloodType       string    //体型
	Animal          string    //属相
	Constellation   string    //星座
	Lang            string    //语言
	Introduction    string    //自我介绍
	Selfjudge       string    //自评
	Education       string    //教育程度
	Income          string    //收入情况
	IncomeMin       int       //收入最低
	IncomeMax       int       //收入最高
	School          string    //毕业学校
	Occupation      string    //职业
	Housing         string    //购房情况
	Carstatus       string    //购车情况
	Speciality      string    //技能
	Marriage        string    //婚姻状况
	Companytype     string    //公司类型
	Companyindustry string    //公司领域
	Nationnality    string    //民族
	Religion        string    //信仰
	Charactor       string    //性格类型
	Hobbies         string    //兴趣爱好
	CityLove        int       //是否接受异地恋: 0(视情况而定) 1(接受) 2(不接受)
	Naken           int       //是否接受婚前性行为: 0(视情况而定) 1(接受) 2(不接受)
	Allow_age       string    //择偶条件:年龄
	Allow_residence string    //择偶条件:居住地
	Allow_height    string    //择偶条件:身高
	Allow_marriage  string    //择偶条件:婚姻状况
	Allow_education string    //择偶条件:教育程度
	Allow_housing   string    //择偶条件:购房情况
	Allow_income    string    //择偶条件:收入
	Allow_kidstatus string    //择偶条件:子女情况
	IconUrl         string    //头像url
	Pics            []string  //照片列表
}
