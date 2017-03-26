package vos

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestCreateAsyncTaskOk(t *testing.T) {
	id, err := CreateAsyncTask("abc", NewDummyAsncyTask(), 1)
	if err != nil {
		t.Errorf("CreateAsyncTask failed, err = %v", err)
		return
	}
	DestroyTask(id)

	_, _, err = FindTask("abc")
	if err == nil {
		t.Errorf("task is still exist after DestroyTask")
		return
	}
}

func TestCreateAsyncTaskWithSameName(t *testing.T) {
	id, err := CreateAsyncTask("abc", NewDummyAsncyTask(), 1)
	if err != nil {
		t.Errorf("CreateAsyncTask failed, err = %v", err)
		return
	}

	_, err = CreateAsyncTask("abc", NewDummyAsncyTask(), 1)
	if err == nil {
		t.Errorf("CreateAsyncTask should be failed")
		return
	}

	DestroyTask(id)

	_, _, err = FindTask("abc")
	if err == nil {
		t.Errorf("task is still exist after DestroyTask")
		return
	}
}

func TestCreateAndRunAsyncTaskOk(t *testing.T) {
	id, err := CreateAsyncTask("abc", NewDummyAsncyTask(), 1)
	if err != nil {
		t.Errorf("CreateAsyncTask failed, err = %v", err)
		return
	}

	err = StartTask(id)
	if err != nil {
		t.Errorf("StartTask failed, err = %v", err)
	}

	DestroyTask(id)

	time.Sleep(time.Second / 1000)

	_, _, err = FindTask("abc")
	if err == nil {
		t.Errorf("task is still exist after DestroyTask")
		return
	}
}

func TestCreateSyncTaskOk(t *testing.T) {
	task := NewDummySyncTask()

	id, err := CreateSyncTask(task.name, task, 1)
	if err != nil {
		t.Errorf("CreateSyncTask failed, err = %v", err)
		return
	}
	DestroyTask(id)

	_, _, err = FindTask(task.name)
	if err == nil {
		t.Errorf("task is still exist after DestroyTask")
		return
	}
}

func TestCreateSyncTaskWithSameName(t *testing.T) {
	task := NewDummySyncTask()

	id, err := CreateSyncTask(task.name, task, 1)
	if err != nil {
		t.Errorf("CreatesyncTask failed, err = %v", err)
		return
	}

	_, err = CreateSyncTask(task.name, task, 1)
	if err == nil {
		t.Errorf("CreateSyncTask should be failed")
		return
	}

	DestroyTask(id)

	_, _, err = FindTask(task.name)
	if err == nil {
		t.Errorf("task is still exist after DestroyTask")
		return
	}
}

func TestCreateAndRunSyncTaskOk(t *testing.T) {
	task := NewDummySyncTask()

	id, err := CreateSyncTask(task.name, task, 1)
	if err != nil {
		t.Errorf("CreateSyncTask failed, err = %v", err)
		return
	}

	err = StartTask(id)
	if err != nil {
		t.Errorf("StartTask failed, err = %v", err)
	}

	time.Sleep(time.Second / 1000)

	msgInfo := &MsgInfo{version: 1, MsgId: 1, Source: id, Dest: id}
	SendMsgToTask(msgInfo, nil)

	DestroyTask(id)

	time.Sleep(time.Second / 1000)

	_, _, err = FindTask(task.name)
	if err == nil {
		t.Errorf("task is still exist after DestroyTask")
		return
	}
}

func TestDestroyTaskWithWrongChecksum(t *testing.T) {
	id, err := CreateAsyncTask("abc", NewDummyAsncyTask(), 1)
	if err != nil {
		t.Errorf("CreateAsyncTask failed, err = %v", err)
		return
	}

	id2 := id
	id2.checksum = 10000

	err = DestroyTask(id2)
	if err == nil {
		t.Errorf("DestroyTask should be failed")
		return
	}

	DestroyTask(id)

	_, _, err = FindTask("abc")
	if err == nil {
		t.Errorf("task is still exist after DestroyTask")
		return
	}
}

type testClient struct {
}

type testServer struct {
}

const (
	EV_TEST_REQ uint32 = 1
	EV_TEST_RSP uint32 = 2
)

type testRequest struct {
	name string
}

type testResponse struct {
	id int
}

func (this *testClient) Run() {
	server, _, _ := FindTask("server")
	client, mailbox, _ := FindTask("client")

	msgInfo := &MsgInfo{version: 1, MsgId: EV_TEST_REQ, Source: client, Dest: server}
	request := &testRequest{name: "send request"}

	SendMsgToTask(msgInfo, request)
	atomic.AddUint64(&testStat.requestSend, 1)

	t1 := time.NewTimer(time.Second / 10)

	for {
		select {
		case msg := <-mailbox:
			msg1, _ := msg.(*Message)

			switch msg1.Msginfo.MsgId {
			case EV_TEST_RSP:
				_, ok := msg1.Payload.(*testResponse)
				if !ok {
					break
				}

				atomic.AddUint64(&testStat.responseRecv, 1)
				return
			}

		case <-t1.C:
			atomic.AddUint64(&testStat.requestTimeut, 1)
			fmt.Println("TestClient timeout")
			return
		}
	}

}

func (this *testServer) RecvMsg(msg *Message) {
	switch msg.Msginfo.MsgId {
	case EV_TEST_REQ:
		_, ok := msg.Payload.(*testRequest)
		if !ok {
			break
		}

		atomic.AddUint64(&testStat.requestRecv, 1)

		msgInfo := &MsgInfo{version: 1, MsgId: EV_TEST_RSP, Source: msg.Msginfo.Dest, Dest: msg.Msginfo.Source}
		response := &testResponse{id: 404}

		SendMsgToTask(msgInfo, response)
		atomic.AddUint64(&testStat.responseSend, 1)

		DestroyTask(msg.Msginfo.Dest) // this line to test deadlock
	}
}

type testEventStat struct {
	requestSend   uint64
	requestRecv   uint64
	requestTimeut uint64
	responseSend  uint64
	responseRecv  uint64
}

var testStat testEventStat

func (this *testEventStat) Clear() {
	this.requestSend = 0
	this.requestRecv = 0
	this.requestTimeut = 0
	this.responseSend = 0
	this.responseRecv = 0
}

func TestSendMsgToTaskOk(t *testing.T) {
	testStat.Clear()

	clientTask := &testClient{}
	serverTask := &testServer{}

	client, _ := CreateSyncTask("client", clientTask, 1)
	defer DestroyTask(client)

	server, _ := CreateAsyncTask("server", serverTask, 1)
	defer DestroyTask(server)

	err := StartTask(server)
	if err != nil {
		t.Errorf("TestSendMsgToTaskOk failed, StartTask \"server\" failed, err = %s", err.Error())
		return
	}

	err = StartTask(client)
	if err != nil {
		t.Errorf("TestSendMsgToTaskOk failed, StartTask \"client\" failed, err = %s", err.Error())
		return
	}

	time.Sleep(time.Second / 5)

	_, _, err = FindTask("client")
	if err == nil {
		t.Errorf("TestSendMsgToTaskOk failed, \"client\" task should have been destroyed")
		return
	}

	if testStat.requestSend != 1 {
		t.Errorf("TestSendMsgToTaskOk failed, wrong requestSend = %v", testStat.requestSend)
		return
	}

	if testStat.requestRecv != 1 {
		t.Errorf("TestSendMsgToTaskOk failed, wrong requestRecv = %v", testStat.requestRecv)
		return
	}

	if testStat.requestTimeut != 0 {
		t.Errorf("TestSendMsgToTaskOk failed, wrong requestTimeut = %v", testStat.requestTimeut)
		return
	}

	if testStat.responseSend != 1 {
		t.Errorf("TestSendMsgToTaskOk failed, wrong responseSend = %v", testStat.responseSend)
		return
	}

	if testStat.responseRecv != 1 {
		t.Errorf("TestSendMsgToTaskOk failed, wrong responseRecv = %v", testStat.responseRecv)
		return
	}
}
