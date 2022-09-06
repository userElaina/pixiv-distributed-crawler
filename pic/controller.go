package pic

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Queue struct {
	lock *sync.Mutex
	buf  []int
}

func (me *Queue) Push(c chan int) {
	for {
		x := <-c
		me.lock.Lock()
		me.buf = append(me.buf, x)
		me.lock.Unlock()
	}
}

func (me *Queue) Pop(c chan int) {
	for {
		for len(me.buf) < 1 {
			time.Sleep(time.Second)
		}
		c <- me.buf[0]
		me.lock.Lock()
		me.buf = me.buf[1:]
		me.lock.Unlock()
	}
}

type InfoMuti struct {
	flag         bool
	lock         *sync.Mutex
	maxThread    int
	maxClientUse int
	dir          string
	popchan      chan int
	waitchan     chan int
}

func NewInfoMuti(maxThread, maxClientUse int, dir string, ch chan int) *InfoMuti {
	q := &Queue{&sync.Mutex{}, []int{}}
	go q.Push(ch)
	popchan := make(chan int)
	go q.Pop(popchan)
	return &InfoMuti{false, &sync.Mutex{}, maxThread, maxClientUse, dir, popchan, make(chan int)}
}

func DefaultInfoMuti(dir string, ch chan int) *InfoMuti {
	return NewInfoMuti(runtime.NumCPU(), 16, dir, ch)
}

func (me *InfoMuti) aThread(n int) {
	fmt.Println(n, "BEGIN")
	for {
		client := &http.Client{}
		fmt.Println(n, "NEW CLIENT")
		for i := 0; i < me.maxClientUse; i++ {
			flg := false
			pid := 0

			me.lock.Lock()
			if me.flag {
				flg = true
			} else {
				pid = <-me.popchan
				if pid <= 0 {
					me.flag = true
					flg = true
				}
			}
			me.lock.Unlock()

			if flg {
				fmt.Println(n, "READY TO END")
				me.waitchan <- n
				// fmt.Println(n, "END")
				return
			}
			fmt.Println(n, "GET", pid)
			result, err := Crawl(pid, me.dir, client)
			if err != nil {
				fmt.Println(n, err)
			} else {
				fmt.Println(n, result)
			}
		}
	}
}

func (me *InfoMuti) MutiThread() {
	for i := 0; i < me.maxThread; i++ {
		go me.aThread(i)
	}
}

func (me *InfoMuti) Kill() {
	me.lock.Lock()
	me.flag = true
	me.lock.Unlock()
}

func (me *InfoMuti) Join() {
	for i := 0; i < me.maxThread; i++ {
		x := <- me.waitchan
		fmt.Println(x, "END")
	}
}
