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

const (
	PUSHMSG_TYPE_GREET     = 0
	PUSHMSG_TYPE_RECOMMEND = 1
	PUSHMSG_TYPE_VISIT     = 2
	PUSHMSG_TYPE_NOTIFYMSG = 3
)

const (
	PUSH_NOTIFYMSG_INVALID = 0
	PUSH_NOTIFYMSG_UNREAD  = 1
)

type pushMsgInfo struct {
	msgtype int
	badge   int
	cid     string
	title   string
	content string
}

var g_pushSock thrift.TTransport
var g_pushclient *pushmsg.PushMsgClient

var g_pushlock sync.Mutex
var g_pushchan chan *pushMsgInfo
var g_pushactive chan int
var g_pushmap map[string]*pushMsgInfo
var g_CountPush int

func InitPush() {
	var err error
	g_pushSock, err = thrift.NewTSocket(config.Conf_GeTuiAddr)
	if nil != err {
		panic(err.Error())
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	g_pushclient = pushmsg.NewPushMsgClientFactory(g_pushSock, protocolFactory)
	if nil == g_pushclient {
		panic("create push client failed")
	}

	g_pushmap = make(map[string]*pushMsgInfo)
	g_pushchan = make(chan *pushMsgInfo, PUSHMSG_CHAN_SIZE)
	g_pushactive = make(chan int)

	go pushRoutine()
	go pushWorkRoutine()
}

func connectServer() bool {
	err := g_pushSock.Open()
	if nil != err {
		g_pushSock.Close()
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

	g_pushlock.Lock()
	info, ok := g_pushmap[key]
	if true == ok {
		if PUSHMSG_TYPE_NOTIFYMSG == pushtype {
			info.content = content
		}
	} else {
		g_pushmap[key] = &pushMsgInfo{msgtype: pushtype, badge: badge, cid: cid, title: title, content: content}
	}
	g_pushlock.Unlock()
}

func DoPush() {
	g_pushactive <- 1
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
		<-g_pushactive

		for key, info := range g_pushmap {
			g_pushlock.Lock()
			delete(g_pushmap, key)
			g_pushlock.Unlock()

			g_pushchan <- info
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
	for info := range g_pushchan {
		if true != g_pushSock.IsOpen() {
			if true != connectServer() {
				fmt.Println("[GETUI] connect push server failed.")
				continue
			}
		}

		g_CountPush = g_CountPush + 1
		err := g_pushclient.Notify(int32(info.badge), info.cid, int32(info.msgtype), info.title, info.content)
		if nil != err {
			if strings.Contains(err.Error(), "broken pipe") {
				g_pushSock.Close()
			}
		}
	}
}

func GetPushNum() int {
	return g_CountPush
}
