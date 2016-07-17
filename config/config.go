package config

import "time"

const DEBUGGING = 0
const Conf_Driver = "mysql"
const Conf_AccessKey = "AtpDdb9Eh642X53CZM5KM7-ncvmgxPq2sFnlgcg5"
const Conf_SecretKey = "f-L1udoQwBf3wQiq-J-nnqX6UUhrZP6ZtYkcO6Ht"
const Toplist_Duration = int64(time.Hour / time.Second * 24 * 30)

var Conf_Dns = []string{"root:Sancho8790@/bh_db", "bhuser:bhpasswd@/bh_db"}[DEBUGGING]
var Conf_QiniuPre = []string{"http://7xjwto.com1.z0.glb.clouddn.com/", "http://7xjwip.com1.z0.glb.clouddn.com/"}[DEBUGGING]
var Conf_QiniuScope = []string{"herefriendpub", "heretest"}[DEBUGGING]
var Conf_GeTuiAddr = []string{"localhost:9090", "192.168.185.141:9090"}[DEBUGGING]
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
