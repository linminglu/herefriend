package push

import (
	"fmt"
	"strings"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
	"herefriend/config"
	"herefriend/lib/push/pushmsg"
)

const (
	PUSHMSG_CHAN_SIZE = 1
)

/*
 * 推送消息类型, getuiserver将采取不同的推送方式
 */
const (
	PUSHMSG_TYPE_GREET     = 0 //普通推送
	PUSHMSG_TYPE_RECOMMEND = 1 //普通推送
	PUSHMSG_TYPE_VISIT     = 2 //普通推送
	PUSHMSG_TYPE_NOTIFYMSG = 3 //透明推送,无弹出框
)

/*
 * 透明推送的子消息类型
 */
const (
	PUSH_NOTIFYMSG_INVALID    = 0
	PUSH_NOTIFYMSG_UNREAD     = 1
	PUSH_NOTIFYMSG_EVALUATION = 2
	PUSH_NOTIFYMSG_RECVGIFT   = 3
)

type pushMsgInfo struct {
	msgtype int
	badge   int
	cid     string
	title   string
	content string
}

var gPushSock thrift.TTransport
var gPushClient *pushmsg.PushMsgClient

var gPushLock sync.Mutex
var gPushChan chan *pushMsgInfo
var gPushActive chan int
var gPushMap map[string]*pushMsgInfo
var gPushCount int

func InitPush() {
	var err error
	gPushSock, err = thrift.NewTSocket(config.Conf_GeTuiAddr)
	if nil != err {
		panic(err.Error())
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	gPushClient = pushmsg.NewPushMsgClientFactory(gPushSock, protocolFactory)
	if nil == gPushClient {
		panic("create push client failed")
	}

	gPushMap = make(map[string]*pushMsgInfo)
	gPushChan = make(chan *pushMsgInfo, PUSHMSG_CHAN_SIZE)
	gPushActive = make(chan int)

	go pushRoutine()
	go pushWorkRoutine()
}

func connectServer() bool {
	err := gPushSock.Open()
	if nil != err {
		gPushSock.Close()
		return false
	}

	return true
}

func Add(badge int, cid string, pushtype, subtype int, title, content string) {
	if "" == cid {
		return
	}

	key := func() string {
		switch pushtype {
		case PUSHMSG_TYPE_NOTIFYMSG:
			return fmt.Sprintf("%s_%d_%d", cid, pushtype, subtype)
		default:
			return fmt.Sprintf("%s_%d_%s_%s", cid, pushtype, title, content)
		}
	}()

	gPushLock.Lock()
	info, ok := gPushMap[key]
	if true == ok {
		if PUSHMSG_TYPE_NOTIFYMSG == pushtype {
			info.content = content
		}
	} else {
		gPushMap[key] = &pushMsgInfo{msgtype: pushtype, badge: badge, cid: cid, title: title, content: content}
	}
	gPushLock.Unlock()
}

func DoPush() {
	gPushActive <- 1
}

/*
 |    Function: pushRoutine
 |      Author: Mr.Sancho
 |        Date: 2016-01-05
 |   Arguments:
 |      Return:
 | Description: 负责将要push的信息放入工作队列的routine
 |
*/
func pushRoutine() {
	for {
		<-gPushActive

		for key, info := range gPushMap {
			gPushLock.Lock()
			delete(gPushMap, key)
			gPushLock.Unlock()

			gPushChan <- info
		}
	}
}

/*
 |    Function: pushWorkRoutine
 |      Author: Mr.Sancho
 |        Date: 2016-01-06
 |   Arguments:
 |      Return:
 | Description: push的工作线程
 |
*/
func pushWorkRoutine() {
	for info := range gPushChan {
		if true != gPushSock.IsOpen() {
			if true != connectServer() {
				fmt.Println("[GETUI] connect push server failed.")
				continue
			}
		}

		gPushCount = gPushCount + 1
		err := gPushClient.Notify(int32(info.badge), info.cid, int32(info.msgtype), info.title, info.content)
		if nil != err {
			if strings.Contains(err.Error(), "broken pipe") {
				gPushSock.Close()
			}
		}
	}
}

func GetPushNum() int {
	return gPushCount
}
