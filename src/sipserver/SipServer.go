package main

import (
	"fmt"
	//"logger"
	"sipparser"
	"sync/atomic"
	//"time"
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

var msg string = "INVITE sip:6135000@24.15.255.4 SIP/2.0\r\n" +
	//"Content-Length: 226\r\n" +
	"Content-Length: 0\r\n" +
	"Via: SIP/2.0/UDP 24.15.255.101:5060\r\n" +
	"From: \"User ID\" <sip:6140000@24.15.255.4>;tag=dab70900252036d7134be-4ec05abe\r\n" +
	"To: <sip:6135000@24.15.255.4>\r\n" +
	"Call-ID: 0009b7da-0352000f-30a69b83-0e7b53d6@24.15.255.101\r\n" +
	"CSeq: 101 INVITE\r\n" +
	//"Expires: 180\r\n" +
	//"User-Agent: Cisco-SIP-IP-Phone/2\r\n" +
	//"Accept: application/sdp\r\n" +
	"Contact: sip:6140000@24.15.255.101:5060\r\n" +
	"Content-Type: application/sdp\r\n" +
	"\r\n"

func main() {

	x := &[]byte{0, 1}

	fmt.Println("x[1] =", (*x)[1])

	y := []byte("12345")

	fmt.Println("x[1] =", &y[0:1])

	//fmt.Println("ABNF_SIP_HDR_TOTAL_NUM =", sipparser.ABNF_SIP_HDR_TOTAL_NUM)
	context := sipparser.NewParseContext()
	context.SetAllocator(sipparser.NewMemAllocator(1024 * 30))
	sipmsg := sipparser.NewSipMsg(context).GetSipMsg(context)
	remain := context.Used()
	msg1 := []byte(msg)

	for i := 0; i < 1; i++ {
		context.ClearAllocNum()
		context.FreePart(remain)
		_, err := sipmsg.Parse(context, msg1, 0)
		if err != nil {
			fmt.Println("parse sip msg failed, err =", err.Error())
			fmt.Println("msg1 = ", string(msg1))
			break
		} //*/
	}

	fmt.Printf("allocator.AllocNum = %d\n", context.GetAllocNum())
	fmt.Printf("allocator.Used = %d\n", context.Used())
	fmt.Printf("len(msg) = %d\n", len(msg))
	fmt.Printf("msg = \n%s\n", msg)

	/*
		fmt.Println("Count1 =", sipparser.Count1)
		fmt.Println("Count2 =", sipparser.Count2)
		fmt.Println("Count3 =", sipparser.Count3)
		fmt.Println("Count4 =", sipparser.Count4)
		//*/

	/*
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
	*/
}
