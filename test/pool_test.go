package test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Task struct {
	f func() error //具体业务逻辑
}

func NewTask(funcArg func() error) *Task {
	return &Task{f: funcArg}
}

type Pool struct {
	RunningWorkers int64      //运行着的worker数量
	Capacity       int64      //协程池worker容量 也就是开几个协程处理
	JobCh          chan *Task //worker任务
	sync.Mutex
}

func NewPool(capacity int64, taskNum int) *Pool {
	return &Pool{
		Capacity: capacity,
		JobCh:    make(chan *Task, taskNum),
	}
}

func (receiver *Pool) GetCap() int64 {
	return receiver.Capacity
}

// 运行中的任务数+1
func (receiver *Pool) incRunning() {
	atomic.AddInt64(&receiver.RunningWorkers, 1)
}

// 结束运行
func (receiver *Pool) decRunning() {
	atomic.AddInt64(&receiver.RunningWorkers, -1)
}

// 获取运行中的任务数
func (receiver *Pool) GetRunningWorkers() int64 {
	return atomic.LoadInt64(&receiver.RunningWorkers)
}

func (receiver *Pool) run() {
	receiver.incRunning()
	go func() {
		defer func() {
			receiver.decRunning()
		}()
		//这边会一直从chan里面取任务，然后chan是安全的
		for task := range receiver.JobCh {
			task.f()
		}
	}()
}

func (receiver *Pool) AddTask(task *Task) {
	receiver.Lock()
	defer receiver.Unlock()
	if receiver.GetRunningWorkers() < receiver.GetCap() {
		receiver.run()
	}
	receiver.JobCh <- task
}

func TestPool(t *testing.T) {
	pool := NewPool(3, 10)
	for i := 0; i < 20; i++ {
		b := i
		fmt.Println(i)
		pool.AddTask(NewTask(func() error {
			fmt.Printf("I am Task %d \n", b)
			time.Sleep(1 * time.Second)
			return nil
		}))
	}
	time.Sleep(50 * time.Second)
}
