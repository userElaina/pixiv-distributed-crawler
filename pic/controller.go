package pic

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	flagLock     *sync.Mutex
	db           *datas
	maxThread    int
	maxClientUse int
	retry        int
	dir          string
	popchan      chan int
	waitchan     chan int
}

func NewMuti(maxThread, maxClientUse, retry int, dir string, ch chan int) *InfoMuti {
	os.MkdirAll(dir, 0644)

	db := newDatas(filepath.Join(dir, "info.db"))
	e := db.initDB()
	if e != nil {
		fmt.Println(e)
	}

	q := &Queue{&sync.Mutex{}, []int{}}
	go q.Push(ch)
	popchan := make(chan int)
	go q.Pop(popchan)

	return &InfoMuti{false, &sync.Mutex{}, db, maxThread, maxClientUse, retry, dir, popchan, make(chan int)}
}

func DefaultMuti(dir string, ch chan int) *InfoMuti {
	return NewMuti(runtime.NumCPU(), 16, 5, dir, ch)
}

func (me *InfoMuti) aThread(n int) {
	fmt.Println(n, "BEGIN")
	for {
		client := &http.Client{}
		fmt.Println(n, "NEW CLIENT")
		for i := 0; i < me.maxClientUse; i++ {
			flg := false
			pid := 0

			me.flagLock.Lock()
			if me.flag {
				flg = true
			} else {
				pid = <-me.popchan
				if pid <= 0 {
					me.flag = true
					flg = true
				}
			}
			me.flagLock.Unlock()

			if flg {
				fmt.Println(n, "READY TO END")
				me.waitchan <- n
				return
			}
			fmt.Println(n, "GET", pid)
			result, e := Crawl(pid, me.retry, me.dir, client)
			if e != nil {
				fmt.Println(n, "ERR", e)
			} else {
				fmt.Println(n, "SUCC", result)
				// me.dbLock.Lock()
				e = me.db.loadInfoPic(&result)
				// me.dbLock.Unlock()
				if e != nil {
					fmt.Println(n, "ERR", e)
				}
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
	me.flagLock.Lock()
	me.flag = true
	me.flagLock.Unlock()
}

func (me *InfoMuti) Join() {
	for i := 0; i < me.maxThread; i++ {
		x := <-me.waitchan
		fmt.Println(x, "END")
	}
	me.db.Close()
}
