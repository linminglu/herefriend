package config

const DEBUGGING = 0
const Conf_Driver = "mysql"
const Conf_AccessKey = "AtpDdb9Eh642X53CZM5KM7-ncvmgxPq2sFnlgcg5"
const Conf_SecretKey = "f-L1udoQwBf3wQiq-J-nnqX6UUhrZP6ZtYkcO6Ht"

var Conf_Dns = []string{"root:Sancho8790@/bh_db", "bhuser:bhpasswd@/bh_db"}[DEBUGGING]
var Conf_QiniuPre = []string{"http://7xjwto.com1.z0.glb.clouddn.com/", "http://7xjwip.com1.z0.glb.clouddn.com/"}[DEBUGGING]
var Conf_QiniuScope = []string{"herefriendpub", "heretest"}[DEBUGGING]
var Conf_GeTuiAddr = []string{"localhost:9090", "192.168.185.141:9090"}[DEBUGGING]
var Conf_AgeMin = 18
var Conf_AgeMax = 85

/*
 * 定期推送消息
 */
var Conf_EvaluationMsgContent = "好评送免费VIP 3个月哦，赶紧去评价一下吧。"
var Conf_EnableEvaluation = false
