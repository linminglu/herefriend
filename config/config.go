package config

import (
	"os"
	"time"
)

var Conf_Driver = "mysql"
var Conf_AccessKey = "AtpDdb9Eh642X53CZM5KM7-ncvmgxPq2sFnlgcg5"
var Conf_SecretKey = "f-L1udoQwBf3wQiq-J-nnqX6UUhrZP6ZtYkcO6Ht"
var Toplist_Duration = int64(time.Hour / time.Second * 24 * 30 * 12)
var Conf_Dns = "root:Sancho8790@/bh_db"
var Conf_QiniuPre = "http://7xjwto.com1.z0.glb.clouddn.com/"
var Conf_QiniuScope = "herefriendpub"
var Conf_GeTuiAddr = "localhost:9090"
var Conf_AgeMin = 18
var Conf_AgeMax = 85

/*
 * 定期推送消息
 */
var Conf_EvaluationMsgContent = "好评送免费VIP两天哦，赶紧去评价一下吧。在您给予5星好评后，将评价内容发送给红娘，我们将人工验证，确认收到5星评论后，我们将会立即赠送您24小时VIP会员。"
var Conf_EnableEvaluation = true

/*
 * 欢迎信息
 */
var Conf_WelcomeMessage = `欢迎你加入寂寞交友！ 
我是您的专职红娘，有什么需求，可以给我留言，我会尽快给您答复； 
平台交友秘诀： 
上传头像————获得更多异性青睐； 
完善资料————答复提高交友成功率； 
主动搭讪————机会把握在自己手中； 
送礼物————最亮眼的勾搭(^o^)/ ； 
喜欢TA，就勾搭一下吧！`

func init() {
	if os.Getenv("DEBUG") == "1" {
		Conf_Dns = "bhuser:bhpasswd@/bh_db"
		Conf_QiniuPre = "http://7xjwip.com1.z0.glb.clouddn.com/"
		Conf_QiniuScope = "heretest"
		Conf_GeTuiAddr = "192.168.185.141:9090"
	}
}
