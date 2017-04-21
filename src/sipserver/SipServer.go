package main

import (
	"fmt"
	"logger"
	"sync/atomic"
	"time"
	"vos"
)

type taskClient struct {
}

type taskServer struct {
}

const (
	EV_TEST_REQ uint32 = 1
	EV_TEST_RSP uint32 = 2
)

type testRequest struct {
	name int
}

type testResponse struct {
	id int
}

func (this *taskClient) Run() {
	server, _, _ := vos.FindTask("server")
	client, mailbox, _ := vos.FindTask("client")

	msgInfo := &vos.MsgInfo{Version: 1, MsgId: EV_TEST_REQ, Source: client, Dest: server}
	//request := &testRequest{name: "send request"}
	request := &testRequest{name: 1}

	vos.SendMsgToTask(msgInfo, request)
	atomic.AddUint64(&stat.requestSend, 1)

	for {
		select {
		case msg := <-mailbox:
			msg1, _ := msg.(*vos.Message)

			switch msg1.Msginfo.MsgId {
			case EV_TEST_RSP:
				_, ok := msg1.Payload.(*testResponse)
				if !ok {
					break
				}

				atomic.AddUint64(&stat.responseRecv, 1)

				vos.SendMsgToTask(msgInfo, request)
				atomic.AddUint64(&stat.requestSend, 1)
			}
		}
	}

}

func (this *taskServer) RecvMsg(msg *vos.Message) {
	switch msg.Msginfo.MsgId {
	case EV_TEST_REQ:
		_, ok := msg.Payload.(*testRequest)
		if !ok {
			break
		}

		atomic.AddUint64(&stat.requestRecv, 1)

		msgInfo := &vos.MsgInfo{Version: 1, MsgId: EV_TEST_RSP, Source: msg.Msginfo.Dest, Dest: msg.Msginfo.Source}
		response := &testResponse{id: 404}

		vos.SendMsgToTask(msgInfo, response)
		atomic.AddUint64(&stat.responseSend, 1)
	}
}

type eventStat struct {
	requestSend   uint64
	requestRecv   uint64
	requestTimeut uint64
	responseSend  uint64
	responseRecv  uint64
}

var stat eventStat

func (this *eventStat) Clear() {
	this.requestSend = 0
	this.requestRecv = 0
	this.requestTimeut = 0
	this.responseSend = 0
	this.responseRecv = 0
}

func (this *eventStat) Print() {
	fmt.Println("requestSend   =", this.requestSend)
	fmt.Println("requestRecv   =", this.requestRecv)
	fmt.Println("requestTimeut =", this.requestTimeut)
	fmt.Println("responseSend  =", this.responseSend)
	fmt.Println("responseRecv  =", this.responseRecv)

}

func main() {
	logger.SetLevel(logger.DEBUG)
	logger.Emergency("test EMERGENCY")
	logger.Alert("test ALERT")
	logger.Critical("test CRITICAL")

	logger.Error("test ERROR")
	logger.Warning("test WARNING")
	logger.Notice("test NOTICE")
	logger.Info("test INFO")
	logger.Debug("test DEBUG")

	logger.Print("test Print")
	logger.Print("test PrintStack")
	logger.PrintStack()
	return

	stat.Clear()

	clientTask := &taskClient{}
	serverTask := &taskServer{}

	client, _ := vos.CreateSyncTask("client", clientTask, 1)
	defer vos.DestroyTask(client)

	server, _ := vos.CreateAsyncTask("server", serverTask, 1)
	defer vos.DestroyTask(server)

	vos.StartTask(server)
	vos.StartTask(client)

	var seconds uint64 = 3

	time.Sleep(time.Second * time.Duration(seconds))

	stat.Print()
}
