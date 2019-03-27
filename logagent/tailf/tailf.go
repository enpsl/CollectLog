package tailf

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

const (
	STATUS_NORMAL = 1
	STATUS_DELETE = -1
)

type Collect struct {
	LogPath string `json:"logPath"`
	Topic   string `json:"topic"`
}

type TextMsg struct {
	Msg   string
	Topic string
}

type TailObjList struct {
	tailItem []*TailObj
	msgChan  chan *TextMsg
	lock     sync.Mutex
}

type TailObj struct {
	tails    *tail.Tail
	conf     Collect
	status   int
	exitChan chan int
}

var (
	tailList *TailObjList
)

func GetOneLine() (msg *TextMsg) {
	msg = <-tailList.msgChan
	return
}

func UpdateConfig(conf []Collect) {
	tailList.lock.Lock()
	defer tailList.lock.Unlock()
	for _, v := range conf {
		var isRunning = false
		for _, obj := range tailList.tailItem {
			if v.LogPath == obj.conf.LogPath {
				isRunning = true
				break
			}
		}
		if isRunning {
			continue
		}
		createNewTask(v)
	}
	for _, obj := range tailList.tailItem {
		obj.status = STATUS_DELETE
		for _, v := range conf {
			if v.LogPath == obj.conf.LogPath {
				obj.status = STATUS_NORMAL
				break
			}
		}
		if obj.status == STATUS_DELETE {
			obj.exitChan <- 1
			continue
		}
	}
	return
}

func createNewTask(collect Collect) {
	obj := &TailObj{
		conf:     collect,
		exitChan: make(chan int, 1),
	}
	tail_obj, errTail := tail.TailFile(collect.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,
		//Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if errTail != nil {
		fmt.Println("tail file err:", errTail)
		return
	}
	obj.tails = tail_obj
	tailList.tailItem = append(tailList.tailItem, obj)
	go ReadFromTail(obj)
}

func InitTailf(collect []Collect, chanSize int) (err error) {
	tailList = &TailObjList{
		msgChan: make(chan *TextMsg, chanSize),
	}
	if len(collect) == 0 {
		logs.Error("no conf to collect: %+v", collect)
		return
	}
	for _, v := range collect {
		createNewTask(v)
	}
	return
}

func ReadFromTail(tailObj *TailObj) {
	for true {
		select {
		case line, ok := <-tailObj.tails.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename:%s", tailObj.tails.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			TextMsg := &TextMsg{
				Msg:   line.Text,
				Topic: tailObj.conf.Topic,
			}
			tailList.msgChan <- TextMsg
		case <-tailObj.exitChan:
			fmt.Println("tail obj will exited, conf:%v", tailObj.conf)
			return
		}
	}
}
