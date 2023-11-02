package _101

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	handle1 := func() error {
		panic("1111")
		return nil
	}
	handle2 := func() error {
		panic("22222")
		return nil
	}
	err := withGoroutine(handle1, handle2)
	if err != nil {
		fmt.Printf("recover22:%v\n", err)
	}
}

func withGoroutine(opts ...func() error) (err error) {
	var wg sync.WaitGroup
	for _, opt := range opts {
		wg.Add(1)
		// 开启goroutine，做并行处理
		go func(handler func() error) {
			defer func() { // 协程内部捕获panic
				if e := recover(); e != nil {
					fmt.Printf("recover:%v\n", e)
				}
				wg.Done()
			}()

			e := handler() // 真正的逻辑调用
			// 取第一个报错的handler调用的错误返回
			if err == nil && e != nil {
				err = e
			}
		}(opt) // 将goroutine的函数逻辑通过封装成的函数变量传入
	}
	wg.Wait()
	return
}

// 无缓冲适合协程内写，协程外等待数据，接收数据，这样才不会卡主
// 有缓冲不阻塞，多协程可以，可能这个协程写，那个协程释放
func TestCh(t *testing.T) {
	ch := make(chan bool, 1)
	var num int
	for i := 1; i <= 100; i++ {
		go add(ch, &num)
	}
	time.Sleep(2 * time.Second)
	fmt.Println(num)
}

func add(ch chan bool, num *int) {
	ch <- true
	*num = *num + 1
	<-ch
}
