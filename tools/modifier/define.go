package main

type randOption struct {
	option string
	dot    int
}

type randValue struct {
	value1 int
	value2 int
	dot    int
}

type valuenode struct {
	value1 int
	value2 int
}

/*
+---vip
	  +---0
	  +---1
	  +---2
	  +---3
非心动用户：无：65%；1级: 5%；2级:10%；3级:20%。
心动用户：无：30%；1级：10%；2级：20%；3级：40%。
*/
var gRandViplevelGirls = []randValue{
	{0, 0, 250},
	{1, 30, 30},
	{1, 60, 10},
	{1, 90, 20},
	{2, 30, 10},
	{2, 60, 5},
	{2, 90, 10},
	{3, 30, 10},
	{3, 60, 10},
	{3, 90, 30},
}

var gRandViplevelGuys = gRandViplevelGirls

var gRandViplevelHeartbeatGirls = []randValue{
	{0, 0, 90},
	{1, 30, 30},
	{1, 60, 10},
	{1, 90, 20},
	{2, 30, 10},
	{2, 60, 5},
	{2, 90, 10},
	{3, 30, 10},
	{3, 60, 10},
	{3, 90, 30},
}

var gRandViplevelHeartbeatGuys = gRandViplevelHeartbeatGirls

/*
+---收入
	  +---小于2000
	  +---2000-5000
	  +---5000-10000
	  +---10000-20000
	  +---20000及以上
*/
var gRandIncomeGirls = []randValue{
	{0, 2000, 60},     //0 ~ 2000
	{2000, 5000, 60},  //2000 ~ 5000
	{5000, 10000, 5},  //5000 ~ 10000
	{10000, 20000, 2}, //10000 ~ 20000
	{20000, 0, 1},     //20000及以上
	{0, 0, 20},        //不限
}

var gRandIncomeGuys = []randValue{
	{0, 2000, 0},       //0 ~ 2000
	{2000, 5000, 20},   //2000 ~ 5000
	{5000, 10000, 40},  //5000 ~ 10000
	{10000, 20000, 20}, //10000 ~ 20000
	{20000, 0, 5},      //20000及以上
	{0, 0, 20},         //不限
}

/*
 * district
 */
var gRandDistrictGirls = []randOption{
	{"北京市", 30},
	{"天津市", 10},
	{"河北省", 10},
	{"山西省", 5},
	{"内蒙古自治区", 5},
	{"辽宁省", 5},
	{"吉林省", 8},
	{"黑龙江省", 10},
	{"上海市", 17},
	{"江苏省", 13},
	{"浙江省", 20},
	{"安徽省", 5},
	{"福建省", 15},
	{"江西省", 5},
	{"山东省", 15},
	{"河南省", 10},
	{"湖北省", 10},
	{"湖南省", 10},
	{"广东省", 30},
	{"广西壮族自治区", 5},
	{"海南省", 12},
	{"重庆市", 14},
	{"四川省", 15},
	{"贵州省", 5},
	{"云南省", 5},
	{"西藏自治区", 3},
	{"陕西省", 5},
	{"甘肃省", 5},
	{"青海省", 3},
	{"宁夏回族自治区", 5},
	{"新疆维吾尔自治区", 5},
	{"台湾省", 5},
}

var gRandDistrictGuys = gRandDistrictGirls

/*
+---职业
	  +---企业职工
	  +---在校学生
	  +---商业贸易
	  +---教育
	  +---科研技术
	  +---医院医疗
	  +---艺术行业
	  +---其它
*/
var gRandOccupationGirls = []randOption{
	{"企业职工", 10},
	{"在校学生", 10},
	{"商业贸易", 10},
	{"教育", 10},
	{"科研技术", 10},
	{"医院医疗", 10},
	{"艺术行业", 10},
	{"其它", 20},
	{"", 30}, //不限
}

var gRandOccupationGuys = gRandOccupationGirls

/*
+---学历
	  +---初中及以下
	  +---高中及中专
	  +---大专
	  +---本科
	  +---硕士及以上
*/
var gRandEducationGirls = []randOption{
	{"初中及以下", 30},
	{"高中及中专", 30},
	{"大专", 10},
	{"本科", 20},
	{"硕士及以上", 5},
	{"", 30}, //不限
}

var gRandEducationGuys = []randOption{
	{"初中及以下", 8},
	{"高中及中专", 10},
	{"大专", 15},
	{"本科", 60},
	{"硕士及以上", 20},
	{"", 30}, //不限
}

/*
+---居住条件
	  +---买房
	  +---租房
	  +---不限
*/
var gRandHousingGirls = []randOption{
	{"买房", 10},
	{"租房", 40},
	{"", 30},
}

var gRandHousingGuys = []randOption{
	{"买房", 40},
	{"租房", 10},
	{"", 30},
}

/*
+---婚姻
	  +---未婚
	  +---已婚
	  +---离异
	  +---丧偶
*/
var gRandMarriageGirls = []randOption{
	{"未婚", 80},
	{"已婚", 10},
	{"离异", 5},
	{"丧偶", 1},
}

var gRandMarriageGuys = []randOption{
	{"未婚", 80},
	{"已婚", 10},
	{"离异", 15},
	{"丧偶", 1},
}

/*
+---性格(girl)
	  +---温柔
	  +---感性
	  +---贤惠
	  +---可爱
	  +---大方
	  +---热情
	  +---成熟
	  +---文静
+---性格(guy)
	  +---幽默
	  +---感性
	  +---体贴
	  +---憨厚
	  +---稳重
	  +---好强
	  +---冷静
	  +---温柔
*/
var gRandCharactorGirls = []randOption{
	{"温柔", 30},
	{"感性", 20},
	{"贤惠", 10},
	{"可爱", 17},
	{"大方", 10},
	{"热情", 20},
	{"成熟", 5},
	{"文静", 20},
	{"", 20},
}

var gRandCharactorGuys = []randOption{
	{"幽默", 20},
	{"感性", 10},
	{"体贴", 20},
	{"憨厚", 9},
	{"稳重", 18},
	{"好强", 10},
	{"冷静", 12},
	{"温柔", 10},
	{"", 20},
}

/*
+---兴趣爱好
	  +---运动
	  +---烹饪
	  +---看电影
	  +---读书
	  +---上网
	  +---听音乐
	  +---养小动物
	  +---旅游
*/
var gRandHobbiesGirls = []randOption{
	{"运动", 10},
	{"烹饪", 20},
	{"看电影", 20},
	{"读书", 10},
	{"上网", 10},
	{"听音乐", 10},
	{"养小动物", 20},
	{"旅游", 20},
	{"", 20},
}

var gRandHobbiesGuys = []randOption{
	{"运动", 30},
	{"烹饪", 10},
	{"看电影", 10},
	{"读书", 10},
	{"上网", 30},
	{"听音乐", 10},
	{"养小动物", 5},
	{"旅游", 20},
	{"", 20},
}

const (
	V_INCOME_GIRLS = iota
	V_INCOME_GUYS
	V_VIPLEVEL_GIRLS
	V_VIPLEVEL_GUYS
	V_HEART_VIPLEVEL_GIRLS
	V_HEART_VIPLEVEL_GUYS
	V_SIZE
)

var gRandValueMax [V_SIZE]int
var gRandValueMap [V_SIZE]map[int]valuenode

const (
	O_DISTRICT_GIRLS = iota
	O_DISTRICT_GUYS
	O_OCCUPATION_GIRLS
	O_OCCUPATION_GUYS
	O_EDUCATION_GIRLS
	O_EDUCATION_GUYS
	O_HOUSING_GIRLS
	O_HOUSING_GUYS
	O_MARRIAGE_GIRLS
	O_MARRIAGE_GUYS
	O_CHARACTOR_GIRLS
	O_CHARACTOR_GUYS
	O_HOBBIES_GIRLS
	O_HOBBIES_GUYS
	O_SIZE
)

var gRandOptionMax [O_SIZE]int
var gRandOptionMap [O_SIZE]map[int]string

/*
 |    Function: getRandOptionResult
 |      Author: Mr.Sancho
 |        Date: 2016-02-03
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func getRandOptionResult(options []randOption) (int, map[int]string) {
	index := 0
	indexmap := make(map[int]string)

	for _, o := range options {
		for i := 0; i < o.dot; i = i + 1 {
			indexmap[index] = o.option
			index = index + 1
		}
	}

	return index, indexmap
}

/*
 |    Function: getRandValueResult
 |      Author: Mr.Sancho
 |        Date: 2016-02-04
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func getRandValueResult(values []randValue) (int, map[int]valuenode) {
	index := 0
	indexmap := make(map[int]valuenode)

	for _, v := range values {
		for i := 0; i < v.dot; i = i + 1 {
			indexmap[index] = valuenode{v.value1, v.value2}
			index = index + 1
		}
	}

	return index, indexmap
}
