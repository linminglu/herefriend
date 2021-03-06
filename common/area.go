package common

import "strings"

type districtInfo struct {
	Province string
	Provcode string
	District string
	Distcode string
}

var gDistrictBrief = [...]districtInfo{
	{"北京", "8611", "东城", "861101"},
	{"北京", "8611", "西城", "861102"},
	{"北京", "8611", "崇文", "861103"},
	{"北京", "8611", "宣武", "861104"},
	{"北京", "8611", "朝阳", "861105"},
	{"北京", "8611", "丰台", "861106"},
	{"北京", "8611", "石景山", "861107"},
	{"北京", "8611", "海淀", "861108"},
	{"北京", "8611", "门头沟", "861109"},
	{"北京", "8611", "房山", "861110"},
	{"北京", "8611", "通州", "861111"},
	{"北京", "8611", "顺义", "861112"},
	{"北京", "8611", "昌平", "861113"},
	{"北京", "8611", "大兴", "861114"},
	{"北京", "8611", "怀柔", "861115"},
	{"北京", "8611", "平谷", "861116"},
	{"北京", "8611", "密云", "861117"},
	{"北京", "8611", "延庆", "861118"},
	{"天津", "8612", "和平", "861201"},
	{"天津", "8612", "河东", "861202"},
	{"天津", "8612", "河西", "861203"},
	{"天津", "8612", "南开", "861204"},
	{"天津", "8612", "河北", "861205"},
	{"天津", "8612", "红桥", "861206"},
	{"天津", "8612", "塘沽", "861207"},
	{"天津", "8612", "汉沽", "861208"},
	{"天津", "8612", "大港", "861209"},
	{"天津", "8612", "东丽", "861210"},
	{"天津", "8612", "西青", "861211"},
	{"天津", "8612", "津南", "861212"},
	{"天津", "8612", "北辰", "861213"},
	{"天津", "8612", "武清", "861214"},
	{"天津", "8612", "宝坻", "861215"},
	{"天津", "8612", "宁河", "861216"},
	{"天津", "8612", "静海", "861217"},
	{"天津", "8612", "蓟县", "861218"},
	{"河北", "8613", "石家庄", "861301"},
	{"河北", "8613", "唐山", "861302"},
	{"河北", "8613", "秦皇岛", "861303"},
	{"河北", "8613", "邯郸", "861304"},
	{"河北", "8613", "邢台", "861305"},
	{"河北", "8613", "保定", "861306"},
	{"河北", "8613", "张家口", "861307"},
	{"河北", "8613", "承德", "861308"},
	{"河北", "8613", "沧州", "861309"},
	{"河北", "8613", "廊坊", "861310"},
	{"河北", "8613", "衡水", "861311"},
	{"山西", "8614", "太原", "861401"},
	{"山西", "8614", "大同", "861402"},
	{"山西", "8614", "阳泉", "861403"},
	{"山西", "8614", "长治", "861404"},
	{"山西", "8614", "晋城", "861405"},
	{"山西", "8614", "朔州", "861406"},
	{"山西", "8614", "晋中", "861407"},
	{"山西", "8614", "运城", "861408"},
	{"山西", "8614", "忻州", "861409"},
	{"山西", "8614", "临汾", "861410"},
	{"山西", "8614", "吕梁", "861411"},
	{"内蒙古自治", "8615", "呼和浩特", "861501"},
	{"内蒙古自治", "8615", "包头", "861502"},
	{"内蒙古自治", "8615", "乌海", "861503"},
	{"内蒙古自治", "8615", "赤峰", "861504"},
	{"内蒙古自治", "8615", "通辽", "861505"},
	{"内蒙古自治", "8615", "鄂尔多斯", "861506"},
	{"内蒙古自治", "8615", "呼伦贝尔", "861507"},
	{"内蒙古自治", "8615", "巴彦淖尔", "861508"},
	{"内蒙古自治", "8615", "乌兰察布", "861509"},
	{"内蒙古自治", "8615", "兴安盟", "861522"},
	{"内蒙古自治", "8615", "锡林郭勒盟", "861525"},
	{"内蒙古自治", "8615", "阿拉善盟", "861529"},
	{"辽宁", "8621", "沈阳", "862101"},
	{"辽宁", "8621", "大连", "862102"},
	{"辽宁", "8621", "鞍山", "862103"},
	{"辽宁", "8621", "抚顺", "862104"},
	{"辽宁", "8621", "本溪", "862105"},
	{"辽宁", "8621", "丹东", "862106"},
	{"辽宁", "8621", "锦州", "862107"},
	{"辽宁", "8621", "营口", "862108"},
	{"辽宁", "8621", "阜新", "862109"},
	{"辽宁", "8621", "辽阳", "862110"},
	{"辽宁", "8621", "盘锦", "862111"},
	{"辽宁", "8621", "铁岭", "862112"},
	{"辽宁", "8621", "朝阳", "862113"},
	{"辽宁", "8621", "葫芦岛", "862114"},
	{"吉林", "8622", "长春", "862201"},
	{"吉林", "8622", "吉林", "862202"},
	{"吉林", "8622", "四平", "862203"},
	{"吉林", "8622", "辽源", "862204"},
	{"吉林", "8622", "通化", "862205"},
	{"吉林", "8622", "白山", "862206"},
	{"吉林", "8622", "松原", "862207"},
	{"吉林", "8622", "白城", "862208"},
	{"吉林", "8622", "延边朝鲜族自治州", "862224"},
	{"黑龙江", "8623", "哈尔滨", "862301"},
	{"黑龙江", "8623", "齐齐哈尔", "862302"},
	{"黑龙江", "8623", "鸡西", "862303"},
	{"黑龙江", "8623", "鹤岗", "862304"},
	{"黑龙江", "8623", "双鸭山", "862305"},
	{"黑龙江", "8623", "大庆", "862306"},
	{"黑龙江", "8623", "伊春", "862307"},
	{"黑龙江", "8623", "佳木斯", "862308"},
	{"黑龙江", "8623", "七台河", "862309"},
	{"黑龙江", "8623", "牡丹江", "862310"},
	{"黑龙江", "8623", "黑河", "862311"},
	{"黑龙江", "8623", "绥化", "862312"},
	{"黑龙江", "8623", "大兴安岭地", "862327"},
	{"上海", "8631", "黄浦", "863101"},
	{"上海", "8631", "卢湾", "863102"},
	{"上海", "8631", "徐汇", "863103"},
	{"上海", "8631", "长宁", "863104"},
	{"上海", "8631", "静安", "863105"},
	{"上海", "8631", "普陀", "863106"},
	{"上海", "8631", "闸北", "863107"},
	{"上海", "8631", "虹口", "863108"},
	{"上海", "8631", "杨浦", "863109"},
	{"上海", "8631", "闵行", "863110"},
	{"上海", "8631", "宝山", "863111"},
	{"上海", "8631", "嘉定", "863112"},
	{"上海", "8631", "浦东新", "863113"},
	{"上海", "8631", "金山", "863114"},
	{"上海", "8631", "松江", "863115"},
	{"上海", "8631", "青浦", "863116"},
	{"上海", "8631", "南汇", "863117"},
	{"上海", "8631", "奉贤", "863118"},
	{"上海", "8631", "崇明", "863119"},
	{"江苏", "8632", "南京", "863201"},
	{"江苏", "8632", "无锡", "863202"},
	{"江苏", "8632", "徐州", "863203"},
	{"江苏", "8632", "常州", "863204"},
	{"江苏", "8632", "苏州", "863205"},
	{"江苏", "8632", "南通", "863206"},
	{"江苏", "8632", "连云港", "863207"},
	{"江苏", "8632", "淮安", "863208"},
	{"江苏", "8632", "盐城", "863209"},
	{"江苏", "8632", "扬州", "863210"},
	{"江苏", "8632", "镇江", "863211"},
	{"江苏", "8632", "泰州", "863212"},
	{"江苏", "8632", "宿迁", "863213"},
	{"浙江", "8633", "杭州", "863301"},
	{"浙江", "8633", "宁波", "863302"},
	{"浙江", "8633", "温州", "863303"},
	{"浙江", "8633", "嘉兴", "863304"},
	{"浙江", "8633", "湖州", "863305"},
	{"浙江", "8633", "绍兴", "863306"},
	{"浙江", "8633", "金华", "863307"},
	{"浙江", "8633", "衢州", "863308"},
	{"浙江", "8633", "舟山", "863309"},
	{"浙江", "8633", "台州", "863310"},
	{"浙江", "8633", "丽水", "863311"},
	{"安徽", "8634", "合肥", "863401"},
	{"安徽", "8634", "芜湖", "863402"},
	{"安徽", "8634", "蚌埠", "863403"},
	{"安徽", "8634", "淮南", "863404"},
	{"安徽", "8634", "马鞍山", "863405"},
	{"安徽", "8634", "淮北", "863406"},
	{"安徽", "8634", "铜陵", "863407"},
	{"安徽", "8634", "安庆", "863408"},
	{"安徽", "8634", "黄山", "863410"},
	{"安徽", "8634", "滁州", "863411"},
	{"安徽", "8634", "阜阳", "863412"},
	{"安徽", "8634", "宿州", "863413"},
	{"安徽", "8634", "巢湖", "863414"},
	{"安徽", "8634", "六安", "863415"},
	{"安徽", "8634", "亳州", "863416"},
	{"安徽", "8634", "池州", "863417"},
	{"安徽", "8634", "宣城", "863418"},
	{"福建", "8635", "福州", "863501"},
	{"福建", "8635", "厦门", "863502"},
	{"福建", "8635", "莆田", "863503"},
	{"福建", "8635", "三明", "863504"},
	{"福建", "8635", "泉州", "863505"},
	{"福建", "8635", "漳州", "863506"},
	{"福建", "8635", "南平", "863507"},
	{"福建", "8635", "龙岩", "863508"},
	{"福建", "8635", "宁德", "863509"},
	{"江西", "8636", "南昌", "863601"},
	{"江西", "8636", "景德镇", "863602"},
	{"江西", "8636", "萍乡", "863603"},
	{"江西", "8636", "九江", "863604"},
	{"江西", "8636", "新余", "863605"},
	{"江西", "8636", "鹰潭", "863606"},
	{"江西", "8636", "赣州", "863607"},
	{"江西", "8636", "吉安", "863608"},
	{"江西", "8636", "宜春", "863609"},
	{"江西", "8636", "抚州", "863610"},
	{"江西", "8636", "上饶", "863611"},
	{"山东", "8637", "济南", "863701"},
	{"山东", "8637", "青岛", "863702"},
	{"山东", "8637", "淄博", "863703"},
	{"山东", "8637", "枣庄", "863704"},
	{"山东", "8637", "东营", "863705"},
	{"山东", "8637", "烟台", "863706"},
	{"山东", "8637", "潍坊", "863707"},
	{"山东", "8637", "济宁", "863708"},
	{"山东", "8637", "泰安", "863709"},
	{"山东", "8637", "威海", "863710"},
	{"山东", "8637", "日照", "863711"},
	{"山东", "8637", "莱芜", "863712"},
	{"山东", "8637", "临沂", "863713"},
	{"山东", "8637", "德州", "863714"},
	{"山东", "8637", "聊城", "863715"},
	{"山东", "8637", "滨州", "863716"},
	{"山东", "8637", "菏泽", "863717"},
	{"河南", "8641", "郑州", "864101"},
	{"河南", "8641", "开封", "864102"},
	{"河南", "8641", "洛阳", "864103"},
	{"河南", "8641", "平顶山", "864104"},
	{"河南", "8641", "安阳", "864105"},
	{"河南", "8641", "鹤壁", "864106"},
	{"河南", "8641", "新乡", "864107"},
	{"河南", "8641", "焦作", "864108"},
	{"河南", "8641", "濮阳", "864109"},
	{"河南", "8641", "许昌", "864110"},
	{"河南", "8641", "漯河", "864111"},
	{"河南", "8641", "三门峡", "864112"},
	{"河南", "8641", "南阳", "864113"},
	{"河南", "8641", "商丘", "864114"},
	{"河南", "8641", "信阳", "864115"},
	{"河南", "8641", "周口", "864116"},
	{"河南", "8641", "驻马店", "864117"},
	{"湖北", "8642", "武汉", "864201"},
	{"湖北", "8642", "黄石", "864202"},
	{"湖北", "8642", "十堰", "864203"},
	{"湖北", "8642", "宜昌", "864205"},
	{"湖北", "8642", "襄樊", "864206"},
	{"湖北", "8642", "鄂州", "864207"},
	{"湖北", "8642", "荆门", "864208"},
	{"湖北", "8642", "孝感", "864209"},
	{"湖北", "8642", "荆州", "864210"},
	{"湖北", "8642", "黄冈", "864211"},
	{"湖北", "8642", "咸宁", "864212"},
	{"湖北", "8642", "随州", "864213"},
	{"湖北", "8642", "恩施土家族苗族自治州", "864228"},
	{"湖北", "8642", "", "864290"},
	{"湖南", "8643", "长沙", "864301"},
	{"湖南", "8643", "株洲", "864302"},
	{"湖南", "8643", "湘潭", "864303"},
	{"湖南", "8643", "衡阳", "864304"},
	{"湖南", "8643", "邵阳", "864305"},
	{"湖南", "8643", "岳阳", "864306"},
	{"湖南", "8643", "常德", "864307"},
	{"湖南", "8643", "张家界", "864308"},
	{"湖南", "8643", "益阳", "864309"},
	{"湖南", "8643", "郴州", "864310"},
	{"湖南", "8643", "永州", "864311"},
	{"湖南", "8643", "怀化", "864312"},
	{"湖南", "8643", "娄底", "864313"},
	{"湖南", "8643", "湘西土家族苗族自治州", "864331"},
	{"广东", "8644", "广州", "864401"},
	{"广东", "8644", "韶关", "864402"},
	{"广东", "8644", "深圳", "864403"},
	{"广东", "8644", "珠海", "864404"},
	{"广东", "8644", "汕头", "864405"},
	{"广东", "8644", "佛山", "864406"},
	{"广东", "8644", "江门", "864407"},
	{"广东", "8644", "湛江", "864408"},
	{"广东", "8644", "茂名", "864409"},
	{"广东", "8644", "肇庆", "864412"},
	{"广东", "8644", "惠州", "864413"},
	{"广东", "8644", "梅州", "864414"},
	{"广东", "8644", "汕尾", "864415"},
	{"广东", "8644", "河源", "864416"},
	{"广东", "8644", "阳江", "864417"},
	{"广东", "8644", "清远", "864418"},
	{"广东", "8644", "东莞", "864419"},
	{"广东", "8644", "中山", "864420"},
	{"广东", "8644", "潮州", "864451"},
	{"广东", "8644", "揭阳", "864452"},
	{"广东", "8644", "云浮", "864453"},
	{"广西自治", "8645", "南宁", "864501"},
	{"广西自治", "8645", "柳州", "864502"},
	{"广西自治", "8645", "桂林", "864503"},
	{"广西自治", "8645", "梧州", "864504"},
	{"广西自治", "8645", "北海", "864505"},
	{"广西自治", "8645", "防城港", "864506"},
	{"广西自治", "8645", "钦州", "864507"},
	{"广西自治", "8645", "贵港", "864508"},
	{"广西自治", "8645", "玉林", "864509"},
	{"广西自治", "8645", "百色", "864510"},
	{"广西自治", "8645", "贺州", "864511"},
	{"广西自治", "8645", "河池", "864512"},
	{"广西自治", "8645", "来宾", "864513"},
	{"广西自治", "8645", "崇左", "864514"},
	{"海南", "8646", "海口", "864601"},
	{"海南", "8646", "三亚", "864602"},
	{"海南", "8646", "", "864690"},
	{"重庆", "8650", "万州", "865001"},
	{"重庆", "8650", "涪陵", "865002"},
	{"重庆", "8650", "渝中", "865003"},
	{"重庆", "8650", "大渡口", "865004"},
	{"重庆", "8650", "江北", "865005"},
	{"重庆", "8650", "沙坪坝", "865006"},
	{"重庆", "8650", "九龙坡", "865007"},
	{"重庆", "8650", "南岸", "865008"},
	{"重庆", "8650", "北碚", "865009"},
	{"重庆", "8650", "万盛", "865010"},
	{"重庆", "8650", "双桥", "865011"},
	{"重庆", "8650", "渝北", "865012"},
	{"重庆", "8650", "巴南", "865013"},
	{"重庆", "8650", "黔江", "865014"},
	{"重庆", "8650", "长寿", "865015"},
	{"重庆", "8650", "綦江", "865016"},
	{"重庆", "8650", "潼南", "865017"},
	{"重庆", "8650", "铜梁", "865018"},
	{"重庆", "8650", "大足", "865019"},
	{"重庆", "8650", "荣昌", "865020"},
	{"重庆", "8650", "璧山", "865021"},
	{"重庆", "8650", "梁平", "865022"},
	{"重庆", "8650", "城口", "865023"},
	{"重庆", "8650", "丰都", "865024"},
	{"重庆", "8650", "垫江", "865025"},
	{"重庆", "8650", "武隆", "865026"},
	{"重庆", "8650", "忠", "865027"},
	{"重庆", "8650", "开", "865028"},
	{"重庆", "8650", "云阳", "865029"},
	{"重庆", "8650", "奉节", "865030"},
	{"重庆", "8650", "巫山", "865031"},
	{"重庆", "8650", "巫溪", "865032"},
	{"重庆", "8650", "石柱土家族自治", "865033"},
	{"重庆", "8650", "秀山土家族苗族自治", "865034"},
	{"重庆", "8650", "酉阳土家族苗族自治", "865035"},
	{"重庆", "8650", "彭水苗族土家族自治", "865036"},
	{"重庆", "8650", "江津", "865037"},
	{"重庆", "8650", "合川", "865038"},
	{"重庆", "8650", "永川", "865039"},
	{"重庆", "8650", "南川", "865040"},
	{"四川", "8651", "成都", "865101"},
	{"四川", "8651", "自贡", "865103"},
	{"四川", "8651", "攀枝花", "865104"},
	{"四川", "8651", "泸州", "865105"},
	{"四川", "8651", "德阳", "865106"},
	{"四川", "8651", "绵阳", "865107"},
	{"四川", "8651", "广元", "865108"},
	{"四川", "8651", "遂宁", "865109"},
	{"四川", "8651", "内江", "865110"},
	{"四川", "8651", "乐山", "865111"},
	{"四川", "8651", "南充", "865113"},
	{"四川", "8651", "眉山", "865114"},
	{"四川", "8651", "宜宾", "865115"},
	{"四川", "8651", "广安", "865116"},
	{"四川", "8651", "达州", "865117"},
	{"四川", "8651", "雅安", "865118"},
	{"四川", "8651", "巴中", "865119"},
	{"四川", "8651", "资阳", "865120"},
	{"四川", "8651", "阿坝藏族羌族自治州", "865132"},
	{"四川", "8651", "甘孜藏族自治州", "865133"},
	{"四川", "8651", "凉山彝族自治州", "865134"},
	{"贵州", "8652", "贵阳", "865201"},
	{"贵州", "8652", "六盘水", "865202"},
	{"贵州", "8652", "遵义", "865203"},
	{"贵州", "8652", "安顺", "865204"},
	{"贵州", "8652", "铜仁地", "865222"},
	{"贵州", "8652", "黔西南布依族苗族自治州", "865223"},
	{"贵州", "8652", "毕节地", "865224"},
	{"贵州", "8652", "黔东南苗族侗族自治州", "865226"},
	{"贵州", "8652", "黔南布依族苗族自治州", "865227"},
	{"云南", "8653", "昆明", "865301"},
	{"云南", "8653", "曲靖", "865303"},
	{"云南", "8653", "玉溪", "865304"},
	{"云南", "8653", "保山", "865305"},
	{"云南", "8653", "昭通", "865306"},
	{"云南", "8653", "丽江", "865307"},
	{"云南", "8653", "思茅", "865308"},
	{"云南", "8653", "临沧", "865309"},
	{"云南", "8653", "楚雄彝族自治州", "865323"},
	{"云南", "8653", "红河哈尼族彝族自治州", "865325"},
	{"云南", "8653", "文山壮族苗族自治州", "865326"},
	{"云南", "8653", "西双版纳傣族自治州", "865328"},
	{"云南", "8653", "大理白族自治州", "865329"},
	{"云南", "8653", "德宏傣族景颇族自治州", "865331"},
	{"云南", "8653", "怒江傈僳族自治州", "865333"},
	{"云南", "8653", "迪庆藏族自治州", "865334"},
	{"西藏自治", "8654", "拉萨", "865401"},
	{"西藏自治", "8654", "昌都地", "865421"},
	{"西藏自治", "8654", "山南地", "865422"},
	{"西藏自治", "8654", "日喀则地", "865423"},
	{"西藏自治", "8654", "那曲地", "865424"},
	{"西藏自治", "8654", "阿里地", "865425"},
	{"西藏自治", "8654", "林芝地", "865426"},
	{"陕西", "8661", "西安", "866101"},
	{"陕西", "8661", "铜川", "866102"},
	{"陕西", "8661", "宝鸡", "866103"},
	{"陕西", "8661", "咸阳", "866104"},
	{"陕西", "8661", "渭南", "866105"},
	{"陕西", "8661", "延安", "866106"},
	{"陕西", "8661", "汉中", "866107"},
	{"陕西", "8661", "榆林", "866108"},
	{"陕西", "8661", "安康", "866109"},
	{"陕西", "8661", "商洛", "866110"},
	{"甘肃", "8662", "兰州", "866201"},
	{"甘肃", "8662", "嘉峪关", "866202"},
	{"甘肃", "8662", "金昌", "866203"},
	{"甘肃", "8662", "白银", "866204"},
	{"甘肃", "8662", "天水", "866205"},
	{"甘肃", "8662", "武威", "866206"},
	{"甘肃", "8662", "张掖", "866207"},
	{"甘肃", "8662", "平凉", "866208"},
	{"甘肃", "8662", "酒泉", "866209"},
	{"甘肃", "8662", "庆阳", "866210"},
	{"甘肃", "8662", "定西", "866211"},
	{"甘肃", "8662", "陇南", "866212"},
	{"甘肃", "8662", "临夏回族自治州", "866229"},
	{"甘肃", "8662", "甘南藏族自治州", "866230"},
	{"青海", "8663", "西宁", "866301"},
	{"青海", "8663", "海东地", "866321"},
	{"青海", "8663", "海北藏族自治州", "866322"},
	{"青海", "8663", "黄南藏族自治州", "866323"},
	{"青海", "8663", "海南藏族自治州", "866325"},
	{"青海", "8663", "果洛藏族自治州", "866326"},
	{"青海", "8663", "玉树藏族自治州", "866327"},
	{"青海", "8663", "海西蒙古族藏族自治州", "866328"},
	{"宁夏自治", "8664", "银川", "866401"},
	{"宁夏自治", "8664", "石嘴山", "866402"},
	{"宁夏自治", "8664", "吴忠", "866403"},
	{"宁夏自治", "8664", "固原", "866404"},
	{"宁夏自治", "8664", "中卫", "866405"},
	{"新疆自治", "8665", "乌鲁木齐", "866501"},
	{"新疆自治", "8665", "克拉玛依", "866502"},
	{"新疆自治", "8665", "吐鲁番地", "866521"},
	{"新疆自治", "8665", "哈密地", "866522"},
	{"新疆自治", "8665", "昌吉回族自治州", "866523"},
	{"新疆自治", "8665", "博尔塔拉蒙古自治州", "866527"},
	{"新疆自治", "8665", "巴音郭楞蒙古自治州", "866528"},
	{"新疆自治", "8665", "阿克苏地", "866529"},
	{"新疆自治", "8665", "克孜勒苏柯尔克孜自治州", "866530"},
	{"新疆自治", "8665", "喀什地", "866531"},
	{"新疆自治", "8665", "和田地", "866532"},
	{"新疆自治", "8665", "伊犁哈萨克自治州", "866540"},
	{"新疆自治", "8665", "塔城地", "866542"},
	{"新疆自治", "8665", "阿勒泰地", "866543"},
	{"新疆自治", "8665", "", "866590"},
	{"台湾", "8671", "新竹", "867101"},
	{"台湾", "8671", "高雄", "867102"},
	{"台湾", "8671", "基隆", "867103"},
	{"台湾", "8671", "台中", "867104"},
	{"台湾", "8671", "台南", "867105"},
	{"台湾", "8671", "台北", "867106"},
	{"台湾", "8671", "桃园", "867107"},
	{"台湾", "8671", "云林", "867108"}}

// CommonDistrcitInfos is the common district infos
var CommonDistrcitInfos = []districtInfo{
	{"北京市", "8611", "东城区", "861101"},
	{"北京市", "8611", "西城区", "861102"},
	{"北京市", "8611", "崇文区", "861103"},
	{"北京市", "8611", "宣武区", "861104"},
	{"北京市", "8611", "朝阳区", "861105"},
	{"北京市", "8611", "丰台区", "861106"},
	{"北京市", "8611", "石景山区", "861107"},
	{"北京市", "8611", "海淀区", "861108"},
	{"北京市", "8611", "门头沟区", "861109"},
	{"北京市", "8611", "房山区", "861110"},
	{"北京市", "8611", "通州区", "861111"},
	{"北京市", "8611", "顺义区", "861112"},
	{"北京市", "8611", "昌平区", "861113"},
	{"北京市", "8611", "大兴区", "861114"},
	{"北京市", "8611", "怀柔区", "861115"},
	{"北京市", "8611", "平谷区", "861116"},
	{"北京市", "8611", "密云县", "861117"},
	{"北京市", "8611", "延庆县", "861118"},
	{"天津市", "8612", "和平区", "861201"},
	{"天津市", "8612", "河东区", "861202"},
	{"天津市", "8612", "河西区", "861203"},
	{"天津市", "8612", "南开区", "861204"},
	{"天津市", "8612", "河北区", "861205"},
	{"天津市", "8612", "红桥区", "861206"},
	{"天津市", "8612", "塘沽区", "861207"},
	{"天津市", "8612", "汉沽区", "861208"},
	{"天津市", "8612", "大港区", "861209"},
	{"天津市", "8612", "东丽区", "861210"},
	{"天津市", "8612", "西青区", "861211"},
	{"天津市", "8612", "津南区", "861212"},
	{"天津市", "8612", "北辰区", "861213"},
	{"天津市", "8612", "武清区", "861214"},
	{"天津市", "8612", "宝坻区", "861215"},
	{"天津市", "8612", "宁河县", "861216"},
	{"天津市", "8612", "静海县", "861217"},
	{"天津市", "8612", "蓟县", "861218"},
	{"河北省", "8613", "石家庄市", "861301"},
	{"河北省", "8613", "唐山市", "861302"},
	{"河北省", "8613", "秦皇岛市", "861303"},
	{"河北省", "8613", "邯郸市", "861304"},
	{"河北省", "8613", "邢台市", "861305"},
	{"河北省", "8613", "保定市", "861306"},
	{"河北省", "8613", "张家口市", "861307"},
	{"河北省", "8613", "承德市", "861308"},
	{"河北省", "8613", "沧州市", "861309"},
	{"河北省", "8613", "廊坊市", "861310"},
	{"河北省", "8613", "衡水市", "861311"},
	{"山西省", "8614", "太原市", "861401"},
	{"山西省", "8614", "大同市", "861402"},
	{"山西省", "8614", "阳泉市", "861403"},
	{"山西省", "8614", "长治市", "861404"},
	{"山西省", "8614", "晋城市", "861405"},
	{"山西省", "8614", "朔州市", "861406"},
	{"山西省", "8614", "晋中市", "861407"},
	{"山西省", "8614", "运城市", "861408"},
	{"山西省", "8614", "忻州市", "861409"},
	{"山西省", "8614", "临汾市", "861410"},
	{"山西省", "8614", "吕梁市", "861411"},
	{"内蒙古自治区", "8615", "呼和浩特市", "861501"},
	{"内蒙古自治区", "8615", "包头市", "861502"},
	{"内蒙古自治区", "8615", "乌海市", "861503"},
	{"内蒙古自治区", "8615", "赤峰市", "861504"},
	{"内蒙古自治区", "8615", "通辽市", "861505"},
	{"内蒙古自治区", "8615", "鄂尔多斯市", "861506"},
	{"内蒙古自治区", "8615", "呼伦贝尔市", "861507"},
	{"内蒙古自治区", "8615", "巴彦淖尔市", "861508"},
	{"内蒙古自治区", "8615", "乌兰察布市", "861509"},
	{"内蒙古自治区", "8615", "兴安盟", "861522"},
	{"内蒙古自治区", "8615", "锡林郭勒盟", "861525"},
	{"内蒙古自治区", "8615", "阿拉善盟", "861529"},
	{"辽宁省", "8621", "沈阳市", "862101"},
	{"辽宁省", "8621", "大连市", "862102"},
	{"辽宁省", "8621", "鞍山市", "862103"},
	{"辽宁省", "8621", "抚顺市", "862104"},
	{"辽宁省", "8621", "本溪市", "862105"},
	{"辽宁省", "8621", "丹东市", "862106"},
	{"辽宁省", "8621", "锦州市", "862107"},
	{"辽宁省", "8621", "营口市", "862108"},
	{"辽宁省", "8621", "阜新市", "862109"},
	{"辽宁省", "8621", "辽阳市", "862110"},
	{"辽宁省", "8621", "盘锦市", "862111"},
	{"辽宁省", "8621", "铁岭市", "862112"},
	{"辽宁省", "8621", "朝阳市", "862113"},
	{"辽宁省", "8621", "葫芦岛市", "862114"},
	{"吉林省", "8622", "长春市", "862201"},
	{"吉林省", "8622", "吉林市", "862202"},
	{"吉林省", "8622", "四平市", "862203"},
	{"吉林省", "8622", "辽源市", "862204"},
	{"吉林省", "8622", "通化市", "862205"},
	{"吉林省", "8622", "白山市", "862206"},
	{"吉林省", "8622", "松原市", "862207"},
	{"吉林省", "8622", "白城市", "862208"},
	{"吉林省", "8622", "延边朝鲜族自治州", "862224"},
	{"黑龙江省", "8623", "哈尔滨市", "862301"},
	{"黑龙江省", "8623", "齐齐哈尔市", "862302"},
	{"黑龙江省", "8623", "鸡西市", "862303"},
	{"黑龙江省", "8623", "鹤岗市", "862304"},
	{"黑龙江省", "8623", "双鸭山市", "862305"},
	{"黑龙江省", "8623", "大庆市", "862306"},
	{"黑龙江省", "8623", "伊春市", "862307"},
	{"黑龙江省", "8623", "佳木斯市", "862308"},
	{"黑龙江省", "8623", "七台河市", "862309"},
	{"黑龙江省", "8623", "牡丹江市", "862310"},
	{"黑龙江省", "8623", "黑河市", "862311"},
	{"黑龙江省", "8623", "绥化市", "862312"},
	{"黑龙江省", "8623", "大兴安岭地区", "862327"},
	{"上海市", "8631", "黄浦区", "863101"},
	{"上海市", "8631", "卢湾区", "863102"},
	{"上海市", "8631", "徐汇区", "863103"},
	{"上海市", "8631", "长宁区", "863104"},
	{"上海市", "8631", "静安区", "863105"},
	{"上海市", "8631", "普陀区", "863106"},
	{"上海市", "8631", "闸北区", "863107"},
	{"上海市", "8631", "虹口区", "863108"},
	{"上海市", "8631", "杨浦区", "863109"},
	{"上海市", "8631", "闵行区", "863110"},
	{"上海市", "8631", "宝山区", "863111"},
	{"上海市", "8631", "嘉定区", "863112"},
	{"上海市", "8631", "浦东新区", "863113"},
	{"上海市", "8631", "金山区", "863114"},
	{"上海市", "8631", "松江区", "863115"},
	{"上海市", "8631", "青浦区", "863116"},
	{"上海市", "8631", "南汇区", "863117"},
	{"上海市", "8631", "奉贤区", "863118"},
	{"上海市", "8631", "崇明县", "863119"},
	{"江苏省", "8632", "南京市", "863201"},
	{"江苏省", "8632", "无锡市", "863202"},
	{"江苏省", "8632", "徐州市", "863203"},
	{"江苏省", "8632", "常州市", "863204"},
	{"江苏省", "8632", "苏州市", "863205"},
	{"江苏省", "8632", "南通市", "863206"},
	{"江苏省", "8632", "连云港市", "863207"},
	{"江苏省", "8632", "淮安市", "863208"},
	{"江苏省", "8632", "盐城市", "863209"},
	{"江苏省", "8632", "扬州市", "863210"},
	{"江苏省", "8632", "镇江市", "863211"},
	{"江苏省", "8632", "泰州市", "863212"},
	{"江苏省", "8632", "宿迁市", "863213"},
	{"浙江省", "8633", "杭州市", "863301"},
	{"浙江省", "8633", "宁波市", "863302"},
	{"浙江省", "8633", "温州市", "863303"},
	{"浙江省", "8633", "嘉兴市", "863304"},
	{"浙江省", "8633", "湖州市", "863305"},
	{"浙江省", "8633", "绍兴市", "863306"},
	{"浙江省", "8633", "金华市", "863307"},
	{"浙江省", "8633", "衢州市", "863308"},
	{"浙江省", "8633", "舟山市", "863309"},
	{"浙江省", "8633", "台州市", "863310"},
	{"浙江省", "8633", "丽水市", "863311"},
	{"安徽省", "8634", "合肥市", "863401"},
	{"安徽省", "8634", "芜湖市", "863402"},
	{"安徽省", "8634", "蚌埠市", "863403"},
	{"安徽省", "8634", "淮南市", "863404"},
	{"安徽省", "8634", "马鞍山市", "863405"},
	{"安徽省", "8634", "淮北市", "863406"},
	{"安徽省", "8634", "铜陵市", "863407"},
	{"安徽省", "8634", "安庆市", "863408"},
	{"安徽省", "8634", "黄山市", "863410"},
	{"安徽省", "8634", "滁州市", "863411"},
	{"安徽省", "8634", "阜阳市", "863412"},
	{"安徽省", "8634", "宿州市", "863413"},
	{"安徽省", "8634", "巢湖市", "863414"},
	{"安徽省", "8634", "六安市", "863415"},
	{"安徽省", "8634", "亳州市", "863416"},
	{"安徽省", "8634", "池州市", "863417"},
	{"安徽省", "8634", "宣城市", "863418"},
	{"福建省", "8635", "福州市", "863501"},
	{"福建省", "8635", "厦门市", "863502"},
	{"福建省", "8635", "莆田市", "863503"},
	{"福建省", "8635", "三明市", "863504"},
	{"福建省", "8635", "泉州市", "863505"},
	{"福建省", "8635", "漳州市", "863506"},
	{"福建省", "8635", "南平市", "863507"},
	{"福建省", "8635", "龙岩市", "863508"},
	{"福建省", "8635", "宁德市", "863509"},
	{"江西省", "8636", "南昌市", "863601"},
	{"江西省", "8636", "景德镇市", "863602"},
	{"江西省", "8636", "萍乡市", "863603"},
	{"江西省", "8636", "九江市", "863604"},
	{"江西省", "8636", "新余市", "863605"},
	{"江西省", "8636", "鹰潭市", "863606"},
	{"江西省", "8636", "赣州市", "863607"},
	{"江西省", "8636", "吉安市", "863608"},
	{"江西省", "8636", "宜春市", "863609"},
	{"江西省", "8636", "抚州市", "863610"},
	{"江西省", "8636", "上饶市", "863611"},
	{"山东省", "8637", "济南市", "863701"},
	{"山东省", "8637", "青岛市", "863702"},
	{"山东省", "8637", "淄博市", "863703"},
	{"山东省", "8637", "枣庄市", "863704"},
	{"山东省", "8637", "东营市", "863705"},
	{"山东省", "8637", "烟台市", "863706"},
	{"山东省", "8637", "潍坊市", "863707"},
	{"山东省", "8637", "济宁市", "863708"},
	{"山东省", "8637", "泰安市", "863709"},
	{"山东省", "8637", "威海市", "863710"},
	{"山东省", "8637", "日照市", "863711"},
	{"山东省", "8637", "莱芜市", "863712"},
	{"山东省", "8637", "临沂市", "863713"},
	{"山东省", "8637", "德州市", "863714"},
	{"山东省", "8637", "聊城市", "863715"},
	{"山东省", "8637", "滨州市", "863716"},
	{"山东省", "8637", "菏泽市", "863717"},
	{"河南省", "8641", "郑州市", "864101"},
	{"河南省", "8641", "开封市", "864102"},
	{"河南省", "8641", "洛阳市", "864103"},
	{"河南省", "8641", "平顶山市", "864104"},
	{"河南省", "8641", "安阳市", "864105"},
	{"河南省", "8641", "鹤壁市", "864106"},
	{"河南省", "8641", "新乡市", "864107"},
	{"河南省", "8641", "焦作市", "864108"},
	{"河南省", "8641", "濮阳市", "864109"},
	{"河南省", "8641", "许昌市", "864110"},
	{"河南省", "8641", "漯河市", "864111"},
	{"河南省", "8641", "三门峡市", "864112"},
	{"河南省", "8641", "南阳市", "864113"},
	{"河南省", "8641", "商丘市", "864114"},
	{"河南省", "8641", "信阳市", "864115"},
	{"河南省", "8641", "周口市", "864116"},
	{"河南省", "8641", "驻马店市", "864117"},
	{"湖北省", "8642", "武汉市", "864201"},
	{"湖北省", "8642", "黄石市", "864202"},
	{"湖北省", "8642", "十堰市", "864203"},
	{"湖北省", "8642", "宜昌市", "864205"},
	{"湖北省", "8642", "襄樊市", "864206"},
	{"湖北省", "8642", "鄂州市", "864207"},
	{"湖北省", "8642", "荆门市", "864208"},
	{"湖北省", "8642", "孝感市", "864209"},
	{"湖北省", "8642", "荆州市", "864210"},
	{"湖北省", "8642", "黄冈市", "864211"},
	{"湖北省", "8642", "咸宁市", "864212"},
	{"湖北省", "8642", "随州市", "864213"},
	{"湖北省", "8642", "恩施土家族苗族自治州", "864228"},
	{"湖北省", "8642", "", "864290"},
	{"湖南省", "8643", "长沙市", "864301"},
	{"湖南省", "8643", "株洲市", "864302"},
	{"湖南省", "8643", "湘潭市", "864303"},
	{"湖南省", "8643", "衡阳市", "864304"},
	{"湖南省", "8643", "邵阳市", "864305"},
	{"湖南省", "8643", "岳阳市", "864306"},
	{"湖南省", "8643", "常德市", "864307"},
	{"湖南省", "8643", "张家界市", "864308"},
	{"湖南省", "8643", "益阳市", "864309"},
	{"湖南省", "8643", "郴州市", "864310"},
	{"湖南省", "8643", "永州市", "864311"},
	{"湖南省", "8643", "怀化市", "864312"},
	{"湖南省", "8643", "娄底市", "864313"},
	{"湖南省", "8643", "湘西土家族苗族自治州", "864331"},
	{"广东省", "8644", "广州市", "864401"},
	{"广东省", "8644", "韶关市", "864402"},
	{"广东省", "8644", "深圳市", "864403"},
	{"广东省", "8644", "珠海市", "864404"},
	{"广东省", "8644", "汕头市", "864405"},
	{"广东省", "8644", "佛山市", "864406"},
	{"广东省", "8644", "江门市", "864407"},
	{"广东省", "8644", "湛江市", "864408"},
	{"广东省", "8644", "茂名市", "864409"},
	{"广东省", "8644", "肇庆市", "864412"},
	{"广东省", "8644", "惠州市", "864413"},
	{"广东省", "8644", "梅州市", "864414"},
	{"广东省", "8644", "汕尾市", "864415"},
	{"广东省", "8644", "河源市", "864416"},
	{"广东省", "8644", "阳江市", "864417"},
	{"广东省", "8644", "清远市", "864418"},
	{"广东省", "8644", "东莞市", "864419"},
	{"广东省", "8644", "中山市", "864420"},
	{"广东省", "8644", "潮州市", "864451"},
	{"广东省", "8644", "揭阳市", "864452"},
	{"广东省", "8644", "云浮市", "864453"},
	{"广西壮族自治区", "8645", "南宁市", "864501"},
	{"广西壮族自治区", "8645", "柳州市", "864502"},
	{"广西壮族自治区", "8645", "桂林市", "864503"},
	{"广西壮族自治区", "8645", "梧州市", "864504"},
	{"广西壮族自治区", "8645", "北海市", "864505"},
	{"广西壮族自治区", "8645", "防城港市", "864506"},
	{"广西壮族自治区", "8645", "钦州市", "864507"},
	{"广西壮族自治区", "8645", "贵港市", "864508"},
	{"广西壮族自治区", "8645", "玉林市", "864509"},
	{"广西壮族自治区", "8645", "百色市", "864510"},
	{"广西壮族自治区", "8645", "贺州市", "864511"},
	{"广西壮族自治区", "8645", "河池市", "864512"},
	{"广西壮族自治区", "8645", "来宾市", "864513"},
	{"广西壮族自治区", "8645", "崇左市", "864514"},
	{"海南省", "8646", "海口市", "864601"},
	{"海南省", "8646", "三亚市", "864602"},
	{"海南省", "8646", "", "864690"},
	{"重庆市", "8650", "万州区", "865001"},
	{"重庆市", "8650", "涪陵区", "865002"},
	{"重庆市", "8650", "渝中区", "865003"},
	{"重庆市", "8650", "大渡口区", "865004"},
	{"重庆市", "8650", "江北区", "865005"},
	{"重庆市", "8650", "沙坪坝区", "865006"},
	{"重庆市", "8650", "九龙坡区", "865007"},
	{"重庆市", "8650", "南岸区", "865008"},
	{"重庆市", "8650", "北碚区", "865009"},
	{"重庆市", "8650", "万盛区", "865010"},
	{"重庆市", "8650", "双桥区", "865011"},
	{"重庆市", "8650", "渝北区", "865012"},
	{"重庆市", "8650", "巴南区", "865013"},
	{"重庆市", "8650", "黔江区", "865014"},
	{"重庆市", "8650", "长寿区", "865015"},
	{"重庆市", "8650", "綦江县", "865016"},
	{"重庆市", "8650", "潼南县", "865017"},
	{"重庆市", "8650", "铜梁县", "865018"},
	{"重庆市", "8650", "大足县", "865019"},
	{"重庆市", "8650", "荣昌县", "865020"},
	{"重庆市", "8650", "璧山县", "865021"},
	{"重庆市", "8650", "梁平县", "865022"},
	{"重庆市", "8650", "城口县", "865023"},
	{"重庆市", "8650", "丰都县", "865024"},
	{"重庆市", "8650", "垫江县", "865025"},
	{"重庆市", "8650", "武隆县", "865026"},
	{"重庆市", "8650", "忠县", "865027"},
	{"重庆市", "8650", "开县", "865028"},
	{"重庆市", "8650", "云阳县", "865029"},
	{"重庆市", "8650", "奉节县", "865030"},
	{"重庆市", "8650", "巫山县", "865031"},
	{"重庆市", "8650", "巫溪县", "865032"},
	{"重庆市", "8650", "石柱土家族自治县", "865033"},
	{"重庆市", "8650", "秀山土家族苗族自治县", "865034"},
	{"重庆市", "8650", "酉阳土家族苗族自治县", "865035"},
	{"重庆市", "8650", "彭水苗族土家族自治县", "865036"},
	{"重庆市", "8650", "江津市", "865037"},
	{"重庆市", "8650", "合川市", "865038"},
	{"重庆市", "8650", "永川市", "865039"},
	{"重庆市", "8650", "南川市", "865040"},
	{"四川省", "8651", "成都市", "865101"},
	{"四川省", "8651", "自贡市", "865103"},
	{"四川省", "8651", "攀枝花市", "865104"},
	{"四川省", "8651", "泸州市", "865105"},
	{"四川省", "8651", "德阳市", "865106"},
	{"四川省", "8651", "绵阳市", "865107"},
	{"四川省", "8651", "广元市", "865108"},
	{"四川省", "8651", "遂宁市", "865109"},
	{"四川省", "8651", "内江市", "865110"},
	{"四川省", "8651", "乐山市", "865111"},
	{"四川省", "8651", "南充市", "865113"},
	{"四川省", "8651", "眉山市", "865114"},
	{"四川省", "8651", "宜宾市", "865115"},
	{"四川省", "8651", "广安市", "865116"},
	{"四川省", "8651", "达州市", "865117"},
	{"四川省", "8651", "雅安市", "865118"},
	{"四川省", "8651", "巴中市", "865119"},
	{"四川省", "8651", "资阳市", "865120"},
	{"四川省", "8651", "阿坝藏族羌族自治州", "865132"},
	{"四川省", "8651", "甘孜藏族自治州", "865133"},
	{"四川省", "8651", "凉山彝族自治州", "865134"},
	{"贵州省", "8652", "贵阳市", "865201"},
	{"贵州省", "8652", "六盘水市", "865202"},
	{"贵州省", "8652", "遵义市", "865203"},
	{"贵州省", "8652", "安顺市", "865204"},
	{"贵州省", "8652", "铜仁地区", "865222"},
	{"贵州省", "8652", "黔西南布依族苗族自治州", "865223"},
	{"贵州省", "8652", "毕节地区", "865224"},
	{"贵州省", "8652", "黔东南苗族侗族自治州", "865226"},
	{"贵州省", "8652", "黔南布依族苗族自治州", "865227"},
	{"云南省", "8653", "昆明市", "865301"},
	{"云南省", "8653", "曲靖市", "865303"},
	{"云南省", "8653", "玉溪市", "865304"},
	{"云南省", "8653", "保山市", "865305"},
	{"云南省", "8653", "昭通市", "865306"},
	{"云南省", "8653", "丽江市", "865307"},
	{"云南省", "8653", "思茅市", "865308"},
	{"云南省", "8653", "临沧市", "865309"},
	{"云南省", "8653", "楚雄彝族自治州", "865323"},
	{"云南省", "8653", "红河哈尼族彝族自治州", "865325"},
	{"云南省", "8653", "文山壮族苗族自治州", "865326"},
	{"云南省", "8653", "西双版纳傣族自治州", "865328"},
	{"云南省", "8653", "大理白族自治州", "865329"},
	{"云南省", "8653", "德宏傣族景颇族自治州", "865331"},
	{"云南省", "8653", "怒江傈僳族自治州", "865333"},
	{"云南省", "8653", "迪庆藏族自治州", "865334"},
	{"西藏自治区", "8654", "拉萨市", "865401"},
	{"西藏自治区", "8654", "昌都地区", "865421"},
	{"西藏自治区", "8654", "山南地区", "865422"},
	{"西藏自治区", "8654", "日喀则地区", "865423"},
	{"西藏自治区", "8654", "那曲地区", "865424"},
	{"西藏自治区", "8654", "阿里地区", "865425"},
	{"西藏自治区", "8654", "林芝地区", "865426"},
	{"陕西省", "8661", "西安市", "866101"},
	{"陕西省", "8661", "铜川市", "866102"},
	{"陕西省", "8661", "宝鸡市", "866103"},
	{"陕西省", "8661", "咸阳市", "866104"},
	{"陕西省", "8661", "渭南市", "866105"},
	{"陕西省", "8661", "延安市", "866106"},
	{"陕西省", "8661", "汉中市", "866107"},
	{"陕西省", "8661", "榆林市", "866108"},
	{"陕西省", "8661", "安康市", "866109"},
	{"陕西省", "8661", "商洛市", "866110"},
	{"甘肃省", "8662", "兰州市", "866201"},
	{"甘肃省", "8662", "嘉峪关市", "866202"},
	{"甘肃省", "8662", "金昌市", "866203"},
	{"甘肃省", "8662", "白银市", "866204"},
	{"甘肃省", "8662", "天水市", "866205"},
	{"甘肃省", "8662", "武威市", "866206"},
	{"甘肃省", "8662", "张掖市", "866207"},
	{"甘肃省", "8662", "平凉市", "866208"},
	{"甘肃省", "8662", "酒泉市", "866209"},
	{"甘肃省", "8662", "庆阳市", "866210"},
	{"甘肃省", "8662", "定西市", "866211"},
	{"甘肃省", "8662", "陇南市", "866212"},
	{"甘肃省", "8662", "临夏回族自治州", "866229"},
	{"甘肃省", "8662", "甘南藏族自治州", "866230"},
	{"青海省", "8663", "西宁市", "866301"},
	{"青海省", "8663", "海东地区", "866321"},
	{"青海省", "8663", "海北藏族自治州", "866322"},
	{"青海省", "8663", "黄南藏族自治州", "866323"},
	{"青海省", "8663", "海南藏族自治州", "866325"},
	{"青海省", "8663", "果洛藏族自治州", "866326"},
	{"青海省", "8663", "玉树藏族自治州", "866327"},
	{"青海省", "8663", "海西蒙古族藏族自治州", "866328"},
	{"宁夏回族自治区", "8664", "银川市", "866401"},
	{"宁夏回族自治区", "8664", "石嘴山市", "866402"},
	{"宁夏回族自治区", "8664", "吴忠市", "866403"},
	{"宁夏回族自治区", "8664", "固原市", "866404"},
	{"宁夏回族自治区", "8664", "中卫市", "866405"},
	{"新疆维吾尔自治区", "8665", "乌鲁木齐市", "866501"},
	{"新疆维吾尔自治区", "8665", "克拉玛依市", "866502"},
	{"新疆维吾尔自治区", "8665", "吐鲁番地区", "866521"},
	{"新疆维吾尔自治区", "8665", "哈密地区", "866522"},
	{"新疆维吾尔自治区", "8665", "昌吉回族自治州", "866523"},
	{"新疆维吾尔自治区", "8665", "博尔塔拉蒙古自治州", "866527"},
	{"新疆维吾尔自治区", "8665", "巴音郭楞蒙古自治州", "866528"},
	{"新疆维吾尔自治区", "8665", "阿克苏地区", "866529"},
	{"新疆维吾尔自治区", "8665", "克孜勒苏柯尔克孜自治州", "866530"},
	{"新疆维吾尔自治区", "8665", "喀什地区", "866531"},
	{"新疆维吾尔自治区", "8665", "和田地区", "866532"},
	{"新疆维吾尔自治区", "8665", "伊犁哈萨克自治州", "866540"},
	{"新疆维吾尔自治区", "8665", "塔城地区", "866542"},
	{"新疆维吾尔自治区", "8665", "阿勒泰地区", "866543"},
	{"新疆维吾尔自治区", "8665", "", "866590"},
	{"台湾省", "8671", "新竹市", "867101"},
	{"台湾省", "8671", "高雄市", "867102"},
	{"台湾省", "8671", "基隆市", "867103"},
	{"台湾省", "8671", "台中市", "867104"},
	{"台湾省", "8671", "台南市", "867105"},
	{"台湾省", "8671", "台北市", "867106"},
	{"台湾省", "8671", "桃园县", "867107"},
	{"台湾省", "8671", "云林县", "867108"}}

// GetDistrictByString get province and district by address string
func GetDistrictByString(addStr string) (string, string) {
	var provcode string
	var distcode string

	if "" != addStr {
		for _, s := range gDistrictBrief {
			if strings.Contains(addStr, s.Province) {
				provcode = s.Provcode
				break
			}
		}

		if "" != provcode {
			for _, s := range gDistrictBrief {
				if (provcode == s.Provcode) && (strings.Contains(addStr, s.District)) {
					distcode = s.Distcode
					break
				}
			}
		}
	}

	var province string
	var district string

	if "" != provcode {
		for _, s := range CommonDistrcitInfos {
			if provcode == s.Provcode {
				province = s.Province
				break
			}
		}
	}

	if "" != distcode {
		for _, s := range CommonDistrcitInfos {
			if distcode == s.Distcode {
				district = s.District
				break
			}
		}
	}

	if "" == province {
		tmp := strings.Split(addStr, " ")

		province = tmp[0]
		if len(tmp) > 1 {
			district = tmp[1]
		}
	}

	return province, district
}
