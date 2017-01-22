package vos

import (
	"errors"
	"fmt"
)

const (
	VOS_EV_INIT uint32 = 1
)

func init() {
	g_taskPool = newTaskPool()
}

type MsgInfo struct {
	Version uint32
	MsgId   uint32
	Source  TaskId
	Dest    TaskId
}

type Message struct {
	Msginfo MsgInfo
	Payload interface{}
}

type TaskId struct {
	id       uint32
	checksum uint32
}

func (this *TaskId) Equal(other *TaskId) bool {
	return this.id == other.id && this.checksum == other.checksum
}

type AsyncTask interface {
	RecvMsg(msg *Message)
}

type SyncTask interface {
	Run()
}

type taskItem struct {
	name        string
	id          TaskId
	isSync      bool
	running     bool
	asyncTask   AsyncTask
	syncTask    SyncTask
	exit        chan bool
	mailboxSize uint32
	mailbox     chan interface{}
}

func newTaskItem(name string, id TaskId, mailboxSize uint32) *taskItem {
	return &taskItem{name: name, id: id, exit: make(chan bool), mailboxSize: mailboxSize, mailbox: make(chan interface{}, mailboxSize)}
}

func newAsyncTaskItem(name string, id TaskId, task AsyncTask, mailboxSize uint32) *taskItem {
	item := newTaskItem(name, id, mailboxSize)
	item.asyncTask = task
	return item
}

func newSyncTaskItem(name string, id TaskId, task SyncTask, mailboxSize uint32) *taskItem {
	item := newTaskItem(name, id, mailboxSize)
	item.isSync = true
	item.syncTask = task
	return item
}

func (this *taskItem) Start(pool *taskPool) {
	this.running = true

	if this.isSync {
		this.startSync(pool)
	} else {
		this.startAsync()
	}
}

func (this *taskItem) startAsync() {
	go func() {
		for {
			select {
			case msg := <-this.mailbox:
				switch data := msg.(type) {
				case *Message:
					this.asyncTask.RecvMsg(data)
				}
			case <-this.exit:
				this.running = false
				return
			}
		}
	}()
}

func (this *taskItem) startSync(pool *taskPool) {
	go func() {
		this.syncTask.Run()
		this.running = false
		pool.delTaskItem(this)
	}()
}

type taskPool struct {
	taskIdMap   map[uint32]*taskItem
	taskNameMap map[string]*taskItem
	idIndex     uint32
	checksum    uint32
}

func newTaskPool() *taskPool {
	return &taskPool{taskIdMap: make(map[uint32]*taskItem), taskNameMap: make(map[string]*taskItem)}
}

func (this *taskPool) AddAsyncTask(name string, task AsyncTask, mailboxSize uint32) (id TaskId, err error) {
	item, found := this.FindByName(name)
	if found {
		return id, errors.New(fmt.Sprintf("AddAsyncTask: \"%s\" is already exist", name))
	}

	id.id = this.idIndex
	id.checksum = this.checksum

	item = newAsyncTaskItem(name, id, task, mailboxSize)

	this.taskIdMap[id.id] = item
	this.taskNameMap[name] = item
	this.idIndex++
	this.checksum++

	return id, nil
}

func (this *taskPool) AddSyncTask(name string, task SyncTask, mailboxSize uint32) (id TaskId, err error) {
	item, found := this.FindByName(name)
	if found {
		return id, errors.New(fmt.Sprintf("AddSyncTask: \"%s\" is already exist", name))
	}

	id.id = this.idIndex
	id.checksum = this.checksum

	item = newSyncTaskItem(name, id, task, mailboxSize)

	this.taskIdMap[id.id] = item
	this.taskNameMap[name] = item
	this.idIndex++
	this.checksum++

	return id, nil
}

func (this *taskPool) DelById(id *TaskId) (err error) {
	item, found := this.FindById(id)
	if !found {
		return errors.New(fmt.Sprintf("DelById: %v is not exist", *id))
	}

	this.delTaskItem(item)
	return nil
}

func (this *taskPool) delTaskItem(item *taskItem) {
	if !item.isSync && item.running {
		go func() {
			item.exit <- true
			delete(this.taskIdMap, item.id.id)
			delete(this.taskNameMap, item.name)
		}()
		return
	}

	delete(this.taskIdMap, item.id.id)
	delete(this.taskNameMap, item.name)
}

func (this *taskPool) FindByName(name string) (item *taskItem, found bool) {
	item, ok := this.taskNameMap[name]
	if ok {
		return item, true
	}
	return nil, false
}

func (this *taskPool) FindById(id *TaskId) (item *taskItem, found bool) {
	item, ok := this.taskIdMap[id.id]
	if ok && id.Equal(&item.id) {
		return item, true
	}
	return nil, false
}

func (this *taskPool) StartTask(id TaskId) (err error) {
	item, found := this.FindById(&id)
	if !found {
		return errors.New(fmt.Sprintf("StartTask: %v is not exist", id))
	}

	item.Start(this)

	return nil
}

var g_taskPool *taskPool

func CreateAsyncTask(name string, task AsyncTask, mailboxSize uint32) (id TaskId, err error) {
	return g_taskPool.AddAsyncTask(name, task, mailboxSize)
}

func CreateSyncTask(name string, task SyncTask, mailboxSize uint32) (id TaskId, err error) {
	return g_taskPool.AddSyncTask(name, task, mailboxSize)
}

func DestroyTask(id TaskId) (err error) {
	return g_taskPool.DelById(&id)
}

func FindTask(name string) (id TaskId, mailbox chan interface{}, err error) {
	item, found := g_taskPool.FindByName(name)
	if !found {
		return id, nil, errors.New(fmt.Sprintf("FindTask: \"%s\" is not exist", name))
	}

	return item.id, item.mailbox, nil
}

func StartTask(id TaskId) (err error) {
	return g_taskPool.StartTask(id)
}

func SendMsgToTask(msgInfo *MsgInfo, data interface{}) (err error) {
	if msgInfo == nil {
		return errors.New("SendMsgToTask: msgInfo is nil")
	}

	item, found := g_taskPool.FindById(&msgInfo.Dest)
	if !found {
		return errors.New(fmt.Sprintf("SendMsgToTask: Dest = %v is not exist", msgInfo.Dest))
	}

	msg := &Message{*msgInfo, data}

	go func() {
		item.mailbox <- msg
	}()

	return nil
}

type DummyAsncyTask struct {
	RecvMsgNum uint32
}

func NewDummyAsncyTask() *DummyAsncyTask {
	return &DummyAsncyTask{}
}

func (this *DummyAsncyTask) RecvMsg(msg *Message) {
	this.RecvMsgNum++
}

type DummySyncTask struct {
	RecvMsgNum uint32
	name       string
}

func NewDummySyncTask() *DummySyncTask {
	return &DummySyncTask{name: "DummySyncTask"}
}

func (this *DummySyncTask) Run() {
	_, mailbox, err := FindTask(this.name)
	if err != nil {
		return
	}

	for {
		select {
		case <-mailbox:
			this.RecvMsgNum++
			return
		}
	}
}
