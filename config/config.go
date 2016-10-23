package config

import (
	"os"
	"time"

	log "qiniupkg.com/x/log.v7"
)

var (
	// ConfDBDriver .
	ConfDBDriver = "mysql"
	// ConfDBDns .
	ConfDBDns = "root:Sancho8790@/bh_db"
	// ConfQiniuAccessKey .
	ConfQiniuAccessKey = "AtpDdb9Eh642X53CZM5KM7-ncvmgxPq2sFnlgcg5"
	// ConfQiniuSecretKey .
	ConfQiniuSecretKey = "f-L1udoQwBf3wQiq-J-nnqX6UUhrZP6ZtYkcO6Ht"
	// ConfQiniuPre .
	ConfQiniuPre = "http://7xjwto.com1.z0.glb.clouddn.com/"
	// ConfQiniuScope .
	ConfQiniuScope = "herefriendpub"
	// ConfGeTuiAddr .
	ConfGeTuiAddr = "localhost:9090"
	// ConfAgeMin .
	ConfAgeMin = 18
	// ConfAgeMax .
	ConfAgeMax = 85
	// ConfToplistDuration .
	ConfToplistDuration = int64(time.Hour / time.Second * 24 * 30 * 12)
	// ConfEvaluationSwitch is the switch of following push
	ConfEvaluationSwitch = true
	// ConfEvaluationMsgContent will push to client automatically
	ConfEvaluationMsgContent = "好评送免费VIP两天哦，赶紧去评价一下吧。在您给予5星好评后，将评价内容发送给红娘，我们将人工验证，确认收到5星评论后，我们将会立即赠送您24小时VIP会员。"
	// ConfWelcomeMessage is the welcome message to return to client
	ConfWelcomeMessage = `欢迎你加入寂寞交友！ 
我是您的专职红娘，有什么需求，可以给我留言，我会尽快给您答复； 
平台交友秘诀： 
上传头像————获得更多异性青睐； 
完善资料————答复提高交友成功率； 
主动搭讪————机会把握在自己手中； 
送礼物————最亮眼的勾搭(^o^)/ ； 
喜欢TA，就勾搭一下吧！`
)

func init() {
	if os.Getenv("DEBUG") == "1" {
		log.Info("Now we start by DEBUG mode")

		ConfDBDns = "bhuser:bhpasswd@/bh_db"
		ConfQiniuPre = "http://7xjwip.com1.z0.glb.clouddn.com/"
		ConfQiniuScope = "heretest"
		ConfGeTuiAddr = "192.168.185.141:9090"
	}

	if os.Getenv("NOEVAL") == "1" {
		log.Info("Now we start without VIP evaluation enabled")

		ConfEvaluationSwitch = false
	}
}
