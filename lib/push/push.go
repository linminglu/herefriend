package push

import (
	"fmt"
	"strings"
	"sync"

	"herefriend/config"
	"herefriend/lib/push/pushmsg"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	gPushChanSize = 1
)

const (
	// PushMsgGreet is greet msg type
	PushMsgGreet = 0
	// PushMsgComment is comment msg type
	PushMsgComment = 1
	// PushMsgVisit is visit msg type
	PushMsgVisit = 2
	// PushMsgNotify is notify msg type
	PushMsgNotify = 3
)

/*
 * 透明推送的子消息类型
 */
const (
	// NotifyMsgInvalid .
	NotifyMsgInvalid = 0
	// NotifyMsgUnRead .
	NotifyMsgUnRead = 1
	// NotifyMsgEvaluation .
	NotifyMsgEvaluation = 2
	// NotifyMsgRecvGift .
	NotifyMsgRecvGift = 3
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

// InitPush init resources for push
func InitPush() {
	var err error
	gPushSock, err = thrift.NewTSocket(config.ConfGeTuiAddr)
	if nil != err {
		panic(err.Error())
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	gPushClient = pushmsg.NewPushMsgClientFactory(gPushSock, protocolFactory)
	if nil == gPushClient {
		panic("create push client failed")
	}

	gPushMap = make(map[string]*pushMsgInfo)
	gPushChan = make(chan *pushMsgInfo, gPushChanSize)
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

// Add push new msg to buffer for later doPush
func Add(badge int, cid string, pushtype, subtype int, title, content string) {
	if "" == cid {
		return
	}

	key := func() string {
		switch pushtype {
		case PushMsgNotify:
			return fmt.Sprintf("%s_%d_%d", cid, pushtype, subtype)
		default:
			return fmt.Sprintf("%s_%d_%s_%s", cid, pushtype, title, content)
		}
	}()

	gPushLock.Lock()
	info, ok := gPushMap[key]
	if true == ok {
		if PushMsgNotify == pushtype {
			info.content = content
		}
	} else {
		gPushMap[key] = &pushMsgInfo{msgtype: pushtype, badge: badge, cid: cid, title: title, content: content}
	}
	gPushLock.Unlock()
}

// DoPush active the push action
func DoPush() {
	gPushActive <- 1
}

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

// GetPushNum returns the push numbers
func GetPushNum() int {
	return gPushCount
}
